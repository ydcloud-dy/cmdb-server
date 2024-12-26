package dao

import (
	"DYCLOUD/global"
	"DYCLOUD/model/configCenter"
	"fmt"
	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/gitea"
	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/driver/gitlab"
	"github.com/drone/go-scm/scm/driver/gogs"
	"github.com/drone/go-scm/scm/transport"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

// 验证仓库源是否能正常连通，若无token，只要地址能通就行；若有token，则必须认证通过
func VerifyRepoConnetion(scmType configCenter.ServiceType, url string, token string) error {
	//  创建对应类型的SCM客户端
	scmClient, err := NewScmProvider(scmType, url, token)
	if err != nil {
		return err
	}
	// 定义分页选项，只请求第一页内容，放置加载过多数据
	op := scm.ListOptions{
		Page: 0, // 设置请求第一页
		Size: 1, // 每页请求一个寄过
	}
	// 尝试获取组织信息，以测试连接性
	_, resp, err := scmClient.Organizations.List(context.Background(), op)
	if resp == nil && err != nil {
		return fmt.Errorf("连接源码仓库失败,错误信息：%s", err)
	}
	// 如果没有提供token，则允许返回401（未授权），表示仓库可连接但没有认证
	if token == "" {
		// 检查响应状态码，200-299 ，401 表示未授权
		if (resp.Status >= 200 && resp.Status <= 299) || resp.Status == 401 {
			return nil
		}
		return fmt.Errorf("连接源码仓库失败,服务器返回：%d", resp.Status)
	} else {
		// 如果提供了 token，要求响应状态码在 200-299 范围内，表示认证通过
		if resp.Status >= 200 && resp.Status <= 299 {
			return nil
		} else if resp.Status == 401 {
			return fmt.Errorf("连接源码仓库成功，但是认证失败,仓库返回：401")
		} else {
			return fmt.Errorf("连接源码仓库失败,服务器返回：%d", resp.Status)
		}
	}
}

// NewScmProvider
//
//	@Description: 创建一个SCM客户端，根据不同的代码管理系统类型（github、gitlab等）初始化不同的客户端
//	@param vcsType
//	@param vcsPath
//	@param token
//	@return *scm.Client
//	@return error
func NewScmProvider(vcsType configCenter.ServiceType, vcsPath, token string) (*scm.Client, error) {
	var err error
	var client *scm.Client

	// 根据 vcsType 使用不同的 SCM 驱动
	switch vcsType {
	case configCenter.GITEA_TYPE:
		client, err = createSCMClient(gitea.New, vcsPath)
	case configCenter.GITLAB_TYPE:
		client, err = createSCMClient(gitlab.New, vcsPath)
	case configCenter.GOGS_TYPE:
		client, err = createSCMClient(gogs.New, vcsPath)
	case configCenter.GITHUB_TYPE:
		client = github.NewDefault()
	case configCenter.GITEE_TYPE:
		client = gitee.NewDefault()
	default:
		err = fmt.Errorf("source code management system not configured")
	}

	// 如果客户端成功创建并提供了 token，则设置 HTTP 客户端，用于处理认证
	if client != nil {
		client.Client = getSCMHttpClient(vcsType, token)
	}

	return client, err
}

// createSCMClient 是一个辅助函数，用于处理创建 SCM 客户端并去除 `.git` 后缀
func createSCMClient(newFunc func(string) (*scm.Client, error), vcsPath string) (*scm.Client, error) {
	// 去掉路径后缀 ".git"，确保 URL 格式一致
	if strings.HasSuffix(vcsPath, ".git") {
		vcsPath = strings.TrimSuffix(vcsPath, ".git")
	}

	// 分割 URL 的 schema 和路径部分
	vcsPathSplit := strings.Split(vcsPath, "://")
	if len(vcsPathSplit) < 2 {
		return nil, fmt.Errorf("invalid VCS path")
	}

	// 使用指定的构造函数创建 SCM 客户端
	return newFunc(vcsPathSplit[0] + "://" + strings.Split(vcsPathSplit[1], "/")[0])
}

// 根据不同类型获取客户端，若有token，则配置对应token；无token,表示公共库，不配置鉴权信息
func getSCMHttpClient(scmType configCenter.ServiceType, token string) *http.Client {
	if token == "" {
		return &http.Client{}
	}

	// 根据 scmType 设置对应的认证方式
	switch scmType {
	case configCenter.GITLAB_TYPE, configCenter.GOGS_TYPE:
		return &http.Client{
			Transport: &transport.PrivateToken{
				Token: token, // 使用私有 token 认证
			},
		}
	case configCenter.GITEA_TYPE, configCenter.GITEE_TYPE, configCenter.GITHUB_TYPE:
		return &http.Client{
			Transport: &transport.BearerToken{
				Token: token, // 使用 Bearer token 认证
			},
		}
	default:
		return nil
	}
}
func GetScmConf(scmType configCenter.ServiceType, config interface{}) configCenter.GitConfig {
	scmCONF := configCenter.GitConfig{}
	switch scmType {
	case configCenter.GITLAB_TYPE, configCenter.GITEA_TYPE, configCenter.GITEE_TYPE, configCenter.GITHUB_TYPE:
		if conf, ok := config.(*configCenter.GitConfig); ok {
			scmCONF.Url = conf.Url
			scmCONF.UserName = conf.UserName
			scmCONF.Token = conf.Token
		} else {
			global.DYCLOUD_LOG.Error(fmt.Sprintf("parse type: %s conf error", configCenter.GITLAB_TYPE))
		}
	}
	return scmCONF
}
