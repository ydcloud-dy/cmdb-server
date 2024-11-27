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
	pipelineTasks := []tektonv1.PipelineTask{}

	// 第一步：拉取代码的任务
	pipelineTasks = append(pipelineTasks, tektonv1.PipelineTask{
		Name: "clone-source",
		TaskSpec: &tektonv1.EmbeddedTask{
			TaskSpec: tektonv1.TaskSpec{
				Steps: []tektonv1.Step{
					{
						Name:    "clone",
						Image:   "swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/alpine/git:v2.45.2",
						Command: []string{"/bin/sh"},
						Args: []string{
							"-c",
							fmt.Sprintf("rm -rf $(workspaces.source.path)/* && git clone -b %s %s $(workspaces.source.path)", req.GitBranch, req.GitUrl),
						},
						WorkingDir: "$(workspaces.source.path)",
					},
				},
			},
		},
		Workspaces: []tektonv1.WorkspacePipelineTaskBinding{
			{Name: "source", Workspace: "source"},
		},
	})

	// 遍历前端传递的阶段数并创建 PipelineTask
	previousTaskName := "clone-source" // 第一个任务是 clone-source，后续任务将依赖于它
	for _, stage := range req.Stages {
		for _, task := range stage.TaskList {
			var pipelineTask tektonv1.PipelineTask

			// 根据任务类型选择不同的逻辑
			switch task.Branch {
			case "1":
				pipelineTask = tektonv1.PipelineTask{
					Name:     task.Name,
					RunAfter: []string{previousTaskName}, // 依赖于前一个任务的完成
					TaskSpec: &tektonv1.EmbeddedTask{
						TaskSpec: tektonv1.TaskSpec{
							Steps: []tektonv1.Step{
								{
									Name:       "execute-script",
									Image:      task.Image,
									Script:     task.Script,
									WorkingDir: "$(workspaces.source.path)",
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "maven-cache",
											MountPath: "/root/.m2", // 缓存 Maven 依赖，如果不是 Maven，也可以灵活替换路径
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "maven-cache",
									VolumeSource: corev1.VolumeSource{
										PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
											ClaimName: "maven-cache-pvc", // 使用你创建的 PVC
										},
									},
								},
							},
						},
					},
				}
			case "2":
				fmt.Printf("Building Image URL: %s:%s", fmt.Sprintf("%s/%s", task.Warehouse, task.SpatialName), task.MirrorTag)
				pipelineTask = tektonv1.PipelineTask{
					Name:     task.Name,
					RunAfter: []string{previousTaskName}, // 依赖于前一个任务的完成
					TaskSpec: &tektonv1.EmbeddedTask{
						TaskSpec: tektonv1.TaskSpec{
							Results: []tektonv1.TaskResult{
								{Name: "built-image-url", Description: "The URL of the built image."},
							},
							Params: []tektonv1.ParamSpec{
								{Name: "dockerfile", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: task.Dockerfile}},
								{Name: "image-url", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: fmt.Sprintf("%s/%s", task.Warehouse, task.SpatialName)}},
								{Name: "image-tag", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: task.MirrorTag}},
							},
							Workspaces: []tektonv1.WorkspaceDeclaration{
								{Name: "source"},
							},
							Steps: []tektonv1.Step{
								{
									Name:    "docker-build",
									Image:   "registry.cn-hangzhou.aliyuncs.com/dyclouds/executor:v1.23.2",
									Command: []string{"/kaniko/executor"},
									Args: []string{
										"--dockerfile=$(params.dockerfile)",
										"--context=$(workspaces.source.path)",
										"--destination=$(params.image-url):$(params.image-tag)",
										"--insecure",
										"--skip-tls-verify",
									},
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      "kaniko-config",
											MountPath: "/kaniko/.docker/config.json",
											SubPath:   ".dockerconfigjson", // 修改为正确的 SubPath
										},
									},
									Env: []corev1.EnvVar{
										{Name: "DOCKER_CONFIG", Value: "/kaniko/.docker"},
									},
								},
								{
									Name:    "output-result",
									Image:   "swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/alpine:latest",
									Command: []string{"sh", "-c"},
									Args: []string{
										"echo $(params.image-url):$(params.image-tag) > /tekton/results/built-image-url",
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: "kaniko-config",
									VolumeSource: corev1.VolumeSource{
										Secret: &corev1.SecretVolumeSource{
											SecretName: "registry-secret", // 正确引用 k8s Secret
											Items: []corev1.KeyToPath{
												{
													Key:  ".dockerconfigjson", // 确保与 Secret 的 Key 匹配
													Path: ".dockerconfigjson", // 挂载为 config.json
												},
											},
										},
									},
								},
							},
						},
					},
				}
			case "4":
				fmt.Println("yaml: ", task.YamlResource)
				pipelineTask = tektonv1.PipelineTask{
					Name:     task.Name,
					RunAfter: []string{previousTaskName}, // 依赖于前一个任务的完成
					Params: []tektonv1.Param{
						{
							Name: "yaml-content",
							Value: tektonv1.ParamValue{
								Type:      tektonv1.ParamTypeString,
								StringVal: task.YamlResource, // 从前端 JSON 获取完整 YAML 内容
							},
						},
						{
							Name: "built-image-url",
							Value: tektonv1.ParamValue{
								Type:      tektonv1.ParamTypeString,
								StringVal: "$(tasks." + previousTaskName + ".results.built-image-url)", // 动态引用结果
							},
						},
					},
					TaskSpec: &tektonv1.EmbeddedTask{
						TaskSpec: tektonv1.TaskSpec{
							Params: []tektonv1.ParamSpec{
								{
									Name:        "yaml-content",
									Type:        tektonv1.ParamTypeString,
									Description: "YAML resource for deployment",
								},
								{
									Name:        "built-image-url",
									Type:        tektonv1.ParamTypeString,
									Description: "Image URL to replace in YAML",
								},
							},
							Steps: []tektonv1.Step{
								{
									Name:  "update-yaml", // 更新 YAML 文件中的镜像信息
									Image: "registry.cn-hangzhou.aliyuncs.com/dyclouds/alpine:latest",
									Script: `
										# 打印调试信息
										echo "Original YAML Content:"
										echo "$(params.yaml-content)"
										echo "Built Image URL:"
										echo "$(params.built-image-url)"
							
										# 写入原始 YAML 内容
										printf "%s" "$(params.yaml-content)" > $(workspaces.source.path)/original-deployment.yaml
										IMAGE_URL=$(echo "$(params.built-image-url)" | sed 's/[&/\]/\\&/g')

										# 使用 sed 替换镜像 URL
										sed "s|__IMAGE__|$IMAGE_URL|g" $(workspaces.source.path)/original-deployment.yaml > $(workspaces.source.path)/updated-deployment.yaml
										#sed 's/__IMAGE__/registry.cn-hangzhou.aliyuncs.com\/dyclouds\/yiyuetong:v1.0.0/g'  $(workspaces.source.path)/original-deployment.yaml > $(workspaces.source.path)/updated-deployment.yaml

										# 打印更新后的 YAML 文件
										cat $(workspaces.source.path)/updated-deployment.yaml
									`,
								},
								{
									Name:    "apply-kubectl", // 应用更新后的 YAML 文件
									Image:   task.Image,      // 使用前端传递的 kubectl 镜像
									Command: []string{"kubectl"},
									Args: []string{
										"apply", "-f", "$(workspaces.source.path)/updated-deployment.yaml",
									},
								},
							},
						},
					},
					Workspaces: []tektonv1.WorkspacePipelineTaskBinding{
						{Name: "source", Workspace: "source"},
					},
				}
			case "3":
				pipelineTask = tektonv1.PipelineTask{
					Name:     task.Name,
					RunAfter: []string{previousTaskName}, // 依赖于前一个任务的完成
					TaskSpec: &tektonv1.EmbeddedTask{
						TaskSpec: tektonv1.TaskSpec{
							Steps: []tektonv1.Step{
								{
									Name:    "upload-to-oss",
									Image:   task.Image,
									Command: []string{"/bin/sh"},
									Args: []string{
										"-c",
										fmt.Sprintf("ossutil cp %s oss://%s/%s -u", task.ProductPath, task.SpatialName, task.ProductName),
									},
									WorkingDir: "$(workspaces.source.path)",
								},
							},
						},
					},
				}
			}

			// 将每个创建的任务添加到 pipelineTasks 列表中
			pipelineTasks = append(pipelineTasks, pipelineTask)

			// 更新 previousTaskName，确保下一个任务依赖于当前任务
			previousTaskName = task.Name
		}
	}

	// 创建 PipelineRun
	pipelineRun := &tektonv1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-pipeline-", req.AppName),
			Namespace:    req.K8SNamespace,
		},
		Spec: tektonv1.PipelineRunSpec{
			TaskRunTemplate: tektonv1.PipelineTaskRunTemplate{
				ServiceAccountName: "deploy-sa", // 使用新创建的 ServiceAccount

			},
			PipelineSpec: &tektonv1.PipelineSpec{
				Tasks: pipelineTasks,
			},
			Workspaces: []tektonv1.WorkspaceBinding{
				{
					Name: "source",
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "shared-workspace-pvc", // 使用你已有的 PVC
					},
					SubPath: fmt.Sprintf("run-%s", metav1.Now().Format("20060102150405")),
				},
				{
					Name: "maven-cache",
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: "maven-cache-pvc", // 缓存 Maven 构建依赖的 PVC
					},
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
