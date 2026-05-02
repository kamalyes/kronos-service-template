# <% .ObjectPrefix %> <% .ServiceName %> 微服务

> 基于 Go gRPC 的 <% .ServiceName %> 微服务，使用 go-rpc-gateway 框架提供 gRPC + HTTP 双协议支持

## ✨ 功能特性

- **gRPC & HTTP 双协议**：同时支持 gRPC 高性能通信和 RESTful HTTP/JSON 接口
- **自动化数据库迁移**：基于 go-sqlbuilder Migrator，支持自动建表、索引和表注释
- **国际化错误处理**：内置 i18n 错误系统，支持 27 种语言的错误消息
- **仓库工厂模式**：统一管理数据访问层，支持泛型 BaseRepository
- **PBMO 模型转换**：Proto ↔ Model 双向自动转换
- **Snowflake ID 生成**：内置分布式短 ID 生成器
- **认证拦截器**：可选的 gRPC 认证拦截器
- **事件订阅**：可选的分布式事件订阅支持

## 📂 项目结构

```bash
.
├── .vscode/                    # VS Code 配置
├── bootstrap/                  # 启动引导
│   ├── database.go             # 数据库初始化与迁移
│   ├── server.go               # gRPC/HTTP 服务启动与生命周期
│   ├── services.go             # 服务注册（gRPC + HTTP Handler）
│   └── subscriber.go           # 事件订阅（可选）
├── constants/                  # 常量定义
│   └── service_name.go         # 服务名称常量
├── docs/                       # 文档目录
├── __ai_home__/                # AI 规范，渲染时改名成 .trae / .claude / .cursor / .vscode
│   ├── codes.go                # 业务错误码常量
│   ├── errors.go               # i18n 错误构造函数
│   └── grpc.go                 # gRPC 错误码映射
├── locales/                    # 国际化文件（27 种语言）
│   ├── zh.json                 # 中文
│   ├── en.json                 # 英文
│   └── ...                     # 其他语言
├── middleware/                  # 中间件
│   ├── auth_interceptor.go     # 认证拦截器（可选）
│   └── public_methods.go       # 公开方法白名单
├── migrations/                  # 数据库迁移
│   ├── auto_migrate.go         # 自动迁移（模型、索引、表注释）
│   └── migrate.go              # 种子数据
├── models/                     # 数据模型
│   ├── base.go                 # 基础类型（BuildPagination, PBConvertible, LocalizedTextList）
│   ├── table_names.go          # 表名常量
│   └── __example_name___model.go  # 示例模型（可选）
├── repository/                  # 数据访问层
│   ├── factory.go              # RepositoryFactory 单例
│   └── __example_name___repository.go  # 示例仓库（可选）
├── resources/                   # 配置文件
│   └── __project_name__-dev.yaml     # 唯一环境配置；新增环境见下方说明
├── service/                     # 业务逻辑层
│   ├── idgen.go                # Snowflake ID 生成器
│   ├── __service_name__.go     # 主服务实现
│   └── __example_name___service.go  # 示例服务（可选）
├── .gitignore
├── .golangci.yml
├── go.mod
└── main.go
```

## 🚀 快速开始

### 1. 环境要求

- **Go**: 版本 `1.22+` 或更高

### 2. 使用 kopy 生成项目

#### 安装kopy
```bash
go install github.com/kamalyes/kopy@latest
```

#### 交互模式

```bash
kopy new git@github.com:kamalyes/kronos-service-template.git -o ./<% .ProjectName %>
```

交互式问答示例：

```
❓ 组织名称（如 kamalyes，用于 Go 模块路径 github.com/{org}/{project}） [kamalyes]:
❓ 组织前缀（如 kronos，用于项目名 {prefix}-{type}-service） [kronos]:
❓ 服务类型标识（如 open/payment/core，用于项目名 {prefix}-{type}-service） [my]: payment
❓ 服务名称（PascalCase，如 PaymentService） [PaymentService]:
❓ 模块名称（小写，如 payment-service） [payment-service]:
❓ 项目名称（如 kronos-payment-service） [kronos-payment-service]:
❓ 作者名称 [kamalyes]:
❓ 作者邮箱 [example@qq.com]:
❓ Go 版本 [1.25]:
❓ 网关配置前缀（如 gateway-payment） [gateway-payment]:
❓ gRPC 端口 [9190]:
❓ HTTP 端口 [8190]:
❓ 健康检查端口 [8191]:
❓ 数据库名称 [payment_service]:
❓ 缓存键前缀 [payment:]:
❓ 监控指标命名空间 [payment]:
❓ 是否使用 Redis [true]:
❓ 是否启用认证拦截器 [true]:
❓ 是否启用事件订阅 [false]:
❓ 是否生成示例 Model/Repository/Service 文件 [true]:
❓ 示例实体名称（PascalCase，如 Tenant） [Demo]: Order
❓ 共享Proto模块路径 [github.com/kamalyes/kronos-share-proto]:
```

#### 非交互模式

使用 `--var` 参数跳过交互式问答：

```bash
kopy new git@github.com:kamalyes/kronos-service-template.git -o ./<% .ProjectName %> \
  --var OrgName=kamalyes \
  --var ObjectPrefix=kronos \
  --var ServiceType=payment \
  --var ServiceName=PaymentService \
  --var ModuleName=payment-service \
  --var ProjectName=kronos-payment-service \
  --var AuthorName=kamalyes \
  --var AuthorEmail=dev@example.com \
  --var GoVersion=1.25 \
  --var GRPCPort=9190 \
  --var HTTPPort=8190 \
  --var HealthPort=8191 \
  --var DBName=payment_service \
  --var CachePrefix=payment: \
  --var MetricsNamespace=payment \
  --var UseRedis=true \
  --var UseAuthInterceptor=true \
  --var UseSubscriber=false \
  --var GenerateExample=true \
  --var ExampleName=Order \
  --var ShareProtoModule=github.com/kamalyes/kronos-share-proto \
  --overwrite
```

#### 预览模式

使用 `--dry-run` 仅预览生成结果，不写入文件：

```bash
kopy new ./kronos-service-template --dry-run --var OrgName=kamalyes --var ObjectPrefix=kronos --var ServiceType=payment
```

### 3. 生成后步骤

```bash
cd <% .ProjectName %>

# 1. 安装依赖
go mod tidy

# 2. 确保 share-proto 项目在同级目录
# go.mod 中的 replace 指令默认指向 ../kronos-share-proto

# 3. 运行服务
go run main.go
```

## 📋 模板变量

| 变量 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `OrgName` | string | `kamalyes` | 组织名称，用于 Go 模块路径 `github.com/{org}/{project}` |
| `ObjectPrefix` | string | `kronos` | 组织前缀，用于项目名 `{prefix}-{type}-service` |
| `ServiceType` | string | `my` | 服务类型标识，如 open/payment/core |
| `ServiceName` | string | `<% .ServiceType \| pascal %>Service` | 服务名称（PascalCase） |
| `ModuleName` | string | `<% .ServiceType %>-service` | 模块名称（小写） |
| `ProjectName` | string | `<% .ObjectPrefix %>-<% .ServiceType %>-service` | 项目名称 |
| `AuthorName` | string | `kamalyes` | 作者名称 |
| `AuthorEmail` | string | `example@qq.com` | 作者邮箱 |
| `GoVersion` | select | `1.25` | Go 版本（1.22/1.23/1.24/1.25） |
| `GRPCPort` | int | `9190` | gRPC 端口 |
| `HTTPPort` | int | `8190` | HTTP 端口 |
| `HealthPort` | int | `8191` | 健康检查端口 |
| `DemoRemoteService` | string | `<% .SecondModule %>-service` | 服务间调用 demo 的下游服务名 |
| `DBName` | string | `<% .ServiceType \| snake %>_service` | 数据库名称 |
| `CachePrefix` | string | `<% .ServiceType \| snake %>:` | 缓存键前缀 |
| `MetricsNamespace` | string | `<% .ServiceType \| snake %>` | 监控指标命名空间 |
| `UseRedis` | bool | `true` | 是否使用 Redis |
| `UseAuthInterceptor` | bool | `true` | 是否启用认证拦截器 |
| `UseSubscriber` | bool | `false` | 是否启用事件订阅 |
| `GenerateExample` | bool | `true` | 是否生成示例 Model/Repository/Service |
| `ExampleName` | string | `Demo` | 示例实体名称（PascalCase） |
| `ShareProtoModule` | string | `github.com/<% .OrgName %>/<% .ObjectPrefix %>-share-proto` | 共享 Proto 模块路径 |

## 🔄 文件重命名规则

| 模式 | 替换为 | 说明 |
|------|--------|------|
| `__project_name__` | `<% .ProjectName %>` | 项目名称（也是配置文件前缀）|
| `__service_name__` | `<% .ServiceName \| snake %>` | 服务名称（蛇形）|
| `__module_name__` | `<% .ModuleName %>` | 模块名称 |
| `__example_module__` | `<% .ExampleModule %>` | 示例模块目录（service/ 下分包）|
| `__second_module__` | `<% .SecondModule %>` | 第二模块目录 |
| `__mixed_module__` | `<% .MixedModule %>` | 混合模块目录 |
| `__example_service__` | `<% .ExampleService \| snake %>` | 示例实体（蛇形）|
| `__second_service__` | `<% .SecondService \| snake %>` | 第二实体（蛇形）|
| `__mixed_service__` | `<% .MixedService \| snake %>` | 混合模块聚合根（蛇形）|
| `__mixed_member__` | `<% .MixedService \| snake %>_member` | 混合模块成员实体 |
| `__ai_home__` | `<% .AiAssistant %>` | AI 规范目录（.trae / .claude / .cursor / .vscode）|

## 🔤 命名转换函数

模板中的 `default` 字段和文件内容支持以下命名转换函数：

| 函数 | 说明 | 示例 |
|------|------|------|
| `snake` | 蛇形命名 | `PaymentService` → `payment_service` |
| `camel` | 驼峰命名 | `payment_service` → `paymentService` |
| `pascal` | 帕斯卡命名 | `payment_service` → `PaymentService` |
| `kebab` | 短横线命名 | `PaymentService` → `payment-service` |
| `lower` | 转小写 | `Hello` → `hello` |
| `upper` | 转大写 | `hello` → `HELLO` |

## 📝 添加新业务实体

当需要添加新的业务实体（如 Order、Product 等）时，需要修改以下文件：

### 1. 创建模型 (`models/xxx_model.go`)

```go
type XxxModel struct {
    XxxId     string    `gorm:"column:xxx_id;type:varchar(36);primaryKey" json:"xxx_id"`
    Name      string    `gorm:"column:name;type:varchar(128);not null" json:"name"`
    Status    int32     `gorm:"column:status;type:tinyint;default:1" json:"status"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (XxxModel) TableName() string { return "xxxs" }
func (XxxModel) TableComment() string { return "XXX表" }
```

### 2. 添加表名常量 (`models/table_names.go`)

```go
var TableNameXxx = XxxModel{}.TableName()
```

### 3. 创建仓库 (`repository/xxx_repository.go`)

继承 `sqlrepo.BaseRepository[T]`，实现 `GetByID`、`ListWithPagination`、`UpdateByID`、`DeleteByID` 等方法。

### 4. 注册仓库 (`repository/factory.go`)

在 `RepositoryFactory` 中添加字段并在 `NewRepositoryFactory()` 中初始化。

### 5. 创建服务 (`service/xxx_service.go`)

封装业务逻辑，调用 Repository 方法。

### 6. 注册服务 (`service/__service_name___service.go`)

在主 ServiceImpl 中添加子服务字段并初始化。

### 7. 数据库迁移 (`migrations/auto_migrate.go`)

在 `GetAllModels()`、`GetAllIndexes()`、`GetAllComments()` 中添加新实体。

### 8. 错误码 (`errors/codes.go`)

添加实体相关的业务错误码。

### 9. 国际化 (`locales/*.json`)

在所有语言文件中添加错误码的翻译。

### 10. gRPC/HTTP 注册 (`bootstrap/services.go`)

注册 gRPC 服务和 HTTP Handler。

## ⚠️ 重要注意事项

- 模型字段名必须与 Proto 字段名一致，以支持 PBMO 自动转换
- 使用 `errors.ToGRPCError()` 返回错误，不要直接使用 `status.Errorf()`
- 使用 `service.NewSnowflakeID()` 生成实体 ID
- i18n 多语言字段使用 `models.LocalizedTextList` 类型
- `go.mod` 中的 `replace` 指令默认指向同级目录的 share-proto 项目
