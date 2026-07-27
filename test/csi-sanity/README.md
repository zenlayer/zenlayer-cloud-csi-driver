# csi-sanity

## Description
csi-sanity is part of the Kubernetes CSI testing framework (csi-test). It issues a series of gRPC calls to check whether a CSI driver complies with the CSI specification, and it provides both unit-level and basic functional testing.

## Install

csi-test v5.5.0 (https://github.com/kubernetes-csi/csi-test/releases)

## Set Up the Test Environment
csi-sanity drives the controller service and the node service through a single endpoint, so the driver must be installed with `mixDriver=true`. That setting makes the controller pod register the node service on its own socket as well. It is intended for this test only and must not be used in production; see the `mixDriver` comment in `chart/values.yaml` for details.

```shell
helm install zeccsi oci://registry-1.docker.io/zenlayer297/zenlayer-cloud-csi-driver --version 1.2.0 --set defaultResourceGroup="" --set defaultZone="" --set maxVolume=6 --set mixDriver=true
```

## Run the Test
Copy the `csi-sanity` binary into the controller pod and run it there. Use the same pod name in both commands:

```shell
POD=$(kubectl get pod -n kube-system -l app=csi-zecplugin-provisioner -o jsonpath='{.items[0].metadata.name}')
kubectl cp ./csi-sanity -n kube-system "$POD":/
kubectl exec -it -n kube-system "$POD" -- /csi-sanity -csi.endpoint unix:///csi/csi.sock -csi.controllerendpoint unix:///csi/csi.sock -csi.testvolumesize 21474836480 > ~/csi-sanity.log
```

## Known Failures
`csi-sanity.log` in this directory records the result of the last run: 60 passed, 2 failed, 1 pending, 33 skipped.

Both failures are caused by the same limitation. The CSI specification requires a driver to accept a name of up to 128 bytes, but a ZEC cloud disk name and a ZEC snapshot name are limited to 64 characters, and the driver currently rejects anything longer:

* `CreateVolume should not fail when creating volume with maximum-length name`
* `CreateSnapshot should succeed when creating snapshot with maximum-length name`

This does not affect Kubernetes workloads. external-provisioner names volumes `zeccsi-pv-<uuid>` (46 characters) and external-snapshotter names snapshots `zeccsi-snapshot-<uuid>` (52 characters), so neither ever reaches the 64-character limit.
