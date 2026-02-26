# Dev Portal — 开发者 API 密钥控制台 需求文档

## 1. 项目概述

模拟 AWS / OpenAI 开发者后台的 **API 密钥管理控制台**。用户可以创建项目（Project），并在项目下创建、查看、启用/禁用和删除 API 密钥（API Key）。系统对安全性和交互体验有较高要求：敏感数据（密钥明文）需做脱敏处理，开关状态需即时反馈。

## 2. 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| 包管理 | pnpm | Monorepo，禁止 npm/yarn |
| 构建编排 | Turborepo | 任务在各 package 的 package.json 定义 |
| 后端 | Go + SQLite | `apps/api`，标准库优先，REST API |
| 前端 | Nuxt 4 + TypeScript + VueUse + Tailwind CSS + DaisyUI | `apps/web`，Composition API + `<script setup>` |

## 3. 数据模型

### 3.1 Project（项目）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | INTEGER | PK, AUTO INCREMENT | 主键 |
| `name` | TEXT | NOT NULL, UNIQUE | 项目名称 |
| `description` | TEXT | | 项目描述 |
| `status` | TEXT | NOT NULL, DEFAULT 'Active' | 项目状态：`Active` / `Archived` |
| `created_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 更新时间 |

### 3.2 ApiKey（密钥）

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | INTEGER | PK, AUTO INCREMENT | 主键 |
| `key_value` | TEXT | NOT NULL, UNIQUE | 密钥值，格式 `sk-<random>`，后端生成 |
| `name` | TEXT | NOT NULL | 密钥名称/备注，方便用户区分 |
| `is_enabled` | BOOLEAN | NOT NULL, DEFAULT TRUE | 是否启用 |
| `last_used_at` | DATETIME | | 最后使用时间，可为空 |
| `project_id` | INTEGER | NOT NULL, FK → Project(id) ON DELETE CASCADE | 所属项目 |
| `created_at` | DATETIME | NOT NULL, DEFAULT CURRENT_TIMESTAMP | 创建时间 |

### 3.3 关系

- **One-to-Many**：一个 Project 拥有多个 ApiKey
- **级联删除**：删除 Project 时，其下所有 ApiKey 自动级联删除（SQLite `ON DELETE CASCADE`）

## 4. API 接口设计

基础路径：`http://localhost:7080`

### 4.1 项目（Project）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/projects` | 获取所有项目列表 |
| POST | `/projects` | 创建新项目 |
| GET | `/projects/:id` | 获取单个项目详情 |
| PATCH | `/projects/:id` | 更新项目（名称、描述、状态） |
| DELETE | `/projects/:id` | 删除项目（级联删除其下所有密钥） |

#### POST /projects — 请求体

```json
{
  "name": "My App",
  "description": "Production application"
}
```

#### GET /projects — 响应体

```json
[
  {
    "id": 1,
    "name": "My App",
    "description": "Production application",
    "status": "Active",
    "key_count": 3,
    "created_at": "2026-02-25T10:00:00Z",
    "updated_at": "2026-02-25T10:00:00Z"
  }
]
```

#### PATCH /projects/:id — 请求体

```json
{
  "name": "Renamed App",
  "description": "Updated description",
  "status": "Archived"
}
```

### 4.2 密钥（ApiKey）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/projects/:id/keys` | 获取某项目下所有密钥（**脱敏返回**） |
| POST | `/projects/:id/keys` | 为项目生成新密钥 |
| PATCH | `/keys/:id` | 更新密钥（切换启用/禁用状态） |
| DELETE | `/keys/:id` | 删除单个密钥 |
| GET | `/keys/:id/reveal` | 获取单个密钥的**明文**（专用接口） |

#### POST /projects/:id/keys — 请求体

```json
{
  "name": "Production Key"
}
```

#### POST /projects/:id/keys — 响应体（仅创建时返回一次明文）

```json
{
  "id": 5,
  "key_value": "sk-a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
  "name": "Production Key",
  "is_enabled": true,
  "last_used_at": null,
  "project_id": 1,
  "created_at": "2026-02-25T12:00:00Z"
}
```

#### GET /projects/:id/keys — 响应体（脱敏）

```json
[
  {
    "id": 5,
    "key_value_masked": "sk-a1b2******a1b2",
    "name": "Production Key",
    "is_enabled": true,
    "last_used_at": null,
    "project_id": 1,
    "created_at": "2026-02-25T12:00:00Z"
  }
]
```

> **安全设计**：列表接口返回脱敏值 `key_value_masked`，不包含明文字段。明文仅通过 `GET /keys/:id/reveal` 单独获取，便于前端按需加载，降低批量泄漏风险。

#### PATCH /keys/:id — 请求体

```json
{
  "is_enabled": false
}
```

#### GET /keys/:id/reveal — 响应体

```json
{
  "id": 5,
  "key_value": "sk-a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
}
```

### 4.3 通用错误响应

```json
{
  "error": "error message here"
}
```

HTTP 状态码约定：
- `400` — 参数错误 / 校验失败
- `404` — 资源不存在
- `500` — 服务端错误

### 4.4 CORS

后端需配置 CORS，允许前端开发服务器（`http://localhost:3000`）跨域访问。

## 5. 安全设计

### 5.1 后端脱敏（核心策略）

- **列表接口不返回明文**：`GET /projects/:id/keys` 仅返回 `key_value_masked`（如 `sk-a1b2******a1b2`），不含 `key_value` 字段
- **明文按需获取**：前端点击"小眼睛"时，调用 `GET /keys/:id/reveal` 获取完整密钥
- **创建时一次性展示**：`POST /projects/:id/keys` 的响应返回完整明文，前端提示用户"此密钥仅显示一次"

### 5.2 脱敏规则

密钥格式为 `sk-<64位hex字符>`（共 67 字符），脱敏规则：保留前缀 `sk-` + 前 4 位 + `******` + 后 4 位。例如：
- 明文：`sk-a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2`
- 脱敏：`sk-a1b2******a1b2`

### 5.3 密钥生成

- 使用 Go 标准库 `crypto/rand` 生成 32 字节随机数据
- Hex 编码后拼接 `sk-` 前缀，保证高强度随机性
- 最终格式：`sk-<64位hex字符>`

## 6. 前端界面要求

### 6.1 整体风格

- **安全感设计**：深灰/科技蓝色调，强调专业性
- 使用 DaisyUI 的深色主题（如 `dark` / `business`）
- 布局：左侧项目导航 + 右侧密钥管理主区域，或顶部项目切换 + 下方密钥列表

### 6.2 页面结构

#### 项目列表页

- 展示所有项目，显示名称、描述、状态、密钥数量
- 支持创建新项目（弹窗/抽屉表单）
- 支持编辑项目信息、归档/激活项目
- 支持删除项目（二次确认弹窗，提示将级联删除所有密钥）

#### 项目详情 / 密钥管理页

- 展示项目基本信息
- 密钥列表表格，包含以下列：
  - 密钥名称
  - 密钥值（脱敏显示）
  - 状态（启用/禁用 Toggle）
  - 创建时间
  - 最后使用时间
  - 操作（查看/复制/删除）

### 6.3 敏感信息交互

| 功能 | 行为 |
|------|------|
| **默认脱敏** | 密钥在列表中显示为 `sk-a1b2******a1b2` |
| **查看明文** | 提供"小眼睛"图标，点击后调用 reveal API 显示完整密钥，再次点击恢复脱敏 |
| **一键复制** | 提供 Copy 按钮，点击后将完整密钥复制到剪贴板，显示"已复制"提示 |
| **创建展示** | 新密钥创建成功后，弹窗展示完整密钥并警告"此密钥仅显示一次，请妥善保存" |

### 6.4 开关控件

- 使用 DaisyUI 的 Toggle 组件控制密钥启用/禁用
- 点击后**立即更新视觉状态**（乐观更新）
- 同时发起 `PATCH /keys/:id` 请求更新后端状态
- 若请求失败，**回滚视觉状态**并显示错误提示

### 6.5 错误与加载状态

- 所有 API 请求需显示 loading 状态（骨架屏或 spinner）
- 网络/API 错误通过 toast 通知用户
- 空状态：项目无密钥时显示引导文案和创建按钮

## 7. 预置数据（Seed）

后端启动时自动检测并填充预置数据：

### 项目 1：Cloud Platform

- 状态：Active
- 密钥：
  - "Production Key" — 启用
  - "Staging Key" — 启用
  - "Legacy Key" — 禁用
  - "Testing Key" — 启用

### 项目 2：Mobile App

- 状态：Active
- 密钥：
  - "iOS Key" — 启用
  - "Android Key" — 启用
  - "Deprecated Key" — 禁用

## 8. 自检标准

- [ ] **安全意识**：列表接口 `GET /projects/:id/keys` 不返回明文密钥，仅返回脱敏值；明文通过独立的 reveal 接口按需获取
- [ ] **交互体验**：Toggle 开关点击后立即更新视觉状态（乐观更新），同时触发 PATCH API 请求；失败时回滚并提示
- [ ] **级联关系**：删除项目后，其下所有密钥在数据库中被级联删除，不存在游离数据
- [ ] **密钥生成**：使用加密安全的随机数生成器（`crypto/rand`），密钥具备足够强度
- [ ] **一键复制**：复制功能正常工作，复制的是完整明文而非脱敏值
- [ ] **创建提醒**：新密钥创建后明确提示用户"仅显示一次"
- [ ] **CORS 配置**：前后端跨域访问正常
- [ ] **错误处理**：API 错误统一格式，前端对错误有可见的用户反馈
