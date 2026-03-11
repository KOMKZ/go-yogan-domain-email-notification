# go-yogan-domain-email-notification

邮件通知领域包，负责模板管理、发送日志、触发点定义与邮件发送编排。

## 目录结构

- `model/`: 领域模型
- `errors/`: 错误定义
- `repository/`: 仓储接口与 MySQL 实现
- `service/`: 领域服务
- `sender/`: 发送器实现
- `provider/do/`: 依赖注入 Provider
- `migrations/`: 数据库迁移脚本

## 关键能力

- 邮件模板 CRUD 与渲染
- 触发点注册与参数定义
- 发送日志记录与查询
- 组件发送器适配（`go-yogan-component-email`）

## 开发命令

```bash
go test ./...
```
