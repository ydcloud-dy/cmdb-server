package cmdb

import (
	"DYCLOUD/global"
	"DYCLOUD/model/cmdb"
	cmdbReq "DYCLOUD/model/cmdb/request"
	utils "DYCLOUD/utils/cmdb"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

type CmdbHostsService struct{}

// SSHTestCmdbHosts
//
//	@Description:  测试本地是否可以免密连接远程主机
//	@receiver cmdbHostsService
//	@param cmdbHosts
//	@return err
func (cmdbHostsService *CmdbHostsService) SSHTestCmdbHosts(req *cmdb.CmdbHosts) (err error) {
	port := strconv.Itoa(*req.Port)
	privateKeyPath, err := utils.GetDefaultPrivateKeyPath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		//return err
	}
	// 检查私钥文件是否存在
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		if err := utils.GenerateSSHKey(privateKeyPath); err != nil {
			fmt.Printf("Error generating SSH key: %v\n", err)
		}
	}
	canLogin, err := utils.CanSSHWithoutPassword(req.ServerHost, port, req.Username, privateKeyPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		//return err
	}
	if canLogin {
		fmt.Println("Can login without password.")
		// 创建 SSH 客户端
		client, err := utils.CreateSSHClient(req.ServerHost, port, req.Username, privateKeyPath)
		if err != nil {
			return fmt.Errorf("创建 SSH 客户端失败: %v", err)
		}
		defer client.Close()
		// 获取主机信息
		if err := utils.PopulateHostInfo(client, req); err != nil {
			return fmt.Errorf("获取主机信息失败: %v", err)
		}
		req.Status = "已验证"
		// 存储到数据库
		if err := global.DYCLOUD_DB.Create(req).Error; err != nil {
			return fmt.Errorf("存储到数据库失败: %v", err)
		}
		return nil
	}
	fmt.Println("Cannot login without password.")
	if strings.Contains(err.Error(), "ssh: handshake failed: ssh: unable to authenticate, attempted methods [none publickey], no supported methods remain") {
		return fmt.Errorf("auth failed")
	}
	return err

}

// CreateCmdbHosts
//
//	@Description: 创建CMDB主机
//	@receiver cmdbHostsService
//	@param cmdbHosts
//	@return err
func (cmdbHostsService *CmdbHostsService) CreateCmdbHosts(req *cmdb.CmdbHosts) (err error) {
	port := strconv.Itoa(*req.Port)

	// 验证 SSH 连接
	client, err := utils.ValidateSSHConnection(req.ServerHost, port, req.Username, req.Password)
	if err != nil {
		return fmt.Errorf("SSH 验证失败: %v", err)
	}
	defer client.Close()

	// 获取主机信息
	if err := utils.PopulateHostInfo(client, req); err != nil {
		return fmt.Errorf("获取主机信息失败: %v", err)
	}
	// 连接成功后，设置免密登录
	if err := utils.EnablePasswordlessSSH(client, req.Username); err != nil {
		return fmt.Errorf("设置免密登录失败: %v", err)
	}
	req.Status = "已验证"

	// 连接成功，存储到数据库
	if err := global.DYCLOUD_DB.Create(req).Error; err != nil {
		return fmt.Errorf("存储到数据库失败: %v", err)
	}
	return nil
}

// DeleteCmdbHosts
//
//	@Description: 删除主机
//	@receiver cmdbHostsService
//	@param ID
//	@param userID
//	@return err
func (cmdbHostsService *CmdbHostsService) DeleteCmdbHosts(ID string, userID uint) (err error) {
	err = global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&cmdb.CmdbHosts{}).Where("id = ?", ID).Update("deleted_by", userID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&cmdb.CmdbHosts{}, "id = ?", ID).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteCmdbHostsByIds
//
//	@Description: 批量删除cmdbHosts表记录
//	@receiver cmdbHostsService
//	@param ID
//	@param userID
//	@return err
func (cmdbHostsService *CmdbHostsService) DeleteCmdbHostsByIds(IDs []string, deleted_by uint) (err error) {
	err = global.DYCLOUD_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&cmdb.CmdbHosts{}).Where("id in ?", IDs).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", IDs).Delete(&cmdb.CmdbHosts{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateCmdbHosts
//
//	@Description: 更新主机信息
//	@receiver cmdbHostsService
//	@param cmdbHosts
//	@return err
func (cmdbHostsService *CmdbHostsService) UpdateCmdbHosts(cmdbHosts cmdb.CmdbHosts) (err error) {
	err = global.DYCLOUD_DB.Model(&cmdb.CmdbHosts{}).Where("id = ?", cmdbHosts.ID).Updates(&cmdbHosts).Error
	return err
}

// GetCmdbHosts
//
//	@Description: 根据ID获取cmdbHosts表记录
//	@receiver cmdbHostsService
//	@param ID
//	@return cmdbHosts
//	@return err
func (cmdbHostsService *CmdbHostsService) GetCmdbHosts(ID string) (cmdbHosts cmdb.CmdbHosts, err error) {
	err = global.DYCLOUD_DB.Where("id = ?", ID).First(&cmdbHosts).Error
	return
}
func (cmdbHostsService *CmdbHostsService) ImportHosts(filePath string, projectId int) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %v", err)
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("无法读取行: %v", err)
	}
	for _, row := range rows[1:] { // 忽略标题行
		name := row[0]
		serverHost := row[1]
		port, _ := strconv.Atoi(row[2])
		username := row[3]
		password := row[4]
		note := row[5]

		host := &cmdb.CmdbHosts{
			Name:       name,
			ServerHost: serverHost,
			Port:       &port,
			Username:   username,
			Password:   password,
			Note:       note,
			Project:    projectId,
		}

		if err := cmdbHostsService.CreateCmdbHosts(host); err != nil {
			fmt.Printf("创建主机 %s 失败: %v\n", name, err)
		} else {
			fmt.Printf("主机 %s 创建成功\n", name)
		}
	}

	return nil
}

// GetCmdbHostsInfoList
//
//	@Description: 分页获取cmdbHosts表记录
//	@receiver cmdbHostsService
//	@param info
//	@return list
//	@return total
//	@return err
func (cmdbHostsService *CmdbHostsService) GetCmdbHostsInfoList(info cmdbReq.CmdbHostsSearch) (list []cmdb.CmdbHosts, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.DYCLOUD_DB.Model(&cmdb.CmdbHosts{})
	var cmdbHostss []cmdb.CmdbHosts
	// 如果有条件搜索 下方会自动创建搜索语句
	//if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
	//db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt).Where("project = ?",info.Project)
	db = db.Where("project = ?", info.Project)
	//}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&cmdbHostss).Error
	return cmdbHostss, total, err
}
func (cmdbHostsService *CmdbHostsService) GetCmdbHostsPublic() {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
