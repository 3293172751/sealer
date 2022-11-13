// Copyright © 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k3s

import (
	"fmt"
	"strings"
	"time"

	"github.com/sealerio/sealer/common"
	"github.com/sealerio/sealer/pkg/registry"
	"github.com/sealerio/sealer/pkg/runtime"
	v2 "github.com/sealerio/sealer/types/api/v2"
	"github.com/sealerio/sealer/utils/ssh"
	"github.com/sirupsen/logrus"
)

// Runtime struct is the runtime interface for k3s

type Runtime struct {
    cluster *v2.Cluster
}

func (k *Runtime) Init() error {
	return k.init()
}

func (k *Runtime) init() error {
	if k.cluster == nil {
        return fmt.Errorf("cluster is nil")
    }
	return nil
	
