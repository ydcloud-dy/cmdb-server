// 自动生成模板K8sVeleroTasks
package velero

import (
	"DYCLOUD/global"
)

// k8sVeleroTasks表 结构体  K8sVeleroTasks
type K8sVeleroTasks struct {
	global.DYCLOUD_MODEL
	ClusterId       *int   `json:"clusterId" form:"clusterId" gorm:"column:cluster_id;comment:集群id;size:10;" binding:"required"`                               //集群id
	ExcludeNs       string `json:"excludeNs" form:"excludeNs" gorm:"column:excludeNs;comment:指定不需要备份的命名空间;size:255;"`                                          //指定不需要备份的命名空间
	ExcludeResource string `json:"excludeResource" form:"excludeResource" gorm:"column:excludeResource;comment:指定不需要备份的资源类型;size:255;"`                        //指定不需要备份的资源类型
	Name            string `json:"name" form:"name" gorm:"column:name;comment:备份名称;size:255;" binding:"required"`                                              //备份名称
	NamespaceType   *int   `json:"namespaceType" form:"namespaceType" gorm:"column:namespace_type;comment:是否备份所有命名空间，0为部分，1为备份所有;size:10;" binding:"required"` //是否备份所有命名空间，0为部分，1为备份所有
	Remark          string `json:"remark" form:"remark" gorm:"column:remark;comment:备注;"`                                                                      //备注
	Resource        string `json:"resource" form:"resource" gorm:"column:resource;comment:备份资源;size:255;"`                                                     //备份资源
	ResourceType    *int   `json:"resourceType" form:"resourceType" gorm:"column:resource_type;comment:是否备份所有资源类型，0为部分，1为备份所有;size:10;"`                       //是否备份所有资源类型，0为部分，1为备份所有
	RetentionTime   string `json:"retentionTime" form:"retentionTime" gorm:"column:retentionTime;comment:备份保留时长;size:255;"`                                    //备份保留时长
	Strategy        string `json:"strategy" form:"strategy" gorm:"column:strategy;comment:备份策略;size:255;" binding:"required"`                                  //备份策略
	CreatedBy       uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy       uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy       uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName k8sVeleroTasks表 K8sVeleroTasks自定义表名 k8s_velero_tasks
func (K8sVeleroTasks) TableName() string {
	return "k8s_velero_tasks"
}
