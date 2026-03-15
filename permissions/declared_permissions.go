package permissions

type DeclaredPermission struct {
	PermissionCode string
	PermissionName string
	PermissionType string
	ResourceCode   string
	GroupCode      string
	Description    string
}

func DeclaredPermissions() []DeclaredPermission {
	return []DeclaredPermission{
		{
			PermissionCode: "email_notification:read",
			PermissionName: "查看邮件通知",
			PermissionType: "READ",
			ResourceCode:   "email_notification",
			GroupCode:      "SYSTEM",
			Description:    "邮件模板与发送记录查看",
		},
		{
			PermissionCode: "email_notification:write",
			PermissionName: "管理邮件通知",
			PermissionType: "WRITE",
			ResourceCode:   "email_notification",
			GroupCode:      "SYSTEM",
			Description:    "邮件模板维护与触发配置",
		},
	}
}
