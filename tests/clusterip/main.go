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
	"github.com/thekubeworld/k8devel/pkg/deployment"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/kubeproxy"
	"github.com/thekubeworld/k8devel/pkg/logschema"
	"github.com/thekubeworld/k8devel/pkg/iptables"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/curl"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/gookit/color"
)

func main() {
	logschema.SetLogrusLogging()
	logrus.Infof("kube-proxy tests has started...")

	logrus.Infof("\n")
	logrus.Infof("Test #1) Pod connect via kube-proxy to a service and Pod")
	logrus.Infof("\n")
	logrus.Infof("           POD                  ")
	logrus.Infof("            |                   ")
	logrus.Infof("         Traffic                ")
	logrus.Infof("            |                   ")
	logrus.Infof("            |                   ")
	logrus.Infof("   +-------------------+        ")
	logrus.Infof("   |     kube-proxy    |        ")
	logrus.Infof("   +-------------------+        ")
	logrus.Infof("       |           |            ")
	logrus.Infof("+------------------------------+")
	logrus.Infof("|     Service (Cluster IP)     |")
	logrus.Infof("|+-----------------------------+")
	logrus.Infof("|       |           |          |")
	logrus.Infof("|   +-------+   +-------+      |")
	logrus.Infof("|   |  Pod  |   |  Pod  |      |")
	logrus.Infof("|   +-------+   +-------+      |")
	logrus.Infof("|                              |")
	logrus.Infof("| kubernetes cluster           |")
	logrus.Infof("+------------------------------+")
	logrus.Infof("\n")

	// Initial set
	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		logrus.Fatal(err)
	}

	namespaceName := "kptesting" + randStr
	labelApp := "kptesting"

        c := client.Client{}
	c.Namespace = namespaceName
	c.NumberMaxOfAttemptsPerTask = 10
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
	logrus.Infof("\t\tnamespace: %s", namespaceKP)
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

	// TODO: use flags
	// START: kube-proxy variables
	KPTestContainerName := kyPods[0]
	KPTestNamespaceName := c.Namespace

	randStr, err = util.GenerateRandomString(6, "lower")
	if err != nil {
		logrus.Fatal(err)
	}
	KPTestServiceName := KPTestNamespaceName +
			"service" +
			randStr

	KPTestNginxDeploymentName := KPTestNamespaceName +
			"nginxdeployment" +
			randStr
	// END: kube-proxy variables

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

	// START: Deployment
	d := deployment.Instance {
		Name: KPTestNginxDeploymentName,
		Namespace: KPTestNamespaceName,
		Replicas: 1,
		LabelKey: "app",
		LabelValue: labelApp,
	}

	d.Pod.Name = "nginx"
	d.Pod.Image = "nginx:1.14.2"
	d.Pod.ContainerPortName = "http"
	d.Pod.ContainerPortProtocol = "TCP"
	d.Pod.ContainerPort = 80

	logrus.Infof("\n")
	err = deployment.Create(&c, &d)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	// END: Deployment

	// START: Service
	s := service.Instance {
		Name: KPTestServiceName,
		Namespace: KPTestNamespaceName,
		LabelKey: "app",
		LabelValue: labelApp,
		SelectorKey: "app",
		SelectorValue: labelApp,
		ClusterIP: "",
		Port: 80,
	}
	err = service.CreateClusterIP(&c, &s)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}

	IPService, err := service.GetIP(
		&c,
		KPTestServiceName,
		KPTestNamespaceName)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	// END: Service

	iptablesStateAfterEndpointCreated, err := iptables.Save(
				&c,
				&iptablesCmd,
				KPTestContainerName,
				"kube-system")
        if err != nil {
		logrus.Fatal(err)
	}

	out, err := util.DiffCommand(iptablesInitialState.Name(),
			iptablesStateAfterEndpointCreated.Name())
	if err != nil {
		logrus.Fatal(err)
	}

	if kpMode == "iptables" {
		if len(string(out)) > 0 {
			logrus.Infof("%s", string(out))
		}
	}
		// END: iptables diff

	// START: Pod
	// PodCommandInitBash struct for running bash command
	containerName := "kptestingnginx"
	p := pod.Instance {
		Name: containerName,
		Namespace: KPTestNamespaceName,
		Image: "nginx",
		LabelKey: "app",
                LabelValue: labelApp,
	}

	logrus.Infof("\n")
	pod.Create(&c, &p)
	// END: Pod

	// START: Execute curl from the pod created to the new service
	ret, err := curl.ExecuteHTTPReqInsideContainer(
			&c,
			containerName,
			KPTestNamespaceName,
			"http://" + IPService)
        if err != nil {
		logrus.Fatal(err)
        }
	logrus.Infof("%s", ret)
	color.Green.Println("[Test #1 PASSED]")
	// END: Execute curl from the pod created to the new service

	namespace.Delete(&c, namespaceName)
}
