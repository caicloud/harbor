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

package target

import (
	"fmt"
	"reflect"

	"github.com/vmware/harbor/src/common/dao"
	"github.com/vmware/harbor/src/common/models"
)

// Manager defines the methods that a target manager should implement
type Manager interface {
	GetTarget(interface{}) (*models.RepTarget, error)
}

// DefaultManager implement the Manager interface
type DefaultManager struct{}

// NewDefaultManager returns an instance of DefaultManger
func NewDefaultManager() *DefaultManager {
	return &DefaultManager{}
}

// GetTarget ...
func (d *DefaultManager) GetTarget(idOrName interface{}) (*models.RepTarget, error) {
	var target *models.RepTarget
	var err error
	switch v := idOrName.(type) {
	case int64:
		target, err = dao.GetRepTarget(v)
	case string:
		target, err = dao.GetRepTargetByName(v)
	default:
		return nil, fmt.Errorf("idOrName should have type string or int64, but got %v", reflect.TypeOf(idOrName))
	}

	if err != nil {
		return nil, err
	}

	if target == nil {
		return nil, fmt.Errorf("target '%v' does not exist", idOrName)
	}

	return target, nil
}