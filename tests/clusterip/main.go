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
	"github.com/thekubeworld/k8devel/pkg/deployment"
	"github.com/thekubeworld/k8devel/pkg/service"
	"github.com/thekubeworld/k8devel/pkg/client"
	//"github.com/thekubeworld/k8devel/pkg/kubeproxy"
	"github.com/thekubeworld/k8devel/pkg/pod"
	"github.com/thekubeworld/k8devel/pkg/curl"
	"github.com/thekubeworld/k8devel/pkg/namespace"
	"github.com/thekubeworld/k8devel/pkg/util"
)

func main() {
	// Initial set
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

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
        c.Connect()

	// Saving current state of firewall in kube-proxy
	//fwInitialState, err := kubeproxy.SaveCurrentFirewallState(
	//	&c,
	//	"kube-proxy",
	//	"kube-proxy",
	//	"kube-system")
	//if err != nil {
	//	fmt.Println(err)
	//}

	// START: kube-proxy variables
	randStr, err = util.GenerateRandomString(6, "lower")
	if err != nil {
		fmt.Println(err)
	}

	Namespace := c.Namespace + randStr
	NameService := "kproxysvc" + randStr
	NameDeployment := "kproxydeploy" + randStr

	// START: Namespace
	err = namespace.Create(&c, Namespace)
	if err != nil {
		fmt.Println(
			"exiting... failed to create: ",
			err)
	}

	// START: Deployment
	d := deployment.Instance {
		Name: NameDeployment,
		Namespace: Namespace,
		Replicas: 1,
		LabelKey: "app",
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
	}
	// END: Deployment

	// START: Service
	s := service.Instance {
		Name: NameService,
		Namespace: Namespace,
		LabelKey: "app",
		LabelValue: labelApp,
		SelectorKey: "app",
		SelectorValue: labelApp,
		ClusterIP: "",
		Port: 80,
	}
	err = service.CreateClusterIP(&c, &s)
	if err != nil {
		fmt.Println("exiting... failed to create: ", err)
	}

	IPService, err := service.GetIP(
		&c,
		NameService,
		Namespace)
	if err != nil {
		fmt.Println("exiting... failed to create: ", err)
	}

	// Save firewall state after service, endpoint created
	// fwAfterEndpointCreated, err := kubeproxy.SaveCurrentFirewallState(
	//	&c,
        //        "kube-proxy",
        //       "kube-proxy",
        //        "kube-system")
        //if err != nil {
	//	fmt.Println(err)
	//}

	// See the difference with diff command
	//out, err := util.DiffCommand(fwInitialState, fwAfterEndpointCreated)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("%s", string(out))

	// START: Pod
	// PodCommandInitBash struct for running bash command
	NameContainer := "kpnginx"
	p := pod.Instance {
		Name: NameContainer,
		Namespace: Namespace,
		Image: "nginx",
		LabelKey: "app",
                LabelValue: labelApp,
	}
	pod.Create(&c, &p)

	// START: Execute curl from the pod created to the new service
	_, err = curl.ExecuteHTTPReqInsideContainer(
			&c,
			NameContainer,
			Namespace,
			"http://" + IPService)
        if err != nil {
		fmt.Println(err)
        }
	// Delete namespace
	namespace.Delete(&c, Namespace)
}
