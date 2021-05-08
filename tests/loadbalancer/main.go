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
	"fmt"
	"os"
	"os/exec"

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/curl"
	"github.com/thekubeworld/k8devel/pkg/diagram"

	"github.com/thekubeworld/k8devel/pkg/kubeproxy"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/util"
)

func main() {

	args := []string{"apply", "-f", "https://raw.githubusercontent.com/metallb/metallb/v0.9.6/manifests/namespace.yaml"}
	cmd := exec.Command("kubectl", args...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("cannot create namespace for metallb\n")
		os.Exit(1)
	}

	args = []string{"apply", "-f", "https://raw.githubusercontent.com/metallb/metallb/v0.9.6/manifests/metallb.yaml"}
	cmd = exec.Command("kubectl", args...)
	_, err = cmd.Output()
	if err != nil {
		fmt.Printf("unable to apply metallb yaml\n")
		os.Exit(1)
	}

	args = []string{"create", "secret", "generic", "-n", "metallb-system", "memberlist", "--from-literal=secretkey=\"$(openssl rand -base64 128)\""}
	cmd = exec.Command("kubectl", args...)
	_, err = cmd.Output()
	if err != nil {
		fmt.Printf("unable to create secret..\n")
		os.Exit(1)
	}

	fmt.Printf("kube-proxy tests has started...\n\n")

	fmt.Printf("Test #3) User's Traffic reach loadbalancer that will\n")
	fmt.Printf("route using kubeproxy/iptables to the right service \n")
	fmt.Printf("that has the backend pod                            \n")
	diagram.LoadBalancer()

	// Initial set
	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	namespaceName := "kptesting"
	labelApp := "kptesting"

	c := client.Client{}
	c.Namespace = namespaceName
	c.TimeoutTaksInSec = 20

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
	c.Connect()

	// START: kube-proxy variables
	Namespace := c.Namespace + randStr
	NameService := "kproxysvc" + randStr

	// Saving current state of firewall in kube-proxy
	fwInitialState, err := kubeproxy.SaveCurrentFirewallState(
		&c,
		"kube-proxy",
		"kube-proxy",
		"kube-system")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// START: Namespace
	err = namespace.Create(&c, Namespace)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	// END: Namespace

	// START: Service
	// nodePort - a static port assigned on each the node
	// port - port exposed internally in the cluster
	// targetPort - the container port to send requests to
	s := service.Instance{
		Name:          NameService,
		Namespace:     Namespace,
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
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	// END: Service

	// Save firewall state after service, endpoint created
	fwAfterEndpointCreated, err := kubeproxy.SaveCurrentFirewallState(
		&c,
		"kube-proxy",
		"kube-proxy",
		"kube-system")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// See the difference with diff command
	out, err := util.DiffCommand(fwInitialState, fwAfterEndpointCreated)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s", string(out))
	os.Remove(fwInitialState)
	os.Remove(fwAfterEndpointCreated)

	// TODO: use library
	args = []string{"apply", "-f", "metallbcfg.yaml"}
	cmd = exec.Command("kubectl", args...)
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// START: Pod
	// Creating a POD Behind the service
	p := pod.Instance{
		Name:       "kpnginxbehindservice",
		Namespace:  Namespace,
		Image:      "nginx:1.14.2",
		LabelKey:   "app",
		LabelValue: labelApp,
	}
	fmt.Printf("\n")
	err = pod.Create(&c, &p)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	// END: Pod
	fmt.Printf("\n")

	// Creating a POD outside the service (No labels)
	// So it will try to connect to pod behind the service
	containerName := "nginxtoconnecttoservice"
	p = pod.Instance{
		Name:       containerName,
		Namespace:  Namespace,
		Image:      "nginx",
		LabelKey:   "app",
		LabelValue: "foobar",
	}
	err = pod.Create(&c, &p)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	// END: Pod
	// START: Execute curl from the pod created to the new service
	ret, err := curl.ExecuteHTTPReqInsideContainer(
		&c,
		containerName,
		Namespace,
		"172.17.255.1")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", ret)
	fmt.Printf("PASSED\n")
	// END: Execute curl from the pod created to the new service

	namespace.Delete(&c, namespaceName)
	namespace.Delete(&c, "metallb-system")
}
