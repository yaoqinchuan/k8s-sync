package model

import v1 "k8s.io/api/core/v1"
import coreV1 "k8s.io/api/core/v1"

type ContainerModel struct {
	Name            string            `json:"name" v:"required"`
	ImageName       string            `json:"imageName" v:"required"`
	Envs            map[string]string `json:"envs" v:"required"`
	PostStart       []string          `json:"postStart" v:"required"`
	Command         []string          `json:"command" v:"required"`
	Args            []string          `json:"args" v:"required"`
	ContainerPort   int32             `json:"containerPort" v:"required"`
	ResourceLimit   string            `json:"resourceLimit" v:"required"`
	ResourceRequest string            `json:"resourceRequest" v:"required"`
	VolumeMount     []v1.VolumeMount  `json:"volumeMount" v:"required"`
	LivenessProbe   coreV1.Probe      `json:"livenessProbe" v:"required"`
	ReadinessProbe  coreV1.Probe      `json:"readinessProbe" v:"required"`
}
