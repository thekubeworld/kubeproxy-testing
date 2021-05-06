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
	"os/exec"
	"github.com/thekubeworld/k8devel/pkg/emoji"
)

func main() {
	e := emoji.LoadEmojis()

	fmt.Printf(
		"%s Creating and validating the following kube-proxy scenarios:\n",
		emoji.Show(e.Robot))

	fmt.Printf("%s ClusterIP service ",
		emoji.Show(e.HourGlassNotDone))

	cmd := exec.Command("go", "run", "tests/clusterip/main.go")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s \n\t%s\n", emoji.Show(e.CrossMark), err)
	} else {
		fmt.Printf("%s\n", emoji.Show(e.CheckMarkButton))
	}

	fmt.Printf("%s NodePort service ",
		emoji.Show(e.HourGlassNotDone))
	cmd = exec.Command("go", "run", "tests/nodeport/main.go")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%s \n\t%s\n", emoji.Show(e.CrossMark), err)
	} else {
		fmt.Printf("%s\n", emoji.Show(e.CheckMarkButton))
	}

	fmt.Printf("%s LoadBalancer service ",
		emoji.Show(e.HourGlassNotDone))
	cmd = exec.Command("go", "run", "tests/loadbalancer/main.go")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%s \n\t%s\n", emoji.Show(e.CrossMark), err)
	} else {
		fmt.Printf("%s\n", emoji.Show(e.CheckMarkButton))
	}

	fmt.Printf("%s ExternalName service ",
		emoji.Show(e.HourGlassNotDone))
	cmd = exec.Command("go", "run", "tests/externalname/main.go")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%s \n\t%s\n", emoji.Show(e.CrossMark), err)
	} else {
		fmt.Printf("%s\n", emoji.Show(e.CheckMarkButton))
	}

}
