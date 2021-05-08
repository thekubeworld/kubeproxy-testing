/*
Copyright 2021 The Kubernetes Authors.

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
	"github.com/thekubeworld/k8devel/pkg/deployment"
	"github.com/thekubeworld/k8devel/pkg/kubeproxy"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/util"
)

const domain = "example.com"

func main() {
	// Initial set
	fmt.Printf("Initializing kubeproxy-testing...\n")
	fmt.Printf("\t- selected ExternalName test...\n")
	randStr, err := util.GenerateRandomString(6, "lower")
	if err != nil {
		fmt.Println(err)
	}

	namespaceName := "kptesting"
	labelApp := "kptesting"

	c := client.Client{}
	c.Namespace = namespaceName
	c.NumberMaxOfAttemptsPerTask = 10
	c.TimeoutTaksInSec = 20
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
	NameDeployment := "kproxydeploy" + randStr

	// START: Namespace
	err = namespace.Create(&c, Namespace)
	if err != nil {
		fmt.Println("exiting... failed to create: ", err)
		os.Exit(1)
	}
	fmt.Printf("\t- namespace created %s...\n", Namespace)

	// START: Deployment
	d := deployment.Instance{
		Name:       NameDeployment,
		Namespace:  Namespace,
		Replicas:   1,
		LabelKey:   "app",
		LabelValue: labelApp,
	}

	d.Pod.Name = "nginx"
	d.Pod.Image = "nginx:1.14.2"
	d.Pod.ContainerPortName = "http"
	d.Pod.ContainerPortProtocol = "TCP"
	d.Pod.ContainerPort = 80

	err = deployment.Create(&c, &d)
	if err != nil {
		fmt.Println("exiting... failed to create: ", err)
		os.Exit(1)
	}
	fmt.Printf("\t- deployment created %s...\n", d.Name)
	// END: Deployment

	// START: Service
	s := service.Instance{
		Name:         NameService,
		Namespace:    Namespace,
		ExternalName: domain,
	}
	if err := service.CreateExternalName(&c, &s); err != nil {
		fmt.Printf("exiting... failed to create: ", err)
		os.Exit(1)
	}
	fmt.Printf("\t- service created %s...\n", s.Name)

	// START: Pod
	NameContainer := "kpnginx"
	p := pod.Instance{
		Name:       NameContainer,
		Namespace:  Namespace,
		Image:      "nginx",
		LabelKey:   "app",
		LabelValue: labelApp,
	}
	pod.Create(&c, &p)
	fmt.Printf("\t- pod created %s...\n", p.Name)

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

	// START: Execute curl from the pod created to the new service
	output, err := curl.ExecuteHTTPReqInsideContainer(
		&c,
		NameContainer,
		Namespace,
		fmt.Sprintf("http://%s", domain))
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", output)
	fmt.Printf("PASSED\n")
	// Delete namespace
	//namespace.Delete(&c, Namespace)
	fmt.Printf("Removed namespace %s...\n", Namespace)
}
