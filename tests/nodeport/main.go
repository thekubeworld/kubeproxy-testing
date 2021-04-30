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
	"github.com/thekubeworld/k8devel"
	"github.com/sirupsen/logrus"
	"github.com/gookit/color"
)

func main() {
	k8devel.SetLogrusLogging()
	logrus.Infof("kube-proxy tests has started...")

	logrus.Infof("\n")
	logrus.Infof("Test #2) Traffic will reach kube-proxy that will    ")
	logrus.Infof("route using iptables to the right node that has the ")
	logrus.Infof("backend pod                                         ")
	logrus.Infof("\n                                                  ")
	logrus.Infof("               +---------------------+              ")
	logrus.Infof("               |      TRAFFIC        |              ")
	logrus.Infof("               |    FROM USERS       |              ")
	logrus.Infof("               +---------------------+              ")
	logrus.Infof("                          |                         ")
	logrus.Infof("                          v                         ")
	logrus.Infof("             +-- kube-proxy/iptables-+              ")
	logrus.Infof("             |                       |              ")
	logrus.Infof("             v                       v              ")
	logrus.Infof("          Node IPs               Node IPs           ")
	logrus.Infof("         and Ports               and ports          ")
        logrus.Infof("+------------|-------+    +-----------|------------+")
	logrus.Infof("| 10.10.50.54|:30001 |    |10.10.50.51|:30001      |")
	logrus.Infof("+------------|--------    +-----------|------------+")
	logrus.Infof("|		   +------------+-----------+            |")
	logrus.Infof("|                         |                        |")
	logrus.Infof("+-----------------------+-V-+----------------------|")
	logrus.Infof("| ClusterIP 10.111.239.7|:80|                      |")
	logrus.Infof("+-----------------------|-|-|----------------------|")
	logrus.Infof("| Service 1             | | |                      |")
	logrus.Infof("|   Selector:           | | |                      |")
	logrus.Infof("|   app: nginx          | | |                      |")
	logrus.Infof("+-----------------------|-|-|----------------------|")
	logrus.Infof("|               +-------|-+-|-----------+          |")
	logrus.Infof("|               |       |   |           |          |")
	logrus.Infof("+---------------|-------|   |-----------|----------|")
	logrus.Infof("| EndpointIP and| ports |   |EndpointIP | and ports|")
	logrus.Infof("|               v       |   |           v          |")
	logrus.Infof("|   10.244.2.2 :80 :8080|   |10.244.2.3 :80  :8080 |")
	logrus.Infof("+---------------|---|---|   |-----------|----|-----|")
	logrus.Infof("|               |   |   |   |           |    |     |")
	logrus.Infof("+---------------|---|---|   |-----------|----|-----|")
	logrus.Infof("| Container port|   |   |   | Container |Port|     |")
	logrus.Infof("|               |   |   |   |           |    |     |")
	logrus.Infof("|             :80 :8080 |   |          :80  :8080  |")
	logrus.Infof("+---------------|---|---|---|-----------|----|     |")
	logrus.Infof("|               |   |   |   |           |    |     |")
	logrus.Infof("|               |   |   |   |           |    |     |")
	logrus.Infof("|               |   |   |   |           |    |     |")
	logrus.Infof("|               v   |   |   |           v    |     |")
	logrus.Infof("|      Container 1  |   |   |  Container 1   |     |")
	logrus.Infof("|                   |   |   |                |     |")
	logrus.Infof("|                   v   |   |                v     |")
	logrus.Infof("|         Container 2   |   |          Container 2 |")
	logrus.Infof("|    Labels: app nginx  |   |    Labels: app nginx |")
	logrus.Infof("|                       |   |                      |")
	logrus.Infof("| Pod 1                 |   | Pod 2                |")
	logrus.Infof("| Node 1                |   | Node 2               |")
	logrus.Infof("+-----------------------+   +----------------------+")
	logrus.Infof("\n")

	// Initial set
        c := k8devel.Client{}
	c.Namespace = "kptesting"
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
	kyPods, kyNumberPods := k8devel.FindPodsWithNameContains(&c,
		KP,
		namespaceKP)
	if kyNumberPods < 0 {
		logrus.Fatal("exiting... unable to find kube-proxy pod..")
	}
	logrus.Infof("Found the following kube-proxy pods:")
	logrus.Infof("\t\tNamespace: %s", namespaceKP)
	logrus.Infof("\t\t%s", kyPods)

	// Detect Kube-proxy mode
	kpMode, err := k8devel.DetectKubeProxyMode(&c,
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
	randStr, err := k8devel.GenerateRandomString(6, "lower")
	if err != nil {
		logrus.Fatal(err)
	}

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


	// Setting Service Name
	KPTestServiceName := KPTestNamespaceName +
			"service" +
			randStr

	// START: Service
	// nodePort - a static port assigned on each the node
        // port - port exposed internally in the cluster
        // targetPort - the container port to send requests to
	s := k8devel.Service {
		Name: KPTestServiceName,
		Namespace: KPTestNamespaceName,
		LabelKey: "app",
		LabelValue: "kptesting",
		SelectorKey: "app",
		SelectorValue: "kptesting",
		PortName: "http",
		PortProtocol: "TCP",
		Port: 80,              // service
		TargetPort: 80,        // container
		NodePort: 30001,       // node
	}
	err = k8devel.CreateNodePortService(&c, &s)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}

	IPNodes, err := k8devel.GetIPFromNodes(&c)
	if err != nil {
		logrus.Fatal("exiting... failed to create: ", err)
	}
	// END: Service


	// START: iptables diff
	iptablesStateAfterEndpointCreated, err := k8devel.IPTablesSaveNatTable(
				&c, &iptablesCmd, KPTestContainerName, "kube-system")
        if err != nil {
		logrus.Fatal(err)
        }

	// Make a diff between two states we collected from iptables
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

	// Creating a POD Behind the service
	p := k8devel.Pod {
		Name: "kpnginxbehindservice",
		Namespace: KPTestNamespaceName,
		Image: "nginx:1.14.2",
		LabelKey: "app",
		LabelValue: "kptesting",
	}
	logrus.Infof("\n")
	err = k8devel.CreatePod(&c, &p)
        if err != nil {
		logrus.Fatal(err)
        }
	// END: Pod
	logrus.Info("\n")

	// Creating a POD outside the service (No labels)
	// So it will try to connect to pod behind the service
	containerName := "nginxtoconnecttoservice"
	p = k8devel.Pod {
		Name: containerName,
		Namespace: KPTestNamespaceName,
		Image: "nginx",
		LabelKey: "app",
		LabelValue: "foobar",
	}
	err = k8devel.CreatePod(&c, &p)
        if err != nil {
		logrus.Fatal(err)
        }
	// END: Pod

	// START: Execute curl from the pod created to the new service
	ret, err := k8devel.ExecuteHTTPReqInsideContainer(
			&c,
			containerName,
			KPTestNamespaceName,
			IPNodes[0] + ":" + fmt.Sprint(s.NodePort))
        if err != nil {
		logrus.Fatal(err)
        }
	logrus.Infof("%s", ret)
	color.Green.Println("[Test #2 PASSED]")
	// END: Execute curl from the pod created to the new service

	// TODO use cleanup function
	k8devel.DeleteNamespace(&c, KPTestNamespaceName)
}