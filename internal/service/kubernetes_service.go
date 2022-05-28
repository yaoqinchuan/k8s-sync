package service

import (
	"context"
	"fmt"
	"k8s-sync/internal/model"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateStatefulSet(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) (*appsV1.StatefulSet, error) {
	var containers []coreV1.Container
	for i := 0; i < len(workspaceSpecModel.Containers); i++ {
		var envs []coreV1.EnvVar
		for k, v := range workspaceSpecModel.Containers[i].Envs {
			env := coreV1.EnvVar{
				Name:  k,
				Value: v,
			}
			envs = append(envs, env)
		}

		request := coreV1.ResourceList{}
		request[coreV1.ResourceMemory] = resource.MustParse(workspaceSpecModel.Containers[i].ResourceRequest)
		limit := coreV1.ResourceList{}
		limit[coreV1.ResourceMemory] = resource.MustParse(workspaceSpecModel.Containers[i].ResourceLimit)
		container := coreV1.Container{
			Env:             envs,
			Command:         workspaceSpecModel.Containers[i].Command,
			Args:            workspaceSpecModel.Containers[i].Args,
			Name:            workspaceSpecModel.Containers[i].Name,
			Image:           workspaceSpecModel.Containers[i].ImageName,
			ImagePullPolicy: coreV1.PullIfNotPresent,
			VolumeMounts:    workspaceSpecModel.Containers[i].VolumeMount,
			Resources: coreV1.ResourceRequirements{
				Requests: request,
				Limits:   limit,
			},
		}
		if 0 != len(workspaceSpecModel.Containers[i].PostStart) {
			container.Lifecycle = &coreV1.Lifecycle{
				PostStart: &coreV1.LifecycleHandler{
					Exec: &coreV1.ExecAction{
						Command: workspaceSpecModel.Containers[i].PostStart,
					},
				},
			}
		}

		if 0 != workspaceSpecModel.Containers[i].ContainerPort {
			container.Ports = []coreV1.ContainerPort{
				{
					Name:          "http",
					Protocol:      coreV1.ProtocolTCP,
					ContainerPort: workspaceSpecModel.Containers[i].ContainerPort,
				},
			}
		}
		containers = append(containers, container)
	}

	var tolerations []coreV1.Toleration
	for k, v := range workspaceSpecModel.Tolerations {
		toleration := coreV1.Toleration{
			Key:      k,
			Operator: coreV1.TolerationOpExists,
			Value:    v,
		}
		tolerations = append(tolerations, toleration)
	}

	var selectors []coreV1.NodeSelectorRequirement
	for k, v := range workspaceSpecModel.NodeLabels {
		selector := coreV1.NodeSelectorRequirement{
			Key:      k,
			Operator: coreV1.NodeSelectorOpIn,
			Values: []string{
				v,
			},
		}
		selectors = append(selectors, selector)
	}
	stateSet := &appsV1.StatefulSet{
		ObjectMeta: metaV1.ObjectMeta{
			Name: fmt.Sprintf("sts-%v", workspaceSpecModel.Name),
		},
		Spec: appsV1.StatefulSetSpec{
			Selector: &metaV1.LabelSelector{
				MatchLabels: workspaceSpecModel.Labels,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Name:   workspaceSpecModel.Name,
					Labels: workspaceSpecModel.Labels,
				},
				Spec: coreV1.PodSpec{
					Volumes:       workspaceSpecModel.Volumes,
					Containers:    containers,
					RestartPolicy: coreV1.RestartPolicyAlways,
					Tolerations:   tolerations,
					Affinity: &coreV1.Affinity{
						NodeAffinity: &coreV1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &coreV1.NodeSelector{
								NodeSelectorTerms: []coreV1.NodeSelectorTerm{
									{
										MatchExpressions: selectors,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return clientSet.AppsV1().StatefulSets(workspaceSpecModel.NameSpace).Create(ctx, stateSet, metaV1.CreateOptions{})
}

func CreateService(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) (*coreV1.Service, error) {
	service := &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name: fmt.Sprintf("svc-%v", workspaceSpecModel.Name),
		},
		Spec: coreV1.ServiceSpec{
			Type:     coreV1.ServiceTypeNodePort,
			Selector: workspaceSpecModel.Labels,
			Ports: []coreV1.ServicePort{
				{
					Name:     "http",
					Port:     workspaceSpecModel.NodePort,
					Protocol: coreV1.ProtocolTCP,
				},
			},
		},
	}
	return clientSet.CoreV1().Services(workspaceSpecModel.NameSpace).Create(ctx, service, metaV1.CreateOptions{})
}

func CreatePVC(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) (*coreV1.PersistentVolumeClaim, error) {
	if "" == workspaceSpecModel.PVCCapacity {
		return nil, nil
	}
	requests := coreV1.ResourceList{}
	requests[coreV1.ResourceMemory] = resource.MustParse(workspaceSpecModel.PVCCapacity)
	accessModes := make([]coreV1.PersistentVolumeAccessMode, 1)
	accessModes = append(accessModes, coreV1.ReadWriteOnce)
	pvc := &coreV1.PersistentVolumeClaim{
		ObjectMeta: metaV1.ObjectMeta{
			Name: fmt.Sprintf("pvc-%v", workspaceSpecModel.Name),
		},
		Spec: coreV1.PersistentVolumeClaimSpec{
			AccessModes: accessModes,
			Resources: coreV1.ResourceRequirements{
				Requests: requests,
			},
			VolumeName:       fmt.Sprintf("pvc-%v", workspaceSpecModel.Name),
			StorageClassName: &workspaceSpecModel.StorageClass,
		},
	}
	return clientSet.CoreV1().PersistentVolumeClaims(workspaceSpecModel.NameSpace).Create(ctx, pvc, metaV1.CreateOptions{})
}

func deleteStatefulSet(ctx context.Context, clientSet *kubernetes.Clientset, workspaceName string, namespace string) error {
	return clientSet.AppsV1().StatefulSets(namespace).Delete(ctx, fmt.Sprintf("sts-%v", workspaceName), metaV1.DeleteOptions{})
}

func deleteService(ctx context.Context, clientSet *kubernetes.Clientset, workspaceName string, namespace string) {
	clientSet.CoreV1().Services(namespace).Delete(ctx, fmt.Sprintf("svc-%v", workspaceName), metaV1.DeleteOptions{})
}
