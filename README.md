[![Go Report Card](https://goreportcard.com/badge/github.com/thekubeworld/kubeproxy-testing)](https://goreportcard.com/report/github.com/thekubeworld/kubeproxy-testing)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# kubeproxy-testing
This repo made progress and now is part of bigger movement for test kubernetes services, see [k8s-service-lb-validator project](https://github.com/K8sbykeshed/k8s-service-lb-validator/)

**Still interested to give it a try**? 

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
