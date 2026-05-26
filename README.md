# MyBlog

一个前后端分离的个人博客项目。后端使用 Go + Gin + GORM + MySQL 提供文章、分类、标签、评论、留言、友链和后台管理接口；前端使用 Astro + React + Tailwind CSS 构建公开博客页面和管理后台。

## 功能概览

- 公开页面：文章列表、文章详情、归档时间线、关于页、友链页、留言板
- 内容管理：文章创建/编辑/删除、发布状态管理、分类和标签管理
- 互动能力：文章评论、留言板、访客身份记录、评论地区展示
- 后台管理：管理员登录、仪表盘、评论/留言管理、友链管理
- 接口文档：后端集成 Swagger，可在本地服务启动后访问

## 技术栈

- 后端：Go 1.24、Gin、GORM、MySQL、Viper、JWT、Swagger
- 前端：Astro 5、React 19、Tailwind CSS、Axios、Marked、Lucide React
- 数据库：MySQL，字符集建议使用 `utf8mb4`

## 目录结构

```text
.
├── backend/          # Go Gin 后端服务
│   ├── config/       # 配置文件与配置读取
│   ├── controllers/  # 接口控制器
│   ├── docs/         # Swagger 文档产物
│   ├── initialize/   # 数据库初始化
│   ├── middlewares/  # 鉴权、CORS、限流等中间件
│   ├── models/       # GORM 模型
│   ├── router/       # 路由注册
│   ├── sql/          # 初始化/修复 SQL
│   └── utils/        # JWT、日志、IP 地区等工具
├── frontend-astro/   # Astro 前端
│   ├── public/       # 静态资源
│   └── src/          # 页面、组件、API 和样式
└── v0-export/        # 早期 UI 导出稿/参考实现
```

## 本地运行

### 1. 准备数据库

创建 MySQL 数据库并导入初始化 SQL：

```bash
mysql -u root -p < backend/sql/init.sql
```

如需修复历史文章评论数，可按需执行：

```bash
mysql -u root -p my_blog < backend/sql/fix_article_comment_count.sql
```

### 2. 配置后端

复制配置模板并按本地环境修改数据库连接和 JWT 密钥：

```bash
cd backend
cp config/config.yml.example config/config.yml
```

关键配置示例：

```yaml
app:
  port: 8081

database:
  dsn: "root:your_password@tcp(127.0.0.1:3306)/my_blog?charset=utf8mb4&parseTime=True&loc=Local"

jwt:
  secret: your_jwt_secret_key_here
  expire_hours: 24
```

### 3. 启动后端

```bash
cd backend
go mod download
go run .
```

启动后可访问：

- 健康检查：`http://127.0.0.1:8081/ping`
- Swagger：`http://127.0.0.1:8081/swagger/index.html`

### 4. 启动前端

前端开发环境默认通过 `PUBLIC_API_BASE_URL` 指向后端：

```bash
cd frontend-astro
npm install
echo "PUBLIC_API_BASE_URL=http://127.0.0.1:8081" > .env.local
npm run dev
```

启动后访问 Astro 提示的本地地址，通常是：

```text
http://localhost:4321
```

## 常用命令

后端：

```bash
cd backend
go run .
go test ./...
```

前端：

```bash
cd frontend-astro
npm run dev
npm run build
npm run preview
```

## API 路由概览

公开接口主要位于 `/api`：

- `GET /api/categories`
- `GET /api/tags`
- `GET /api/articles/timeline`
- `GET /api/articles/search`
- `GET /api/articles/:id`
- `POST /api/comments/add`
- `GET /api/friendlinks`

后台接口主要位于 `/api/admin`，登录后通过 `Authorization: Bearer <token>` 访问：

- `POST /api/admin/login`
- `GET /api/admin/dashboard`
- `GET /api/admin/articles`
- `POST /api/admin/articles/create`
- `PUT /api/admin/articles/update/:id`
- `DELETE /api/admin/articles/delete/:id`
- `GET /api/admin/comments`
- `GET /api/admin/friendlinks`

更完整的接口说明以 Swagger 文档为准。

## 部署说明

前端生产构建后会默认使用同源 API。常见部署方式是将前端静态资源交给 Nginx 托管，并把 `/api` 与 `/swagger` 反向代理到后端服务。

如果后端部署在 Nginx 后面，为了让评论和留言能获取真实客户端 IP，需要透传以下请求头：

```nginx
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
```

## 注意事项

- `backend/config/config.yml`、`frontend-astro/.env.local` 等本地配置不建议提交到公开仓库。
- `frontend-astro/dist/`、`frontend-astro/node_modules/`、`backend/myblog-server` 属于构建或依赖产物，通常不需要纳入版本管理。
- 当前后端配置读取路径为 `backend/config/config.yml`，请在 `backend` 目录下启动服务。
