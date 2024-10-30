package kubernetes

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes"
	cluster2 "DYCLOUD/model/kubernetes/cluster"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	authV1 "k8s.io/api/authorization/v1"
	v1 "k8s.io/api/authorization/v1"
	certv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"sync"
)

// KubernetesService
// @Description: 定义K8s服务相关接口
type KubernetesService interface {
	//
	// CheckPermissions
	//  @Description: 检查当前用户权限
	//  @return err
	//
	CheckPermissions() (err error)
	//
	// HasPermission
	//  @Description: 检查当前用户是否具有指定的资源权限
	//  @return kubernetes.PermissionCheckResult  权限检查结果
	//  @return error
	//
	HasPermission(attributes v1.ResourceAttributes) (kubernetes.PermissionCheckResult, error)
	//
	// Config
	//  @Description: 返回k8s客户端配置
	//  @return config
	//  @return err
	//
	Config() (config *rest.Config, err error)
	//
	// CheckRequiredPermissionsclient
	//  @Description: 所需权限的映射，映射键为资源类型，值为权限列表
	//  @return string
	//  @return error
	//
	CheckRequiredPermissionsclient(requiredPermissions map[string][]string) (string, error)
	//
	// Client
	//  @Description: 返回k8s客户端
	//  @return *k8s.Clientset
	//  @return error
	//
	Client() (*k8s.Clientset, error)
	//
	// GetKubeCaCrt
	//  @Description: 获取k8s集群的CA证书
	//  @return string  base64 编码的CA证书
	//  @return error
	//
	GetKubeCaCrt() (string, error)
	//
	// CreateDefaultClusterRoles
	//  @Description: 创建默认的集群角色
	//  @return err
	//
	CreateDefaultClusterRoles() (err error)
	//
	// CleanAllRBACResource
	//  @Description: 清理所有RBAC资源
	//  @return err
	//
	CleanAllRBACResource() (err error)
	//
	// KubeconfigJson
	//  @Description: 根据提供的集群信息、用户名、私钥和公钥证书生成kubeConfig的JSON格式字符串
	//  @return yaml   生成的kubeConfig JSON格式的字符串
	//  @return err
	//
	KubeconfigJson(cluster cluster2.K8sCluster, username, privateKey, publicCert string) (yaml string, err error)
	//
	// CleanManagedRoleBinding
	//  @Description: 清理管理的角色绑定
	//  @return error
	//
	CleanManagedRoleBinding(username string) error
	//
	// CleanManagedClusterRole
	//  @Description: 清理管理的集群角色
	//  @return error
	//
	CleanManagedClusterRole() error
	//
	// CleanManagedClusterRoleBinding
	//  @Description: 清理管理的集群角色绑定
	//  @return error
	//
	CleanManagedClusterRoleBinding(username string) error
	//
	// CreateClientCertificate
	//  @Description: 创建客户端证书
	//  @return privateKey 返回私钥
	//  @return publicCert 返回公钥
	//  @return err
	//
	CreateClientCertificate(username string) (privateKey, publicCert string, err error)
	//
	// ClusterRoles
	//  @Description: 根据指定的角色类型获取集群角色列表
	//  @return roles  返回包含集群角色的列表
	//  @return err
	//
	ClusterRoles(roleType string) (roles []rbacV1.ClusterRole, err error)
	//
	// CanVisitAllNamespace
	//  @Description: 检查用户是否可以访问所有命名空间
	//  @return bool
	//  @return error
	//
	CanVisitAllNamespace(username string) (bool, error)
	//
	// CreateOrUpdateRolebinding
	//  @Description: 创建或更新角色绑定
	//  @return error
	//
	CreateOrUpdateRolebinding(namespace string, clusterRoleName string, username string, builtIn bool) error
	//
	// CreateOrUpdateClusterRoleBinding
	//  @Description: 创建或更新集群角色绑定
	//  @return error
	//
	CreateOrUpdateClusterRoleBinding(clusterRoleName string, username string, builtIn bool) error
	//
	// GetUserNamespaceNames
	//  @Description:  获取用户命名空间名称列表
	//  @return []string   命名空间名称列表
	//  @return error
	//
	GetUserNamespaceNames(username string) ([]string, error)
	//
	// NodeMatchesKeyword
	//  @Description: 检查节点是否匹配关键字
	//  @return bool
	//
	NodeMatchesKeyword(node corev1.Node, keyword string) bool
}

func GenerateTLSTransport(c *cluster2.K8sCluster, u *cluster2.User, isadmin bool) (http.RoundTripper, error) {
	kube := NewKubernetes(c, u, isadmin)
	kubeconfig, err := kube.Config()
	if err != nil {
		return nil, err
	}

	return rest.TransportFor(kubeconfig)
}

// Kubernetes
// @Description: 结构体包含 k8sCluster集群对象、集群用户信息和是否为管理员的标志
type Kubernetes struct {
	*cluster2.K8sCluster
	*cluster2.User
	IsAdmin bool
}

// NewKubernetes
//
//	@Description: 创建一个新的kubernetes实例
//	@return *Kubernetes
func NewKubernetes(cluster *cluster2.K8sCluster, user *cluster2.User, isAdmin bool) *Kubernetes {
	return &Kubernetes{cluster, user, isAdmin}

}

// CheckPermissions
//
//	@Description: 检查当前用户在k8s集群的权限
//	@receiver k
//	@return err
func (k *Kubernetes) CheckPermissions() (err error) {
	// 定义需要检查的权限列表
	permissions := map[string][]string{
		"namespaces":       {"get", "post", "delete"},
		"clusterroles":     {"get", "post", "delete"},
		"clusterrolebings": {"get", "post", "delete"},
		"roles":            {"get", "post", "delete"},
		"rolebindings":     {"get", "post", "delete"},
		"nodes":            {"get", "post", "delete"},
	}
	// 检查所需的权限
	notAllowed, err := k.CheckRequiredPermissionsclient(permissions)
	if err != nil {
		global.DYCLOUD_LOG.Error(fmt.Sprintf("notAllowed faile: %s", err.Error()))
		return errors.New(fmt.Sprintf("notAllowed faile: %s", err.Error()))
	}
	// 如果有权限未被允许，返回错误信息
	if notAllowed != "" {
		return errors.New(fmt.Sprintf("permission %s required", notAllowed))
	}

	return err
}

// CheckRequiredPermissionsclient
//
//	@Description: 检查所需的权限
//	@receiver k
//	@return string
//	@return error
func (k *Kubernetes) CheckRequiredPermissionsclient(requiredPermissions map[string][]string) (string, error) {
	wg := sync.WaitGroup{}                                  // 等待所有协程完成
	errCh := make(chan error)                               // 用于接受错误信息
	resultCh := make(chan kubernetes.PermissionCheckResult) // 用于接受权限检查结果
	doneCh := make(chan struct{})                           // 用于通知所有协程已完成
	// 遍历所需的权限并启动协程进行检查
	for key := range requiredPermissions {
		for i := range requiredPermissions[key] {
			wg.Add(1) // 增加等待组计数
			i := i    // 避免闭包捕获循环变量
			go func(key string, index int) {
				// 检查单个权限
				rs, err := k.HasPermission(authV1.ResourceAttributes{
					Verb:     requiredPermissions[key][i],
					Resource: key,
				})
				if err != nil {
					errCh <- err // 发送错误信息到通道
					wg.Done()    // 减少等待组计数
					return
				}
				resultCh <- rs // 发送权限检查结果到通道
				wg.Done()      // 减少等待组计数
			}(key, i)
		}
	}
	// 等待所有协程完成
	go func() {
		wg.Wait()
		doneCh <- struct{}{} // 发送成功信号
	}()
	// 处理检查结果
	for {
		select {
		case <-doneCh:
			goto end // 所有协程完成后跳转到end标签
		case err := <-errCh:
			return "", err // 返回错误
		case b := <-resultCh:
			if !b.Allowed {
				return b.Resource.Resource, nil
			} // 返回未被允许的资源
		}
	}
end:
	return "", nil
}

// HasPermission
//
//	@Description: 检查单个权限
//	@receiver k
//	@return kubernetes.PermissionCheckResult
//	@return error
func (k *Kubernetes) HasPermission(attributes v1.ResourceAttributes) (kubernetes.PermissionCheckResult, error) {
	// 获取k8s客户端
	clientset, err := k.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("clientset init failed: " + err.Error())
		return kubernetes.PermissionCheckResult{}, err
	}
	// 创建
	resp, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), &v1.SelfSubjectAccessReview{
		Spec: v1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &attributes,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return kubernetes.PermissionCheckResult{}, err
	}

	return kubernetes.PermissionCheckResult{
		Resource: attributes,
		Allowed:  resp.Status.Allowed,
	}, nil

}

// Client
//
//	@Description: 初始化k8s客户端
//	@receiver k
//	@return *k8s.Clientset
//	@return error
func (k *Kubernetes) Client() (*k8s.Clientset, error) {
	cfg, err := k.Config()
	if err != nil {
		global.DYCLOUD_LOG.Error("Config init failed: ", zap.Error(err))
		return nil, err
	}

	clientset, err := k8s.NewForConfig(cfg)
	if err != nil {
		global.DYCLOUD_LOG.Error("clientset init failed: ", zap.Error(err))
		return nil, err
	}

	return clientset, err
}

// Config
//
//	@Description: 初始化k8s配置
//	@receiver k
//	@return config
//	@return err
func (k *Kubernetes) Config() (config *rest.Config, err error) {
	// 判断用户创建k8s集群的方式，1为kubeconfig，2为token
	if k.KubeType == 1 {
		if k.IsAdmin {
			config, err = clientcmd.RESTConfigFromKubeConfig([]byte(k.K8sCluster.KubeConfig))
			if err != nil {
				global.DYCLOUD_LOG.Error("RESTConfigFromKubeConfig init failed: ", zap.Error(err))
				return nil, err
			}
		} else {
			config, err = clientcmd.RESTConfigFromKubeConfig([]byte(k.User.KubeConfig))
			if err != nil {
				global.DYCLOUD_LOG.Error("RESTConfigFromKubeConfig init failed: ", zap.Error(err))
				return nil, err
			}
		}
		return config, err
	}
	if k.KubeType == 2 {
		if k.IsAdmin {
			return &rest.Config{
				Host:            k.ApiAddress,
				BearerToken:     k.K8sCluster.KubeConfig,
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			}, err
		} else {
			return &rest.Config{
				Host:            k.ApiAddress,
				BearerToken:     k.User.KubeConfig,
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			}, err
		}
	}

	return
}

// GetKubeCaCrt
//
//	@Description: 获取k8s集群的ca证书
//	@receiver k
//	@return cacrt 返回base64编码的CA证书
//	@return err
func (k *Kubernetes) GetKubeCaCrt() (cacrt string, err error) {
	clientset, err := k.Client()
	if err != nil {
		return "", err
	}

	User := fmt.Sprintf("%s-%s-%d", "random", "devops", time.Now().Unix())

	serviceaccount := &corev1.ServiceAccount{TypeMeta: metav1.TypeMeta{}, ObjectMeta: metav1.ObjectMeta{Name: User}}
	if _, err = clientset.CoreV1().ServiceAccounts("default").Create(context.TODO(), serviceaccount, metav1.CreateOptions{}); err != nil {
		return "", errors.New("ServiceAccounts 创建失败: " + err.Error())
	}

	secret := &corev1.Secret{TypeMeta: metav1.TypeMeta{}, ObjectMeta: metav1.ObjectMeta{Name: User,
		Annotations: map[string]string{"kubernetes.io/service-account.name": User}},
		Type: "kubernetes.io/service-account-token",
	}
	if _, err = clientset.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{}); err != nil {
		return "", errors.New("Secrets 创建失败: " + err.Error())
	}

	for i := 1; i <= 150; i++ {
		time.Sleep(3 * time.Second)

		secrets, err := clientset.CoreV1().Secrets("default").Get(context.TODO(), User, metav1.GetOptions{})
		if err != nil {
			continue
		}

		secretCaCrt := fmt.Sprintf("%s", secrets.Data["ca.crt"])
		if secretCaCrt != "" {
			if err = clientset.CoreV1().Secrets("default").Delete(context.TODO(), User, metav1.DeleteOptions{}); err != nil {
				return "", err
			}

			if err = clientset.CoreV1().ServiceAccounts("default").Delete(context.TODO(), User, metav1.DeleteOptions{}); err != nil {
				return "", err
			}
		}

		return base64.StdEncoding.EncodeToString([]byte(secretCaCrt)), err
	}

	return "", err
}

// CreateDefaultClusterRoles
//
//	@Description: 创建默认的集群角色
//	@receiver k
//	@return error
func (k *Kubernetes) CreateDefaultClusterRoles() error {
	clientset, err := k.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("clientset init failed: " + err.Error())
		return err
	}
	// 遍历并创建或更新默认的集群角色
	for i := range InitClusterRoles {
		instance, err := clientset.RbacV1().ClusterRoles().Get(context.TODO(), InitClusterRoles[i].Name, metav1.GetOptions{})
		if err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "not found") {
				return err
			}
		}

		// 如果集群角色不存在，则创建；否则更新
		if instance == nil || instance.Name == "" {
			_, err = clientset.RbacV1().ClusterRoles().Create(context.TODO(), &InitClusterRoles[i], metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			_, err = clientset.RbacV1().ClusterRoles().Update(context.TODO(), &InitClusterRoles[i], metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CleanAllRBACResource
//
//	@Description: 清理所有RBAC资源
//	@receiver k
//	@return error
func (k *Kubernetes) CleanAllRBACResource() error {
	// 清理集群资源
	if err := k.CleanManagedClusterRole(); err != nil {
		return err
	}
	// 清理集群角色绑定
	if err := k.CleanManagedClusterRoleBinding(""); err != nil {
		return err
	}
	// 清理角色绑定
	if err := k.CleanManagedRoleBinding(""); err != nil {
		return err
	}

	return nil
}

// CleanManagedClusterRole
//
//	@Description: 清理管理的集群角色
//	@receiver k
//	@return error
func (k *Kubernetes) CleanManagedClusterRole() error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	// 构建标签选择器
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "devops"),
	}
	// 删除符合标签选择器的集群角色
	return client.RbacV1().ClusterRoles().DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	})
}

// CleanManagedRoleBinding
//
//	@Description: 清理管理的角色绑定
//	@receiver k
//	@return error
func (k *Kubernetes) CleanManagedRoleBinding(username string) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	// 列出所有命名空间
	nss, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	// 构建标签选择器
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "devops"),
	}
	if username != "" {
		labels = append(labels, fmt.Sprintf("%s=%s", LabelUsername, username))
	}
	// 遍历所有命名空间并删除符合条件的角色绑定
	for i := range nss.Items {
		if err := client.RbacV1().RoleBindings(nss.Items[i].Name).DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
			LabelSelector: strings.Join(labels, ","),
		}); err != nil {
			return err
		}
	}
	return nil
}
func (k *Kubernetes) CleanManagedClusterRoleBinding(username string) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "devops"),
		fmt.Sprintf("%s=%s", LabelClusterId, k.K8sCluster.UUID),
	}
	if username != "" {
		labels = append(labels, fmt.Sprintf("%s=%s", LabelUsername, username))
	}
	return client.RbacV1().ClusterRoleBindings().DeleteCollection(context.TODO(), metav1.DeleteOptions{}, metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	})
}

// CreateClientCertificate
//
//	@Description: 创建集群证书
//	@receiver k
//	@return privateKey
//	@return publicCert
//	@return err
func (k *Kubernetes) CreateClientCertificate(username string) (privateKey, publicCert string, err error) {
	// 获取clientSet客户端
	clientset, err := k.Client()
	if err != nil {
		return privateKey, publicCert, err
	}

	// 生成一个私钥
	RandprivateKey, err := GeneratePrivateKey()
	if err != nil {
		return privateKey, publicCert, err
	}

	// 生成用户证书申请
	cert, err := CreateClientCertificateRequest(username, RandprivateKey)
	if err != nil {
		return privateKey, publicCert, err
	}
	// 创建k8s证书签名请求(CSR)
	csr := certv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: username,
		},
		Spec: certv1.CertificateSigningRequestSpec{
			SignerName: "kubernetes.io/kube-apiserver-client", // 指定签名者
			Request:    cert,                                  // 设置CSR 请求的证书
			Groups: []string{
				"system:authenticated", // 指定用户组
			},
			Usages: []certv1.KeyUsage{
				"client auth", // 指定证书用途为客户端认证
			},
		},
	}
	// 在k8s集群中创建CSR
	createResp, err := clientset.CertificatesV1().CertificateSigningRequests().Create(context.TODO(), &csr, metav1.CreateOptions{})
	if err != nil {
		return privateKey, publicCert, err
	}

	// 审批CSR证书
	createResp.Status.Conditions = append(createResp.Status.Conditions, certv1.CertificateSigningRequestCondition{
		Reason:         "Approved by Devops",
		Type:           certv1.CertificateApproved,
		LastUpdateTime: metav1.Now(),
		Status:         "True",
	})
	// 更新CSR的审批状态
	updateResp, err := clientset.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), createResp.Name, createResp, metav1.UpdateOptions{})
	if err != nil {
		return privateKey, publicCert, err
	}
	// 循环等待证书颁发
	for {
		max_num := 0
		for {
			// 每次循环等待3秒
			time.Sleep(3 * time.Second)
			if max_num > 150 {
				break
			}
			// 获取CSR的状态
			getResp, err := clientset.CertificatesV1().CertificateSigningRequests().Get(context.TODO(), updateResp.Name, metav1.GetOptions{})
			if err != nil {
				max_num += 1
				global.DYCLOUD_LOG.Warn("CertificateSigningRequests Get failed: " + err.Error())
				continue
			}
			// 检查CSR是否已经颁发证书
			if getResp.Status.Certificate != nil {
				// 删除已经颁发证书的CSR
				if err = clientset.CertificatesV1().CertificateSigningRequests().Delete(context.TODO(), username, metav1.DeleteOptions{}); err != nil {
					max_num += 1
					global.DYCLOUD_LOG.Warn("CertificateSigningRequests Get failed: " + err.Error())
					continue
				}
				// 返回私钥和公钥证书(编码为Base64)
				return base64.StdEncoding.EncodeToString(
					pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: RandprivateKey}),
				), base64.StdEncoding.EncodeToString(getResp.Status.Certificate), err
			}

			max_num += 1
			continue
		}
	}
}

// CreateClientCertificateRequest
//
//	@Description: 生成客户端证书CSR
//	@return []byte
//	@return error
func CreateClientCertificateRequest(userName string, key []byte, org ...string) ([]byte, error) {
	// 解析PEM 编码的RSA私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, err
	}
	// 构建证书请求的主题信息，包括用户名等
	subj := pkix.Name{
		CommonName: userName,
	}
	for i := range org {
		subj.Organization = append(subj.Organization, org[i])
	}
	// 将主题信息转换为RDN序列
	rawSubj := subj.ToRDNSequence()
	// 使用ASN.1 编码主题
	asn1Subj, err := asn1.Marshal(rawSubj)
	if err != nil {
		return nil, err
	}
	// 创建x509的证书请求对象
	csr := &x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	// 生成证书CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, privateKey)
	if err != nil {
		return nil, err
	}
	// 将CSR转为PEM格式
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}), nil
}

// GeneratePrivateKey
//
//	@Description: 生成私钥
//	@return []byte
//	@return error
func GeneratePrivateKey() ([]byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return x509.MarshalPKCS1PrivateKey(privateKey), nil
}

// KubeconfigJson
//
//	@Description: 根据提供的集群信息、用户名、私钥和公钥证书生成kubeConfig的json格式字符串C
//	@receiver k
//	@return yaml   生成的kubeConfig JSON格式字符串
//	@return err
func (k *Kubernetes) KubeconfigJson(cluster cluster2.K8sCluster, username, privateKey, publicCert string) (yaml string, err error) {
	// 创建Kubeconfig结构
	kubeconfig := kubernetes.Kubeconfig{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []kubernetes.ClusterEntry{
			{
				Name:    cluster.Name,
				Cluster: kubernetes.KubeCluster{CertificateAuthorityData: cluster.KubeCaCrt, Server: cluster.ApiAddress},
			},
		},
		Contexts: []kubernetes.ContextEntry{
			{
				Name: cluster.Name,
				Context: kubernetes.Context{
					Cluster: cluster.Name,
					User:    username,
				},
			},
		},
		CurrentContext: cluster.Name,
		Preferences:    struct{}{},
		Users: []kubernetes.UserEntry{
			{
				Name: username,
				User: kubernetes.KubeUser{
					ClientCertificateData: publicCert,
					ClientKeyData:         privateKey,
				},
			},
		},
	}

	// 将Kubeconfig结构转换为Json格式的字节
	jsonData, err := json.Marshal(&kubeconfig)
	if err != nil {
		return "", err
	}

	return string(jsonData), err
}

// ClusterRoles
//
//	@Description: 根据指定的角色类型获取集群角色列表
//	@receiver k
//	@return roles  返回包含集群角色的列表
//	@return err
func (k *Kubernetes) ClusterRoles(roleType string) (roles []rbacV1.ClusterRole, err error) {
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	// 构建标签选择器
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "devops"),
		fmt.Sprintf("%s=%s", "devops/role-type", roleType),
	}
	// 列出符合标签选择器的集群角色
	roleList, err := client.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	})

	if err != nil {
		return nil, err
	}

	return roleList.Items, err

}

// ServerGroupsAndResources
//
//	@Description:  根据指定的API 类型(namespace或者cluster)获取服务器上的API组和资源
//	@receiver k
//	@return groups 返回包含API组和资源的列表
//	@return err
func (k *Kubernetes) ServerGroupsAndResources(api_type string) (groups []kubernetes.ApiGroupOption, err error) {
	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	// 获取服务器上的 API 组和资源
	_, rss, err := client.ServerGroupsAndResources()
	if err != nil {
		return nil, err
	}
	// 创建一个映射，用于跟踪已经处理过的组
	var (
		groupMap = make(map[string]struct{})
	)
	// 定义一个通用的资源选项，表示对所有资源的所有操作
	itemResource := kubernetes.ApiResourceOption{
		Resource: "*",
		Verbs:    []string{"create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"},
	}
	// 判断 API 类型，如果是 "namespace"，则添加一个通用的 API 组选项
	//if api_type == "namespace" {
	//	groups = append(groups, kubernetes.ApiGroupOption{Group: "*", Resources: []kubernetes.ApiResourceOption{itemResource}})
	//} else {
	groups = append(groups, kubernetes.ApiGroupOption{Group: "*", Resources: []kubernetes.ApiResourceOption{itemResource}})
	//}
	// 遍历所有 API 组和资源

	for _, group := range rss {
		// 提取组版本中的组名
		if strings.Contains(group.GroupVersion, "/") {
			group.GroupVersion = group.GroupVersion[0:strings.Index(group.GroupVersion, "/")]
		}

		name := group.GroupVersion
		if name == "v1" {
			name = ""
		}
		// 如果组名已经处理过，则跳过
		if _, ok := groupMap[name]; ok {
			continue
		}
		// 标记组名为已处理
		groupMap[name] = struct{}{}
		// 创建一个 API 组选项
		itemGroup := kubernetes.ApiGroupOption{Group: name, Resources: []kubernetes.ApiResourceOption{itemResource}}
		// 遍历组中的所有资源
		for _, resource := range group.APIResources {
			// 根据 API 类型过滤资源
			if !(api_type == "namespace" && resource.Namespaced) && !(api_type == "cluster" && !resource.Namespaced) {
				continue
			}
			// 将符合条件的资源添加到 API 组选项中
			itemGroup.Resources = append(itemGroup.Resources, kubernetes.ApiResourceOption{Resource: resource.Name, Verbs: resource.Verbs})
		}
		// 如果组中有资源，添加到结果列表中
		if len(itemGroup.Resources) > 1 {
			groups = append(groups, itemGroup)
		}

	}

	return groups, err

}

// CreateOrUpdateClusterRoleBinding
//
//	@Description: 创建或更新k8s ClusterRoleBinding
//	@receiver k
//	@return error
func (k *Kubernetes) CreateOrUpdateClusterRoleBinding(clusterRoleName string, username string, builtIn bool) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	// 生成 ClusterRoleBinding 的名称
	name := fmt.Sprintf("%s:%s:%s", username, clusterRoleName, k.K8sCluster.UUID.String())
	labels := map[string]string{
		LabelManageKey: "devops",
		LabelClusterId: k.K8sCluster.UUID.String(),
		LabelUsername:  username,
	}
	annotations := map[string]string{
		"built-in":   strconv.FormatBool(builtIn),
		"created-at": time.Now().Format("2006-01-02 15:04:05"),
	}
	// 创建 ClusterRoleBinding 对象
	item := rbacV1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacV1.RoleRef{
			Kind: "ClusterRole",
			Name: clusterRoleName,
		},
	}
	// 获取现有的 ClusterRoleBinding
	baseItem, err := client.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	if baseItem != nil && baseItem.Name != "" {
		_, err := client.RbacV1().ClusterRoleBindings().Create(context.TODO(), &item, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := client.RbacV1().ClusterRoleBindings().Update(context.TODO(), &item, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateOrUpdateRolebinding
//
//	@Description: 创建或更新k8s RoleBinding
//	@receiver k
//	@return error
func (k *Kubernetes) CreateOrUpdateRolebinding(namespace string, clusterRoleName string, username string, builtIn bool) error {
	client, err := k.Client()
	if err != nil {
		return err
	}
	// 设置标签
	labels := map[string]string{
		LabelManageKey: "devops",
		LabelClusterId: k.K8sCluster.UUID.String(),
		LabelUsername:  username,
	}
	// 设置注解
	annotations := map[string]string{
		"built-in":   strconv.FormatBool(builtIn),
		"created-at": time.Now().Format("2006-01-02 15:04:05"),
	}
	// 生成RoleBinding的名称
	name := fmt.Sprintf("%s:%s:%s:%s", namespace, username, clusterRoleName, k.K8sCluster.UUID)
	// 创建RoleBinding对象
	item := rbacV1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
			Namespace:   namespace,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacV1.RoleRef{
			Kind: "ClusterRole",
			Name: clusterRoleName,
		},
	}
	// 获取现有的 RoleBinding
	baseItem, err := client.RbacV1().RoleBindings(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}
	// 如果 RoleBinding 不存在，则创建新的 RoleBinding
	if baseItem != nil && baseItem.Name != "" {
		_, err := client.RbacV1().RoleBindings(namespace).Create(context.TODO(), &item, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		// 如果 RoleBinding 存在，则更新现有的 RoleBinding
		_, err := client.RbacV1().RoleBindings(namespace).Update(context.TODO(), &item, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserNamespaceNames
//
//	@Description: 获取指定用户可以访问的命名空间列表
//	@receiver k
//	@return []string 返回可以访问的命名空间名称列表
//	@return error
func (k *Kubernetes) GetUserNamespaceNames(username string) ([]string, error) {

	client, err := k.Client()
	if err != nil {
		return nil, err
	}
	// 检查用户是否具有访问所有命名空间的权限
	all, err := k.CanVisitAllNamespace(username)
	if err != nil {
		return nil, err
	}
	// 创建新的字符串集合，用于存储命名空间名称
	namespaceSet := NewStringSet()
	if all {
		// 如果用户具有访问所有命名空间的权限
		// 列出所有命名空间
		ns, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		// 遍历命名空间列表，并将状态为 "Active" 的命名空间名称添加到集合中
		for i := range ns.Items {
			if ns.Items[i].Status.Phase == "Active" {
				namespaceSet.Add(ns.Items[i].Name)
			}
		}
	} else {
		// 如果用户不具有访问所有命名空间的权限
		// 列出所有 RoleBindings
		rbs, err := client.RbacV1().RoleBindings("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		// 遍历 RoleBindings 列表，并将用户为指定用户名的命名空间名称添加到集合中
		for i := range rbs.Items {
			for j := range rbs.Items[i].Subjects {
				if rbs.Items[i].Subjects[j].Kind == "User" && rbs.Items[i].Subjects[j].Name == username {
					namespaceSet.Add(rbs.Items[i].Namespace)
				}
			}
		}
	}
	// 将集合转换为切片，并进行排序
	result := namespaceSet.ToSlice()
	sort.Strings(result)
	return result, nil
}

// CanVisitAllNamespace
//
//	@Description: 检查指定用户是否具有访问所有命名空间的权限
//	@receiver k
//	@return bool 如果用户可以访问所有命名空间，返回true，反之返回false
//	@return error
func (k *Kubernetes) CanVisitAllNamespace(username string) (bool, error) {
	client, err := k.Client()
	if err != nil {
		return false, err
	}
	// 创建新的字符串集合，用于存储角色名称
	roleSet := NewStringSet()
	// 构建标签选择器，用于过滤 ClusterRoleBindings
	labels := []string{
		fmt.Sprintf("%s=%s", LabelManageKey, "devops"),
		fmt.Sprintf("%s=%s", LabelClusterId, k.K8sCluster.UUID),
		fmt.Sprintf("%s=%s", LabelUsername, username),
	}
	// 列出所有符合标签选择器的 ClusterRoleBindings
	clusterrolebindings, err := client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{
		LabelSelector: strings.Join(labels, ","),
	})
	if err != nil {
		return false, err
	}
	// 将每个 ClusterRoleBinding 中的角色名称添加到 roleSet 中
	for i := range clusterrolebindings.Items {
		roleSet.Add(clusterrolebindings.Items[i].RoleRef.Name)
	}
	// 遍历 roleSet 中的每个角色名称
	for _, roleName := range roleSet.ToSlice() {
		// 获取 ClusterRole 对象
		role, err := client.RbacV1().ClusterRoles().Get(context.TODO(), roleName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		// 检查 ClusterRole 中的规则，判断是否具有访问所有命名空间的权限
		for i := range role.Rules {
			// 如果 APIGroups 和 Resources 均包含 "*"，则表示具有全局访问权限
			if IndexOfStringSlice(role.Rules[i].APIGroups, "*") != -1 && IndexOfStringSlice(role.Rules[i].Resources, "*") != -1 {
				return true, nil
			}
		}
	}
	return false, nil
}

// IndexOfStringSlice 返回字符串在切片中的索引，如果未找到则返回 -1
func IndexOfStringSlice(s []string, target string) int {
	for i := range s {
		if s[i] == target {
			return i
		}
	}
	return -1
}

// NodeMatchesKeyword
//
//	@Description: 根据用户上传的字段查询node列表
//	@receiver k
//	@return bool
func (k *Kubernetes) NodeMatchesKeyword(node corev1.Node, keyword string) bool {
	if strings.Contains(node.Name, keyword) {
		return true
	}

	for _, address := range node.Status.Addresses {
		if strings.Contains(address.Address, keyword) {
			return true
		}
	}

	return false
}

func (k *Kubernetes) DeploymentMatchesKeyword(deployment appsv1.Deployment, keyword string) bool {
	if strings.Contains(deployment.Name, keyword) {
		return true
	}

	return false
}
