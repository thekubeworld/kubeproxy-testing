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

	"github.com/thekubeworld/k8devel/pkg/client"
	"github.com/thekubeworld/k8devel/pkg/curl"
	"github.com/thekubeworld/k8devel/pkg/diagram"
	"github.com/thekubeworld/k8devel/pkg/kubeproxy"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/node"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/util"
)

func main() {
	fmt.Printf("kube-proxy tests has started...\n")
	diagram.NodePort()

	// Initial set
	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	namespaceName := "kptesting" + randStr
	labelApp := "kptesting"

	c := client.Client{}
	c.Namespace = namespaceName
	c.TimeoutTaskInSec = 20

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
	c.Connect()

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

	// START: kube-proxy variables
	Namespace := c.Namespace + randStr
	NameService := "kproxysvc" + randStr

	// START: Namespace
	err = namespace.Create(&c, Namespace)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
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
		Port:          80,    // service
		TargetPort:    80,    // container
		NodePort:      30001, // node
	}
	err = service.CreateNodePort(&c, &s)
	if err != nil {
		fmt.Printf("exiting... failed to create: %s\n", err)
		os.Exit(1)
	}

	IPNodes, err := node.GetIPFromNodes(&c)
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

	// START: Pod

	// Creating a POD Behind the service
	p := pod.Instance{
		Name:       "kpnginxbehindservice",
		Namespace:  Namespace,
		Image:      "nginx:1.14.2",
		LabelKey:   "app",
		LabelValue: labelApp,
	}
	err = pod.Create(&c, &p)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Pod behind the service created...\n")
	// END: Pod

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
	fmt.Printf("Pod to connect to service created...\n")
	// END: Pod

	// START: Execute curl from the pod created to the new service
	ret, err := curl.ExecuteHTTPReqInsideContainer(
		&c,
		containerName,
		Namespace,
		IPNodes[0]+":"+fmt.Sprint(s.NodePort))
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", ret)
	fmt.Printf("PASSED\n")

	namespace.Delete(&c, Namespace)
	fmt.Printf("namespace deleted: %s\n", Namespace)
}
