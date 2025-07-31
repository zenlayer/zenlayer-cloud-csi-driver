# csi-sanity  

## description
Csi-sanity is a part of the Kubernetes CSI testing framework (CSI-Test), used to check whether the CSI driver complies with the specification through a series of gRPC calls. It is a way of unit and basic functional testing

## install

csi-test v5.3.1 (https://github.com/kubernetes-csi/csi-test/releases)     

## build env
* helm install zeccsi oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.0.0 --set defaultResourceGroup="" --set defaultZone="" --set maxVolume=6 --set mixDriver=true

## run test
* kubectl cp ./csi-sanity -n=kube-system csi-zecplugin-provisioner-7f577b964c-86mhh:/         
* kubectl -n=kube-system exec -it csi-zecplugin-provisioner-7f577b964c-92zb6 -- /csi-sanity -csi.endpoint unix:///csi/csi.sock -csi.controllerendpoint unix:///csi/csi.sock -csi.testvolumesize 21474836480  > csi-sanity.log      