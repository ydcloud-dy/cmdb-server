## GVA 组织管理功能
#### 开发者：奇淼
### 本插件支持自动安装

### 图例

![img](https://qmplusimg.henrongyi.top/plugin%2Forganization1.jpg)

![img](https://qmplusimg.henrongyi.top/plugin%2Forganization2.jpg)

![img](https://qmplusimg.henrongyi.top/plugin%2Forganization3.jpg)

![img](https://qmplusimg.henrongyi.top/plugin%2Forganization4.jpg)

![img](https://qmplusimg.henrongyi.top/plugin%2Forganization5.jpg)

### 手动安装方法

    1.解压zip获得organization文件夹
    2.将 organization/web/plugin/organization 放置在web/plugin下
    3.将 organization/server/plugin/organization 放置在server/plugin下

#### 执行如下注册方法

### 注册（手动自动都需要）

#### 1. 前往GVA主程序下的initialize/router.go 在Routers 方法最末尾按照你需要的及安全模式添加本插件
    PluginInit(PublicGroup, organization.CreateOrganizationPlug())
    到gva系统，角色管理，分配角色的菜单权限和api权限即可，插件会自动注册菜单和api，需要手动分配。
### 2. 配置说明

#### 2-1 全局配置结构体说明

    无配置

#### 2-2 结构说明
    
    // 组织基础结构    
    type Organization struct {
	global.DYCLOUD_MODEL
	Name     string `json:"name" form:"name" gorm:"column:name;comment:;"`
	ParentID uint   `json:"parentID" form:"parentID" gorm:"column:parent_id;comment:父节点ID;"`
    }   

    // 组织用户关系表
    type OrgUser struct {
    Organization   Organization   `json:"organization"`
    OrganizationID uint           `json:"organizationID,omitempty" form:"organizationID" `
    SysUserID      uint           `json:"sysUserID,omitempty" form:"sysUserID"`
    IsAdmin        bool           `json:"isAdmin" form:"isAdmin"`
    SysUser        system.SysUser `json:"sysUser"`
    }
    
    // 组织内用户操作结构
    type OrgUserReq struct {
    OrganizationID   uint   `json:"organizationID,omitempty"`
    ToOrganizationID uint   `json:"toOrganizationID,omitempty"`
    SysUserIDS       []uint `json:"sysUserIDS,omitempty"`
    }


    // 搜索组织功能（未使用）
    type OrganizationSearch struct {
    organization.Organization
    request.PageInfo
    }
    
    // 搜索组织内用户功能
    type OrgUserSearch struct {
    organization.OrgUser
    UserName string `json:"userName" form:"userName"`
    request.PageInfo
    }



### 3. 方法API
    
    无，后续维护会增加便捷查询当前用户所属组织，所属组织子组织以及成员等便捷方法，根据反馈添加。

### 4. 可直接调用的接口

    POST  /org/createOrganization // 新建Organization
    入参为 Organization 结构体

    DELETE  /org/deleteOrganization // 删除Organization
    入参为 ID int


    DELETE  /org/deleteOrganizationByIds // 批量删除Organization
    入参为 IDS:[]int


    PUT  /org/updateOrganization // 更新Organization
    入参为 Organization 结构体


    POST  /org/createOrgUser // 人员入职
    入参为 OrgUserReq 结构体


    PUT  /org/setOrgUserAdmin // 管理员设置
    入参为 {
        sysUserID: 用户id,
        isAdmin: 是否管理员 bool
    }


    GET  /org/findOrganization // 根据ID获取Organization
    入参为 { ID: 组织id }


    GET  /org/getOrganizationList // 获取Organization列表
    无需入参


    GET  /org/findOrgUserAll // 获取当前组织下所有用户ID
    入参为 { organizationID: 组织id }


    GET  /org/findOrgUserList // 获取当前组织下所有用户（分页）
    入参为 OrgUserSearch 结构体


    DELETE  /org/deleteOrgUser // 删除当前组织下选中用户
    入参为 { sysUserIDS:用户id []int, organizationID: 当前组织id }


    PUT  /org/transferOrgUser // 用户转移组织
       入参为  sysUserIDS: 需要操作的用户的ids []int,
        organizationID: 原始组织ID  int,
        toOrganizationID: 目标组织ID int,

