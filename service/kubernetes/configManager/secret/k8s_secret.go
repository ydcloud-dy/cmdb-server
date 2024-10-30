package secret

import (
	"DYCLOUD/global"
	"DYCLOUD/model/kubernetes/secret"
	"DYCLOUD/utils/kubernetes"
	"DYCLOUD/utils/kubernetes/paginate"
	"context"
	"encoding/json"
	"github.com/gofrs/uuid/v5"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type K8sSecretService struct {
	kubernetes.BaseService
}

func (k *K8sSecretService) GetSecretList(req secret.GetSecretList, uuid uuid.UUID) (*[]v1.Secret, int, error) {
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
	data, err := client.CoreV1().Secrets(req.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("获取失败:" + err.Error())
		return nil, 0, err
	}
	var secretList []v1.Secret
	if req.Keyword != "" {
		for _, secret := range data.Items {
			if strings.Contains(secret.Name, req.Keyword) {
				secretList = append(secretList, secret)
			}
		}
	} else {
		secretList = data.Items
	}

	result, err := paginate.Paginate(secretList, req.Page, req.PageSize)

	return result.(*[]v1.Secret), len(secretList), nil
}
func (k *K8sSecretService) DescribeSecret(req secret.DescribeSecretReq, uuid uuid.UUID) (*v1.Secret, error) {
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
	SecretIns, err := client.CoreV1().Secrets(req.Namespace).Get(context.TODO(), req.SecretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return SecretIns, nil
}
func (k *K8sSecretService) UpdateSecret(req secret.UpdateSecretReq, uuid uuid.UUID) (*v1.Secret, error) {
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
	var SecretIns *v1.Secret
	err = json.Unmarshal([]byte(tmp), &SecretIns)
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	result, err := client.CoreV1().Secrets(req.Namespace).Update(context.TODO(), SecretIns, metav1.UpdateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("更新失败:" + err.Error())
		return nil, err
	}
	return result, nil
}
func (k *K8sSecretService) DeleteSecret(req secret.DeleteSecretReq, uuid uuid.UUID) error {
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

	err = client.CoreV1().Secrets(req.Namespace).Delete(context.TODO(), req.SecretName, metav1.DeleteOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())

		return err
	}
	return nil
}
func (k *K8sSecretService) CreateSecret(req secret.CreateSecretReq, uuid uuid.UUID) (*v1.Secret, error) {
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
	var SecretIns *v1.Secret
	tmp := string(data)
	json.Unmarshal([]byte(tmp), &SecretIns)
	ins, err := client.CoreV1().Secrets(req.Namespace).Create(context.TODO(), SecretIns, metav1.CreateOptions{})
	if err != nil {
		global.DYCLOUD_LOG.Error("删除失败:" + err.Error())
		return nil, err
	}
	return ins, nil
}
