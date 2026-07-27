# Helm is the official package manager for Kubernetes        
## Chart      
A Helm package is distributed as a TAR archive. Much like a DEB package for APT or an RPM package for YUM, it contains a set of YAML files that define Kubernetes resources.       

## Repository     
A Helm repository is essentially a web server that hosts a set of chart packages for users to download, together with an index file that lists the charts available in that repository.      
Helm can manage several repositories at the same time.       

## Release        
A chart deployed into a Kubernetes cluster with `helm install` is called a release. You can think of a release as one application instance that Helm deployed from a chart package. The same chart can usually be installed more than once in a cluster, and each installation creates a new release.        


## Install Helm 3      
[Helm](https://github.com/helm/helm/tags)       

## cmd
```shell
Helm Commands:

helm pull oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.2.0
helm show all oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.2.0
```