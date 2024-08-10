package enum

import (
	"errors"
	"github.com/fengyuan-liang/GoKit/collection/maps"
)

type Permission uint64

const (
	// PermissionUserRead 用户信息读权限
	PermissionUserRead Permission = 1 << 0

	// PermissionUserWrite 用户信息写权限
	PermissionUserWrite Permission = 1 << 1

	// PermissionUserReadWrite 用户信息读写权限
	PermissionUserReadWrite = PermissionUserRead | PermissionUserWrite

	// PermissionAppRead 应用信息读权限
	PermissionAppRead Permission = 1 << 2

	// PermissionAppWrite 应用信息写权限
	PermissionAppWrite Permission = 1 << 3

	// PermissionAppReadWrite 应用信息读写权限
	PermissionAppReadWrite Permission = PermissionAppRead | PermissionAppWrite

	// PermissionAdminRead 管理员信息读权限
	PermissionAdminRead Permission = 1 << 62

	// PermissionAdminWrite 管理员信息写权限
	PermissionAdminWrite Permission = 1 << 63

	// PermissionAdminReadWrite 管理员信息读写权限
	PermissionAdminReadWrite Permission = PermissionAdminRead | PermissionAdminWrite

	// PermissionAll 所有权限
	PermissionAll Permission = 0xFFFFFFFFFFFFFFFF
)

const (
	RoleTsPermission            Permission = PermissionUserRead | PermissionAppRead | PermissionAdminRead
	RoleManagerPermission       Permission = PermissionUserReadWrite | PermissionAppReadWrite
	RoleAdministratorPermission Permission = PermissionAll
)

// RoleType 角色类型
type RoleType string

const (
	// RoleTS 技术支持
	RoleTS RoleType = "ts"

	// RoleManager 管理员
	RoleManager RoleType = "manager"

	// RoleAdministrator 系统管理员
	RoleAdministrator RoleType = "administrator"

	// RoleWxUser 微信用户
	RoleWxUser RoleType = "wx_user"
	// RoleAssistant 助教 打手
	RoleAssistant RoleType = "assistant"
)

var rolePermissionMap = func() maps.IMap[RoleType, Permission] {
	linkedHashMap := maps.NewLinkedHashMap[RoleType, Permission]()
	linkedHashMap.PutAll([]*maps.Pair[RoleType, Permission]{
		{RoleTS, RoleTsPermission},
		{RoleManager, RoleManagerPermission},
		{RoleAdministrator, RoleAdministratorPermission},
	})
	return linkedHashMap
}()

var RoleDisPlayNameMap = func() maps.IMap[RoleType, string] {
	linkedHashMap := maps.NewLinkedHashMap[RoleType, string]()
	linkedHashMap.PutAll([]*maps.Pair[RoleType, string]{
		{RoleTS, "技术支持"},
		{RoleManager, "管理员"},
		{RoleAdministrator, "系统管理员"},
		{RoleWxUser, "微信用户"},
		{RoleAssistant, "助教 打手"},
	})
	return linkedHashMap
}()

func (r RoleType) Permission() Permission {
	if per, ok := rolePermissionMap.Get(r); ok {
		return per
	} else {
		return 0
	}
}

func (r RoleType) DisPlayName() string {
	return RoleDisPlayNameMap.MustGet(r)
}

func (r RoleType) CheckPermission(requiredPermission Permission) error {
	if r.Permission()&requiredPermission != requiredPermission {
		return errors.New("权限不够")
	}
	return nil
}
