// Copyright © 2021 Alibaba Group Holding Ltd.
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

package processor

import (
	"fmt"

	"github.com/alibaba/sealer/pkg/clusterfile"

	"github.com/alibaba/sealer/pkg/filesystem/cloudfilesystem"
	"github.com/alibaba/sealer/pkg/filesystem/cloudimage"

	"github.com/alibaba/sealer/utils"

	"github.com/alibaba/sealer/pkg/plugin"

	"github.com/alibaba/sealer/common"

	"github.com/alibaba/sealer/pkg/filesystem"
	"github.com/alibaba/sealer/pkg/runtime"
	v2 "github.com/alibaba/sealer/types/api/v2"
)

type DeleteProcessor struct {
	cloudImageMounter cloudimage.Interface
	ClusterFile       clusterfile.Interface
}

func (d DeleteProcessor) Reset(cluster *v2.Cluster) error {
	runTime, err := runtime.NewDefaultRuntime(cluster, d.ClusterFile.GetKubeadmConfig())
	if err != nil {
		return fmt.Errorf("failed to init runtime, %v", err)
	}

	return runTime.Reset()
}

func (d DeleteProcessor) GetPipeLine() ([]func(cluster *v2.Cluster) error, error) {
	var todoList []func(cluster *v2.Cluster) error
	todoList = append(todoList,
		d.Reset,
		d.ApplyCleanPlugin,
		d.UnMountRootfs,
		d.UnMountImage,
		d.CleanFS,
	)
	return todoList, nil
}

func (d DeleteProcessor) UnMountRootfs(cluster *v2.Cluster) error {
	hosts := append(cluster.GetMasterIPList(), cluster.GetNodeIPList()...)
	config := runtime.GetRegistryConfig(common.DefaultTheClusterRootfsDir(cluster.Name), runtime.GetMaster0Ip(cluster))
	if utils.NotIn(config.IP, hosts) {
		hosts = append(hosts, config.IP)
	}
	fs, err := filesystem.NewFilesystem(common.DefaultTheClusterRootfsDir(cluster.Name))
	if err != nil {
		return err
	}

	return fs.UnMountRootfs(cluster, hosts)
}

func (d DeleteProcessor) UnMountImage(cluster *v2.Cluster) error {
	return d.cloudImageMounter.UnMountImage(cluster)
}

func (d DeleteProcessor) ApplyCleanPlugin(cluster *v2.Cluster) error {
	plugins := plugin.NewPlugins(cluster)
	if err := plugins.Dump(d.ClusterFile.GetPlugins()); err != nil {
		return err
	}
	if err := plugins.Load(); err != nil {
		return err
	}
	return plugins.Run(cluster.GetAllIPList(), plugin.PhasePostClean)
}

func (d DeleteProcessor) CleanFS(cluster *v2.Cluster) error {
	return cloudfilesystem.CleanFilesystem(cluster.Name)
}

func NewDeleteProcessor(clusterFile clusterfile.Interface) (Processor, error) {
	mounter, err := filesystem.NewCloudImageMounter()
	if err != nil {
		return nil, err
	}

	return DeleteProcessor{
		ClusterFile:       clusterFile,
		cloudImageMounter: mounter,
	}, nil
}