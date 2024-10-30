package deployment

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/deployment"
	"DYCLOUD/service/kubernetes/cluster"
	"DYCLOUD/utils/kubernetes"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	k8s "k8s.io/client-go/kubernetes"
	cache2 "k8s.io/client-go/tools/cache"
	"log"
	"strings"
)

type K8sDeploymentService struct {
	kubernetes.BaseService
	kubernetes.KubernetesService
	DeploymentInformer cache2.SharedIndexInformer
}

// GetDeploymentList 获取deployment列表
//
// @Description 获取deployment列表
// @Tags Deployment
// @Accept json
// @Produce json
// @Param req query deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Param uuid path uuid.UUID true "集群用户的uuid"
// @Success 200 {object} []v1.Deployment

// @Router /deployment/list [get]
func (k *K8sDeploymentService) GetDeploymentList(req deployment.GetDeploymentListReq,
	uuid uuid.UUID) (deploymentList *[]v1.Deployment, total int, err error) {

	// 获取k8s客户端实例
	client, err := InformerFactory(k, req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, 0, err
	}
	// 获取所有或特定命名空间的 Deployment 列表
	var allDeploymentsInterface interface{}
	deploymentListObj, total, err := GetResourceList(allDeploymentsInterface, client, req.Namespace, &req)

	return deploymentListObj, total, nil

}

// GetResourceList 获取资源列表
//
// @Description 获取指定命名空间中的 Deployment 列表，支持分页和关键词过滤
// @Tags Deployment
// @Accept json
// @Produce json
// @Param allDeploymentsInterface body interface{} true "所有部署接口"
// @Param client body k8s.Clientset true "Kubernetes 客户端"
// @Param namespace path string true "命名空间"
// @Param req body deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Success 200 {object} []v1.Deployment

// @Router /resource/list [get]
func GetResourceList(allDeploymentsInterface interface{}, client *k8s.Clientset,
	namespace string, req *deployment.GetDeploymentListReq) (*[]v1.Deployment, int, error) {

	var err error
	service := kubernetes.K8sResourceService{}

	allDeploymentsInterface, err = service.GetK8sResourceList(client, "deployments", namespace,
		req.Keyword, func(namespace string, opts metav1.ListOptions) (interface{}, error) {

			return client.AppsV1().Deployments(namespace).List(context.TODO(), opts)
		})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, 0, err
	}

	deploymentListResult, ok := allDeploymentsInterface.(*v1.DeploymentList)
	if !ok {
		deploymentListResult := allDeploymentsInterface.(*[]v1.Deployment)
		total := len(*deploymentListResult)
		return paginateDeployments(*deploymentListResult, req.Page, req.PageSize), total, nil
	}
	total := len(deploymentListResult.Items)
	return paginateDeployments(deploymentListResult.Items, req.Page, req.PageSize), total, nil

}

// getDeploymentsFromAllNamespaces 获取所有命名空间的部署
//
// @Description 获取所有命名空间的部署
// @Tags Deployment
// @Accept json
// @Produce json
// @Param client body k8s.Clientset true "Kubernetes 客户端"
// @Param req query deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {object} []v1.Deployment

// @Router /deployment/all/namespaces [get]
func (k *K8sDeploymentService) getDeploymentsFromAllNamespaces(client *k8s.Clientset, req deployment.GetDeploymentListReq, uuid uuid.UUID) ([]v1.Deployment, error) {
	var allDeployments []v1.Deployment
	c := cluster.K8sClusterService{}

	nsList, err := c.GetClusterUserNamespace(req.ClusterId, uuid)
	if err != nil {
		return nil, err
	}

	for _, namespace := range nsList {
		deployments, err := k.getDeploymentsFromNamespace(client, namespace, req.Keyword)
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		allDeployments = append(allDeployments, deployments...)
	}

	return allDeployments, nil
}

// getDeploymentsFromNamespace 获取指定命名空间的部署
//
// @Description 获取指定命名空间的部署
// @Tags Deployment
// @Accept json
// @Produce json
// @Param client body k8s.Clientset true "Kubernetes 客户端"
// @Param namespace path string true "命名空间"
// @Param keyword query string true "关键词"
// @Success 200 {object} []v1.Deployment

// @Router /deployment/namespace [get]
func (k *K8sDeploymentService) getDeploymentsFromNamespace(client *k8s.Clientset, namespace, keyword string) ([]v1.Deployment, error) {
	cacheKey := fmt.Sprintf("deployments:%s:%s", namespace, keyword)
	deploymentsJSON, err := global.DYCLOUD_REDIS.Get(context.TODO(), cacheKey).Result()

	if err == redis.Nil {
		return k.fetchAndCacheDeployments(client, namespace, keyword, cacheKey)
	} else if err != nil {
		return nil, err
	}

	var deployments []v1.Deployment
	err = json.Unmarshal([]byte(deploymentsJSON), &deployments)
	return deployments, err
}

// fetchAndCacheDeployments 从 Kubernetes API 获取并缓存部署
//
// @Description 从 Kubernetes API 获取并缓存部署
// @Tags Deployment
// @Accept json
// @Produce json
// @Param client body k8s.Clientset true "Kubernetes 客户端"
// @Param namespace path string true "命名空间"
// @Param keyword query string true "关键词"
// @Param cacheKey query string true "缓存键"
// @Success 200 {object} []v1.Deployment

// @Router /deployment/fetch/cache [get]
func (k *K8sDeploymentService) fetchAndCacheDeployments(client *k8s.Clientset, namespace, keyword, cacheKey string) ([]v1.Deployment, error) {
	deployList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var filteredDeployments []v1.Deployment
	if keyword != "" {
		for _, deploy := range deployList.Items {
			if k.DeploymentMatchesKeyword(deploy, keyword) {
				filteredDeployments = append(filteredDeployments, deploy)
			}
		}
	} else {
		filteredDeployments = deployList.Items
	}

	deploymentsJSON, err := json.Marshal(filteredDeployments)
	if err == nil {
		err = global.DYCLOUD_REDIS.Set(context.TODO(), cacheKey, deploymentsJSON, 0).Err()
	}

	return filteredDeployments, err
}

// paginateDeployments 分页处理部署列表
func paginateDeployments(allDeployments []v1.Deployment, page, pageSize int) *[]v1.Deployment {
	total := len(allDeployments)
	startIndex := (page - 1) * pageSize
	if startIndex >= total {
		return &[]v1.Deployment{}
	}

	endIndex := startIndex + pageSize
	if endIndex > total {
		endIndex = total
	}

	pagedDeployments := allDeployments[startIndex:endIndex]
	return &pagedDeployments
}
func (k *K8sDeploymentService) DeploymentMatchesKeyword(deployment v1.Deployment, keyword string) bool {
	if strings.Contains(deployment.Name, keyword) {
		return true
	}

	return false
}

func InformerFactory(k *K8sDeploymentService, req deployment.GetDeploymentListReq, uuid uuid.UUID) (*k8s.Clientset, error) {

	k8s, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}

	client, err := k8s.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	factory := informers.NewSharedInformerFactory(client, 0)
	k.DeploymentInformer = factory.Apps().V1().Deployments().Informer()

	// 设置事件处理程序
	k.DeploymentInformer.AddEventHandler(cache2.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			k.handleAddUpdate(obj, req, uuid)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			k.handleAddUpdate(newObj, req, uuid)
		},
		DeleteFunc: func(obj interface{}) {
			k.handleDelete(obj, req, uuid)
		},
	})
	// 启动 informer
	stopCh := make(chan struct{})
	go k.DeploymentInformer.Run(stopCh)
	if !cache2.WaitForCacheSync(stopCh, k.DeploymentInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return nil, fmt.Errorf("Timed out waiting for caches to sync")
	}

	return client, nil
}

// invalidateCache 使缓存失效
//
// @Description 使缓存失效
// @Tags Deployment
// @Accept json
// @Produce json
// @Param cacheKey path string true "缓存键"
// @Success 200 {string} string "缓存已失效"
// @Router /deployment/cache/invalidate [post]
func (k *K8sDeploymentService) invalidateCache(cacheKey string) {
	err := global.DYCLOUD_REDIS.Del(context.TODO(), cacheKey).Err()
	if err != nil {
		log.Println("Error invalidating Redis cache:", err)
	}
}

// DescribeDeploymentInfo 查看部署详情
//
// @Description 查看部署详情
// @Tags Deployment
// @Accept json
// @Produce json
// @Param req query deployment.DescribeDeploymentInfoReq true "部署详情请求参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {object} v1.Deployment

// @Router /deployment/describe [get]
func (k *K8sDeploymentService) DescribeDeploymentInfo(req deployment.DescribeDeploymentInfoReq, uuid uuid.UUID) (*v1.Deployment, error) {

	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		return &v1.Deployment{}, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1.Deployment{}, err
	}
	deploy, err := client.AppsV1().Deployments(req.Namespace).Get(context.TODO(), req.DeploymentName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return &v1.Deployment{}, err

	}
	return deploy, nil
}

// UpdateDeploymentInfo 更新部署信息
//
// @Description 更新部署信息
// @Tags Deployment
// @Accept json
// @Produce json
// @Param req body deployment.UpdateDeploymentInfoReq true "更新部署信息请求参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {object} v1.Deployment

// @Router /deployment/update [put]
func (k *K8sDeploymentService) UpdateDeploymentInfo(req deployment.UpdateDeploymentInfoReq, uuid uuid.UUID) (*v1.Deployment, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	c1, err := json.Marshal(req.Content)
	if err != nil {
		return nil, err
	}
	content := string(c1)
	var deployIns v1.Deployment
	err = json.Unmarshal([]byte(content), &deployIns)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	data, err := client.AppsV1().Deployments(req.Namespace).Update(context.TODO(), &deployIns, metav1.UpdateOptions{})
	if err != nil {
		return nil, err

	}
	return data, nil

}

// CreateDeployment 创建部署
//
// @Description 创建部署
// @Tags Deployment
// @Accept json
// @Produce json
// @Param req body deployment.CreateDeploymentReq true "创建部署请求参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {object} v1.Deployment

// @Router /deployment/create [post]
func (k *K8sDeploymentService) CreateDeployment(req deployment.CreateDeploymentReq, uuid uuid.UUID) (*v1.Deployment, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	c1, err := json.Marshal(req.Content)
	if err != nil {
		return nil, err
	}
	content := string(c1)
	var deployIns v1.Deployment
	err = json.Unmarshal([]byte(content), &deployIns)
	json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	data, err := client.AppsV1().Deployments(req.Namespace).Create(context.TODO(), &deployIns, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	//cacheKeyNamespace := fmt.Sprintf("deployments:%s:%s", req.Namespace, "")
	//k.updateCache(cacheKeyNamespace, req.Namespace, "", deployment.GetDeploymentListReq{
	//	ClusterId: req.ClusterId,
	//	Namespace: req.Namespace,
	//}, uuid)
	//
	//// 更新所有命名空间的缓存
	//cacheKeyAllNamespaces := "deployments::"
	//k.updateAllNamespacesCache(cacheKeyAllNamespaces, deployment.GetDeploymentListReq{
	//	ClusterId: req.ClusterId,
	//	Namespace: req.Namespace,
	//}, uuid)
	return data, nil
}

// DeleteDeployment 删除部署
//
// @Description 删除部署
// @Tags Deployment
// @Accept json
// @Produce json
// @Param req body deployment.DeleteDeploymentReq true "删除部署请求参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {string} string "删除成功"

// @Router /deployment/delete [post]
func (k *K8sDeploymentService) DeleteDeployment(req deployment.DeleteDeploymentReq, uuid uuid.UUID) error {
	k8s, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}
	client, err := k8s.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return err
	}

	err = client.AppsV1().Deployments(req.Namespace).Delete(context.TODO(), req.DeploymentName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	service := kubernetes.K8sResourceService{}
	// 删除缓存
	// 删除所有缓存
	err = service.DeleteAllK8sResourceCache()
	if err != nil {
		return err
	}
	//err = service.DeleteK8sResourceCache("deployments", req.Namespace, "")
	//if err != nil {
	//	return err
	//}

	return nil
}

// RollBackDeployment 回滚部署
//
// @Description 回滚部署
// @Tags Deployment
// @Accept json
// @Produce json
// @Param req body deployment.RollBackDeployment true "回滚部署请求参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {object} v1.Deployment

// @Router /deployment/rollback [post]
func (k *K8sDeploymentService) RollBackDeployment(req deployment.RollBackDeployment, uuid uuid.UUID) (*v1.Deployment, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	// 序列化用户传入的模板
	c1, err := json.Marshal(req.Content)
	if err != nil {
		return nil, err
	}
	content := string(c1)
	// 反序列化到Deployment
	var targetTemplate v1.Deployment
	err = json.Unmarshal([]byte(content), &targetTemplate)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	// 回滚deployment
	deployment, err := updateDeploymentWithTemplate(client, req.Namespace, req.DeploymentName, targetTemplate)
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	return deployment, nil

}

// updateDeploymentWithTemplate 更新部署模板
//
// @Description 更新部署模板
// @Tags Deployment
// @Accept json
// @Produce json
// @Param clientset body k8s.Clientset true "Kubernetes 客户端"
// @Param namespace path string true "命名空间"
// @Param deploymentName path string true "部署名称"
// @Param targetTemplate body v1.Deployment true "目标模板"
// @Success 200 {object} v1.Deployment

// @Router /deployment/update/template [put]
func updateDeploymentWithTemplate(clientset *k8s.Clientset, namespace, deploymentName string, targetTemplate v1.Deployment) (*v1.Deployment, error) {

	// 获取当前的 Deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}
	// 更新 Deployment 的模板
	deployment.Spec.Template = targetTemplate.Spec.Template
	// 应用更新后的 Deployment
	deployment, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error(err.Error())
		return nil, err
	}

	return deployment, nil
}

// updateAllNamespacesCache 更新所有命名空间的缓存
//
// @Description 更新所有命名空间的缓存
// @Tags Deployment
// @Accept json
// @Produce json
// @Param cacheKey path string true "缓存键"
// @Param req body deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {string} string "缓存已更新"

// @Router /deployment/cache/update/all [post]
func (k *K8sDeploymentService) updateAllNamespacesCache(cacheKey string, req deployment.GetDeploymentListReq, uuid uuid.UUID) {
	k8s, err := k.Generic(&req, uuid)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	client, err := k8s.Client()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	namespaceList, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error fetching namespaces:", err)
		return
	}

	var allDeployments []v1.Deployment
	for _, ns := range namespaceList.Items {
		deployList, err := client.AppsV1().Deployments(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Println("Error fetching deployments in namespace", ns.Name, ":", err)
			continue
		}
		allDeployments = append(allDeployments, deployList.Items...)
	}

	deploymentsJSON, err := json.Marshal(allDeployments)
	if err != nil {
		log.Println("Error marshalling deployments:", err)
		return
	}

	err = global.DYCLOUD_REDIS.Set(context.TODO(), cacheKey, deploymentsJSON, 0).Err()
	if err != nil {
		log.Println("Error setting Redis cache:", err)
	}
}

// handleAddUpdate 处理添加或更新事件
//
// @Description 处理添加或更新事件
// @Tags Deployment
// @Accept json
// @Produce json
// @Param obj body interface{} true "对象"
// @Param req body deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {string} string "处理成功"

// @Router /deployment/event/add-update [post]
func (k *K8sDeploymentService) handleAddUpdate(obj interface{}, req deployment.GetDeploymentListReq, uuid uuid.UUID) {
	deploy, ok := obj.(*v1.Deployment)
	if !ok {
		return
	}
	namespace := deploy.Namespace
	keyword := "" // 这里假设没有关键词过滤，如果有需要可以根据需求调整
	cacheKey := fmt.Sprintf("deployments:%s:%s", namespace, keyword)
	k.invalidateCache(cacheKey)
	k.updateCache(cacheKey, namespace, keyword, req, uuid)
	//k.updateAllNamespacesCache(cacheKey, req, uuid)
}

// handleDelete 处理删除事件
//
// @Description 处理删除事件
// @Tags Deployment
// @Accept json
// @Produce json
// @Param obj body interface{} true "对象"
// @Param req body deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {string} string "处理成功"

// @Router /deployment/event/delete [post]
func (k *K8sDeploymentService) handleDelete(obj interface{}, req deployment.GetDeploymentListReq, uuid uuid.UUID) {
	deploy, ok := obj.(*v1.Deployment)
	if !ok {
		return
	}
	namespace := deploy.Namespace
	keyword := "" // 这里假设没有关键词过滤，如果有需要可以根据需求调整
	cacheKey := fmt.Sprintf("deployments:%s:%s", namespace, keyword)
	k.invalidateCache(cacheKey)
	k.updateCache(cacheKey, namespace, keyword, req, uuid)
	//k.updateAllNamespacesCache(cacheKey, req, uuid)

}

// updateCache 更新缓存
//
// @Description 更新缓存
// @Tags Deployment
// @Accept json
// @Produce json
// @Param cacheKey path string true "缓存键"
// @Param namespace path string true "命名空间"
// @Param keyword query string true "关键词"
// @Param req body deployment.GetDeploymentListReq true "获取列表传入的参数"
// @Param uuid path uuid.UUID true "集群用户的 uuid"
// @Success 200 {string} string "缓存已更新"

// @Router /deployment/cache/update [post]
func (k *K8sDeploymentService) updateCache(cacheKey, namespace, keyword string, req deployment.GetDeploymentListReq, uuid uuid.UUID) {

	//if namespace == "" {
	//	// 处理所有命名空间的情况
	//	//k.updateAllNamespacesCache(cacheKey, req, uuid)
	//	return
	//}

	k8s, err := k.Generic(&req, uuid)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	client, err := k8s.Client()
	if err != nil {
		log.Println("Error:", err)
		return
	}
	deployList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error fetching deployments:", err)
		return
	}
	var filteredDeployments []v1.Deployment
	if keyword != "" {
		for _, deploy := range deployList.Items {
			if k8s.DeploymentMatchesKeyword(deploy, keyword) {
				filteredDeployments = append(filteredDeployments, deploy)
			}
		}
	} else {
		filteredDeployments = deployList.Items
	}
	deployList.Items = filteredDeployments
	deploymentsJSON, err := json.Marshal(deployList)
	if err != nil {
		log.Println("Error marshalling deployments:", err)
		return
	}

	err = global.DYCLOUD_REDIS.Set(context.TODO(), cacheKey, deploymentsJSON, 0).Err()
	if err != nil {
		log.Println("Error setting Redis cache:", err)
	}
}
