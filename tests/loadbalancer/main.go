/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"os/exec"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/curl"
	"github.com/thekubeworld/k8devel/pkg/diagram"
	"github.com/thekubeworld/k8devel/pkg/iptables"
	"github.com/thekubeworld/k8devel/pkg/kubeproxy"
	"github.com/thekubeworld/k8devel/pkg/logschema"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/util"
)

func main() {
	logschema.SetLogrusLogging()

	args := []string{"apply", "-f", "https://raw.githubusercontent.com/metallb/metallb/v0.9.6/manifests/namespace.yaml"}
	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		logrus.Infof("%s", err)
	}

	args = []string{"apply", "-f", "https://raw.githubusercontent.com/metallb/metallb/v0.9.6/manifests/metallb.yaml"}
	cmd = exec.Command("kubectl", args...)
	out, err = cmd.Output()
	if err != nil {
		logrus.Infof("%s", err)
	}

	args = []string{"create", "secret", "generic", "-n", "metallb-system", "memberlist", "--from-literal=secretkey=\"$(openssl rand -base64 128)\""}
	cmd = exec.Command("kubectl", args...)
	out, err = cmd.Output()
	if err != nil {
		logrus.Infof("%s", err)
	}

	logrus.Infof("kube-proxy tests has started...")

	logrus.Infof("\n")
	logrus.Infof("\n")
	logrus.Infof("Test #3) User's Traffic reach loadbalancer that will")
	logrus.Infof("route using kubeproxy/iptables to the right service ")
	logrus.Infof("that has the backend pod                            ")
	diagram.LoadBalancer()

	// Initial set
	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		logrus.Fatal(err)
	}

	namespaceName := "kptesting" + randStr
	labelApp := "kptesting"

	c := client.Client{}
	c.Namespace = namespaceName
	c.TimeoutTaksInSec = 20

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	// kube-proxy is a DaemonSet, it will replicate the pod
	// data to all nodes. Let's find one pod name to use and
	// collect data.
	KP := "kube-proxy"
	namespaceKP := "kube-system"
	kyPods, kyNumberPods := pod.FindPodsWithNameContains(&c,
		KP,
		namespaceKP)
	if kyNumberPods < 0 {
		logrus.Fatal("exiting... unable to find kube-proxy pod..")
	}
	logrus.Infof("Found the following kube-proxy pods:")
	logrus.Infof("\t\tNamespace: %s", namespaceKP)
	logrus.Infof("\t\t%s", kyPods)

	// Detect Kube-proxy mode
	kpMode, err := kubeproxy.DetectKubeProxyMode(&c,
		KP,
		namespaceKP)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("\n")
	logrus.Infof("Detected kube-proxy mode: %s", kpMode)

	// Setting ContainerName and Namespace
	KPTestContainerName := kyPods[0]
	KPTestNamespaceName := c.Namespace
	randStr, err = util.GenerateRandomString(6, "lower")
	if err != nil {
		logrus.Fatal(err)
	}

	// TODO: Just load the iptables commands if kube-proxy
	// is IPTABLES or return error
	// Loading some iptables
	iptablesCmd := iptables.LoadPreDefinedCommands()
	if err != nil {
		logrus.Fatal(err)
	}

	// iptables saving initial state
	iptablesInitialState, err := iptables.Save(
		&c,
		&iptablesCmd,
		KPTestContainerName,
		"kube-system")
	if err != nil {
		logrus.Fatal(err)
	}

	// START: Namespace
	_, err = namespace.Exists(&c,
		KPTestNamespaceName)
	if err != nil {
		err = namespace.Create(&c,
			KPTestNamespaceName)
		if err != nil {
			logrus.Fatal("exiting... failed to create: ", err)
		}
	}
	// END: Namespace

	// Setting Service Name
	KPTestServiceName := KPTestNamespaceName +
		"service" +
		randStr

	// START: Service
	// nodePort - a static port assigned on each the node
	// port - port exposed internally in the cluster
	// targetPort - the container port to send requests to
	s := service.Instance{
		Name:          KPTestServiceName,
		Namespace:     KPTestNamespaceName,
		LabelKey:      "app",
		LabelValue:    labelApp,
		SelectorKey:   "app",
		SelectorValue: labelApp,
		PortName:      "http",
		PortProtocol:  "TCP",
		Port:          80,
	}
	err = service.CreateLoadBalancer(&c, &s)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}

	// END: Service

	// START: iptables diff
	iptablesStateAfterEndpointCreated, err := iptables.Save(
		&c, &iptablesCmd, KPTestContainerName, "kube-system")
	if err != nil {
		logrus.Fatal(err)
	}

	// Make a diff between two states we collected from iptables
	out, err = util.DiffCommand(iptablesInitialState.Name(),
		iptablesStateAfterEndpointCreated.Name())
	if err != nil {
		logrus.Fatal(err)
	}

	if len(string(out)) > 0 {
		logrus.Infof("%s", string(out))
	}
	// END: iptables diff

	// TODO: use library
	args = []string{"apply", "-f", "metallbcfg.yaml"}
	cmd = exec.Command("kubectl", args...)
	out, err = cmd.Output()
	if err != nil {
		logrus.Infof("%s", err)
	}

	/*
		ExternalIPService, err := k8devel.GetExternalIPFromService(
	                &c,
	                KPTestServiceName,
	                KPTestNamespaceName)
	        if err != nil {
	                logrus.Fatal("exiting... failed to create: ", err)
	        }*/

	// START: Pod

	// Creating a POD Behind the service
	p := pod.Instance{
		Name:       "kpnginxbehindservice",
		Namespace:  KPTestNamespaceName,
		Image:      "nginx:1.14.2",
		LabelKey:   "app",
		LabelValue: labelApp,
	}
	logrus.Infof("\n")
	err = pod.Create(&c, &p)
	if err != nil {
		logrus.Fatal(err)
	}
	// END: Pod
	logrus.Info("\n")

	// Creating a POD outside the service (No labels)
	// So it will try to connect to pod behind the service
	containerName := "nginxtoconnecttoservice"
	p = pod.Instance{
		Name:       containerName,
		Namespace:  KPTestNamespaceName,
		Image:      "nginx",
		LabelKey:   "app",
		LabelValue: "foobar",
	}
	err = pod.Create(&c, &p)
	if err != nil {
		logrus.Fatal(err)
	}
	// END: Pod
	// START: Execute curl from the pod created to the new service
	ret, err := curl.ExecuteHTTPReqInsideContainer(
		&c,
		containerName,
		KPTestNamespaceName,
		"172.17.255.1")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("%s", ret)
	color.Green.Println("[Test #2 PASSED]")
	// END: Execute curl from the pod created to the new service

	namespace.Delete(&c, namespaceName)
}
