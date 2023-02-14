/*
Copyright 2022 SUSE.
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package ignition aggregates all Ignition flavors into a single package to be consumed
// by the bootstrap provider by exposing an API similar to 'internal/cloudinit' package.
package ignition

import (
	"fmt"

	bootstrapv1 "github.com/rancher-sandbox/cluster-api-provider-rke2/bootstrap/api/v1alpha1"
	"github.com/rancher-sandbox/cluster-api-provider-rke2/bootstrap/internal/cloudinit"
	"github.com/rancher-sandbox/cluster-api-provider-rke2/bootstrap/internal/ignition/clc"
)

// NodeInput defines the context to generate a node user data.
type NodeInput struct {
	*cloudinit.BaseUserData
	Ignition *bootstrapv1.IgnitionSpec
}

// ControlPlaneJoinInput defines context to generate controlplane instance user data for control plane node join.
type ControlPlaneJoinInput struct {
	*cloudinit.ControlPlaneInput
	Ignition *bootstrapv1.IgnitionSpec
}

// ControlPlaneInput defines the context to generate a controlplane instance user data.
type ControlPlaneInitInput struct {
	*cloudinit.ControlPlaneInput
	Ignition *bootstrapv1.IgnitionSpec
}

// NewNode returns Ignition configuration for new worker node joining the cluster.
func NewNode(input *NodeInput) ([]byte, string, error) {
	if input == nil {
		return nil, "", fmt.Errorf("input can't be nil")
	}

	if input.BaseUserData == nil {
		return nil, "", fmt.Errorf("node input can't be nil")
	}

	input.WriteFiles = append(input.WriteFiles, input.WriteFiles...)
	if input.AirGapped {
		input.RKE2Command = "INSTALL_RKE2_ARTIFACT_PATH=/opt/rke2-artifacts INSTALL_RKE2_TYPE=\"agent\" sh /opt/install.sh"
	} else {
		input.RKE2Command = fmt.Sprintf("curl -sfL https://get.rke2.io | INSTALL_RKE2_VERSION=%[1]s INSTALL_RKE2_TYPE=\"agent\" sh -s -", input.RKE2Version)
	}

	return render(input.BaseUserData, input.Ignition)
}

// NewJoinControlPlane returns Ignition configuration for new controlplane node joining the cluster.
func NewJoinControlPlane(input *ControlPlaneJoinInput) ([]byte, string, error) {
	if input == nil {
		return nil, "", fmt.Errorf("input can't be nil")
	}

	if input.ControlPlaneInput == nil {
		return nil, "", fmt.Errorf("controlplane join input can't be nil")
	}

	input.WriteFiles = input.Certificates.AsFiles()
	input.WriteFiles = append(input.WriteFiles, input.WriteFiles...)
	if input.AirGapped {
		input.RKE2Command = "INSTALL_RKE2_ARTIFACT_PATH=/opt/rke2-artifacts INSTALL_RKE2_TYPE=\"server\" sh /opt/install.sh"
	} else {
		input.RKE2Command = fmt.Sprintf("curl -sfL https://get.rke2.io | INSTALL_RKE2_VERSION=%[1]s INSTALL_RKE2_TYPE=\"server\" sh -s -", input.RKE2Version)
	}

	return render(&input.BaseUserData, input.Ignition)
}

// NewInitControlPlane returns Ignition configuration for bootstrapping new cluster.
func NewInitControlPlane(input *ControlPlaneInitInput) ([]byte, string, error) {
	if input == nil {
		return nil, "", fmt.Errorf("input can't be nil")
	}

	if input.ControlPlaneInput == nil {
		return nil, "", fmt.Errorf("controlplane input can't be nil")
	}

	input.WriteFiles = input.Certificates.AsFiles()
	input.WriteFiles = append(input.WriteFiles, input.WriteFiles...)
	if input.AirGapped {
		input.RKE2Command = "INSTALL_RKE2_ARTIFACT_PATH=/opt/rke2-artifacts INSTALL_RKE2_TYPE=\"server\" sh /opt/install.sh"
	} else {
		input.RKE2Command = fmt.Sprintf("curl -sfL https://get.rke2.io | INSTALL_RKE2_VERSION=%[1]s INSTALL_RKE2_TYPE=\"server\" sh -s -", input.RKE2Version)
	}
	//kubeadmConfig := fmt.Sprintf("%s\n---\n%s", input.ClusterConfiguration, input.InitConfiguration)

	return render(&input.BaseUserData, input.Ignition)
}

func render(input *cloudinit.BaseUserData, ignitionConfig *bootstrapv1.IgnitionSpec) ([]byte, string, error) {
	clcConfig := &bootstrapv1.ContainerLinuxConfig{}
	if ignitionConfig != nil && ignitionConfig.ContainerLinuxConfig != nil {
		clcConfig = ignitionConfig.ContainerLinuxConfig
	}

	return clc.Render(input, clcConfig)
}
