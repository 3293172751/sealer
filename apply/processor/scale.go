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

	"github.com/alibaba/sealer/common"
	"github.com/alibaba/sealer/pkg/clusterfile"
	"github.com/alibaba/sealer/pkg/config"
	"github.com/alibaba/sealer/pkg/plugin"

	"github.com/alibaba/sealer/pkg/filesystem/cloudfilesystem"

	"github.com/alibaba/sealer/utils"

	"github.com/alibaba/sealer/pkg/filesystem"
	"github.com/alibaba/sealer/pkg/runtime"
	v2 "github.com/alibaba/sealer/types/api/v2"
)

type ScaleProcessor struct {
	fileSystem      cloudfilesystem.Interface
	ClusterFile     clusterfile.Interface
	Runtime         runtime.Interface
	KubeadmConfig   *runtime.KubeadmConfig
	Config          config.Interface
	Plugins         plugin.Plugins
	MastersToJoin   []string
	MastersToDelete []string
	NodesToJoin     []string
	NodesToDelete   []string
	IsScaleUp       bool
}

func (s *ScaleProcessor) PreProcess(cluster *v2.Cluster) error {
	runTime, err := runtime.NewDefaultRuntime(cluster, s.KubeadmConfig)
	if err != nil {
		return fmt.Errorf("failed to init runtime, %v", err)
	}
	s.Runtime = runTime
	s.Config = config.NewConfiguration(cluster)
	if s.IsScaleUp {
		if err = utils.SaveClusterInfoToFile(cluster, cluster.Name); err != nil {
			return err
		}
	}
	return s.initPlugin(cluster)
}

func (s *ScaleProcessor) GetPipeLine() ([]func(cluster *v2.Cluster) error, error) {
	if s.IsScaleUp {
		return s.ScaleUpPipeLine()
	}
	return s.ScaleDownPipeLine()
}

func (s *ScaleProcessor) ScaleUp(cluster *v2.Cluster) error {
	pipLine, err := s.ScaleUpPipeLine()
	if err != nil {
		return err
	}

	for _, f := range pipLine {
		if err = f(cluster); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScaleProcessor) ScaleDown(cluster *v2.Cluster) error {
	pipLine, err := s.ScaleDownPipeLine()
	if err != nil {
		return err
	}

	for _, f := range pipLine {
		if err = f(cluster); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScaleProcessor) ScaleUpPipeLine() ([]func(cluster *v2.Cluster) error, error) {
	var todoList []func(cluster *v2.Cluster) error
	todoList = append(todoList,
		s.PreProcess,
		s.RunConfig,
		s.MountRootfs,
		s.GetPhasePluginFunc(plugin.PhasePreJoin),
		s.Join,
		s.GetPhasePluginFunc(plugin.PhasePostJoin),
	)
	return todoList, nil
}

func (s *ScaleProcessor) ScaleDownPipeLine() ([]func(cluster *v2.Cluster) error, error) {
	var todoList []func(cluster *v2.Cluster) error
	todoList = append(todoList,
		s.PreProcess,
		s.Delete,
		s.ApplyCleanPlugin,
		s.UnMountRootfs,
	)
	return todoList, nil
}

func (s *ScaleProcessor) initPlugin(cluster *v2.Cluster) error {
	s.Plugins = plugin.NewPlugins(cluster)
	return s.Plugins.Dump(s.ClusterFile.GetPlugins())
}

func (s *ScaleProcessor) GetPhasePluginFunc(phase plugin.Phase) func(cluster *v2.Cluster) error {
	return func(cluster *v2.Cluster) error {
		if phase == plugin.PhasePreInit {
			if err := s.Plugins.Load(); err != nil {
				return err
			}
		}
		return s.Plugins.Run(append(s.MastersToJoin, s.NodesToJoin...), phase)
	}
}

func (s *ScaleProcessor) RunConfig(cluster *v2.Cluster) error {
	return s.Config.Dump(s.ClusterFile.GetConfigs())
}

func (s *ScaleProcessor) MountRootfs(cluster *v2.Cluster) error {
	return s.fileSystem.MountRootfs(cluster, append(s.MastersToJoin, s.NodesToJoin...), true)
}

func (s *ScaleProcessor) UnMountRootfs(cluster *v2.Cluster) error {
	return s.fileSystem.UnMountRootfs(cluster, append(s.MastersToDelete, s.NodesToDelete...))
}

func (s *ScaleProcessor) Join(cluster *v2.Cluster) error {
	err := s.Runtime.JoinMasters(s.MastersToJoin)
	if err != nil {
		return err
	}
	return s.Runtime.JoinNodes(s.NodesToJoin)
}

func (s *ScaleProcessor) Delete(cluster *v2.Cluster) error {
	err := s.Runtime.DeleteMasters(s.MastersToDelete)
	if err != nil {
		return err
	}
	return s.Runtime.DeleteNodes(s.NodesToDelete)
}

func (s *ScaleProcessor) ApplyCleanPlugin(cluster *v2.Cluster) error {
	plugins := plugin.NewPlugins(cluster)
	if err := plugins.Dump(s.ClusterFile.GetPlugins()); err != nil {
		return err
	}
	if err := plugins.Load(); err != nil {
		return err
	}
	return plugins.Run(cluster.GetAllIPList(), plugin.PhasePostClean)
}

func NewScaleProcessor(kubeadmConfig *runtime.KubeadmConfig, clusterFile clusterfile.Interface, masterToJoin, masterToDelete, nodeToJoin, nodeToDelete []string) (Processor, error) {
	var up bool
	// only scale up or scale down at a time
	if len(masterToJoin) > 0 || len(nodeToJoin) > 0 {
		up = true
	}
	fs, err := filesystem.NewFilesystem(common.DefaultTheClusterRootfsDir(clusterFile.GetCluster().Name))
	if err != nil {
		return nil, err
	}
	return &ScaleProcessor{
		MastersToDelete: masterToDelete,
		MastersToJoin:   masterToJoin,
		NodesToDelete:   nodeToDelete,
		NodesToJoin:     nodeToJoin,
		KubeadmConfig:   kubeadmConfig,
		ClusterFile:     clusterFile,
		IsScaleUp:       up,
		fileSystem:      fs,
	}, nil
}