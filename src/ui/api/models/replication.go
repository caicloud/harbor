// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"github.com/astaxie/beego/validation"
)

const (
	ReplicateSucceed = "Succeed"
	ReplicateFailed = "Failed"
	ReplicateProcessing = "Processing"
)

// Replication defines the properties of model used in replication API
type Replication struct {
	PolicyID int64 `json:"policy_id"`
}

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

// Valid ...
func (r *Replication) Valid(v *validation.Validation) {
	if r.PolicyID <= 0 {
		v.SetError("policy_id", "invalid value")
	}
}
