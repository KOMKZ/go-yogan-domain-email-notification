# go-yogan-domain-email-notification

> 邮件通知领域模型 - 可配置的邮件模板与发送服务

## 功能

- **Trigger 注册**：代码管理触发点，不入库
- **模板管理**：CRUD 邮件模板
- **多语言支持**：同一 Trigger 支持多语言模板
- **参数体系**：通用参数 + Trigger 专属参数
- **发送日志**：记录每次发送

## 安装

```bash
go get github.com/KOMKZ/go-yogan-domain-email-notification
```

## 快速使用

### 1. 创建注册表

```go
registry := email_notification.NewTriggerRegistry()

// 注入通用参数
registry.SetCommonParams([]email_notification.Param{
    {Name: "AppName", Type: "string", Description: "应用名称"},
})

// 注册触发点
registry.Register("user.registered", "用户注册", "...", []email_notification.Param{
    {Name: "UserName", Type: "string", Required: true},
})
```

### 2. 创建服务

```go
svc := email_notification.NewService(db, emailComp, registry)
```

### 3. 发送邮件

```go
err := svc.Send(ctx, email_notification.SendInput{
    TriggerCode: "user.registered",
    Recipient:   "user@example.com",
    Params: map[string]any{
        "AppName":  "Yogan",
        "UserName": "张三",
    },
})
```

## License

MIT
