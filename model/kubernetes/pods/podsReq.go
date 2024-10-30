package pods

import (
	"DYCLOUD/model/common/request"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"io"
	"time"
)

type PodRecord struct {
	ID            uint           `json:"id" gorm:"not null;unique;primary_key"`
	Cluster       string         `json:"cluster" gorm:"comment:集群名称"`
	Namespace     string         `json:"namespace" gorm:"comment:命名空间"`
	PodName       string         `json:"pod_name" gorm:"comment:Pod名称"`
	ContainerName string         `json:"container_name" gorm:"comment:容器名称"`
	UUID          uuid.UUID      `json:"uuid" gorm:"comment:UUID"`
	Username      string         `json:"userName" gorm:"comment:操作用户"`
	NickName      string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`
	Records       []byte         `json:"records" gorm:"type:longblob;comment:'操作记录(二进制存储)';size:128"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
type PodsFilesRequest struct {
	ClusterId     int       `json:"cluster_id" form:"cluster_id" validate:"required"`
	Folder        string    `json:"folder" form:"folder"`
	PodName       string    `json:"podName" form:"podName" validate:"required"`
	ContainerName string    `json:"containerName" form:"containerName"`
	Namespace     string    `json:"namespace" form:"namespace" validate:"required"`
	Path          string    `json:"path" form:"path"`
	OldPath       string    `json:"oldPath" form:"oldPath"`
	Commands      []string  `json:"-" form:"Commands"`
	Stdin         io.Reader `json:"-" form:"Stdin"`
	Content       string    `json:"content" form:"content"`
	FilePath      string    `json:"filePath" form:"filePath"`
	XToken        string    `json:"x-token" form:"x-token"`
}

func (PodRecord) TableName() string {
	return "k8s_pod_records"
}

type PodListReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	LabelSelector string `json:"labelSelector" form:"labelSelector"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
	request.PageInfo
}
type PodMetricsReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
}

type DescribePodInfo struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	PodName   string `json:"podName" form:"podName"`
}

type CreatePodReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	Content   interface{} `json:"content" form:"content"`
}

type DeletePodReq struct {
	ClusterId int    `json:"cluster_id" form:"cluster_id"`
	Namespace string `json:"namespace" form:"namespace"`
	PodName   string `json:"podName" form:"podName"`
}
type UpdatePodReq struct {
	ClusterId int         `json:"cluster_id" form:"cluster_id"`
	Namespace string      `json:"namespace" form:"namespace"`
	PodName   string      `json:"podName" form:"podName"`
	Content   interface{} `json:"content" form:"content"`
}
type PodEventsReq struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	PodName       string `json:"podName" form:"podName"`
	FieldSelector string `json:"fieldSelector" form:"fieldSelector"`
}
type ListPodFiles struct {
	ClusterId     int    `json:"cluster_id" form:"cluster_id"`
	Namespace     string `json:"namespace" form:"namespace"`
	PodName       string `json:"podName" form:"podName"`
	ContainerName string `json:"containerName" form:"containerName"`
	Path          string `json:"path" form:"path"`
}
