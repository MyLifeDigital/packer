package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/hashicorp/packer/helper/multistep"
)

type StepAttachVolume struct {
	CaptureVolume bool
}

func (s *StepAttachVolume) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// Proceed only if block storage volume is used.
	if !s.CaptureVolume {
		return multistep.ActionContinue
	}

	config := state.Get("config").(*Config)
	computeClient, _ := config.computeV2Client()
	// ui := state.Get("ui").(packer.Ui)

	volume := state.Get("volume_id").(string)
	server := state.Get("server").(*servers.Server)

	createOpts := volumeattach.CreateOpts{
		Device:   "/dev/vdb",
		VolumeID: volume,
	}

	_, err := volumeattach.Create(computeClient, server.ID, createOpts).Extract()
	if err != nil {
		err = fmt.Errorf("Error attaching block storage volume: %s", err)
		state.Put("error", err)
		return multistep.ActionHalt
	}

	// Wait for volume to become available.
	// ui.Say(fmt.Sprintf("Waiting for volume %s (volume id: %s) to become attached...", config.VolumeName, volume))
	// if err := WaitForVolume(blockStorageClient, volume); err != nil {
	// 	err := fmt.Errorf("Error waiting for volume: %s", err)
	// 	state.Put("error", err)
	// 	ui.Error(err.Error())
	// 	return multistep.ActionHalt
	// }

	return multistep.ActionContinue
}

func (s *StepAttachVolume) Cleanup(multistep.StateBag) {
	// No cleanup.
}
