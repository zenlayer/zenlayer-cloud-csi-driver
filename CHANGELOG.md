# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-07-24
### Added
* Base feature (Create/Delete/Attach/Detach/Resize Volume)
* Multi Zenlayer region Topology Support
* ReadOnly Support
* Storage-class Supported parameters
``` yaml
provisioner:     
parameters:
  fsType:                    
  type:                     
  zoneID:             
  placeGroupID:          
reclaimPolicy:                  
allowVolumeExpansion:                
volumeBindingMode: 
mountOptions:
```
* PVC Supported parameters
``` yaml
accessModes:        
volumeMode: 
resources:      
  requests:     
    storage:
storageClassName:
```
* Support csi-sanity and e2e.test ut-test
* Support helm pull and helm install for Quick installation
### Changed
-
### Fixed
- 