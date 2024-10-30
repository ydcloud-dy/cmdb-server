package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/utils/cicd"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// 定义service类型
type ServiceType int

const (
	// 0 为k8s
	KUBERNETES_TYPE ServiceType = iota
	// 1 为镜像仓库
	REGISTRY_TYPE
	// 2 为jenkins
	JENKINS_TYPE
	// 3 为git仓库
	GITLAB_TYPE
	GITHUB_TYPE
	GITEE_TYPE
	GITEA_TYPE
	GOGS_TYPE
)

// 定义K8s连接类型
type K8sType int

const (
	// 0 为config
	KUBERNETES_CONFIG K8sType = iota
	// 1 为token
	KUBERNETES_TOKEN
)

type GitType int

const (
	GITLAB GitType = iota
	GITHUB
	GITEE
	GITEA
	GOGS
)

// 私有还是开源
type PrivacyType int

const (
	GIT_PRIVATE PrivacyType = iota
	GIT_PUBLIC
)

// ServiceIntegration
// @Description: 服务列表
type ServiceIntegration struct {
	Name      string          `json:"name" form:"name"`
	Desc      string          `json:"desc" form:"desc"`
	Config    json.RawMessage `json:"config" form:"config"`
	ConfigStr string          `json:"-" gorm:"column:config"`
	Type      ServiceType     `json:"type" form:"type"`
	global.DYCLOUD_MODEL
	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

func (s *ServiceIntegration) TableName() string {
	return "cicd_services"
}
func (s *ServiceIntegration) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

// crypto
//
//	@Description: 将传入的config加密
//	@receiver s
//	@param raw
//	@return string
func (s *ServiceIntegration) crypto(raw string) string {
	plainText := []byte(raw)
	return base64.StdEncoding.EncodeToString(cicd.AesEny(plainText))
}

// decrypt
//
//	@Description: 解密config
//	@receiver s
//	@return string
func (s *ServiceIntegration) decrypt() (string, error) {
	// 从 `Config` 中取出 base64 编码的加密数据并解码
	encryptedData, err := base64.StdEncoding.DecodeString(string(s.Config))
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 config: %v", err)
	}

	// 解密数据（假设 `cicd.AesEny` 是对称解密函数）
	decryptedData := cicd.AesEny(encryptedData)
	return string(decryptedData), nil
}

func (s *ServiceIntegration) CryptoConfig(raw string) {
	// 加密并进行 base64 编码
	encryptedData := s.crypto(raw)
	s.Config = json.RawMessage([]byte(encryptedData))
}
func (s *ServiceIntegration) DecryptConfig() error {

	decrypted, err := s.decrypt()
	if err != nil {
		return err
	}
	s.Config = json.RawMessage(decrypted)
	fmt.Println("解密后的config", string(s.Config))
	fmt.Println("================================================")

	return nil
}

// ChooseConfig
//
//	@Description: 判断当前的config类型
//	@receiver s
//	@return string
//	@return error
func (s *ServiceIntegration) ChooseConfig() (string, error) {
	var config interface{}
	// 将 s.Config 转为字节
	configBytes, err := json.Marshal(s.Config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %v", err)
	}
	switch s.Type {
	case KUBERNETES_TYPE:
		k := &KubernetesConfig{}
		if err := json.Unmarshal(configBytes, &k); err != nil {
			return "", fmt.Errorf("failed to unmarshal KubernetesConfig", err)
		}
		config = k
	case JENKINS_TYPE:
		j := &JenkinsConfig{}
		if err := json.Unmarshal(configBytes, &j); err != nil {
			return "", fmt.Errorf("failed to unmarshal KubernetesConfig", err)
		}
		config = j
	case REGISTRY_TYPE:
		r := &RegistryConfig{}
		if err := json.Unmarshal(configBytes, &r); err != nil {
			return "", fmt.Errorf("failed to unmarshal KubernetesConfig", err)
		}
		config = r
	case GITLAB_TYPE, GITEA_TYPE, GITHUB_TYPE, GITEE_TYPE, GOGS_TYPE:
		g := &GitConfig{}
		if err := json.Unmarshal(configBytes, &g); err != nil {
			return "", fmt.Errorf("failed to unmarshal KubernetesConfig", err)
		}
		config = g
	default:
		return "", errors.New("type not supported")
	}
	configStr, err := json.Marshal(&config)
	if err != nil {
		return "", err
	}
	return string(configStr), nil
}

type Config struct {
}

// Registry
// @Description: registry服务配置
type RegistryConfig struct {
	Url      string `json:"url" form:"url"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	IsHttps  bool   `json:"isHttps" form:"isHttps"`
}

// KubernetesConfig
// @Description:  k8s 配置
type KubernetesConfig struct {
	Type K8sType `json:"type" form:"type"`
	Url  string  `json:"url" form:"url"`
	Conf string  `json:"conf" form:"conf"`
}

// JenkinsConfig
// @Description:  Jenkins配置
type JenkinsConfig struct {
	Url      string `json:"url" form:"url"`
	Username string `json:"username" form:"username"`
	// 用户token
	Token string `json:"token" form:"token"`
	// agent的工作目录，默认是/home/jenkins/agent
	WorkSpace string `json:"workspace" form:"workspace"`
	// agent命名空间，默认是devops
	Namespace string `json:"namespace" form:"namespace"`
}

type GitConfig struct {
	Type     PrivacyType `json:"type" form:"type"`
	Url      string      `json:"url" form:"url"`
	Token    string      `json:"token" form:"token"`
	UserName string      `json:"username" form:"username"`
}
