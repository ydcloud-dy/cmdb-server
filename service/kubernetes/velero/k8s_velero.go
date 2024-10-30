package velero

import (
	"DYCLOUD/global"
	veleroReq "DYCLOUD/model/velero/request"
	"DYCLOUD/utils/kubernetes"
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/utils/pointer"
	"log"
	"strconv"
	"time"
)

type K8sVeleroService struct {
	kubernetes.BaseService
}

func (k *K8sVeleroService) CreateVelero(req *veleroReq.VeleroModel, uuid uuid.UUID) error {
	k8s, err := k.Generic(req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}

	client, err := k8s.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	//config, err := k8s.Config()
	//if err != nil {
	//	global.DYCLOUD_LOG.Error(err.Error())
	//	return err
	//}
	//veleroClientset, err := veleroversioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Velero clientset: %v", err)
	}
	namespace := createNamespace()
	serviceAccount := createServiceAccount()
	//clusterRole := createClusterRole()
	clusterRoleBindg := createClusterRolebinding()
	configMap := createConfigMap(req, namespace.Name)
	job := createJobSpec(req)
	// Create resources
	ctx := context.TODO()
	ns, err := client.CoreV1().Namespaces().Get(ctx, namespace.Name, metav1.GetOptions{})
	if err == nil {
		updateFrontendStatus("正在删除 Namespace")
		options := metav1.DeleteOptions{
			GracePeriodSeconds: pointer.Int64Ptr(0),
		}
		k.DeleteVeleroNamespace(namespace.Name)
		//err = client.RbacV1().ClusterRoles().Delete(ctx, clusterRole.Name, options)
		err = client.RbacV1().ClusterRoleBindings().Delete(ctx, clusterRoleBindg.Name, options)
		err = client.CoreV1().Namespaces().Delete(ctx, namespace.Name, options)

		time.Sleep(10 * time.Second)
		updateFrontendStatus("velero命名空间已存在，正在清除旧velero环境")
	}
	fmt.Println(ns)
	_, err = client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {

		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	_, err = client.CoreV1().ConfigMaps(namespace.Name).Create(ctx, configMap, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	_, err = client.CoreV1().ServiceAccounts(namespace.Name).Create(ctx, serviceAccount, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	//_, err = client.RbacV1().ClusterRoles().Create(ctx, clusterRole, metav1.CreateOptions{})
	//if err != nil {
	//	global.DYCLOUD_LOG.Error(err.Error())
	//	return err
	//}
	_, err = client.RbacV1().ClusterRoleBindings().Create(ctx, clusterRoleBindg, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}

	_, err = client.BatchV1().Jobs(namespace.Name).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}

	fmt.Println("Velero resources created successfully.")
	return nil
}
func updateFrontendStatus(status string) {
	// 在这里实现更新前端状态的逻辑
	// 这可以是通过 WebSocket, HTTP 请求, 或者任何其他方法
	fmt.Println(status)
}
func createNamespace() *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "velero",
		},
	}
}
func createConfigMap(req *veleroReq.VeleroModel, namespace string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "velero-install" + strconv.Itoa(req.ClusterId),
			Namespace: namespace,
		},
		Data: map[string]string{
			"s3-secret": fmt.Sprintf("[default]\r\naws_access_key_id = %s\r\naws_secret_access_key = %s", req.S3Key, req.S3Secret),
		},
	}
}
func createServiceAccount() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "velero",
			Namespace: "velero",
		},
	}
}

func createClusterRolebinding() *rbacV1.ClusterRoleBinding {
	return &rbacV1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "velero",
		},
		Subjects: []rbacV1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "velero",
				Namespace: "velero",
			},
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
	}
}
func createJobSpec(req *veleroReq.VeleroModel) *v1.Job {
	return &v1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "velero-install" + strconv.Itoa(req.ClusterId),
			Namespace: "velero",
			Labels: map[string]string{
				"kubeasy.com/create": "true",
			},
			Annotations: map[string]string{
				"kubeasy.com/create": "true",
			},
		},
		Spec: v1.JobSpec{
			Parallelism:    int32Ptr(1),
			Completions:    int32Ptr(1),
			BackoffLimit:   int32Ptr(6),
			ManualSelector: boolPtr(false),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"job-name": "velero-install" + strconv.Itoa(req.ClusterId),
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "s3-secret",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "velero-install" + strconv.Itoa(req.ClusterId),
									},
									DefaultMode: int32Ptr(420),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "install",
							Image: req.VeleroImage,
							Args: []string{
								"/velero",
								"install",
								"--provider",
								"aws",
								"credentials-velero",
								"--image",
								req.VeleroImage,
								"--plugins",
								req.PluginImage,
								"--bucket",
								req.S3Bucket,
								"--prefix",
								"",
								"--secret-file",
								"/tmp/credentials/s3-secret",
								"--use-volume-snapshots",
								"false",
								"--backup-location-config",
								fmt.Sprintf("region=%s,s3ForcePathStyle=true,s3Url=%s", req.S3Region, req.S3Address),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "s3-secret",
									MountPath: "/tmp/credentials/",
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: corev1.TerminationMessageReadFile,
							ImagePullPolicy:          corev1.PullAlways,
						},
					},
					RestartPolicy:                 corev1.RestartPolicyNever,
					TerminationGracePeriodSeconds: int64Ptr(30),
					DNSPolicy:                     corev1.DNSClusterFirst,
					ServiceAccountName:            "velero",
					SecurityContext:               &corev1.PodSecurityContext{},
					SchedulerName:                 "default-scheduler",
				},
			},
			//CompletionModeNonIndexed
			CompletionMode: NonIndexedCompletionModePtr(),
			Suspend:        boolPtr(false),
		},
	}

}
func (k *K8sVeleroService) DeleteVeleroNamespace(namespace string) error {
	k8s, err := k.Generic(nil, uuid.UUID{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}

	client, err := k8s.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	config, err := k8s.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	ctx := context.TODO()
	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	// 创建 API 扩展客户端
	apiextensionsClientset, err := apiextensionsv1.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating apiextensions clientset: %v", err)
	}
	// 获取命名空间中的所有 CRD 实例
	crdList, err := apiextensionsClientset.CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	for _, crd := range crdList.Items {
		gvr := schema.GroupVersionResource{
			Group:    crd.Spec.Group,
			Version:  crd.Spec.Versions[0].Name,
			Resource: crd.Spec.Names.Plural,
		}

		// 列出 CRD 实例
		crdInstances, err := dynamicClient.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing CRD instances: %v", err)
		}

		// 删除 CRD 实例
		for _, instance := range crdInstances.Items {
			err := dynamicClient.Resource(gvr).Namespace(namespace).Delete(ctx, instance.GetName(), metav1.DeleteOptions{})
			if err != nil {
				log.Fatalf("Error deleting CRD instance: %v", err)
			}
			fmt.Printf("Deleted CRD instance: %s\n", instance.GetName())
		}
	}
	// 删除命名空间
	err = client.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}

	return nil
}

func NonIndexedCompletionModePtr() *v1.CompletionMode {
	mode := v1.NonIndexedCompletion
	return &mode
}
func int32Ptr(i int32) *int32 {
	return &i
}
func int64Ptr(i int64) *int64 { return &i }
func boolPtr(b bool) *bool    { return &b }
