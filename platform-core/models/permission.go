package models

type MenuItems struct {
	Id               string `json:"id" xorm:"id"`                               // 唯一标识
	ParentCode       string `json:"parentCode" xorm:"parent_code"`              // 所属菜单栏
	Code             string `json:"code" xorm:"code"`                           // 编码
	Source           string `json:"source" xorm:"source"`                       // 来源
	Description      string `json:"description" xorm:"description"`             // 描述
	LocalDisplayName string `json:"localDisplayName" xorm:"local_display_name"` // 显示名
	MenuOrder        int    `json:"menuOrder" xorm:"menu_order"`                // 菜单排序
}

type RoleMenu struct {
	Id       string `json:"id" xorm:"id"`              // 唯一标识
	RoleName string `json:"roleName" xorm:"role_name"` // 角色
	MenuCode string `json:"menuCode" xorm:"menu_code"` // 菜单编码
}