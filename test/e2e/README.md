# Zenlayer CSI Driver E2E Test

## Description

This directory contains scripts and config templates used to run [Kubernetes external storage e2e test](https://github.com/kubernetes/kubernetes/tree/master/test/e2e/storage/external).

## Prerequisites

The test can only be run on ZEC nodes with an installed Kubernetes cluster, because it really does create and attach volumes by calling the ZEC IaaS API.   

Make sure no volume is attached to any node; otherwise the volume limits test will be disrupted.    

## Install

1. Prefer a zone with good connectivity to the global internet, because the first run downloads the e2e test packages from Google. If that is not possible, download the package elsewhere and upload it to the test node:    

```bash
curl -L https://storage.googleapis.com/kubernetes-release/release/v1.28.2/kubernetes-test-linux-amd64.tar.gz --output e2e-tests.tar.gz
tar -xf e2e-tests.tar.gz --directory=./zenlayer-cloud-csi-driver/test/e2e
```

2. Use more than one node if you can, because some tests check that a volume does not drift from one node to another. On a single-node cluster those tests are skipped.      

## Run 

Adjust the `-focus` and `-skip` regular expressions below to select which tests to run:

```bash
./ginkgo -focus='External.Storage.*' -skip='(.*Disruptive.*|.*stress.*)' ./e2e.test -- -storage.testdriver=driver.yaml -kubeconfig=/etc/kubernetes/admin.conf --ginkgo.timeout="12h"  > ~/e2e.log     
```

These are plain regular expressions; refer to the [Ginkgo docs](https://onsi.github.io/ginkgo/#the-ginkgo-cli) for more detail.