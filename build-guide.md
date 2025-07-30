# Developer Documentation
* If you need to manually compile this project, please refer to the following steps for reference only. If necessary, please contact after-sales support
## use docker build image
* Install the go environment
* Install docker
* Install make
* Directly executing 'make image' in the project root directory will generate an image to the docker image space, which can be seen using 'docker images'
* run cmd 'docker save -o zeccsi.tar zenlayer297/zeccsi:v1.0.0' export image
* Upload the image package to all nodes of the k8s cluster and manually import it into the k8s namespace 'ctr -n=k8s.io image import ./zeccsi.tar'
* After load the image, then you can continue install csi pod refer to [ZecCSI Installation guide](./doc/install-guide.md)

## use buildkit build image
* Install the go environment
* Install buildctl and buildkitd  [Install](https://github.com/moby/buildkit/releases)     
* Generate conf-file /etc/buildkit/buildkitd.toml
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
* run buildkitd: buildkitd --config /etc/buildkit/buildkitd.toml &
* run cmd in the project root directory: buildctl --addr tcp://127.0.0.1:1234 build --frontend=dockerfile.v0 --local context=.  --local dockerfile=./deploy/buildkit/  --output type=image,name=docker.io/zenlayer297/zeccsi:v1.0.0
* check image: ctr -n=k8s.io image ls -q 
* After load the image, then you can continue install csi pod refer to [ZecCSI Installation guide](./doc/install-guide.md)