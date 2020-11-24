package model

import (
	"github.com/astaxie/beego/validation"
)

const (
	// ImagesReplicateSucceed represents the images replication task succeed
	ImagesReplicateSucceed = "Succeed"
	// ImagesReplicateFailed represents the images replication task failed
	ImagesReplicateFailed = "Failed"
	// ImagesReplicateProcessing represents the images replication task is under processing
	ImagesReplicateProcessing = "Processing"
)

// ImagesReplication describes an images replication request
type ImagesReplication struct {
	Images  []string `json:"images"`
	Targets []string `json:"targets"`
}

// ImagesReplicationRsp describes response of an images replication, it defines
// UUID of the operation, which can be used to retrieve replication status
type ImagesReplicationRsp struct {
	UUID string `json:"uuid"`
}

// ImagesReplicationStatus describes image replication status
type ImagesReplicationStatus struct {
	Status     string           `json:"status"`
	JobsStatus []ImageRepStatus `json:"jobs_status"`
}

// ImageRepStatus describes replication status of an image
type ImageRepStatus struct {
	Image  string `json:"image"`
	Status string `json:"status"`
}

// Valid the ImagesReplication
func (p *ImagesReplication) Valid(v *validation.Validation) {
	if len(p.Images) == 0 {
		v.SetError("images", "cannot be empty")
	}

	if len(p.Targets) == 0 {
		v.SetError("targets", "cannot be empty")
	}

	if len(p.Targets) > 1 {
		v.SetError("targets", "only support 1 target")
	}

}
