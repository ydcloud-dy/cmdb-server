package kubernetes

import (
	"DYCLOUD/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// K8sResourceService 定义一个通用的 Kubernetes 资源服务
type K8sResourceService struct{}

// GetK8sResourceList 获取通用的 Kubernetes 资源列表
func (k *K8sResourceService) GetK8sResourceList(client *kubernetes.Clientset, resourceType, namespace,
	keyword string, listFunc func(string, metav1.ListOptions) (interface{}, error)) (interface{}, error) {

	// 构建缓存键
	cacheKey := fmt.Sprintf("%s:%s:%s", resourceType, namespace, keyword)
	//if namespace == "" {
	//	cacheKey = fmt.Sprintf("%s::%s", resourceType, namespace)
	//}
	// 从 Redis 缓存中获取数据
	resourceJSON, err := global.DYCLOUD_REDIS.Get(context.TODO(), cacheKey).Result()
	if err == redis.Nil {
		// 如果缓存中没有数据，从 Kubernetes API 获取资源并缓存
		return k.fetchAndCacheResource(client, resourceType, namespace, keyword, cacheKey, listFunc)
	} else if err != nil {
		// 如果发生错误，返回错误信息
		return nil, err
	}
	// 根据 resourceType 创建相应的空列表对象
	var resources interface{}
	var resourcesNew interface{}
	switch resourceType {
	case "deployments":
		//if namespace == "" {
		//	resources = &[]v1.Deployment{}
		//} else {
		resources = &v1.DeploymentList{}
		resourcesNew = &[]v1.Deployment{}
		//}
	// 其他资源类型的处理（如 statefulsets, configmaps, services 等）
	// 可以在这里添加更多 case 语句
	default:
		return nil, fmt.Errorf("unsupported resource type")
	}
	// 将 JSON 数据解码为相应的资源列表
	err = json.Unmarshal([]byte(resourceJSON), resources)
	if err != nil {
		err = json.Unmarshal([]byte(resourceJSON), resourcesNew)
		if err != nil {
			return nil, err
		}
		return resourcesNew, nil
	}
	// 返回解码后的资源列表
	return resources, nil
}

func (k *K8sResourceService) DeleteK8sResourceCache(resourceType, namespace, keyword string) error {
	//var cacheKey string
	//if namespace == "" {
	// 如果命名空间为空，删除所有命名空间的缓存数据
	cacheKeyPattern := fmt.Sprintf("%s:*:%s", resourceType, keyword)
	keys, err := global.DYCLOUD_REDIS.Keys(context.TODO(), cacheKeyPattern).Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		err = global.DYCLOUD_REDIS.Del(context.TODO(), key).Err()
		if err != nil {
			return err
		}
	}
	//} else {
	//	// 如果命名空间不为空，只删除特定命名空间的缓存数据
	//	cacheKey = fmt.Sprintf("%s:%s:%s", resourceType, namespace, keyword)
	//	err := global.DYCLOUD_REDIS.Del(context.TODO(), cacheKey).Err()
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}
func (k *K8sResourceService) DeleteAllK8sResourceCache() error {
	// 使用 Redis 的 FlushAll 清空所有缓存
	err := global.DYCLOUD_REDIS.FlushAll(context.TODO()).Err()
	if err != nil {
		return err
	}
	return nil
}
func (k *K8sResourceService) fetchAndCacheResource(client *kubernetes.Clientset, resourceType,
	namespace, keyword, cacheKey string, listFunc func(string, metav1.ListOptions) (interface{}, error)) (interface{}, error) {
	var resourceList interface{}
	var err error

	if namespace == "" {
		// 处理所有命名空间的情况
		namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var allResources []interface{}
		for _, ns := range namespaceList.Items {
			nsResourceList, err := listFunc(ns.Name, metav1.ListOptions{})
			if err != nil {
				return nil, err
			}
			allResources = append(allResources, nsResourceList)
		}
		resourceList = allResources
	} else {
		// 处理单个命名空间的情况
		resourceList, err = listFunc(namespace, metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
	}
	// 过滤资源列表，如果需要关键词过滤
	filteredResources, err := k.filterResourcesByKeyword(resourceList, keyword)
	if err != nil {
		return nil, err
	}

	// 将过滤后的资源列表序列化为 JSON 字符串
	resourceJSON, err := json.Marshal(filteredResources)
	if err == nil {
		// 如果序列化成功，将数据存储到 Redis 缓存中
		err = global.DYCLOUD_REDIS.Set(context.TODO(), cacheKey, resourceJSON, 0).Err()
	}

	// 返回过滤后的资源列表和可能发生的错误
	return filteredResources, err
}

// filterResourcesByKeyword 根据关键词过滤资源列表
func (k *K8sResourceService) filterResourcesByKeyword(resourceList interface{}, keyword string) (interface{}, error) {
	// 如果没有关键词，则不需要过滤
	if keyword == "" {
		return resourceList, nil
	}

	// 根据不同资源类型进行关键词匹配过滤
	switch resources := resourceList.(type) {
	case *v1.DeploymentList:
		var filtered []v1.Deployment
		for _, resource := range resources.Items {
			if strings.Contains(resource.Name, keyword) {
				filtered = append(filtered, resource)
			}
		}
		resources.Items = filtered
		return resources, nil
	case *v1.StatefulSetList:
		var filtered []v1.StatefulSet
		for _, resource := range resources.Items {
			if strings.Contains(resource.Name, keyword) {
				filtered = append(filtered, resource)
			}
		}
		resources.Items = filtered

		return resources, nil
	case *corev1.ConfigMapList:
		var filtered []corev1.ConfigMap
		for _, resource := range resources.Items {
			if strings.Contains(resource.Name, keyword) {
				filtered = append(filtered, resource)
			}
		}
		resources.Items = filtered
		return resources, nil
	case *corev1.ServiceList:
		var filtered []corev1.Service
		for _, resource := range resources.Items {
			if strings.Contains(resource.Name, keyword) {
				filtered = append(filtered, resource)
			}
		}
		resources.Items = filtered

		return resources, nil
	default:
		return nil, fmt.Errorf("unsupported resource type")
	}
}
