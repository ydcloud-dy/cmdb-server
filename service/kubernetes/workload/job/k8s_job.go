package job

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/job"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sJobService struct {
	kubernetes.BaseService
}

func (k *K8sJobService) GetJobList(req job.GetJobListReq, uuid uuid.UUID) (*[]v1.Job, int, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	options := metav1.ListOptions{LabelSelector: req.LabelSelector, FieldSelector: req.FieldSelector}

	data, err := client.BatchV1().Jobs(req.Namespace).List(context.TODO(), options)

	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterJobs []v1.Job
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterJobs = append(filterJobs, item)
			}
		}
	} else {
		filterJobs = data.Items
	}

	result, err := paginate.Paginate(filterJobs, req.Page, req.PageSize)

	return result.(*[]v1.Job), len(filterJobs), nil
}

func (k *K8sJobService) DescribeJob(req job.DescribeJobReq, uuid uuid.UUID) (*v1.Job, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, err
	}
	JobIns, err := client.BatchV1().Jobs(req.Namespace).Get(context.TODO(), req.JobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return JobIns, nil
}
func (k *K8sJobService) UpdateJob(req job.UpdateJobReq, uuid uuid.UUID) (*v1.Job, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}

	tmp := string(data)
	var JobIns *v1.Job
	err = json.Unmarshal([]byte(tmp), &JobIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.BatchV1().Jobs(req.Namespace).Update(context.TODO(), JobIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sJobService) DeleteJob(req job.DeleteJobReq, uuid uuid.UUID) error {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}

	err = client.BatchV1().Jobs(req.Namespace).Delete(context.TODO(), req.JobName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}
	return nil
}
func (k *K8sJobService) CreateJob(req job.CreateJobReq, uuid uuid.UUID) (*v1.Job, error) {
	kubernetes, err := k.Generic(&req, uuid)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	client, err := kubernetes.Client()
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	data, err := json.Marshal(req.Content)
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	var Job *v1.Job
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &Job)
	ins, err := client.BatchV1().Jobs(req.Namespace).Create(context.TODO(), Job, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
