package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library-study/app/model"
	"net/http"
)

// GetRoles 获取角色列表
func GetRoles(c *gin.Context) {
	var roles []model.Roles
	if err := model.MySQLDB.Table("roles").Find(&roles).Error; err != nil {
		fmt.Printf("Retrieve roles error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, roles)
}

// GetPermissions 获取权限列表
func GetPermissions(c *gin.Context) {
	var permissions []model.Permissions
	if err := model.MySQLDB.Table("permissions").Find(&permissions).Error; err != nil {
		fmt.Printf("Retrieve permissions error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, permissions)
}

// GetRolePre 获取角色权限列表
func GetRolePre(c *gin.Context) {
	var RolePre []model.RolePermissions
	if err := model.MySQLDB.Table("role_permissions").Find(&RolePre).Error; err != nil {
		fmt.Printf("Retrieve role_permissions error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, RolePre)
}

// GetUerRoles 获取用户角色列表
func GetUerRoles(c *gin.Context) {
	var UerRoles []model.UserRoles
	if err := model.MySQLDB.Table("user_roles").Find(&UerRoles).Error; err != nil {
		fmt.Printf("Retrieve user_roles error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// 将处理后的数据返回到前端
	c.JSON(http.StatusOK, UerRoles)
}

// UpdateRolePre 修改用户角色
func UpdateRolePre(c *gin.Context) {
	// 创建一个结构体用于接收请求体中的数据
	var updateInfo struct {
		UserId int `json:"userid"`
		RoleId int `json:"roleid"`
	}
	fmt.Println("c", c)
	// 从请求中解析数据到 updateInfo 对象中
	if err := c.BindJSON(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 在数据库中更新用户的角色
	if err := model.MySQLDB.Table("user_roles").
		Where("user_id = ?", updateInfo.UserId).
		Update("role_id", updateInfo.RoleId).Error; err != nil {
		fmt.Printf("Update user_roles error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user role"})
		return
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// AddRolePermissions 添加角色权限
func AddRolePermissions(c *gin.Context) {
	// 创建一个结构体用于接收请求体中的数据
	var addInfo struct {
		RoleId       int `json:"roleid"`
		PermissionId int `json:"permissionid"`
	}

	// 从请求中解析数据到 addInfo 对象中
	if err := c.BindJSON(&addInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println("add----", addInfo)

	// 在尝试插入之前检查是否存在相应的角色权限记录
	var count int64
	model.MySQLDB.Model(&model.RolePermissions{}).Where("role_id = ? AND permission_id = ?", addInfo.RoleId, addInfo.PermissionId).Count(&count)

	// 如果count不为0，说明数据库中已存在该角色权限记录
	if count != 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Role permission already exists"})
		return
	}

	// 插入新角色权限记录
	rolePermission := model.RolePermissions{
		RoleId:       int64(addInfo.RoleId),
		PermissionId: int64(addInfo.PermissionId),
	}
	if err := model.MySQLDB.Create(&rolePermission).Error; err != nil {
		fmt.Printf("Add role_permissions error: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding role permission"})
		return
	}

	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"message": "Role permission added successfully"})
}
