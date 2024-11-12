package cicd

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cicd"
	request "DYCLOUD/model/cicd/request"
	"fmt"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelinesService struct{}

// GetPipelinesList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *PipelinesService) GetPipelinesList(req *request.ApplicationRequest) (envList *[]cicd.Pipelines, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DYCLOUD_DB.Model(&cicd.Pipelines{})

	// 创建db
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		db = db.Where("name like ?", keyword).Or("id = ?", req.Keyword)
	}
	if !req.StartCreatedAt.IsZero() && !req.EndCreatedAt.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartCreatedAt, req.EndCreatedAt)
		db = db.Where("name = ?", req.Keyword)
	}
	var data []cicd.Pipelines
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&data).Error
	if err != nil {
		return nil, 0, nil
	}
	return &data, total, nil
}

// DescribePipelines
//
//	@Description: 查看应用详情
//	@receiver e
//	@param id
//	@return *cicd.Pipelines
//	@return error
func (e *PipelinesService) DescribePipelines(id int) (*cicd.Pipelines, error) {
	var data *cicd.Pipelines
	if err := global.DYCLOUD_DB.Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// CreatePipelines
//
//	@Description: 创建应用
//	@receiver e
//	@param req
//	@return error
func (e *PipelinesService) CreatePipelines(clientSet *tektonclient.Clientset, req *cicd.Pipelines) error {
	pipelineRun := &tektonv1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-pipeline-", req.AppName),
			Namespace:    req.K8SNamespace,
		},
		Spec: tektonv1.PipelineRunSpec{
			TaskRunTemplate: tektonv1.PipelineTaskRunTemplate{
				ServiceAccountName: "default", // 在这里指定 ServiceAccount
			},
			PipelineSpec: &tektonv1.PipelineSpec{
				Tasks: []tektonv1.PipelineTask{
					{
						Name: "clone-source",
						TaskSpec: &tektonv1.EmbeddedTask{
							TaskSpec: tektonv1.TaskSpec{
								Steps: []tektonv1.Step{
									{
										Name:    "clone",
										Image:   "192.168.31.41:82/cicd/git:v2.45.2",
										Command: []string{"/bin/sh"},
										Args: []string{
											"-c",
											// 清空目录后执行 git clone
											"rm -rf $(workspaces.source.path)/* && git clone -b " + req.GitBranch + " " + req.GitUrl + " $(workspaces.source.path)",
										},
										WorkingDir: "$(workspaces.source.path)",
									},
								},
							},
						},
						Workspaces: []tektonv1.WorkspacePipelineTaskBinding{
							{Name: "source", Workspace: "source"},
						},
					},
					{
						Name:     "build-project",
						RunAfter: []string{"clone-source"},
						TaskSpec: &tektonv1.EmbeddedTask{
							TaskSpec: tektonv1.TaskSpec{
								Steps: []tektonv1.Step{
									{
										Name:   "build",
										Image:  "192.168.31.41:82/cicd/maven:3.8.7-eclipse-temurin-11-alpine",
										Script: req.BuildScript,
										//Args: []string{
										//	"-c",
										//	// 下载 settings.xml 并执行构建
										//	"curl -sL https://gitee.com/mageedu/spring-boot-helloWorld/raw/main/maven/settings.xml -o /usr/share/maven/conf/settings.xml && cd $(workspaces.source.path) && mvn clean install",
										//},
										WorkingDir: "$(workspaces.source.path)",
										VolumeMounts: []corev1.VolumeMount{
											{
												Name:      "maven-cache", // 与 volumes 中的名称一致
												MountPath: "/root/.m2",   // Maven 缓存目录
											},
										},
									},
								},
								Volumes: []corev1.Volume{
									{
										Name: "maven-cache",
										VolumeSource: corev1.VolumeSource{
											PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
												ClaimName: "maven-cache", // 这里是您已经创建的 PVC 名称
											},
										},
									},
								},
							},
						},
						Workspaces: []tektonv1.WorkspacePipelineTaskBinding{
							{Name: "source", Workspace: "source"},
						},
					},
					{
						Name:     "docker-build",
						RunAfter: []string{"build-project"},
						TaskSpec: &tektonv1.EmbeddedTask{
							TaskSpec: tektonv1.TaskSpec{
								Steps: []tektonv1.Step{
									{
										Name:    "kaniko",
										Image:   "192.168.31.41:82/cicd/executor:v1.23.2",
										Command: []string{"/kaniko/executor"},
										Args: []string{
											"--dockerfile", req.DockerfilePath,
											"--destination", fmt.Sprintf("%s/%s:%s", req.RegistryURL, req.ImageName, req.ImageTag),
											"--context", "$(workspaces.source.path)",
											"--insecure",
											"--skip-tls-verify",
											"--skip-tls-verify-pull",
											"--skip-push-permission-check",
										},
										Env: []corev1.EnvVar{
											{Name: "DOCKER_CONFIG", Value: "/kaniko/.docker"},
											{Name: "DOCKER_USERNAME", Value: req.RegistryUser},
											{Name: "DOCKER_PASSWORD", Value: req.RegistryPass},
										},
										VolumeMounts: []corev1.VolumeMount{
											{
												Name:      "kaniko-config",
												MountPath: "/kaniko/.docker/config.json",
												SubPath:   "config.json",
											},
										},
									},
								},
								Volumes: []corev1.Volume{
									{
										Name: "kaniko-config",
										VolumeSource: corev1.VolumeSource{
											Secret: &corev1.SecretVolumeSource{
												SecretName: "registry-secret", // 引用 generic Secret
												Items: []corev1.KeyToPath{
													{
														Key:  ".dockerconfigjson",
														Path: "config.json",
													},
												},
											},
										},
									},
								},
							},
						},
						Workspaces: []tektonv1.WorkspacePipelineTaskBinding{
							{Name: "source", Workspace: "source"},
						},
					},
					// 添加 deploy-using-kubectl 任务
					{
						Name:     "deploy-using-kubectl",
						RunAfter: []string{"docker-build"},
						TaskSpec: &tektonv1.EmbeddedTask{
							TaskSpec: tektonv1.TaskSpec{
								Params: []tektonv1.ParamSpec{
									{Name: "deploy-config-file", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: "03-deployment.yaml"}},
									{Name: "image-url", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: fmt.Sprintf("%s/%s", req.RegistryURL, req.ImageName)}},
									{Name: "image-tag", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: req.ImageTag}},
								},
								Steps: []tektonv1.Step{
									{
										Name:    "update-yaml",
										Image:   "192.168.31.41:82/cicd/alpine:latest",
										Command: []string{"sed"},
										Args: []string{
											"-i", "-e",
											"s@__IMAGE__@$(params.image-url):$(params.image-tag)@g",
											"$(workspaces.source.path)/deploy/$(params.deploy-config-file)",
										},
									},
									{
										Name:    "run-kubectl",
										Image:   "192.168.31.41:82/cicd/kubectl",
										Command: []string{"kubectl"},
										Args: []string{
											"apply", "-f",
											"$(workspaces.source.path)/deploy/$(params.deploy-config-file)",
										},
									},
								},
							},
						},
						Workspaces: []tektonv1.WorkspacePipelineTaskBinding{
							{Name: "source", Workspace: "source"},
						},
					},
				},
				Workspaces: []tektonv1.PipelineWorkspaceDeclaration{
					{Name: "source"},
				},
			},
			Workspaces: []tektonv1.WorkspaceBinding{
				{
					Name: "source",
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "shared-workspace-pvc", // 使用您已有的PVC
					},
					SubPath: fmt.Sprintf("run-%s", metav1.Now().Format("20060102150405")),
				},
			},
		},
	}

	// 创建 PipelineRun
	createdPipelineRun, err := clientSet.TektonV1().PipelineRuns(req.K8SNamespace).Create(context.Background(), pipelineRun, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create PipelineRun: %v", err)
	}
	// 将创建的 PipelineRun 信息保存到数据库
	newPipeline := &cicd.Pipelines{
		Name:           createdPipelineRun.Name,
		AppName:        req.AppName,
		EnvName:        req.EnvName,
		BuildScript:    req.BuildScript,
		K8SNamespace:   req.K8SNamespace,
		BaseImage:      req.BaseImage,
		DockerfilePath: req.DockerfilePath,
		ImageName:      req.ImageName,
		ImageTag:       req.ImageTag,
		RegistryURL:    req.RegistryURL,
		RegistryUser:   req.RegistryUser,
		RegistryPass:   req.RegistryPass,
		GitUrl:         req.GitUrl,
		GitBranch:      req.GitBranch,
		GitCommitId:    req.GitCommitId,
	}

	// 保存到数据库
	err = global.DYCLOUD_DB.Create(newPipeline).Error
	if err != nil {
		return fmt.Errorf("failed to save PipelineRun to database: %v", err)
	}
	return nil
}

// SyncBranches
//
//	@Description: 同步应用的分支信息
//	@receiver e
//	@param id
//	@return error
func (e *PipelinesService) CreateAppBranchIfNotExist(branch *cicd.AppBranch) (int, error) {
	result := global.DYCLOUD_DB.Where("branch_name = ? and app_id = ?", branch.BranchName, branch.AppID).FirstOrCreate(branch)
	if result.Error != nil {
		return 0, result.Error
	}
	// 检查是否是创建的新记录，或者已存在
	if result.RowsAffected == 0 {
		return int(branch.ID), fmt.Errorf("branch_name: %v already exists in app branch table", branch.BranchName)
	}
	return int(branch.ID), nil

}
func (e *PipelinesService) UpdateAppBranch(branch *cicd.AppBranch) error {
	err := global.DYCLOUD_DB.Model(&cicd.AppBranch{}).Where("id = ?", branch.ID).Updates(branch).Error
	return err
}

// UpdatePipelines
//
//	@Description: 更新应用
//	@receiver e
//	@param req
//	@return *cicd.Pipelines
//	@return error
func (e *PipelinesService) UpdatePipelines(req *cicd.Pipelines) (*cicd.Pipelines, error) {
	//fmt.Println(req)
	//data, err := e.DescribePipelines(int(req.ID))
	//if err != nil {
	//	return nil, err
	//}
	//data = req
	//if err = global.DYCLOUD_DB.Model(&cicd.Pipelines{}).Where("id = ?", req.ID).Omit("ID").Updates(&req).Error; err != nil {
	//	return nil, err
	//}
	return nil, nil
}

// DeletePipelines
//
//	@Description: 删除应用
//	@receiver e
//	@param id
//	@return error
func (e *PipelinesService) DeletePipelines(id int) error {
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Model(&cicd.Pipelines{}).Where("id = ?", id).Delete(&cicd.Pipelines{}).Error; err != nil {
		return err
	}
	return nil
}

// DeletePipelinesByIds
//
//	@Description: 批量删除应用
//	@receiver e
//	@param ids
//	@return error
func (e *PipelinesService) DeletePipelinesByIds(ids *request.DeleteApplicationByIds) error {
	fmt.Println(ids)
	if err := global.DYCLOUD_DB.Model(&cicd.Pipelines{}).Where("id in ?", ids.Ids).Delete(&cicd.Pipelines{}).Error; err != nil {
		return err
	}
	return nil
}
