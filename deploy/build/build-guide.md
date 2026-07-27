# Developer Documentation
* If you need to build this project yourself, the following steps are for reference only. Contact after-sales support if you need help.
## Build the image with Docker
* Install the Go toolchain
* Install Docker
* Install make
* Run `make image` in the project root directory. It builds the image into the local Docker image store, where you can list it with `docker images`.
* Export the image: `docker save -o zeccsi.tar zenlayer297/zeccsi:v1.2.0`
* Upload the image archive to every node of the Kubernetes cluster and import it into the containerd `k8s.io` namespace: `ctr -n=k8s.io image import ./zeccsi.tar`
* Once the image is loaded, continue installing the CSI pods as described in the [ZecCSI installation guide](./doc/install-guide.md)

## Build the image with BuildKit
* Install the Go toolchain
* Install buildctl and buildkitd  [Install](https://github.com/moby/buildkit/releases)     
* Create the configuration file `/etc/buildkit/buildkitd.toml`
``` shell
[worker]
  [worker.oci]
    enabled = false
  [worker.containerd]
    address = "/run/containerd/containerd.sock"
    enabled = true
    platforms = ["linux/amd64"]
    namespace = "k8s.io"
    gc = true
    gckeepstorage = 9000

[grpc]
  address = ["tcp://0.0.0.0:1234"]
  uid = 0
  gid = 0
  debug = false

[registry]
  [registry."registry.opsxlab.cn"]
    http = true
    insecure = false
```
* Start buildkitd: `buildkitd --config /etc/buildkit/buildkitd.toml &`
* Run the build from the project root directory: `buildctl --addr tcp://127.0.0.1:1234 build --frontend=dockerfile.v0 --local context=.  --local dockerfile=./deploy/buildkit/  --output type=image,name=docker.io/zenlayer297/zeccsi:v1.2.0`
* Verify the image: `ctr -n=k8s.io image ls -q`
* Once the image is loaded, continue installing the CSI pods as described in the [ZecCSI installation guide](./doc/install-guide.md)