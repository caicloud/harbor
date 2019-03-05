package api

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/uuid"
	"github.com/vmware/harbor/src/ui/api/models"
	common_models "github.com/vmware/harbor/src/common/models"
	"github.com/vmware/harbor/src/common/notifier"
	"github.com/vmware/harbor/src/replication/event/notification"
	"github.com/vmware/harbor/src/replication/event/topic"
	rep_models "github.com/vmware/harbor/src/replication/models"
	"github.com/vmware/harbor/src/replication"
	"github.com/vmware/harbor/src/common/dao"
)

// ImageReplicateAPI handles API calls for images replication
type ImageReplicateAPI struct {
	BaseController
}

// Prepare does authentication and authorization works
func (r *ImageReplicateAPI) Prepare() {
	r.BaseController.Prepare()
	if !r.SecurityCtx.IsAuthenticated() {
		r.HandleUnauthorized()
		return
	}

	if !r.SecurityCtx.IsSysAdmin() && !r.SecurityCtx.IsSolutionUser() {
		r.HandleForbidden(r.SecurityCtx.GetUsername())
		return
	}
}

// Replicate replicates a list of images to remote targets
func (r *ImageReplicateAPI) Replicate() {
	imgReplication := &models.ImagesReplication{}
	r.DecodeJSONReqAndValidate(imgReplication)

	var items []rep_models.FilterItem
	for _, img := range imgReplication.Images {
		if !isValidImage(img) {
			r.HandleBadRequest(fmt.Sprintf("malfold image '%s'", img))
			return
		}
		items = append(items, rep_models.FilterItem{
			Kind: replication.FilterItemKindTag,
			Value: img,
			Operation: common_models.RepOpTransfer,
		})
	}

	opUUID := strings.Replace(uuid.Generate().String(), "-", "", -1)
	err := notifier.Publish(topic.StartReplicationTopic,
		notification.StartReplicationNotification{
			PolicyID: -1,
			Metadata: map[string]interface{}{
				"uuid": opUUID,
				"independent": true,
				"candidates": items,
				"targets": imgReplication.Targets,
			},
		})

	if err != nil {
		r.HandleInternalServerError(fmt.Sprintf("publish replication event error: %v", err))
		return
	}

	r.Data["json"] = models.ImagesReplicationRsp{
		UUID: opUUID,
	}
	r.ServeJSON()
}

// Status check status of an image replication request
func (r *ImageReplicateAPI) Status() {
	uuid := r.GetString(":uuid")
	query := &common_models.RepJobQuery{
		OpUUID: uuid,
	}

	jobs, err := dao.GetRepJobs(query)
	if err != nil {
		r.HandleInternalServerError(fmt.Sprintf("failed to get repository jobs, query: %v :%v", query, err))
		return
	}

	if len(jobs) == 0 {
		r.HandleNotFound(fmt.Sprintf("no jobs found for %s", uuid))
		return
	}

	statuses := models.ImagesReplicationStatus{
		Status: models.ReplicateSucceed,
	}
	var hasProcessing, hasFailed bool
	for _, job := range jobs {
		statuses.JobsStatus = append(statuses.JobsStatus, models.ImageRepStatus{
			Image: fmt.Sprintf("%s:%s", job.Repository, job.Tags),
			Status: job.Status,
		})
		s := convertStatus(job.Status)
		if s == models.ReplicateProcessing {
			hasProcessing = true
		} else if s == models.ReplicateFailed {
			hasFailed = true
		}
	}
	if hasProcessing {
		statuses.Status = models.ReplicateProcessing
	} else if hasFailed {
		statuses.Status = models.ReplicateFailed
	}

	r.Data["json"] = statuses
	r.ServeJSON()
}

func convertStatus(status string) string {
	if status == common_models.JobRunning ||
		status == common_models.JobPending ||
		status == common_models.JobRetrying {
		return models.ReplicateProcessing
	}

	if status == common_models.JobError ||
		status == common_models.JobStopped ||
		status == common_models.JobCanceled {
		return models.ReplicateFailed
	}

	return models.ReplicateSucceed
}

func isValidImage(img string) bool {
	_, err := common_models.ParseImage(img)
	return err == nil
}