package model

import v1 "k8s.io/api/core/v1"

type WorkspaceSpecModel struct {
	Name         string            `json:"name" v:"required"`
	NameSpace    string            `json:"namespace" v:"required"`
	StorageClass string            `json:"storageClass" v:"required"`
	NodePort     int32             `json:"nodePort" v:"required"`
	Labels       map[string]string `json:"labels" v:"required"`
	Annotations  map[string]string `json:"annotations" v:"required"`
	NodeLabels   map[string]string `json:"nodeLabels" v:"required"`
	Tolerations  map[string]string `json:"tolerations" v:"required"`
	Volumes      []v1.Volume       `json:"volumes" v:"required"`
	PVCCapacity  string            `json:"PVCCapacity" v:"required"`
	Containers   []ContainerModel  `json:"containers" v:"required"`
}
