package service

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-sync/internal/model"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

func GenerateStatefulSet(workspaceSpecModel *model.WorkspaceSpecModel) *appsV1.StatefulSet {
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
		TypeMeta: metaV1.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      fmt.Sprintf("sts-%v", workspaceSpecModel.Name),
			Namespace: workspaceSpecModel.NameSpace,
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
	return stateSet
}

func doStartWorkspace(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) error {
	_, err := clientSet.AppsV1().StatefulSets(workspaceSpecModel.NameSpace).Create(ctx, GenerateStatefulSet(workspaceSpecModel), metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	if pvc := GeneratePVC(workspaceSpecModel); pvc != nil {
		_, err := clientSet.CoreV1().PersistentVolumeClaims(workspaceSpecModel.NameSpace).Create(ctx, pvc, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	_, err = clientSet.CoreV1().Services(workspaceSpecModel.NameSpace).Create(ctx, GenerateService(workspaceSpecModel), metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func doRestoreWorkspace(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) error {
	_, err := clientSet.AppsV1().StatefulSets(workspaceSpecModel.NameSpace).Create(ctx, GenerateStatefulSet(workspaceSpecModel), metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func CheckWorkspaceRunning(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) (bool, error) {
	_, err := clientSet.CoreV1().Pods(workspaceSpecModel.NameSpace).Get(ctx, fmt.Sprintf("sts-%v-0", workspaceSpecModel.Name), metaV1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func RestoreWorkspace(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) error {
	_, err := clientSet.AppsV1().StatefulSets(workspaceSpecModel.NameSpace).Create(ctx, GenerateStatefulSet(workspaceSpecModel), metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	_, err = clientSet.CoreV1().Services(workspaceSpecModel.NameSpace).Create(ctx, GenerateService(workspaceSpecModel), metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func doDeleteWorkspace(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) error {
	err := clientSet.AppsV1().StatefulSets(workspaceSpecModel.NameSpace).Delete(ctx, fmt.Sprintf("sts-%v", workspaceSpecModel.Name), metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().PersistentVolumeClaims(workspaceSpecModel.NameSpace).Delete(ctx, fmt.Sprintf("pvc-%v", workspaceSpecModel.Name), metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().Services(workspaceSpecModel.NameSpace).Delete(ctx, fmt.Sprintf("svc-%v", workspaceSpecModel.Name), metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func doStopWorkspace(ctx context.Context, clientSet *kubernetes.Clientset, workspaceSpecModel *model.WorkspaceSpecModel) error {
	err := clientSet.AppsV1().StatefulSets(workspaceSpecModel.NameSpace).Delete(ctx, fmt.Sprintf("sts-%v", workspaceSpecModel.Name), metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().Services(workspaceSpecModel.NameSpace).Delete(ctx, fmt.Sprintf("svc-%v", workspaceSpecModel.Name), metaV1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GenerateWorkspace(workspaceSpecModel *model.WorkspaceSpecModel) (*coreV1.List, error) {
	statefulSet, err := json.Marshal(GenerateStatefulSet(workspaceSpecModel))
	if nil != err {
		return nil, err
	}

	service, err := json.Marshal(GenerateService(workspaceSpecModel))
	if nil != err {
		return nil, err
	}
	var pvc []byte
	if pvcGenerated := GeneratePVC(workspaceSpecModel); pvcGenerated != nil {
		pvc, err = json.Marshal(pvcGenerated)
		if nil != err {
			return nil, err
		}
	}

	workspaceSpec := coreV1.List{
		TypeMeta: metaV1.TypeMeta{Kind: "List", APIVersion: "v1"},
		Items: []runtime.RawExtension{
			{Raw: statefulSet}, {
				Raw: service,
			}},
	}
	if 0 != len(pvc) {
		workspaceSpec.Items = append(workspaceSpec.Items, runtime.RawExtension{
			Raw: service,
		})
	}
	return &workspaceSpec, nil
}

func GenerateService(workspaceSpecModel *model.WorkspaceSpecModel) *coreV1.Service {
	service := &coreV1.Service{
		TypeMeta: metaV1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      fmt.Sprintf("svc-%v", workspaceSpecModel.Name),
			Namespace: workspaceSpecModel.NameSpace,
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
	return service
}

func GeneratePVC(workspaceSpecModel *model.WorkspaceSpecModel) *coreV1.PersistentVolumeClaim {
	if "" == workspaceSpecModel.PVCCapacity {
		return nil
	}
	requests := coreV1.ResourceList{}
	requests[coreV1.ResourceMemory] = resource.MustParse(workspaceSpecModel.PVCCapacity)
	accessModes := make([]coreV1.PersistentVolumeAccessMode, 1)
	accessModes = append(accessModes, coreV1.ReadWriteOnce)
	pvc := &coreV1.PersistentVolumeClaim{
		TypeMeta: metaV1.TypeMeta{Kind: "PersistentVolumeClaim", APIVersion: "v1"},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      fmt.Sprintf("pvc-%v", workspaceSpecModel.Name),
			Namespace: workspaceSpecModel.NameSpace,
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
	return pvc
}
