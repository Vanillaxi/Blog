# MyBlog Apifox 接口测试指南

本文档基于当前 Go Gin 后端路由和 Swagger 注释整理，用于 Apifox 导入、环境变量配置和接口联调测试。

## 1. 项目接口测试总说明

- 后端默认从 `backend/config/config.yml` 读取配置。
- 本文档推荐 Apifox 环境变量按需求设置为：
  - `base_url = http://localhost:8080`
  - `admin_token = 登录接口返回的 token`
- 如果本地 `config.yml` 中实际端口不是 `8080`，例如当前配置常见为 `8081`，请将 Apifox 的 `base_url` 改成实际启动地址。
- 管理员接口统一使用请求头：

```http
Authorization: Bearer {{admin_token}}
```

- 需要登录后保存 token，再测试后台管理接口。
- 文章相关前台接口通常要求文章已发布且未删除。
- 本项目 Swagger 文件位置：
  - `backend/docs/swagger.json`
  - `backend/docs/swagger.yaml`
  - `backend/docs/docs.go`

## 2. Apifox 导入 Swagger 步骤

1. 打开 Apifox，进入目标项目。
2. 选择「导入」。
3. 选择 OpenAPI / Swagger。
4. 选择本地文件：

```text
backend/docs/swagger.json
```

5. 导入完成后，创建或选择测试环境。
6. 添加环境变量：

```text
base_url = http://localhost:8080
admin_token =
category_id =
tag_id =
article_id =
comment_id =
friendlink_id =
```

7. 先执行「管理员登录」，将响应中的 `data.token` 保存到 `admin_token`。

Apifox 登录接口后置脚本示例：

```javascript
const json = pm.response.json();
if (json.data && json.data.token) {
  pm.environment.set("admin_token", json.data.token);
}
```

## 3. 推荐接口测试顺序

1. 管理员登录
2. 创建分类，保存 `category_id`
3. 查询分类
4. 创建标签，保存 `tag_id`
5. 创建文章，保存 `article_id`
6. 发布文章
7. 查询文章列表
8. 查询文章详情
9. 更新文章
10. 删除或下架文章
11. 创建评论/留言，保存 `comment_id`
12. 管理员删除评论/留言
13. 友链增删改查

## 4. 健康检查

### 4.1 健康检查

- 请求方法：`GET`
- 请求路径：`{{base_url}}/ping`
- 是否需要 Authorization：否
- Query 参数：无
- Body JSON 示例：无

成功响应示例：

```json
{
  "message": "pong"
}
```

常见失败情况：

```json
{
  "error": "服务未启动或端口配置不正确"
}
```

## 5. 管理员登录

### 5.1 管理员注册

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/register`
- 是否需要 Authorization：否
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json"
}
```

Body JSON 示例：

```json
{
  "username": "admin",
  "password": "123456"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "Welcome!"
}
```

常见失败情况：

```json
{
  "code": 409,
  "msg": "用户名已被注册"
}
```

前置数据：无。若数据库已存在管理员账号，可跳过注册。

### 5.2 管理员登录

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/login`
- 是否需要 Authorization：否
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json"
}
```

Body JSON 示例：

```json
{
  "username": "admin",
  "password": "123456"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOi...",
    "username": "admin",
    "userID": 1
  }
}
```

常见失败情况：

```json
{
  "code": 401,
  "msg": "用户或密码错误"
}
```

前置数据：需要已有管理员账号。成功后保存 `data.token` 到 `admin_token`。

### 5.3 修改管理员用户名

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/username`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "old_password": "123456",
  "new_username": "newadmin"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "用户名修改成功，请重新登录"
}
```

常见失败情况：

```json
{
  "code": 401,
  "msg": "用户未登录"
}
```

前置数据：需要登录 token。当前鉴权中间件和 controller 上下文取值是否完全匹配，需要根据实际 model/controller 确认。

### 5.4 修改管理员密码

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/password`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "密码修改成功，请重新登录"
}
```

常见失败情况：

```json
{
  "code": 401,
  "msg": "旧密码错误"
}
```

前置数据：需要登录 token。当前鉴权中间件和 controller 上下文取值是否完全匹配，需要根据实际 model/controller 确认。

## 6. 分类管理

### 6.1 创建分类

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/categories/create`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "category_name": "Go",
  "slug": "go",
  "sort": 10
}
```

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "id": 1
  },
  "msg": "创建分类成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "参数格式错误"
}
```

前置数据：需要 `admin_token`。成功后保存 `data.id` 到 `category_id`。

### 6.2 查询后台分类

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/categories`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无
- Body JSON 示例：无

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

成功响应示例：

```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "category_name": "Go",
      "sort": 10,
      "slug": "go",
      "status": 1,
      "article_count": 0
    }
  ],
  "msg": "获取分类列表成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取分类列表失败"
}
```

前置数据：需要 `admin_token`。

### 6.3 更新分类状态

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/categories/{{category_id}}/status`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "status": 1
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "更新分类状态成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "分类ID不合法"
}
```

前置数据：需要 `admin_token` 和 `category_id`。当前 controller 使用 `PostForm("id")` 获取 id，和路由 path 参数不一致，是否能按示例直接通过需要根据实际 model/controller 确认。

### 6.4 更新分类排序

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/categories/{{category_id}}/sort`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "sort": 20
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "更新分类排序成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "分类ID不合法"
}
```

前置数据：需要 `admin_token` 和 `category_id`。当前 controller 使用 `PostForm("id")` 获取 id，和路由 path 参数不一致，是否能按示例直接通过需要根据实际 model/controller 确认。

## 7. 文章管理

### 7.1 创建标签

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/tag/create`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "tag_name": "Gin",
  "slug": "gin",
  "sort": 10
}
```

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "id": 1
  },
  "msg": "创建标签成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "参数格式错误"
}
```

前置数据：需要 `admin_token`。成功后保存 `data.id` 到 `tag_id`。

### 7.2 查询后台标签

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/admin/tags`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无
- Body JSON 示例：无

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

成功响应示例：

```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "tag_name": "Gin",
      "article_count": 0,
      "status": 1,
      "slug": "gin",
      "sort": 10
    }
  ],
  "msg": "获取标签列表成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取标签列表失败"
}
```

前置数据：需要 `admin_token`。

### 7.3 更新标签状态

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/tag/{{tag_id}}/status`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "status": 1
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "更新标签状态成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "标签ID不合法"
}
```

前置数据：需要 `admin_token` 和 `tag_id`。当前 controller 使用 `PostForm("id")` 获取 id，和路由 path 参数不一致，是否能按示例直接通过需要根据实际 model/controller 确认。

### 7.4 更新标签排序

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/tag/{{tag_id}}/sort`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "sort": 20
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "更新标签排序成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "标签ID不合法"
}
```

前置数据：需要 `admin_token` 和 `tag_id`。当前 controller 使用 `PostForm("id")` 获取 id，和路由 path 参数不一致，是否能按示例直接通过需要根据实际 model/controller 确认。

### 7.5 创建文章

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/articles/create`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "title": "Go Gin 博客接口测试",
  "category_id": {{category_id}},
  "summary": "这是一篇用于 Apifox 测试的文章摘要",
  "content": "这是一篇用于 Apifox 测试的文章正文。",
  "is_top": 0,
  "cover_url": "https://example.com/cover.jpg",
  "tag_ids": [{{tag_id}}]
}
```

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "id": 1
  },
  "msg": "创建文章成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "分类不存在"
}
```

前置数据：需要 `admin_token`、`category_id`、`tag_id`。成功后保存 `data.id` 到 `article_id`。

### 7.6 查询后台文章列表

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/admin/articles`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

Query 参数：

```text
page=1
pageSize=10
status=-1
is_deleted=0
category_id=0
keyword=
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "size": 10
  },
  "msg": "获取后台文章列表成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "文章状态参数错误"
}
```

前置数据：需要 `admin_token`。

### 7.7 更新文章

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/articles/update/{{article_id}}`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "title": "Go Gin 博客接口测试 - 已更新",
  "category_id": {{category_id}},
  "summary": "更新后的摘要",
  "content": "更新后的正文内容。",
  "is_top": 1,
  "cover_url": "https://example.com/cover-updated.jpg",
  "tag_ids": [{{tag_id}}]
}
```

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "id": 1
  },
  "msg": "修改文章成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "文章ID不合法"
}
```

前置数据：需要 `admin_token`、`article_id`、`category_id`、`tag_id`。

### 7.8 发布或下架文章

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/articles/updateStatus/{{article_id}}`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "status": 1
}
```

说明：

- `status = 1`：发布
- `status = 2`：下架

成功响应示例：

```json
{
  "code": 200,
  "msg": "文章发布成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "文章状态不合法,只允许发布或下架"
}
```

前置数据：需要 `admin_token` 和 `article_id`。

### 7.9 删除文章

- 请求方法：`DELETE`
- 请求路径：`{{base_url}}/api/admin/articles/delete/{{article_id}}`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无
- Body JSON 示例：无

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "删除文章成果"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "文章ID不合法"
}
```

前置数据：需要 `admin_token` 和 `article_id`。

## 8. 前台文章展示

### 8.1 查询前台分类

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/categories`
- 是否需要 Authorization：否
- Query 参数：无
- Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "category_name": "Go",
      "sort": 10,
      "slug": "go",
      "status": 1,
      "article_count": 0
    }
  ],
  "msg": "获取分类列表成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取分类列表失败"
}
```

前置数据：无。

### 8.2 查询前台标签

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/tags`
- 是否需要 Authorization：否
- Query 参数：无
- Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "tag_name": "Gin",
      "article_count": 0,
      "status": 1,
      "slug": "gin",
      "sort": 10
    }
  ],
  "msg": "获取标签列表成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取标签列表失败"
}
```

前置数据：无。

### 8.3 按分类查询文章

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/categories/{{category_id}}/articles`
- 是否需要 Authorization：否

Query 参数：

```text
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "list": [],
    "total": 0
  },
  "msg": "获取文章列表成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "分类ID不合法"
}
```

前置数据：需要 `category_id`，且该分类下存在已发布文章。当前路由参数名和 controller 读取参数名不一致，是否能按示例直接通过需要根据实际 model/controller 确认。

### 8.4 按标签查询文章

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/tags/{{tag_id}}/articles`
- 是否需要 Authorization：否

Query 参数：

```text
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "list": [],
    "total": 0
  },
  "msg": "获取标签文章列表成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "标签ID不合法"
}
```

前置数据：需要 `tag_id`，且该标签下存在已发布文章。

### 8.5 查询文章详情

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/articles/{{article_id}}`
- 是否需要 Authorization：否
- Query 参数：无
- Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "title": "Go Gin 博客接口测试",
    "content": "这是一篇用于 Apifox 测试的文章正文。",
    "summary": "这是一篇用于 Apifox 测试的文章摘要",
    "cover_url": "https://example.com/cover.jpg",
    "category_id": 1,
    "status": 1
  },
  "msg": "查询成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "查询失败或文章不存在"
}
```

前置数据：需要 `article_id`，文章必须已发布且未删除。

## 9. 搜索/时间轴等

### 9.1 查询文章时间轴

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/articles/timeline`
- 是否需要 Authorization：否

Query 参数：

```text
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "list": [],
    "total": 0
  },
  "msg": "获取文章时间轴成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取文章时间轴失败"
}
```

前置数据：需要存在已发布、未删除、且有 `published_time` 的文章。

### 9.2 搜索文章

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/articles/search`
- 是否需要 Authorization：否

Query 参数：

```text
keyword=Gin
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "size": 10
  },
  "msg": "搜索成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "搜索文章失败"
}
```

前置数据：需要存在已发布且未删除的文章。

## 10. 留言/评论管理

### 10.1 创建评论/留言

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/comments/add`
- 是否需要 Authorization：否
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json"
}
```

Body JSON 示例：

```json
{
  "target_type": 1,
  "target_id": {{article_id}},
  "parent_id": 0,
  "nickname": "Apifox Tester",
  "email": "tester@example.com",
  "website": "https://example.com",
  "content": "这是一条 Apifox 测试评论"
}
```

说明：

- `target_type = 1`：文章评论
- `target_type = 2`：留言板，相关业务含义需要根据实际 model/controller 确认

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "target_type": 1,
    "target_id": 1,
    "nickname": "Apifox Tester",
    "content": "这是一条 Apifox 测试评论"
  },
  "msg": "评论成功"
}
```

常见失败情况：

```json
{
  "code": 400,
  "msg": "参数错误"
}
```

前置数据：文章评论需要 `article_id`。成功后保存 `data.id` 到 `comment_id`。

### 10.2 查询评论/留言

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/cmments/get`
- 是否需要 Authorization：否

Query 参数：

```text
target_type=1
target_id={{article_id}}
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "list": [],
    "total": 0,
    "page": 1,
    "page_size": 10
  },
  "msg": "获取评论成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取评论失败"
}
```

前置数据：文章评论需要 `article_id`。注意当前路由路径是 `/api/cmments/get`，不是 `/api/comments/get`。

### 10.3 管理员删除评论/留言

- 请求方法：`DELETE`
- 请求路径：`{{base_url}}/api/admin/comments/delete/{{comment_id}}`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无
- Body JSON 示例：无

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "删除评论成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "删除评论失败"
}
```

前置数据：需要 `admin_token` 和 `comment_id`。

## 11. 友链管理

### 11.1 创建友链

- 请求方法：`POST`
- 请求路径：`{{base_url}}/api/admin/friendlinks/add`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "name": "Example Blog",
  "url": "https://example.com",
  "logo": "https://example.com/logo.png",
  "sort": 10,
  "status": 1
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "创建成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "创建失败"
}
```

前置数据：需要 `admin_token`。

### 11.2 查询后台友链

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/admin/friendlinks`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

Query 参数：

```text
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "links": [
      {
        "id": 1,
        "name": "Example Blog",
        "url": "https://example.com",
        "logo": "https://example.com/logo.png",
        "sort": 10,
        "status": 1
      }
    ],
    "total": 1,
    "page": 1,
    "size": 10
  },
  "msg": "获取成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取友链失败"
}
```

前置数据：需要 `admin_token`。从响应中保存 `links[0].id` 到 `friendlink_id`。

### 11.3 更新友链

- 请求方法：`PUT`
- 请求路径：`{{base_url}}/api/admin/friendlinks/{{friendlink_id}}/update`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无

Headers：

```json
{
  "Content-Type": "application/json",
  "Authorization": "Bearer {{admin_token}}"
}
```

Body JSON 示例：

```json
{
  "id": {{friendlink_id}},
  "name": "Example Blog Updated",
  "url": "https://example.com",
  "logo": "https://example.com/logo.png",
  "sort": 20,
  "status": 1
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "更新成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "更新失败"
}
```

前置数据：需要 `admin_token` 和 `friendlink_id`。当前 controller 未读取 path 中的 `id`，更新时 body 中的 `id` 必须带上。

### 11.4 删除友链

- 请求方法：`DELETE`
- 请求路径：`{{base_url}}/api/admin/friendlinks/{{friendlink_id}}/delete`
- 是否需要 Authorization：是，`Authorization: Bearer {{admin_token}}`
- Query 参数：无
- Body JSON 示例：无

Headers：

```json
{
  "Authorization": "Bearer {{admin_token}}"
}
```

成功响应示例：

```json
{
  "code": 200,
  "msg": "删除成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "删除失败"
}
```

前置数据：需要 `admin_token` 和 `friendlink_id`。

### 11.5 查询前台友链

- 请求方法：`GET`
- 请求路径：`{{base_url}}/api/friendlinks`
- 是否需要 Authorization：否

Query 参数：

```text
page=1
pageSize=10
```

Body JSON 示例：无

成功响应示例：

```json
{
  "code": 200,
  "data": {
    "links": [],
    "total": 0,
    "page": 1,
    "size": 10
  },
  "msg": "获取成功"
}
```

常见失败情况：

```json
{
  "code": 500,
  "msg": "获取友链失败"
}
```

前置数据：无。若要看到数据，需要先创建 `status = 1` 的友链。

## 12. 需要特别确认的问题

- `base_url` 按需求写为 `http://localhost:8080`，但实际运行端口以 `backend/config/config.yml` 为准。
- `/api/cmments/get` 当前拼写为 `cmments`，Apifox 中请按当前路由填写。
- `/api/categories/:id/articles` 的 controller 读取 `category_id`，和路由 `:id` 不一致，按分类查文章可能失败。
- 分类/标签的状态和排序接口读取 `PostForm("id")`，不是 path 参数，按 JSON + path 测试可能失败。
- 修改用户名、修改密码接口依赖上下文中的 `user`，是否能由当前鉴权中间件正确注入，需要根据实际 model/controller 确认。
- 友链更新接口 path 中有 `friendlink_id`，但 controller 主要依赖 body 中的 `id`。
