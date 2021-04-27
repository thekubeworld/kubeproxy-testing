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
	"github.com/dougsland/k8devel"
	"github.com/sirupsen/logrus"
	"github.com/gookit/color"
)

func main() {
	k8devel.SetLogrusLogging()
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

	// Initial set
        c := k8devel.Client{}
	c.Namespace = "kptesting"
	c.NumberMaxOfAttemptsPerTask = 5
        c.TimeoutTaksInSec = 1

	// Connect to cluster from:
	//	- $HOME/kubeconfig (Linux)
	//	- os.Getenv("USERPROFILE") (Windows)
        c.Connect()

	// kube-proxy is a DaemonSet, it will replicate the pod
	// data to all nodes. Let's find one pod name to use and
	// collect data.
	kyPods, kyNumberPods := k8devel.FindPodsWithNameContains(&c,
		"kube-proxy",
		"kube-system")
	if kyNumberPods < 0 {
		logrus.Fatal("exiting... unable to find kube-proxy pod..")
	}

	// TODO: use flags
	// START: kube-proxy variables
	KPTestContainerName := kyPods[0]
	KPTestNamespaceName := c.Namespace

	randStr, err := k8devel.GenerateRandomString(6, "lower")
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

	// TODO: Just load the iptables commands if kube-proxy
	// is IPTABLES or return error
	// Loading some iptables
	iptablesCmd := k8devel.IPTablesLoadPreDefinedCommands()
	if err != nil {
		logrus.Fatal(err)
	}

	// iptables saving initial state
	iptablesInitialState, err := k8devel.IPTablesSaveNatTable(
				&c,
				&iptablesCmd,
				KPTestContainerName,
				"kube-system")
        if err != nil {
		logrus.Fatal(err)
        }

	// START: Namespace
	_, err = k8devel.ExistsNamespace(&c,
			KPTestNamespaceName)
	if err != nil {
		err = k8devel.CreateNamespace(&c,
			KPTestNamespaceName)
		if err != nil {
			logrus.Fatal("exiting... failed to create: ", err)
		}
	}
	// END: Namespace

	// START: Deployment
	d := k8devel.Deployment {
		Name: KPTestNginxDeploymentName,
		Namespace: KPTestNamespaceName,
		Replicas: 1,
		LabelKey: "app",
		LabelValue: "kptesting",
	}

	d.Pod.Name = "nginx"
	d.Pod.Image = "nginx:1.14.2"
	d.Pod.ContainerPortName = "http"
	d.Pod.ContainerPortProtocol = "TCP"
	d.Pod.ContainerPort = 80

	err = k8devel.CreateDeployment(&c, &d)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	// END: Deployment

	// START: Service
	s := k8devel.Service {
		Name: KPTestServiceName,
		Namespace: KPTestNamespaceName,
		LabelKey: "app",
		LabelValue: "kptesting",
		SelectorKey: "app",
		SelectorValue: "kptesting",
		ClusterIP: "",
		Port: 80,
	}
	err = k8devel.CreateService(&c, &s)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}

	IPService, err := k8devel.GetIPFromService(
		&c,
		KPTestServiceName,
		KPTestNamespaceName)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	// END: Service

	// START: iptables diff
	iptablesStateAfterEndpointCreated, err := k8devel.IPTablesSaveNatTable(
				&c,
				&iptablesCmd,
				KPTestContainerName,
				"kube-system")
        if err != nil {
		logrus.Fatal(err)
        }

	out, err := k8devel.DiffCommand(iptablesInitialState.Name(),
			iptablesStateAfterEndpointCreated.Name())
        if err != nil {
		logrus.Fatal(err)
        }

	if len(string(out)) > 0 {
		logrus.Infof("%s", string(out))
	}
	// END: iptables diff

	// START: Pod
	// PodCommandInitBash struct for running bash command
	PodCommandInitBash := []string {
		"/bin/bash",
	}

	SleepOneDay := []string {
		"-c",
		"sleep 5000",
	}

	containerName := "kptestingnginx"
	p := k8devel.Pod {
		Name: containerName,
		Namespace: KPTestNamespaceName,
		Image: "nginx",
		Command: PodCommandInitBash,
		CommandArgs: SleepOneDay,
	}

	k8devel.CreatePod(&c, &p)
	// END: Pod

	// START: Execute curl from the pod created to the new service
	ret, err := k8devel.ExecuteHTTPReqInsideContainer(
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

}
