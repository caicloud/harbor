// Copyright Project Harbor Authors
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

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	common_dao "github.com/goharbor/harbor/src/common/dao"
	common_http "github.com/goharbor/harbor/src/common/http"
	"github.com/goharbor/harbor/src/replication"
	"github.com/goharbor/harbor/src/replication/dao/models"
	"github.com/goharbor/harbor/src/replication/event"
	"github.com/goharbor/harbor/src/replication/model"
)

// ReplicationOperationAPI handles the replication operation requests
type ReplicationOperationAPI struct {
	BaseController
	execution *models.Execution
	task      *models.Task
}

// Prepare ...
func (r *ReplicationOperationAPI) Prepare() {
	r.BaseController.Prepare()
	// As we delegate the jobservice to trigger the scheduled replication,
	// we need to allow the jobservice to call the API
	if !(r.SecurityCtx.IsSysAdmin() || r.SecurityCtx.IsSolutionUser()) {
		if !r.SecurityCtx.IsAuthenticated() {
			r.SendUnAuthorizedError(errors.New("UnAuthorized"))
			return
		}
		r.SendForbiddenError(errors.New(r.SecurityCtx.GetUsername()))
		return
	}

	// check the existence of execution if execution ID is provided in the request path
	executionIDStr := r.GetStringFromPath(":id")
	if len(executionIDStr) > 0 {
		executionID, err := strconv.ParseInt(executionIDStr, 10, 64)
		if err != nil || executionID <= 0 {
			r.SendBadRequestError(fmt.Errorf("invalid execution ID: %s", executionIDStr))
			return
		}
		execution, err := replication.OperationCtl.GetExecution(executionID)
		if err != nil {
			r.SendInternalServerError(fmt.Errorf("failed to get execution %d: %v", executionID, err))
			return
		}
		if execution == nil {
			r.SendNotFoundError(fmt.Errorf("execution %d not found", executionID))
			return
		}
		r.execution = execution

		// check the existence of task if task ID is provided in the request path
		taskIDStr := r.GetStringFromPath(":tid")
		if len(taskIDStr) > 0 {
			taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
			if err != nil || taskID <= 0 {
				r.SendBadRequestError(fmt.Errorf("invalid task ID: %s", taskIDStr))
				return
			}
			task, err := replication.OperationCtl.GetTask(taskID)
			if err != nil {
				r.SendInternalServerError(fmt.Errorf("failed to get task %d: %v", taskID, err))
				return
			}
			if task == nil || task.ExecutionID != executionID {
				r.SendNotFoundError(fmt.Errorf("task %d not found", taskID))
				return
			}
			r.task = task
		}
	}
}

// The API is open only for system admin currently, we can use
// the code commentted below to make the API available to the
// users who have permission for all projects that the policy
// refers
/*
func (r *ReplicationOperationAPI) authorized(policy *model.Policy, resource rbac.Resource, action rbac.Action) bool {

	projects := []string{}
	// pull mode
	if policy.SrcRegistryID != 0 {
		projects = append(projects, policy.DestNamespace)
	} else {
		// push mode
		projects = append(projects, policy.SrcNamespaces...)
	}

	for _, project := range projects {
		resource := rbac.NewProjectNamespace(project).Resource(resource)
		if !r.SecurityCtx.Can(action, resource) {
			r.HandleForbidden(r.SecurityCtx.GetUsername())
			return false
		}
	}

	return true
}
*/

// ListExecutions ...
func (r *ReplicationOperationAPI) ListExecutions() {
	query := &models.ExecutionQuery{
		Trigger: r.GetString("trigger"),
	}

	if len(r.GetString("status")) > 0 {
		query.Statuses = []string{r.GetString("status")}
	}
	if len(r.GetString("policy_id")) > 0 {
		policyID, err := r.GetInt64("policy_id")
		if err != nil || policyID <= 0 {
			r.SendBadRequestError(fmt.Errorf("invalid policy_id %s", r.GetString("policy_id")))
			return
		}
		query.PolicyID = policyID
	}
	page, size, err := r.GetPaginationParams()
	if err != nil {
		r.SendBadRequestError(err)
		return
	}
	query.Page = page
	query.Size = size

	total, executions, err := replication.OperationCtl.ListExecutions(query)
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to list executions: %v", err))
		return
	}
	r.SetPaginationHeader(total, query.Page, query.Size)
	r.WriteJSONData(executions)
}

// CreateExecution starts a replication
func (r *ReplicationOperationAPI) CreateExecution() {
	execution := &models.Execution{}
	if err := r.DecodeJSONReq(execution); err != nil {
		r.SendBadRequestError(err)
		return
	}

	policy, err := replication.PolicyCtl.Get(execution.PolicyID)
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to get policy %d: %v", execution.PolicyID, err))
		return
	}
	if policy == nil {
		r.SendNotFoundError(fmt.Errorf("policy %d not found", execution.PolicyID))
		return
	}
	if !policy.Enabled {
		r.SendBadRequestError(fmt.Errorf("the policy %d is disabled", execution.PolicyID))
		return
	}
	if err = event.PopulateRegistries(replication.RegistryMgr, policy); err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to populate registries for policy %d: %v", execution.PolicyID, err))
		return
	}

	trigger := r.GetString("trigger", string(model.TriggerTypeManual))
	executionID, err := replication.OperationCtl.StartReplication(policy, nil, model.TriggerType(trigger))
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to start replication for policy %d: %v", execution.PolicyID, err))
		return
	}
	r.Redirect(http.StatusCreated, strconv.FormatInt(executionID, 10))
}

// CreateImagesExecution starts an images list replication
func (r *ReplicationOperationAPI) CreateImagesExecution() {
	imagesRep := &model.ImagesReplication{}
	isValid, err := r.DecodeJSONReqAndValidate(imagesRep)
	if !isValid {
		r.SendBadRequestError(err)
		return
	}

	registry, err := replication.RegistryMgr.GetByName(imagesRep.Targets[0])
	if err != nil {
		r.SendNotFoundError(fmt.Errorf("targets registry %v not found, err: %v", imagesRep.Targets[0], err))
		return
	}

	if registry == nil {
		r.SendNotFoundError(fmt.Errorf("targets registry %v not found(nil), please integrate the target registry firstly", imagesRep.Targets[0]))
		return
	}

	policy := &model.Policy{
		// ID:          0,
		// Name:        "images-replication",
		// Description: "will not store the policy to database",
		Creator:     "cargo",
		SrcRegistry: nil,
		DestRegistry: &model.Registry{
			ID:   registry.ID,
			Name: imagesRep.Targets[0],
			Type: model.RegistryTypeHarbor,
		},
		DestNamespace: "",
		Filters:       []*model.Filter{},
		Trigger: &model.Trigger{
			Type: model.TriggerTypeManual,
		},
		Deletion: false,
		// If override the image tag
		Override: true,
		Enabled:  true,
	}
	execution := &models.Execution{}
	if err := r.DecodeJSONReq(execution); err != nil {
		r.SendBadRequestError(err)
		return
	}

	if policy == nil {
		r.SendNotFoundError(fmt.Errorf("policy %d not found", execution.PolicyID))
		return
	}
	if !policy.Enabled {
		r.SendBadRequestError(fmt.Errorf("the policy %d is disabled", execution.PolicyID))
		return
	}
	if err = event.PopulateRegistries(replication.RegistryMgr, policy); err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to populate registries for policy %d: %v", execution.PolicyID, err))
		return
	}
	resources := make([]*model.Resource, 0)

	for _, image := range imagesRep.Images {
		recource, err := imageToResource(image, registry)
		if err != nil {
			r.SendNotFoundError(fmt.Errorf("image %s is malformed: %v", image, err))
			return
		}
		resources = append(resources, recource)
	}

	trigger := r.GetString("trigger", string(model.TriggerTypeManual))
	executionID, err := replication.OperationCtl.StartReplicationWithResources(policy, model.TriggerType(trigger), resources...)
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to start replication for policy %d: %v", execution.PolicyID, err))
		return
	}

	rsp := model.ImagesReplicationRsp{
		UUID: strconv.FormatInt(executionID, 10),
	}
	r.WriteJSONData(rsp)
}

func imageToResource(image string, registry *model.Registry) (*model.Resource, error) {
	var repository, tag string
	if i := strings.LastIndex(image, ":"); i != -1 {
		repository = image[:i]
		tag = image[i+1:]
	}

	repoMeta := make(map[string]interface{})

	repoRecord, err := common_dao.GetRepositoryByName(repository)
	if err != nil {
		return nil, err
	}

	projectMeta, err := common_dao.GetProjectMetadata(repoRecord.ProjectID)
	if err != nil {
		return nil, err
	}

	for _, v := range projectMeta {
		repoMeta[v.Name] = v.Value
	}

	return &model.Resource{
		Type:     model.ResourceTypeImage,
		Registry: registry,
		Metadata: &model.ResourceMetadata{
			Repository: &model.Repository{
				Name:     repository,
				Metadata: repoMeta,
			},
			Vtags: []string{tag},
		},
	}, nil

}

// GetImagesExecutionStatus get status of images execution
func (r *ReplicationOperationAPI) GetImagesExecutionStatus() {
	query := &models.TaskQuery{
		ExecutionID:  r.execution.ID,
		ResourceType: r.GetString("resource_type"),
	}
	status := r.GetString("status")
	if len(status) > 0 {
		query.Statuses = []string{status}
	}
	page, size, err := r.GetPaginationParams()
	if err != nil {
		r.SendBadRequestError(err)
		return
	}
	query.Page = page
	query.Size = size

	_, tasks, err := replication.OperationCtl.ListTasks(query)
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to list tasks: %v", err))
		return
	}

	rsp := model.ImagesReplicationStatus{
		Status: model.ImagesReplicateSucceed,
	}
	var hasProcessing, hasFailed bool

	for _, task := range tasks {
		rsp.JobsStatus = append(rsp.JobsStatus, model.ImageRepStatus{
			// convert 'library/alpine:[3.7]' to 'library/alpine:3.7'
			Image: strings.ReplaceAll(strings.ReplaceAll(task.SrcResource, "[", ""), "]", ""),
			// compatible with Harbor 1.7
			Status: changeStatus(task.Status),
		})

		s := convertStatus(task.Status)
		if s == model.ImagesReplicateProcessing {
			hasProcessing = true
		} else if s == model.ImagesReplicateFailed {
			hasFailed = true
		}
	}

	if hasProcessing {
		rsp.Status = model.ImagesReplicateProcessing
	} else if hasFailed {
		rsp.Status = model.ImagesReplicateFailed
	}

	r.WriteJSONData(rsp)
}

// changeStatus changes the task status in Harbor 1.8.6
//		[https://github.com/goharbor/harbor/blob/v1.8.6/src/replication/dao/models/execution.go#L27-L33]
//	to job status in Harbor 1.7.7
//      [https://github.com/goharbor/harbor/blob/v1.7.7/src/common/models/job.go#L17-L36]
func changeStatus(taskStatus string) string {
	switch taskStatus {
	case models.TaskStatusFailed:
		return "error"
	case models.TaskStatusInitialized:
		return "scheduled"
	case models.TaskStatusInProgress:
		return "running"
	case models.TaskStatusPending:
		return "pending"
	case models.TaskStatusStopped:
		return "stopped"
	case models.TaskStatusSucceed:
		return "finished"
	default:
		return taskStatus
	}
	return taskStatus
}

func convertStatus(taskStatus string) string {
	switch taskStatus {
	case models.TaskStatusFailed, models.TaskStatusStopped:
		return model.ImagesReplicateFailed
	case models.TaskStatusInitialized, models.TaskStatusInProgress, models.TaskStatusPending:
		return model.ImagesReplicateProcessing
	case models.TaskStatusSucceed:
		return model.ImagesReplicateSucceed
	default:
		return taskStatus
	}
	return taskStatus
}

// GetExecution gets one execution of the replication
func (r *ReplicationOperationAPI) GetExecution() {
	r.WriteJSONData(r.execution)
}

// StopExecution stops one execution of the replication
func (r *ReplicationOperationAPI) StopExecution() {
	if err := replication.OperationCtl.StopReplication(r.execution.ID); err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to stop execution %d: %v", r.execution.ID, err))
		return
	}
}

// ListTasks ...
func (r *ReplicationOperationAPI) ListTasks() {
	query := &models.TaskQuery{
		ExecutionID:  r.execution.ID,
		ResourceType: r.GetString("resource_type"),
	}
	status := r.GetString("status")
	if len(status) > 0 {
		query.Statuses = []string{status}
	}
	page, size, err := r.GetPaginationParams()
	if err != nil {
		r.SendBadRequestError(err)
		return
	}
	query.Page = page
	query.Size = size
	total, tasks, err := replication.OperationCtl.ListTasks(query)
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to list tasks: %v", err))
		return
	}
	r.SetPaginationHeader(total, query.Page, query.Size)
	r.WriteJSONData(tasks)
}

// GetTaskLog ...
func (r *ReplicationOperationAPI) GetTaskLog() {
	logBytes, err := replication.OperationCtl.GetTaskLog(r.task.ID)
	if err != nil {
		if httpErr, ok := err.(*common_http.Error); ok {
			if ok && httpErr.Code == http.StatusNotFound {
				r.SendNotFoundError(fmt.Errorf("the log of task %d not found", r.task.ID))
				return
			}
		}
		r.SendInternalServerError(fmt.Errorf("failed to get log of task %d: %v", r.task.ID, err))
		return
	}
	r.Ctx.ResponseWriter.Header().Set(http.CanonicalHeaderKey("Content-Length"), strconv.Itoa(len(logBytes)))
	r.Ctx.ResponseWriter.Header().Set(http.CanonicalHeaderKey("Content-Type"), "text/plain")
	_, err = r.Ctx.ResponseWriter.Write(logBytes)
	if err != nil {
		r.SendInternalServerError(fmt.Errorf("failed to write log of task %d: %v", r.task.ID, err))
		return
	}
}
