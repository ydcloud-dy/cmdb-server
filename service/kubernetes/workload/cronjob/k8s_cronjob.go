package cronjob

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/cronjob"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sCronJobService struct {
	kubernetes.BaseService
}

func (k *K8sCronJobService) GetCronJobList(req cronjob.GetCronJobListReq, uuid uuid.UUID) (*[]v1.CronJob, int, error) {
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

	data, err := client.BatchV1().CronJobs(req.Namespace).List(context.TODO(), options)
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var filterCronJobs []v1.CronJob
	if req.Keyword != "" {
		for _, item := range data.Items {
			if strings.Contains(item.Name, req.Keyword) {
				filterCronJobs = append(filterCronJobs, item)
			}
		}
	} else {
		filterCronJobs = data.Items
	}

	result, err := paginate.Paginate(filterCronJobs, req.Page, req.PageSize)

	return result.(*[]v1.CronJob), len(filterCronJobs), nil
}

func (k *K8sCronJobService) DescribeCronJob(req cronjob.DescribeCronJobReq, uuid uuid.UUID) (*v1.CronJob, error) {
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
	CronJobIns, err := client.BatchV1().CronJobs(req.Namespace).Get(context.TODO(), req.CronJobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return CronJobIns, nil
}
func (k *K8sCronJobService) UpdateCronJob(req cronjob.UpdateCronJobReq, uuid uuid.UUID) (*v1.CronJob, error) {
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
	var CronJobIns *v1.CronJob
	err = json.Unmarshal([]byte(tmp), &CronJobIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.BatchV1().CronJobs(req.Namespace).Update(context.TODO(), CronJobIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sCronJobService) DeleteCronJob(req cronjob.DeleteCronJobReq, uuid uuid.UUID) error {
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

	err = client.BatchV1().CronJobs(req.Namespace).Delete(context.TODO(), req.CronJobName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return err
	}
	return nil
}
func (k *K8sCronJobService) CreateCronJob(req cronjob.CreateCronJobReq, uuid uuid.UUID) (*v1.CronJob, error) {
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
	var CronJob *v1.CronJob
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &CronJob)
	ins, err := client.BatchV1().CronJobs(req.Namespace).Create(context.TODO(), CronJob, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
