# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-08-05

### Added
* CI workflow: gofmt, build, vet and test on Linux, plus an advisory golangci-lint job.
* Unit tests for volume id parsing, storage class zone handling and volume size rounding.
* Every instrumented call now logs its duration: `ExitFunction` reports `Elapsed[<n>ms]` next to the
  hash that already correlates it with the matching entry line.

### Changed
* Updated zenlayercloud-sdk-go from v0.2.23 to v0.2.49.
* Volume expansion is reported as `ONLINE` instead of `OFFLINE`, so a volume can be expanded
  without stopping the pod that uses it.
* Unified the node-side mounter on `mount-utils` `New("")`. Mounts created by earlier versions are
  only corrected once the pod is recreated.
* `GetRequiredVolumeSizeByte` honours `limit_bytes` and rounds up to whole GiB.
* `DeleteVolume` releases a disk twice so it leaves the recycle bin, and can resume if an earlier
  attempt stopped in between.
* `ControllerPublishVolume` returns `FailedPrecondition` when the target instance is not `RUNNING`.
* Internal wait timeouts became named constants, documented against the sidecar `--timeout` values
  they must stay below, and those sidecar timeouts were raised to match: provisioner 150s to 600s,
  attacher 60s to 240s, resizer and snapshotter 150s to 240s. A sidecar that gave up before the
  driver did would retry straight into the volume lock the first call still held and get `Aborted`.
* `mkfs` no longer discards blocks when formatting a new volume: `-E nodiscard` for ext3 and ext4,
  `-K` for xfs. A newly created ZEC cloud disk is already empty, so a full-device discard gains
  nothing and only makes `NodeStageVolume` slower the larger the disk is. Discard/TRIM on the
  mounted filesystem is unaffected and can still be requested through `mountOptions`.
* Access keys are masked in the startup log.
* Chart: added a liveness probe to the node DaemonSet, fixed `imagePullPolicy`, pinned the
  controller Deployment rollout to `maxSurge: 0` / `maxUnavailable: 1` (with `hostNetwork` and pod
  anti-affinity a surge pod can never be scheduled), and documented `mixDriver` as csi-sanity only.
* Docs: grammar and terminology fixes throughout, and fixed the `featureGates` example, which was
  missing `--set`.

### Removed
* Go Report Card badge and its refresh workflow; the service was sunset on 2026-07-01.
* The broken `golangci-lint` workflow, the Dependabot config and a duplicated `storageclasses` RBAC
  rule.
* The `--retry-detach-times-max` flag and the `retryLimiter` behind it. See the matching entry under
  Fixed for why. The chart never passed the flag, but a hand-written manifest that still does will
  now fail to start.

### Fixed
* `ControllerUnpublishVolume` no longer gives up on a volume permanently. The retry limiter counted
  failed detaches per volume id with no way to reset, so once a volume hit
  `--retry-detach-times-max` (10 by default) every later attempt was refused with
  `Internal: exceeds max retry times` for the lifetime of the controller pod, even after the
  underlying problem was gone. That blocked detach, and with it rescheduling the workload.
* `ControllerPublishVolume` no longer detaches a volume found attached to an unexpected instance.
  That disk is in use, so detaching it could corrupt the other instance's filesystem.
* `ControllerPublishVolume` rejects an access mode the driver does not advertise with
  `InvalidArgument` rather than attaching the volume anyway. A ZEC disk can only be attached to a
  single instance, so a `MULTI_NODE_*` request must not be served.
* `ParseCsiVolId` rejects a malformed volume id instead of passing a truncated disk id to the cloud
  API, and splits at the first `-`. Teardown RPCs stay idempotent; the rest return `NotFound`.
* `CreateVolume` no longer returns a volume id assembled from an empty disk id or serial. The cloud
  API may still report the serial as null just after a disk is created, and the resulting id could
  never be parsed back, so the call now fails with `Aborted` and the CO retries.
* `FormatAndMount` no longer appends `defaults` after the caller's mount options. `mount` lets the
  last option win and `defaults` implies `rw`, so a StorageClass asking for `ro` was staged
  read-write while `fsck` was skipped as though it were read-only.
* `CreateSnapshot` accepts several `createTime` formats and falls back to the current time instead
  of failing. The snapshot already exists in the cloud by that point, so returning `Internal` left
  it permanently not ready to use and leaked it, with every idempotent retry failing the same way.
* `NodePublishVolume` detects an already published block volume with `IsMountPoint`.
  `IsLikelyNotMountPoint` compares a directory with its parent and cannot see a device file
  bind-mount, so a repeated call stacked a second mount on the target path. It also no longer
  resolves a device path for filesystem volumes, where only the staging path is bind-mounted.
* `NodeExpandVolume` detects block volumes from `volume_path` when the CO omits
  `volume_capability`, instead of running `resize2fs` on a raw device.
* `IsValidTopology` compares topologies with `proto.Equal`. `reflect.DeepEqual` also compared
  protobuf internals, which broke `CreateVolume` idempotency on retry.
* `NodeUnstageVolume` and `ControllerUnpublishVolume` are idempotent, including when the node is
  already gone.
* `ValidateVolumeCapabilities` reports an unsupported access mode as a success response without
  `confirmed`, as the specification requires.
* Cloud calls check for a nil client so the node driver no longer panics, and `StringToType` no
  longer panics on an unknown disk category.
* `FindVolumeByName` skips disks in the recycle bin, paged lookups accumulate results with `append`,
  and `DetachDisks` failures reported in `failedDiskIds` are detected.
* Serials read from sysfs and access keys read from secret files are trimmed.
* `UpdateParmsZone` matches the `zoneID` key case-insensitively.
* The node liveness probe listens on `:29633` instead of `localhost:29633`. The DaemonSet does not
  use `hostNetwork`, so kubelet probes the pod IP and could never reach a loopback-only listener.
* Typos in identifiers: `VmType.IsVaild` is now `IsValid`, and the `nodeNmae` flag variable is now
  `nodeName`. The `--nodename` flag itself is unchanged.

## [1.1.0] - 2026-03-20

### Added
* StorageClass support for configuring cloud disk performance bursting (`burstEnable`).

### Changed
* Updated zenlayercloud-sdk-go from v0.2.0 to v0.2.23, and refreshed the other go.mod dependencies.
* Documentation updates.

## [1.0.0] - 2025-10-20

### Added
* Base features: create, delete, attach, detach and resize volumes.
* Topology support across multiple Zenlayer regions.
* Filesystem-mode PVCs, including read-only mounts.
* Block-mode PVCs.
* StorageClass and PVC support.
* VolumeSnapshotClass and VolumeSnapshot support.
* csi-sanity and Kubernetes external storage e2e test support.
* Helm chart for quick installation.

[1.2.0]: https://github.com/zenlayer/zenlayer-cloud-csi-driver/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/zenlayer/zenlayer-cloud-csi-driver/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/zenlayer/zenlayer-cloud-csi-driver/releases/tag/v1.0.0
