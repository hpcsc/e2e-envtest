# e2e-envtest

An example to demonstrate running e2e tests for an application that is supposed to be deployed to Kubernetes.

This application is also dependent on Kubernetes API (use Kubernetes client to query API server)

## Run

(Require [Taskfile](https://taskfile.dev))

```shell
task e2e
```

## Test set up

Before the e2e test is run:

- [kubebuilder envtest](https://book.kubebuilder.io/reference/envtest.html) binaries (which contain etcd, kube-apiserver, kubectl) are downloaded and extracted to `./bin/envtest`
- main application is built and placed at `./bin/e2e-envtest`

When e2e test is run:

- It starts envtest using binaries from `./bin/envtest`
- It writes kubernetes rest config file (returned from envtest setup) to a temporary file in filesystem (`./bin/kubeconfig`). This file will be used by main application to interact with envtest
- It sets up test data (create 2 test namespaces in envtest in this case)
- It uses Go `exec.Command()` to execute main application binary
- It verifies main application output and any side effect
- It shuts down envtest and delete temporary kubeconfig file
