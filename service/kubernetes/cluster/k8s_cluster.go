package cluster

import (
	"DYCLOUD/global"
	"DYCLOUD/model/common/request"
	kubernetes2 "DYCLOUD/model/kubernetes"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	clusterReq "DYCLOUD/model/kubernetes/cluster/request"
	"DYCLOUD/model/system"
	"DYCLOUD/utils/kubernetes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type K8sClusterService struct{}

// CreateK8sCluster
//
// @Description: 创建集群，生成集群ca文件
// @receiver k8sClusterService
// @param k8sCluster body cluster2.K8sCluster true "创建k8s集群的请求参数"
// @return err
func (k8sClusterService *K8sClusterService) CreateK8sCluster(k8sCluster *cluster2.K8sCluster) (err error) {
	if !errors.Is(global.DYCLOUD_DB.Where("name = ?", k8sCluster.Name).First(&cluster2.K8sCluster{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("已存在相同集群，请勿重复创建")
	}
	client := kubernetes.NewKubernetes(k8sCluster, &cluster2.User{}, true)
	if err := client.CheckPermissions(); err != nil {
		global.DYCLOUD_LOG.Error("perrmission failed: " + err.Error())
	}
	k8sCluster.UUID = uuid.Must(uuid.NewV4())
	if err := global.DYCLOUD_DB.Create(&k8sCluster).Error; err != nil {
		return err
	}

	// 异步获取 ca.crt 存储入库
	go func(k8sCluster *cluster2.K8sCluster) {
		caCrt, err := client.GetKubeCaCrt()
		if err != nil {
			global.DYCLOUD_LOG.Error("get Ca.crt failed: " + err.Error())
		}

		if err = global.DYCLOUD_DB.Model(&cluster2.K8sCluster{}).Where("id = ?", k8sCluster.ID).Update("kube_ca_crt", caCrt).Error; err != nil {
			global.DYCLOUD_LOG.Error("update cluster filed: " + err.Error())
		}
	}(k8sCluster)
	// 异步创建集群默认角色
	go func() {
		if err := client.CreateDefaultClusterRoles(); err != nil {
			global.DYCLOUD_LOG.Error("create default cluster roles filed: " + err.Error())
		}
	}()

	return err
}

// DeleteK8sCluster
//
// @Description: 删除k8sCluster表记录
// @receiver k8sClusterService
// @param ID path int true "集群ID"
// @return err
func (k8sClusterService *K8sClusterService) DeleteK8sCluster(ID int) (err error) {
	var cl cluster2.K8sCluster
	if err = global.DYCLOUD_DB.Where("id = ?", ID).First(&cl).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	client := kubernetes.NewKubernetes(&cl, &cluster2.User{}, true)
	_ = client.CleanAllRBACResource()
	if err = global.DYCLOUD_DB.Delete(&cl).Error; err != nil {
		return err
	}
	return err
}

// DeleteK8sClusterByIds
//
// @Description: 批量删除k8sCluster表记录
// @receiver k8sClusterService
// @param ids body request.IdsReq true "批量删除的ID列表"
// @return err
func (k8sClusterService *K8sClusterService) DeleteK8sClusterByIds(ids request.IdsReq) (err error) {
	for _, id := range ids.Ids {
		go func(id int) {
			var cl cluster2.K8sCluster
			if err = global.DYCLOUD_DB.Where("id = ?", id).First(&cl).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				global.DYCLOUD_LOG.Error("search  cluster filed: " + err.Error())
			}

			// delete  all rbac resource
			client := kubernetes.NewKubernetes(&cl, &cluster2.User{}, true)
			_ = client.CleanAllRBACResource()

			if err = global.DYCLOUD_DB.Delete(&cl).Error; err != nil {
				global.DYCLOUD_LOG.Error("delete cluster filed: " + err.Error())
			}
		}(id)
	}
	return
}

// UpdateK8sCluster
//
// @Description: 更新k8sCluster表记录
// @receiver k8sClusterService
// @param k8sCluster body cluster2.K8sCluster true "更新k8s集群的请求参数"
// @return err
func (k8sClusterService *K8sClusterService) UpdateK8sCluster(k8sCluster cluster2.K8sCluster) (err error) {
	err = global.DYCLOUD_DB.Model(&cluster2.K8sCluster{}).Where("id = ?", k8sCluster.ID).Updates(&k8sCluster).Error
	return err
}

// GetK8sCluster
//
// @Description: 根据ID获取k8sCluster表记录
// @receiver k8sClusterService
// @param id path int true "集群ID"
// @return k8sCluster
// @return err
func (k8sClusterService *K8sClusterService) GetK8sCluster(id int) (k8sCluster cluster2.K8sCluster, err error) {
	fmt.Println(id)
	err = global.DYCLOUD_DB.Where("id = ?", id).Preload("Users").First(&k8sCluster).Error
	return
}

// CreateCredential
//
// @Description: 创建集群授权文件（用户）
// @receiver K8sClusterService
// @param clusterId path int true "集群ID"
// @param userId path uint true "用户ID"
// @return error
func (K8sClusterService *K8sClusterService) CreateCredential(clusterId int, userId uint) error {
	// 根据当前平台登录的用户id获取用户信息
	var user system.SysUser
	if err := global.DYCLOUD_DB.Where("id = ?", userId).First(&user).Error; err != nil {
		return err
	}
	// 根据当前平台登录的用户的uuid，和上传的集群id查询集群信息，同时预加载相关的用户信息
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", clusterId).Preload("Users", "uuid = ?", user.UUID).First(&clusterIns).Error; err != nil {
		return err
	}
	if !errors.Is(global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", clusterId, user.UUID).First(&cluster2.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("该账号集群凭据已存在")
	}
	// 判断当前平台登录的用户是否是admin，如果是admin则设置权限为true，如果是其他用户则设置为false
	isAdmin := true
	// 构建kubernetes 实例
	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, isAdmin)
	// 生成对应的客户端证书，返回私钥和公钥
	privateKey, publicCert, err := k.CreateClientCertificate(user.Username)
	if err != nil {
		return errors.New("CreateClientCertificate, error:" + err.Error())
	}

	// To JSON
	kubeConfig, err := k.KubeconfigJson(clusterIns, user.Username, privateKey, publicCert)
	if err != nil {
		return errors.New("To YAML, error:" + err.Error())
	}

	// save User
	if err = global.DYCLOUD_DB.Save(&cluster2.User{
		UUID:       user.UUID,
		Username:   user.Username,
		NickName:   user.NickName,
		KubeConfig: kubeConfig,
		ClusterId:  uint(clusterId),
	}).Error; err != nil {
		return err
	}

	return err
}

// GetClusterUserById
//
// @Description: 根据id查询集群用户
// @receiver k8sClusterService
// @param clusterId path int true "集群ID"
// @param userId path uint true "用户ID"
// @return user
// @return err
func (k8sClusterService *K8sClusterService) GetClusterUserById(clusterId int, userId uint) (user []cluster2.User, err error) {
	var userObj system.SysUser
	if err := global.DYCLOUD_DB.Where("id = ?", userId).First(&userObj).Error; err != nil {
		return user, err
	}

	var cluster cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", clusterId).Preload("Users", "uuid = ?", userObj.UUID).First(&cluster).Error; err != nil {
		return user, err
	}
	return cluster.Users, err
}

// GetClusterByUserUUID
//
// @Description: 根据用户UUID查询集群
// @receiver k8sClusterService
// @param id path int true "集群ID"
// @param uuid path uuid.UUID true "用户UUID"
// @return user
// @return err
func (k8sClusterService *K8sClusterService) GetClusterByUserUUID(id int, uuid uuid.UUID) (user cluster2.User, err error) {
	err = global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", id, uuid).First(&user).Error
	return
}

// GetClusterRoles
//
// @Description: 获取集群角色
// @receiver k8sClusterService
// @param roleType body clusterReq.ClusterRoleType true "角色类型请求参数"
// @return roles
// @return err
func (k8sClusterService *K8sClusterService) GetClusterRoles(roleType clusterReq.ClusterRoleType) (roles []rbacV1.ClusterRole, err error) {
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", roleType.ClusterId).First(&clusterIns).Error; err != nil {
		return roles, err
	}

	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)

	return k.ClusterRoles(roleType.RoleType)
}

// GetClusterApiGroups
//
// @Description: 获取集群api资源组
// @receiver k8sClusterService
// @param apiGroups body clusterReq.ClusterApiGroups true "API资源组请求参数"
// @return groups
// @return err
func (k8sClusterService *K8sClusterService) GetClusterApiGroups(apiGroups clusterReq.ClusterApiGroups) (groups []kubernetes2.ApiGroupOption, err error) {
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", apiGroups.ClusterId).First(&clusterIns).Error; err != nil {
		return nil, err
	}

	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	groups, err = k.ServerGroupsAndResources(apiGroups.ApiType)
	if err != nil {
		global.DYCLOUD_LOG.Error("ServerGroupsAndResources get resources failed:" + err.Error())
		return nil, err
	}

	return groups, err
}

// CreateClusterRole
//
// @Description: 创建k8s集群角色
// @receiver k8sClusterService
// @param role body cluster2.RoleData true "角色数据"
// @return err
func (k8sClusterService *K8sClusterService) CreateClusterRole(role cluster2.RoleData) (err error) {
	// 根据角色中的 ClusterId 查找对应的 Kubernetes 集群实例
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", role.ClusterId).First(&clusterIns).Error; err != nil {
		return err
	}
	// 遍历角色规则中的 API 组，将 "core" 替换为空字符串
	for i := range role.Rules {
		for j := range role.Rules[i].APIGroups {
			if role.Rules[i].APIGroups[j] == "core" {
				role.Rules[i].APIGroups[j] = ""
			}
		}
	}
	// 创建 Kubernetes 客户端实例
	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	client, err := k.Client()
	if err != nil {
		return err
	}
	// 设置角色的注解和标签
	role.Annotations["builtin"] = "false"
	role.Annotations["created-at"] = time.Now().Format("2006-01-02 15:04:05")
	role.Labels[kubernetes.LabelManageKey] = "devops"
	role.Labels[kubernetes.LabelClusterId] = clusterIns.UUID.String()
	// 使用 Kubernetes 客户端创建集群角色
	_, err = client.RbacV1().ClusterRoles().Create(context.TODO(), &role.ClusterRole, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

// UpdateClusterRole
//
// @Description: 更新集群角色
// @receiver k8sClusterService
// @param role body cluster2.RoleData true "角色数据"
// @return err
func (k8sClusterService *K8sClusterService) UpdateClusterRole(role cluster2.RoleData) (err error) {
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", role.ClusterId).First(&clusterIns).Error; err != nil {
		return err
	}

	for i := range role.Rules {
		for j := range role.Rules[i].APIGroups {
			if role.Rules[i].APIGroups[j] == "core" {
				role.Rules[i].APIGroups[j] = ""
			}
		}
	}

	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	client, err := k.Client()
	if err != nil {
		return err
	}

	instance, err := client.RbacV1().ClusterRoles().Get(context.TODO(), role.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	instance.Annotations["description"] = role.Annotations["description"]
	instance.Rules = role.Rules
	_, err = client.RbacV1().ClusterRoles().Update(context.TODO(), instance, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return err
}

// DeleteClusterRole
//
// @Description: 删除集群角色
// @receiver k8sClusterService
// @param role body cluster2.RoleData true "角色数据"
// @return err
func (k8sClusterService *K8sClusterService) DeleteClusterRole(role cluster2.RoleData) (err error) {
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", role.ClusterId).First(&clusterIns).Error; err != nil {
		return err
	}

	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	client, err := k.Client()
	if err != nil {
		return err
	}

	return client.RbacV1().ClusterRoles().Delete(context.TODO(), role.Name, metav1.DeleteOptions{})
}

// CreateUser
//
// @Description: 在k8s中为给定的角色创建用户
// @receiver k8sClusterService
// @param role body clusterReq.CreateClusterRole true "创建集群角色的请求参数"
// @return err
func (k8sClusterService *K8sClusterService) CreateUser(role clusterReq.CreateClusterRole) (err error) {
	// 遍历 role.UserUuids 列表，处理每个用户的创建操作
	for _, userUuid := range role.UserUuids {
		// 启动一个新的 goroutine 来并行处理每个用户的创建
		go func(u string, r clusterReq.CreateClusterRole) {
			var user system.SysUser
			// 根据用户的 UUID 从数据库中查找对应的用户记录
			if err = global.DYCLOUD_DB.Where("uuid = ?", u).First(&user).Error; err != nil {
				global.DYCLOUD_LOG.Warn("User Not Found failed:" + err.Error())
				return
			}
			// 检查该用户是否已经存在于指定的集群中
			if !errors.Is(global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", r.ClusterId, user.UUID).First(&cluster2.User{}).Error, gorm.ErrRecordNotFound) {
				global.DYCLOUD_LOG.Warn(fmt.Sprintf("The User: %s cluster credentials for this account already exist", user.NickName))
				return
			}
			// 调用 CreateClusterUser 方法来为用户创建集群凭证
			if err = k8sClusterService.CreateClusterUser(r, user); err != nil {
				global.DYCLOUD_LOG.Error("CreateClusterUser failed: " + err.Error())
				return
			}

			return
		}(userUuid, role)
	}

	return err
}

// CreateClusterUser
//
// @Description: 为给定角色和平台用户在k8s中创建集群用户
// @receiver k8sClusterService
// @param role body clusterReq.CreateClusterRole true "创建集群角色的请求参数"
// @param user body system.SysUser true "用户信息"
// @return err
func (k8sClusterService *K8sClusterService) CreateClusterUser(role clusterReq.CreateClusterRole, user system.SysUser) (err error) {
	// 根据角色中的 ClusterId 查找对应的 Kubernetes 集群实例
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", role.ClusterId).First(&clusterIns).Error; err != nil {
		return err
	}
	// 为用户创建集群凭据
	if err = k8sClusterService.CreateCredential(role.ClusterId, user.ID); err != nil {
		return errors.New("为用户创建集群凭据失败" + err.Error())
	}
	// 创建 Kubernetes 客户端实例
	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	// 查找用户在集群中的数据
	var userData cluster2.User
	if err = global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", role.ClusterId, user.UUID).First(&userData).Error; err != nil {
		return err
	}
	// 如果角色包含 ClusterRoles，则进行相应处理
	if len(role.ClusterRoles) > 0 {
		// 将 ClusterRoles 转换为 JSON 字符串
		clusterRolesString, err := json.Marshal(role.ClusterRoles)
		if err != nil {
			return err
		}
		// 启动一个新的 goroutine 处理 ClusterRoleBinding 的创建和更新
		go func(id uint, clusterRolesString string) {
			// 创建或更新 ClusterRoleBinding
			for r := range role.ClusterRoles {
				if err = k.CreateOrUpdateClusterRoleBinding(role.ClusterRoles[r], user.Username, false); err != nil {
					global.DYCLOUD_LOG.Error("CreateOrUpdateClusterRoleBinding failed:" + err.Error())
				}
			}

			// 更新集群授权
			if err = global.DYCLOUD_DB.Model(&cluster2.User{}).Where("id = ?", id).Update("cluster_roles", clusterRolesString).Error; err != nil {
				global.DYCLOUD_LOG.Error("CreateOrUpdateClusterRoleBinding Cluster Role failed:" + err.Error())
			}
		}(userData.ID, string(clusterRolesString))
	}
	// 如果角色包含 NamespaceRoles，则进行相应处理
	if len(role.NamespaceRoles) > 0 {
		// 将 NamespaceRoles 转换为 JSON 字符串
		namespaceRolesString, err := json.Marshal(role.NamespaceRoles)
		if err != nil {
			return err
		}
		// 启动一个新的 goroutine 处理 Rolebinding 的创建和更新
		go func(id uint, namespaceRolesString string) {
			// 创建或更新 Rolebinding
			for r := range role.NamespaceRoles {
				for j := range role.NamespaceRoles[r].Roles {
					if err := k.CreateOrUpdateRolebinding(role.NamespaceRoles[r].Namespace, role.NamespaceRoles[r].Roles[j], user.Username, false); err != nil {
						global.DYCLOUD_LOG.Error("CreateOrUpdateRolebinding Namespace role failed:" + err.Error())
					}

				}
			}

			// 更新命名空间授权
			if err = global.DYCLOUD_DB.Model(&cluster2.User{}).Where("id = ?", id).Update("cluster_roles", namespaceRolesString).Error; err != nil {
				global.DYCLOUD_LOG.Error("CreateOrUpdateClusterRoleBinding Cluster Role failed:" + err.Error())
			}
		}(userData.ID, string(namespaceRolesString))
	}

	return err
}

// UpdateClusterUser
//
// @Description: 更新集群用户
// @receiver k8sClusterService
// @param role body clusterReq.CreateClusterRole true "创建集群角色的请求参数"
// @return err
func (k8sClusterService *K8sClusterService) UpdateClusterUser(role clusterReq.CreateClusterRole) (err error) {
	var user system.SysUser
	if err := global.DYCLOUD_DB.Where("uuid = ?", role.UUID).First(&user).Error; err != nil {
		return err
	}

	var userData cluster2.User
	if errors.Is(global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", role.ClusterId, user.UUID).First(&userData).Error, gorm.ErrRecordNotFound) {
		return errors.New("该账号集群凭据不存在")
	}

	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", role.ClusterId).Preload("Users", "uuid = ?", role.UUID).First(&clusterIns).Error; err != nil {
		return err
	}

	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)

	// 清理clusterrolebinding
	if err := k.CleanManagedClusterRoleBinding(user.Username); err != nil {
		global.DYCLOUD_LOG.Error("CleanManagedClusterRoleBinding failed:" + err.Error())
	}

	// 清理Rolebinding
	if err := k.CleanManagedRoleBinding(user.Username); err != nil {
		global.DYCLOUD_LOG.Error("CleanManagedRoleBinding failed:" + err.Error())
	}

	if len(role.ClusterRoles) > 0 {
		clusterRolesString, err := json.Marshal(role.ClusterRoles)
		if err != nil {
			return err
		}

		go func(id uint, clusterRolesString string) {
			// 创建clusterrolebinding
			for r := range role.ClusterRoles {
				if err = k.CreateOrUpdateClusterRoleBinding(role.ClusterRoles[r], user.Username, false); err != nil {
					global.DYCLOUD_LOG.Error("CreateOrUpdateClusterRoleBinding failed:" + err.Error())
				}
			}

			// 更新集群授权
			if err = global.DYCLOUD_DB.Model(&cluster2.User{}).Where("id = ?", id).Update("cluster_roles", clusterRolesString).Error; err != nil {
				global.DYCLOUD_LOG.Error("CreateOrUpdateClusterRoleBinding Cluster Role failed:" + err.Error())
			}

		}(userData.ID, string(clusterRolesString))

	} else {
		// 集群授权清空
		if err = global.DYCLOUD_DB.Model(&cluster2.User{}).Where("id = ?", userData.ID).Update("cluster_roles", "").Error; err != nil {
			return err
		}

	}

	if len(role.NamespaceRoles) > 0 {
		namespaceRolesString, err := json.Marshal(role.NamespaceRoles)
		if err != nil {
			return err
		}

		go func(id uint, namespaceRolesString string) {
			// 创建Rolebinding
			for r := range role.NamespaceRoles {
				for j := range role.NamespaceRoles[r].Roles {
					if err := k.CreateOrUpdateRolebinding(role.NamespaceRoles[r].Namespace, role.NamespaceRoles[r].Roles[j], user.Username, false); err != nil {
						global.DYCLOUD_LOG.Error("CreateOrUpdateRolebinding failed:" + err.Error())
					}
				}
			}

			// 更新命名空间授权
			if err = global.DYCLOUD_DB.Model(&cluster2.User{}).Where("id = ?", id).Update("namespace_roles", namespaceRolesString).Error; err != nil {
				global.DYCLOUD_LOG.Error("CreateOrUpdateRolebinding  Namespace Role failed:" + err.Error())
			}

		}(userData.ID, string(namespaceRolesString))

	} else {
		// 命名空间授权清空
		if err = global.DYCLOUD_DB.Model(&cluster2.User{}).Where("id = ?", userData.ID).Update("namespace_roles", "").Error; err != nil {
			return err
		}
	}

	return err
}

// DeleteClusterUser
//
// @Description: 删除集群用户
// @receiver k8sClusterService
// @param role body clusterReq.DeleteClusterRole true "删除集群角色的请求参数"
// @return err
func (k8sClusterService *K8sClusterService) DeleteClusterUser(role clusterReq.DeleteClusterRole) (err error) {
	for _, userUuid := range role.UserUuids {
		var user cluster2.User
		if err := global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", role.ClusterId, userUuid).First(&user).Error; err != nil {
			global.DYCLOUD_LOG.Warn("Cluster Not Found User")
			continue
		}
		if err := global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
			var clusterIns cluster2.K8sCluster
			if err := global.DYCLOUD_DB.Where("id = ?", role.ClusterId).First(&clusterIns).Error; err != nil {
				global.DYCLOUD_LOG.Warn("Cluster Search failed: " + err.Error())
				return err
			}
			// 清理集群角色绑定
			k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
			if err := k.CleanManagedClusterRoleBinding(user.Username); err != nil {
				return err
			}

			// 清理命名空间角色绑定
			if err := k.CleanManagedRoleBinding(user.Username); err != nil {
				return err
			}

			TxErr := tx.Where("cluster_id = ? and uuid = ?", role.ClusterId, userUuid).Delete(&cluster2.User{}).Error
			if TxErr != nil {
				return TxErr
			}

			// 返回 nil 提交事务
			return nil
		}); err != nil {
			return err
		}

	}

	return err
}

// GetClusterUserNamespace
//
// @Description: 获取集群用户的namespace列表
// @receiver k8sClusterService
// @param clusterId path int true "集群ID"
// @param uuid path uuid.UUID true "用户UUID"
// @return nsList
// @return err
func (k8sClusterService *K8sClusterService) GetClusterUserNamespace(clusterId int, uuid uuid.UUID) (nsList []string, err error) {
	var user cluster2.User
	if err := global.DYCLOUD_DB.Where("cluster_id = ? and uuid = ?", clusterId, uuid).First(&user).Error; err != nil {
		return nil, err
	}

	var cluster cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", clusterId).First(&cluster).Error; err != nil {
		return nil, err
	}

	k := kubernetes.NewKubernetes(&cluster, &user, true)

	return k.GetUserNamespaceNames(user.Username)
}

// GetClusterListNamespace
//
// @Description: 获取集群namespace列表
// @receiver k8sClusterService
// @param clusterId path int true "集群ID"
// @return nsList
// @return err
func (k8sClusterService *K8sClusterService) GetClusterListNamespace(clusterId int) (nsList []corev1.Namespace, err error) {
	var clusterIns cluster2.K8sCluster
	if err := global.DYCLOUD_DB.Where("id = ?", clusterId).First(&clusterIns).Error; err != nil {
		return nil, err
	}

	k := kubernetes.NewKubernetes(&clusterIns, &cluster2.User{}, true)
	client, err := k.Client()
	if err != nil {
		return nil, err
	}

	ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return ns.Items, err

}

// GetK8sClusterInfoList
//
// @Description: 获取集群列表
// @receiver k8sClusterService
// @param info query clusterReq.K8sClusterSearch true "集群搜索信息"
// @return list
// @return total
// @return err
func (k8sClusterService *K8sClusterService) GetK8sClusterInfoList(info clusterReq.K8sClusterSearch) (list []cluster2.K8sCluster, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DYCLOUD_DB.Model(&cluster2.K8sCluster{})
	var k8sClusters []cluster2.K8sCluster
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&k8sClusters).Error
	return k8sClusters, total, err
}
