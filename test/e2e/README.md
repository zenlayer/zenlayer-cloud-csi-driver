# Zenlayer CSI Driver E2E Test

## Description

This directory contains scripts and config templates used to run [Kubernetes external storage e2e test](https://github.com/kubernetes/kubernetes/tree/master/test/e2e/storage/external).

## Prerequisites

The test can only be run on Zec nodes with a installed Kubernetes cluster, as it actually creates/attaches volumes by calling the Zec IAAS API.   

Make sure no volume is attached to any node, which will mess up with the volume limits test.    

## Install

1. The zones that have a good network connectivity to global internet is preferred, like `ap2a`, as it will download the e2e test packages from google when it runs for the first time. If that's impossible, you can download the package somewhere else, and upload it to the test node:    

```bash
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.28.2/kubernetes-test-linux-amd64.tar.gz --output e2e-tests.tar.gz
tar -xf e2e-tests.tar.gz --directory=./zenlayer-cloud-csi-driver/test/e2e
```

2. Multiple nodes are preferred, if that's possible, as some tests will check against volume drifting from one node to another. If there is only one, those tests will be skipped by the script.      

## Run 

You can skip some tests/only run some tests, by editing the following lines of the script:

```bash
./ginkgo -focus='External.Storage.*' -skip='(.*Disruptive.*|.*stress.*)' ./e2e.test -- -storage.testdriver=driver.yaml -kubeconfig=/etc/kubernetes/admin.conf --ginkgo.timeout="12h"  > ~/e2e.log     
```

It's simple regex, refer to the [Ginkgo docs](https://onsi.github.io/ginkgo/#the-ginkgo-cli) for more detail.