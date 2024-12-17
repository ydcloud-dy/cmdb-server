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
	rbacV1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type PipelinesService struct{}

// GetPipelinesList
//
//	@Description: 获取应用环境列表
//	@receiver e
//	@param req
//	@return envList
//	@return err
func (e *PipelinesService) GetPipelinesList(req *request.PipelinesRequest) (envList *[]cicd.Pipelines, total int64, err error) {
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
	if req.AppCode != "" {
		db.Where("app_name = ?", req.AppCode)
	}
	if req.EnvCode != "" {
		db.Where("env_code = ?", req.EnvCode)
	}
	var data []cicd.Pipelines
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	// 使用 Preload 加载关联的 stages 和 tasks
	err = db.Preload("Stages.TaskList").Find(&data).Error
	if err != nil {
		return nil, 0, err
	}
	fmt.Println(data)
	return &data, total, nil
}
func (e *PipelinesService) GetPipelinesStatus(client *tektonclient.Clientset, req *request.PipelinesRequest) (*request.PipelineRunStatus, error) {
	var data *cicd.Pipelines
	if err := global.DYCLOUD_DB.Where("app_name = ? and env_name = ?", req.AppCode, req.EnvCode).First(&data).Error; err != nil {
		return nil, nil
	}
	pipelineRun, err := client.TektonV1().PipelineRuns(req.Namespace).Get(context.TODO(), data.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 提取所需的 PipelineRun 状态信息
	// 检查 Conditions 字段来获取状态
	status := ""
	if len(pipelineRun.Status.Conditions) > 0 {
		// 这里假设第一个 Condition 即为我们关注的状态，实际中可能需要检查具体的 Condition 类型或状态
		status = string(pipelineRun.Status.Conditions[0].Reason)
	}
	// 获取最近的开始时间和完成时间
	lastRunTime := ""
	duration := ""

	if pipelineRun.Status.StartTime != nil && pipelineRun.Status.CompletionTime != nil {
		// 从 metav1.Time 转换为 time.Time
		startTime := pipelineRun.Status.StartTime.Time
		completionTime := pipelineRun.Status.CompletionTime.Time

		// 计算持续时间
		lastRunTime = startTime.String()                  // 最近运行时间
		duration = completionTime.Sub(startTime).String() // 耗时
	}
	// 构造返回结果
	result := request.PipelineRunStatus{
		Name:        pipelineRun.Name, // PipelineRun 名称
		Status:      status,           // 状态
		User:        data.CreatedName, // 假设从数据库中获取到的创建者
		Branch:      data.GitBranch,   // 假设从数据库中获取到的 Git 分支
		LastRunTime: lastRunTime,      // 最近运行时间
		Duration:    duration,         // 耗时
	}
	fmt.Println(result)
	return &result, nil
}

// DescribePipelines
//
//	@Description: 查看应用详情
//	@receiver e
//	@param id
//	@return *cicd.Pipelines
//	@return error
func (e *PipelinesService) DescribePipelines(id int) (*cicd.Pipelines, error) {
	var data cicd.Pipelines
	if err := global.DYCLOUD_DB.
		Preload("Stages.TaskList"). // 预加载 Stages 及其 TaskList
		Preload("Stages.Params").
		Where("id = ?", id).
		First(&data).Error; err != nil {
		return nil, err
	}
	fmt.Println(data)
	return &data, nil
}

// CreatePipelines
//
//	@Description: 创建应用
//	@receiver e
//	@param req
//	@return error
func (e *PipelinesService) CreatePipelines(k8sClient *kubernetes.Clientset, clientSet *tektonclient.Clientset, req *cicd.Pipelines) error {
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
				fmt.Printf("Building Image URL: %s:%s", fmt.Sprintf("%s/%s/%s", task.Warehouse, task.SpatialName, req.AppName), task.MirrorTag)
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
								{Name: "image-url", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: fmt.Sprintf("%s/%s/%s", task.Warehouse, task.SpatialName, req.AppName)}},
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
				fmt.Println("appName：", req.AppName)
				fmt.Println("envName：", req.EnvName)
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
						{
							Name: "app-name",
							Value: tektonv1.ParamValue{
								Type:      tektonv1.ParamTypeString,
								StringVal: req.AppName, // 从请求中获取 AppName
							},
						},
						{
							Name: "env-name",
							Value: tektonv1.ParamValue{
								Type:      tektonv1.ParamTypeString,
								StringVal: req.EnvName, // 从请求中获取 EnvName
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
								{
									Name:        "app-name",
									Type:        tektonv1.ParamTypeString,
									Description: "app name",
								},
								{
									Name:        "env-name",
									Type:        tektonv1.ParamTypeString,
									Description: "env name",
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

										# 动态修改 labels 部分
										APP_NAME="$(params.app-name)"
										ENV_NAME="$(params.env-name)"
										sed -i "s|__APP_ENV_NAME__|$APP_NAME-$ENV_NAME|g" $(workspaces.source.path)/updated-deployment.yaml

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
	// 创建 ServiceAccount
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.K8SNamespace,
		},
	}
	_, err := k8sClient.CoreV1().ServiceAccounts(req.K8SNamespace).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create ServiceAccount: %v", err)
	}

	// 创建 Role
	role := &rbacV1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name + "-role",
			Namespace: req.K8SNamespace,
		},
		Rules: []rbacV1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"get", "list", "watch", "create", "delete", "update", "patch"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments"},
				Verbs:     []string{"get", "list", "watch", "create", "delete", "update", "patch"},
			},
		},
	}
	_, err = k8sClient.RbacV1().Roles(req.K8SNamespace).Create(context.Background(), role, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create Role: %v", err)
	}

	// 创建 RoleBinding
	roleBinding := &rbacV1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name + "-rolebinding",
			Namespace: req.K8SNamespace,
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     req.Name + "-role",
		},
		Subjects: []rbacV1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      req.Name,
				Namespace: req.K8SNamespace,
			},
		},
	}
	_, err = k8sClient.RbacV1().RoleBindings(req.K8SNamespace).Create(context.Background(), roleBinding, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create RoleBinding: %v", err)
	}
	// 创建 PipelineRun
	//pipelineRun := &tektonv1.PipelineRun{
	//	ObjectMeta: metav1.ObjectMeta{
	//		//GenerateName: fmt.Sprintf("%s-pipeline-", req.Name),
	//		Name:      req.Name,
	//		Namespace: req.K8SNamespace,
	//	},
	//	Spec: tektonv1.PipelineRunSpec{
	//		TaskRunTemplate: tektonv1.PipelineTaskRunTemplate{
	//			ServiceAccountName: "deploy-sa", // 使用新创建的 ServiceAccount
	//
	//		},
	//		PipelineSpec: &tektonv1.PipelineSpec{
	//			Tasks: pipelineTasks,
	//		},
	//		Workspaces: []tektonv1.WorkspaceBinding{
	//			{
	//				Name: "source",
	//				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	//					ClaimName: "shared-workspace-pvc", // 使用你已有的 PVC
	//				},
	//				SubPath: fmt.Sprintf("run-%s", metav1.Now().Format("20060102150405")),
	//			},
	//			{
	//				Name: "maven-cache",
	//				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
	//					ClaimName: "maven-cache-pvc", // 缓存 Maven 构建依赖的 PVC
	//				},
	//			},
	//		},
	//	},
	//}
	//
	//// 创建 PipelineRun
	//createdPipelineRun, err := clientSet.TektonV1().PipelineRuns(req.K8SNamespace).Create(context.Background(), pipelineRun, metav1.CreateOptions{})
	//if err != nil {
	//	return fmt.Errorf("failed to create PipelineRun: %v", err)
	//}
	// 创建 Pipeline 模板
	pipeline := &tektonv1.Pipeline{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.K8SNamespace,
		},

		Spec: tektonv1.PipelineSpec{
			Tasks: pipelineTasks,

			Workspaces: []tektonv1.PipelineWorkspaceDeclaration{
				{
					Name: "source",
				},
				{
					Name: "maven-cache",
				},
			},
		},
	}

	// 创建 Pipeline
	_, err = clientSet.TektonV1().Pipelines(req.K8SNamespace).Create(context.Background(), pipeline, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Pipeline: %v", err)
	}

	// 开始事务
	tx := global.DYCLOUD_DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出以触发上层的错误处理
		}
	}()

	// 创建 Pipeline
	newPipeline := &cicd.Pipelines{
		Name:           pipeline.Name,
		AppName:        req.AppName,
		EnvName:        req.EnvName,
		BuildScript:    req.BuildScript,
		K8SNamespace:   req.K8SNamespace,
		K8SClusterName: req.K8SClusterName,
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
		CreatedBy:      req.CreatedBy,
		CreatedName:    req.CreatedName,
	}

	if err := tx.Create(newPipeline).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save Pipeline to database: %v", err)
	}

	// 保存 Stages 和 Tasks 到数据库
	for _, stage := range req.Stages {
		newStage := &cicd.Stage{
			PipelineID: newPipeline.ID, // 使用保存后的 Pipeline ID
			Name:       stage.Name,
		}

		if err := tx.Create(newStage).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save Stage to database: %v", err)
		}

		// 保存 Stage 的 Params
		for _, param := range stage.Params {
			fmt.Println(param)
			newParam := &cicd.Param{
				StageID:      newStage.ID, // 使用保存后的 Stage ID
				PipelineID:   newPipeline.ID,
				Name:         param.Name,
				DefaultValue: param.DefaultValue,
			}

			if err := tx.Create(newParam).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to save Param to database: %v", err)
			}
		}
		for _, task := range stage.TaskList {
			newTask := &cicd.Task{
				StageID:      newStage.ID, // 使用保存后的 Stage ID
				PipelineID:   newPipeline.ID,
				Name:         task.Name,
				Branch:       task.Branch,
				Image:        task.Image,
				Plugin:       task.Plugin,
				Script:       task.Script,
				SpatialName:  task.SpatialName,
				Warehouse:    task.Warehouse,
				MirrorTag:    task.MirrorTag,
				Dockerfile:   task.Dockerfile,
				ContextPath:  task.ContextPath,
				ProductName:  task.ProductName,
				ProductPath:  task.ProductPath,
				Version:      task.Version,
				Resource:     task.Resource,
				YamlResource: task.YamlResource,
				GoalResource: task.GoalResource,
				OpenScript:   task.OpenScript,
			}

			if err := tx.Create(newTask).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to save Task to database: %v", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
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
func (e *PipelinesService) UpdatePipelines(k8sClient *kubernetes.Clientset, clientSet *tektonclient.Clientset, req *cicd.Pipelines) (*cicd.Pipelines, error) {
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
				fmt.Printf("Building Image URL: %s:%s", fmt.Sprintf("%s/%s/%s", task.Warehouse, task.SpatialName, req.AppName), task.MirrorTag)
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
								{Name: "image-url", Type: tektonv1.ParamTypeString, Default: &tektonv1.ParamValue{Type: tektonv1.ParamTypeString, StringVal: fmt.Sprintf("%s/%s/%s", task.Warehouse, task.SpatialName, req.AppName)}},
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
				fmt.Println("appName：", req.AppName)
				fmt.Println("envName：", req.EnvName)
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
						{
							Name: "app-name",
							Value: tektonv1.ParamValue{
								Type:      tektonv1.ParamTypeString,
								StringVal: req.AppName, // 从请求中获取 AppName
							},
						},
						{
							Name: "env-name",
							Value: tektonv1.ParamValue{
								Type:      tektonv1.ParamTypeString,
								StringVal: req.EnvName, // 从请求中获取 EnvName
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
								{
									Name:        "app-name",
									Type:        tektonv1.ParamTypeString,
									Description: "app name",
								},
								{
									Name:        "env-name",
									Type:        tektonv1.ParamTypeString,
									Description: "env name",
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

										# 动态修改 labels 部分
										APP_NAME="$(params.app-name)"
										ENV_NAME="$(params.env-name)"
										sed -i "s|__APP_ENV_NAME__|$APP_NAME-$ENV_NAME|g" $(workspaces.source.path)/updated-deployment.yaml

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

	// 获取已有的 Pipeline
	existingPipeline, err := clientSet.TektonV1().Pipelines(req.K8SNamespace).Get(context.Background(), req.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get existing Pipeline: %v", err)
	}
	fmt.Println(existingPipeline)
	// 修改 Pipeline 的内容
	existingPipeline.Spec.Tasks = pipelineTasks
	existingPipeline.Spec.Workspaces = []tektonv1.PipelineWorkspaceDeclaration{
		{
			Name: "source",
		},
		{
			Name: "maven-cache",
		},
	}
	// 更新 Pipeline
	_, err = clientSet.TektonV1().Pipelines(req.K8SNamespace).Update(context.Background(), existingPipeline, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update Pipeline: %v", err)
	}

	// 开始事务
	tx := global.DYCLOUD_DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出以触发上层的错误处理
		}
	}()

	// 获取并更新 Pipeline
	var pipelines *cicd.Pipelines
	fmt.Println(req.ID)

	// 获取 pipeline 数据
	pipelines, err = e.DescribePipelines(int(req.ID))
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find Pipeline: %v", err)
	}

	// 更新 Pipeline 字段
	pipelines.AppName = req.AppName
	pipelines.EnvName = req.EnvName
	pipelines.BuildScript = req.BuildScript
	pipelines.K8SNamespace = req.K8SNamespace
	pipelines.BaseImage = req.BaseImage
	pipelines.DockerfilePath = req.DockerfilePath
	pipelines.ImageName = req.ImageName
	pipelines.ImageTag = req.ImageTag
	pipelines.RegistryURL = req.RegistryURL
	pipelines.RegistryUser = req.RegistryUser
	pipelines.RegistryPass = req.RegistryPass
	pipelines.GitUrl = req.GitUrl
	pipelines.GitBranch = req.GitBranch
	pipelines.GitCommitId = req.GitCommitId
	pipelines.CreatedBy = req.CreatedBy
	pipelines.CreatedName = req.CreatedName
	fmt.Println("Updating Pipeline:", pipelines)

	if err := tx.Save(&pipelines).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update Pipeline in database: %v", err)
	}
	fmt.Println("Pipeline updated successfully:", pipelines)

	// 获取数据库中的所有 Stage
	var existingStages []cicd.Stage
	if err := tx.Where("pipeline_id = ?", pipelines.ID).Find(&existingStages).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find Stages: %v", err)
	}

	// 创建一个 map 用于标记前端请求中的 Stage
	stageNamesFromRequest := make(map[string]bool)
	for _, stage := range req.Stages {
		stageNamesFromRequest[stage.Name] = true
	}

	// 处理 Stage 更新和删除
	for _, existingStage := range existingStages {
		// 如果这个 Stage 不在请求中，说明它被删除了
		if _, exists := stageNamesFromRequest[existingStage.Name]; !exists {
			// 删除这个 Stage
			if err := tx.Where("id = ?", existingStage.ID).Delete(&cicd.Stage{}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to delete Stage from database: %v", err)
			}
			fmt.Println("Deleted Stage:", existingStage.Name)
		}
	}

	// 处理请求中的每个 Stage
	for _, stage := range req.Stages {
		var existingStage cicd.Stage
		fmt.Println("Received Stage:", stage)

		// 查找数据库中的 Stage
		if err := tx.Where("id = ?", stage.ID).First(&existingStage).Error; err != nil {
			// 如果没有找到这个 Stage，说明是新增的，创建新的 Stage
			if err := tx.Create(&cicd.Stage{
				PipelineID: pipelines.ID,
				Name:       stage.Name,
			}).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create Stage in database: %v", err)
			}

			// 创建新的 Stage 后，重新获取其 ID
			if err := tx.Where(" id = ?", stage.ID).First(&existingStage).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to retrieve newly created Stage: %v", err)
			}
		}

		// 更新 Stage 数据
		updatedStage := cicd.Stage{
			Name: stage.Name, // 使用请求中传递的 name
		}
		fmt.Println("Updating Stage:", updatedStage)
		if err := tx.Model(&existingStage).Where("id = ?", existingStage.ID).Updates(updatedStage).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update Stage in database: %v", err)
		}

		// 更新该 Stage 下的 Task 数据
		for _, task := range stage.TaskList {
			var existingTask cicd.Task
			fmt.Println("Current Task from DB:", task)

			// 查找数据库中对应的 Task，使用 stage_id 和 task.name 来查找
			if err := tx.Where("stage_id = ? AND id = ?", existingStage.ID, task.ID).First(&existingTask).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to find Task: %v", err)
			}

			// 更新 Task 数据
			updatedTask := cicd.Task{
				Branch:       task.Branch,
				Image:        task.Image,
				Plugin:       task.Plugin,
				Script:       task.Script,
				SpatialName:  task.SpatialName,
				Warehouse:    task.Warehouse,
				MirrorTag:    task.MirrorTag,
				Dockerfile:   task.Dockerfile,
				ContextPath:  task.ContextPath,
				ProductName:  task.ProductName,
				ProductPath:  task.ProductPath,
				Version:      task.Version,
				Resource:     task.Resource,
				YamlResource: task.YamlResource,
				GoalResource: task.GoalResource,
				OpenScript:   task.OpenScript,
			}
			fmt.Println("Updating Task:", updatedTask)

			// 更新 Task 数据
			if err := tx.Model(&existingTask).Where("id = ?", existingTask.ID).Updates(updatedTask).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update Task in database: %v", err)
			}
		}

		// 更新该 Stage 下的 Param 数据
		for _, param := range stage.Params {
			var existingParam cicd.Param
			fmt.Println("Current Param from DB:", param)

			// 查找数据库中对应的 Param，使用 stage_id 和 param.name 来查找
			if err := tx.Where("stage_id = ? AND name = ?", existingStage.ID, param.ID).First(&existingParam).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to find Param: %v", err)
			}

			// 更新 Param 数据
			updatedParam := cicd.Param{
				Name:         param.Name,
				DefaultValue: param.DefaultValue,
			}
			fmt.Println("Updating Param:", updatedParam)

			// 更新 Param 数据
			if err := tx.Model(&existingParam).Where("id = ?", existingParam.ID).Updates(updatedParam).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update Param in database: %v", err)
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}
	fmt.Println("Transaction committed successfully")

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

func (e *PipelinesService) GetPipelinesNotice(id int) (*cicd.Notice, error) {
	var result = cicd.Notice{}
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Where("pipeline_id = ? and enable = 1", id).First(&result).Error; err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}

func (e *PipelinesService) ClosePipelineNotice(notice *request.ClosePipelineNotice, pipelineID int) error {
	fmt.Println(notice)
	fmt.Println(pipelineID)
	if err := global.DYCLOUD_DB.Model(&cicd.Notice{}).Where("id = ?", pipelineID).Update("enable", notice.Enable).Error; err != nil {
		return err
	}
	return nil
}

func (e *PipelinesService) CreatePipelinesNotice(req *cicd.Notice) error {
	var notice cicd.Notice

	// 查找是否已存在记录
	result := global.DYCLOUD_DB.Where("pipeline_id = ?", req.PipelineID).First(&notice)

	if result.RowsAffected == 0 {
		// 如果不存在，创建新记录
		if err := global.DYCLOUD_DB.Create(req).Error; err != nil {
			return err
		}
	} else {
		// 如果存在，更新记录
		notice.Enable = req.Enable
		notice.NoticeEvent = req.NoticeEvent
		notice.Webhook = req.Webhook
		notice.NoticeType = req.NoticeType
		if err := global.DYCLOUD_DB.Model(&notice).Updates(notice).Error; err != nil {
			return err
		}
	}
	return nil
}

func (e *PipelinesService) GetPipelinesCache(id int) (*cicd.Cache, error) {
	var result = cicd.Cache{}
	fmt.Println(id)
	if err := global.DYCLOUD_DB.Where("pipeline_id = ? and enable = 1", id).First(&result).Error; err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}

func (e *PipelinesService) ClosePipelineCache(notice *request.ClosePipelineCache, pipelineID int) error {
	fmt.Println(notice)
	fmt.Println(pipelineID)
	if err := global.DYCLOUD_DB.Model(&cicd.Cache{}).Where("id = ?", pipelineID).Update("enable", notice.Enable).Error; err != nil {
		return err
	}
	return nil
}

func (e *PipelinesService) CreatePipelinesCache(req *cicd.Cache) error {
	var cache cicd.Cache

	// 查找是否已存在记录
	result := global.DYCLOUD_DB.Where("pipeline_id = ?", req.PipelineID).First(&cache)

	if result.RowsAffected == 0 {
		// 如果不存在，创建新记录
		if err := global.DYCLOUD_DB.Create(req).Error; err != nil {
			return err
		}
	} else {
		// 如果存在，更新记录
		cache.Enable = req.Enable
		cache.CacheDir = req.CacheDir
		cache.CacheOption = req.CacheOption
		if err := global.DYCLOUD_DB.Model(&cache).Updates(cache).Error; err != nil {
			return err
		}
	}
	return nil
}
