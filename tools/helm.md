# Helm is an officially provided package manager for kubernetes        
## Chart      
The Helm software package is in TAR format. Similar to the DEB package of APT or the RPM package of YUM, it contains a set of YAML files that define Kubernetes resources       

## Repository     
Helm's software Repository, Repository, is essentially a Web server that stores a series of Chart software packages for users to download and provides a manifest file of the Chart packages in this Repository for query.      
Helm can manage multiple different repositories simultaneously.       

## Release        
A Chart deployed in a Kubernetes cluster using the helm install command is called a Release. It can be understood as an application instance deployed by Helm using the Chart package. A chart can usually be installed multiple times in the same cluster. Each installation will create a new release.        


## install Helm3      
[Helm](https://github.com/helm/helm/tags)       

## cmd
```shell
Helm Commands:

helm pull oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.0.0
helm show all oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.0.0
helm template <my-release> oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.0.0
helm install <my-release> oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.0.0
helm upgrade <my-release> oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version <new-version>
```