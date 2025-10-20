/*
Copyright (C) 2025 Zenlayer, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this work except in compliance with the License.
You may obtain a copy of the License in the LICENSE file, or at:

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/util/mount"
	k8smount "k8s.io/mount-utils"
	utilexec "k8s.io/utils/exec"
)

func NewSafeMounter() *mount.SafeFormatAndMount {
	realMounter := mount.New("")
	realExec := mount.NewOsExec()
	return &mount.SafeFormatAndMount{
		Interface: realMounter,
		Exec:      realExec,
	}
}

// IsFileExisting check file exist in volume driver
func IsFileExisting(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	// Notice: this err may be is not dictionary error, it will returns true
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDirEmpty check whether the given directory is empty
func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)
	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func FormatAndMount(diskMounter *k8smount.SafeFormatAndMount, source string, target string, fstype string, mkfsOptions []string, mountOptions []string, omitFsCheck bool) error {
	klog.Infof("formatAndMount: mount options [%+v]", mountOptions)
	readOnly := false
	for _, option := range mountOptions {
		if option == "ro" {
			readOnly = true
			break
		}
	}

	// check device fs
	mountOptions = append(mountOptions, "defaults")
	if !readOnly || !omitFsCheck {
		// Run fsck on the disk to fix repairable issues, only do this for volumes requested as rw.
		args := []string{"-a", source}

		out, err := diskMounter.Exec.Command("fsck", args...).CombinedOutput()
		if err != nil {
			ee, isExitError := err.(utilexec.ExitError)
			switch {
			case err == utilexec.ErrExecutableNotFound:
				klog.Warningf("'fsck' not found on system; continuing mount without running 'fsck'.")
			case isExitError && ee.ExitStatus() == 1: //fsckErrorsCorrected
				klog.Infof("Device [%s] has errors which were corrected by fsck.", source)
			case isExitError && ee.ExitStatus() == 4: //fsckErrorsUncorrected
				return fmt.Errorf("'fsck' found errors on device %s but could not correct them: %s", source, string(out))
			case isExitError && ee.ExitStatus() > 4: //fsckErrorsUncorrected
			}
		}
	}

	// Try to mount the disk
	mountErr := diskMounter.Interface.Mount(source, target, fstype, mountOptions)
	if mountErr != nil {
		// Mount failed. This indicates either that the disk is unformatted or
		// it contains an unexpected filesystem.
		existingFormat, err := diskMounter.GetDiskFormat(source)
		if err != nil {
			return err
		}
		if existingFormat == "" {
			return FormatNewDisk(readOnly, source, fstype, target, mkfsOptions, mountOptions, diskMounter)
		}
		// Disk is already formatted and failed to mount
		if len(fstype) == 0 || fstype == existingFormat {
			// This is mount error
			return mountErr
		}
		// Detect if an encrypted disk is empty disk, since atari type partition is detected by blkid.
		if existingFormat == "unknown data, probably partitions" {
			klog.Infof("FormatAndMount: enter special partition logics")
			fsType, ptType, _ := GetDiskFStypePTtype(source)
			if fsType == "" && ptType == "atari" {
				return FormatNewDisk(readOnly, source, fstype, target, mkfsOptions, mountOptions, diskMounter)
			}
		}
		// Block device is formatted with unexpected filesystem, let the user know
		return fmt.Errorf("failed to mount the volume as %q, it already contains %s. Mount error: %v", fstype, existingFormat, mountErr)
	}

	return mountErr
}
func FormatNewDisk(readOnly bool, source, fstype, target string, mkfsOptions, mountOptions []string, diskMounter *k8smount.SafeFormatAndMount) error {
	if readOnly {
		// Don't attempt to format if mounting as readonly, return an error to reflect this.
		return errors.New("failed to mount unformatted volume as read only")
	}

	// Use 'ext4' as the default
	if len(fstype) == 0 {
		fstype = "ext4"
	}

	args := mkfsDefaultArgs(fstype, source)

	// add mkfs options
	if len(mkfsOptions) != 0 {
		args = []string{}
		args = append(args, mkfsOptions...)
		args = append(args, source)
	}

	klog.Infof("Disk [%s] appears to be unformatted, attempting to format as type [%s] with options [%v]", source, fstype, args)
	startT := time.Now()

	pvName := filepath.Base(source)
	_, err := diskMounter.Exec.Command("mkfs."+fstype, args...).CombinedOutput()
	if err == nil {
		// the disk has been formatted successfully try to mount it again.
		klog.Infof("Disk format succeeded, pvName: %s elapsedTime: %+v ms", pvName, time.Since(startT).Milliseconds())
		return diskMounter.Interface.Mount(source, target, fstype, mountOptions)
	}
	klog.Errorf("format of disk %q failed: type:(%q) target:(%q) options:(%q) error:(%v)", source, fstype, target, args, err)
	return err
}

// GetDiskFStypePTtype uses 'blkid' to see if the given disk is unformatted
func GetDiskFStypePTtype(disk string) (fstype string, pttype string, err error) {
	args := []string{"-p", "-s", "TYPE", "-s", "PTTYPE", "-o", "export", disk}

	dataOut, err := utilexec.New().Command("blkid", args...).CombinedOutput()
	output := string(dataOut)

	if err != nil {
		if exit, ok := err.(utilexec.ExitError); ok {
			if exit.ExitStatus() == 2 {
				// Disk device is unformatted.
				// For `blkid`, if the specified token (TYPE/PTTYPE, etc) was
				// not found, or no (specified) devices could be identified, an
				// exit code of 2 is returned.
				return "", "", nil
			}
		}
		klog.Errorf("Could not determine if disk %q is formatted (%v)", disk, err)
		return "", "", err
	}

	lines := strings.Split(output, "\n")
	for _, l := range lines {
		if len(l) <= 0 {
			// Ignore empty line.
			continue
		}
		cs := strings.Split(l, "=")
		if len(cs) != 2 {
			return "", "", fmt.Errorf("blkid returns invalid output: %s", output)
		}
		// TYPE is filesystem type, and PTTYPE is partition table type, according
		// to https://www.kernel.org/pub/linux/utils/util-linux/v2.21/libblkid-docs/.
		if cs[0] == "TYPE" {
			fstype = cs[1]
		} else if cs[0] == "PTTYPE" {
			pttype = cs[1]
		}
	}

	if len(pttype) > 0 {
		// Returns a special non-empty string as filesystem type, then kubelet
		// will not format it.
		return fstype, pttype, nil
	}

	return fstype, "", nil
}
func mkfsDefaultArgs(fstype, source string) (args []string) {
	// default args
	if fstype == "ext4" || fstype == "ext3" {
		args = []string{
			"-F",  // Force flag
			"-m0", // Zero blocks reserved for super-user
			source,
		}
	} else if fstype == "xfs" {
		args = []string{
			"-f",
			source,
		}
	}
	return
}
