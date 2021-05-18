[![Go Report Card](https://goreportcard.com/badge/github.com/thekubeworld/kubeproxy-testing)](https://goreportcard.com/report/github.com/thekubeworld/kubeproxy-testing)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# kubeproxy-testing
A repo for testing kube-proxy and related bits

To see the logs/output from each job, follow these steps:

- Select the artifact from the list above.
- Click on the workflow/job
- Go to Artifacts and download logs.zip. It will contain the output logs from the job.

**Looking to run the tests locally**?  

**1)** Download the [k8devel Go framework for Kubernetes](https://github.com/thekubeworld/k8devel)
```
$ go get github.com/thekubeworld/k8devel
```

**2)** Clone this repo and run the tests
```
$ git clone https://github.com/thekubeworld/kubeproxy-testing
$ cd kubeproxy-testing
$ go run run.go 
🤖 Creating and validating the following kube-proxy scenarios:
⏳ ClusterIP service ✅
⏳ NodePort service ✅
⏳ LoadBalancer service ✅
⏳ ExternalName service ✅
```
