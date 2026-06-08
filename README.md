# mini-issue Part 1

`mini-issue` 是一个用 Go 写的小型 issue 系统。当前 Part 1 先完成用户模块的基础能力：

- 用户注册
- 用户登录
- JWT 鉴权
- 获取当前登录用户信息
- 分层结构：Router / Controller / Service / DAO / Model
- 基础单元测试和路由测试

## 技术栈

- Go
- Gin：HTTP Web 框架
- MySQL：数据存储
- sqlx：数据库访问
- bcrypt：密码哈希
- golang-jwt：JWT 生成和解析
- go-sqlmock：测试中 mock 数据库

## 项目结构

```text
mini-issue
├── cmd
│   └── main.go                  # 程序入口，组装 DB、DAO、Service、Controller、Router
├── internal
│   ├── controller
│   │   └── user_controller.go    # HTTP 请求/响应处理
│   ├── dao
│   │   ├── user_dao.go           # 用户表数据库操作
│   │   └── user_dao_test.go
│   ├── middleware
│   │   ├── auth.go               # JWT 鉴权中间件
│   │   └── auth_test.go
│   ├── model
│   │   └── user.go               # 数据结构、请求/响应结构
│   ├── router
│   │   ├── router.go             # 路由注册
│   │   └── router_test.go
│   └── service
│       ├── user_service.go       # 用户业务逻辑
│       └── user_service_test.go
├── pkg
│   ├── db
│   │   └── mysql.go              # MySQL 连接
│   ├── jwt
│   │   ├── jwt.go                # JWT 生成和解析
│   │   └── jwt_test.go
│   └── response
│       └── response.go           # 统一 JSON 响应
├── go.mod
└── go.sum
```

## 请求流程

一次请求进入项目后，大致按这个顺序流动：

```text
HTTP Request
  ↓
Router
  ↓
Controller
  ↓
Service
  ↓
DAO
  ↓
MySQL
```

例如登录：

```text
POST /login
  ↓
router.NewRouter 负责把 /login 分配给 UserController.Login
  ↓
UserController.Login 解析 JSON 请求体
  ↓
UserService.Login 校验用户名、密码，并生成 JWT
  ↓
UserDAO.GetByUsername 查询数据库里的用户
  ↓
返回统一 JSON 响应
```

## 分层说明

这几个类型看起来都在处理 user，但它们关注的问题不同。

### UserController

负责 HTTP 层：

- 从 `gin.Context` 读取请求体、请求头、URL 参数
- 调用 Service
- 决定 HTTP 状态码
- 返回 JSON

典型代码位置：[internal/controller/user_controller.go](internal/controller/user_controller.go)

### UserService

负责业务规则：

- 注册时检查用户名是否为空
- 检查密码长度
- 检查用户是否已存在
- 登录时校验密码
- 登录成功后生成 token

典型代码位置：[internal/service/user_service.go](internal/service/user_service.go)

### UserDAO

负责数据库访问：

- 写 SQL
- 执行查询或插入
- 把数据库结果转换成 model

典型代码位置：[internal/dao/user_dao.go](internal/dao/user_dao.go)

简单记法：

```text
Controller：处理 HTTP
Service：处理业务
DAO：处理数据库
```

## API

### 注册

```http
POST /register
Content-Type: application/json
```

请求体：

```json
{
  "username": "alice",
  "password": "123456"
}
```

成功响应：

```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

### 登录

```http
POST /login
Content-Type: application/json
```

请求体：

```json
{
  "username": "alice",
  "password": "123456"
}
```

成功响应：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "<jwt-token>"
  }
}
```

### 获取当前用户

```http
GET /me
Authorization: Bearer <jwt-token>
```

成功响应：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "alice",
    "created_at": "2026-06-08T12:00:00Z"
  }
}
```

## 数据库

当前代码默认连接：

```text
root:@tcp(127.0.0.1:3306)/test?parseTime=true
```

可以先创建数据库和用户表：

```sql
create database if not exists test;

use test;

create table if not exists `user` (
  id bigint primary key auto_increment,
  username varchar(64) not null unique,
  password varchar(255) not null,
  created_at datetime not null
);
```

## 运行项目

启动 MySQL 后，在项目根目录运行：

```bash
go run ./cmd
```

服务默认监听：

```text
http://localhost:8080
```

## 测试

运行全部测试：

```bash
go test ./...
```

当前测试覆盖：

- JWT 生成和解析
- 鉴权中间件
- UserDAO 数据库访问逻辑
- UserService 业务逻辑
- Router + Controller 的 HTTP 流程

DAO 和路由测试使用 `go-sqlmock`，不会连接真实 MySQL。

## Part 1 当前边界

Part 1 目前只做用户认证相关能力，还没有真正的 issue 功能。

后续 Part 可以继续扩展：

- issue 创建
- issue 列表
- issue 详情
- issue 状态流转
- 用户和 issue 的权限关系
- 更完善的错误码和配置管理
