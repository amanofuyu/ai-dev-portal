# Dev Portal — 开发文档

> 本文档是 [REQUIREMENTS.md](./REQUIREMENTS.md) 的实现指南，包含架构设计、目录规划、数据库 Schema、接口实现方案、前端组件拆分和开发任务清单。

## 1. 架构总览

```
┌──────────────┐        HTTP/JSON         ┌──────────────┐
│              │  ←────────────────────→   │              │
│   Nuxt 4     │     localhost:7080        │   Go API     │
│   (SSR/CSR)  │                          │   (net/http)  │
│   :3000      │                          │   :7080       │
│              │                          │              │
└──────────────┘                          └──────┬───────┘
                                                 │
                                                 │ database/sql
                                                 ▼
                                          ┌──────────────┐
                                          │   SQLite      │
                                          │   dev.db      │
                                          └──────────────┘
```

- **前端** (`apps/web`) — Nuxt 4，端口 3000，通过 `runtimeConfig.public.apiBaseUrl` 指向后端
- **后端** (`apps/api`) — Go 标准库 `net/http`，端口 7080，SQLite 持久化
- **通信** — RESTful JSON，后端配置 CORS 允许 `localhost:3000`

## 2. 目录结构（规划）

```
dev-portal/
├── apps/
│   ├── api/                          # Go 后端
│   │   ├── cmd/
│   │   │   └── api/
│   │   │       └── main.go           # 入口：初始化 DB、注册路由、启动服务
│   │   ├── internal/
│   │   │   ├── db/
│   │   │   │   ├── db.go             # 数据库初始化、迁移、seed
│   │   │   │   └── seed.go           # 预置数据
│   │   │   ├── handler/
│   │   │   │   ├── project.go        # Project CRUD handlers
│   │   │   │   ├── apikey.go         # ApiKey handlers（含 reveal）
│   │   │   │   └── response.go       # writeJSON / writeError 通用响应
│   │   │   ├── middleware/
│   │   │   │   └── cors.go           # CORS 中间件
│   │   │   └── model/
│   │   │       └── model.go          # Go struct 定义（Project, ApiKey）
│   │   ├── go.mod
│   │   ├── go.sum
│   │   └── package.json              # Turborepo 任务入口
│   │
│   └── web/                          # Nuxt 4 前端
│       ├── app/
│       │   ├── app.vue
│       │   ├── app.config.ts
│       │   ├── layouts/
│       │   │   └── default.vue       # 全局布局（导航栏 + 容器）
│       │   ├── pages/
│       │   │   ├── index.vue         # 项目列表页
│       │   │   └── projects/
│       │   │       └── [id].vue      # 项目详情 / 密钥管理页
│       │   ├── components/
│       │   │   ├── ProjectCard.vue        # 项目卡片
│       │   │   ├── ProjectForm.vue        # 项目创建/编辑表单（Modal）
│       │   │   ├── DeleteConfirm.vue      # 删除确认弹窗
│       │   │   ├── KeyTable.vue           # 密钥列表表格
│       │   │   ├── KeyRow.vue             # 单行密钥（含脱敏/明文切换、toggle、复制）
│       │   │   ├── KeyCreateModal.vue     # 密钥创建弹窗 + 明文展示
│       │   │   ├── ToggleSwitch.vue       # 启用/禁用开关（封装乐观更新逻辑）
│       │   │   └── AppToast.vue           # Toast 通知
│       │   ├── composables/
│       │   │   ├── useProjects.ts         # 项目 CRUD composable
│       │   │   ├── useApiKeys.ts          # 密钥 CRUD + reveal composable
│       │   │   ├── useToast.ts            # Toast 状态管理（provide/inject）
│       │   │   └── useClipboard.ts        # 剪贴板（可直接用 VueUse）
│       │   ├── types/
│       │   │   └── index.ts               # TypeScript 接口定义
│       │   └── assets/
│       │       └── (静态资源)
│       ├── nuxt.config.ts
│       ├── tailwind.css
│       ├── eslint.config.mjs
│       ├── tsconfig.json
│       └── package.json
│
├── turbo.json
├── package.json
├── pnpm-workspace.yaml
├── design.md              # API 契约与数据模型权威来源
├── REQUIREMENTS.md
├── DEVELOPMENT.md
└── TASKS.md
```

## 3. 数据库设计

### 3.1 DDL（SQLite）

```sql
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS projects (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL UNIQUE,
    description TEXT    NOT NULL DEFAULT '',
    status      TEXT    NOT NULL DEFAULT 'Active' CHECK (status IN ('Active', 'Archived')),
    created_at  TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS api_keys (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    key_value   TEXT    NOT NULL UNIQUE,
    name        TEXT    NOT NULL,
    is_enabled  INTEGER NOT NULL DEFAULT 1,
    last_used_at TEXT,
    project_id  INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_api_keys_project_id ON api_keys(project_id);
```

> SQLite 需要 `PRAGMA foreign_keys = ON` 才能启用外键约束。在 `db.go` 中打开连接后立即执行。

### 3.2 Go Model 定义

```go
package model

import "time"

type Project struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    KeyCount    int       `json:"key_count"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ApiKey struct {
    ID             int64      `json:"id"`
    KeyValue       string     `json:"key_value,omitempty"`
    KeyValueMasked string     `json:"key_value_masked,omitempty"`
    Name           string     `json:"name"`
    IsEnabled      bool       `json:"is_enabled"`
    LastUsedAt     *time.Time `json:"last_used_at"`
    ProjectID      int64      `json:"project_id"`
    CreatedAt      time.Time  `json:"created_at"`
}

type ApiKeyRevealed struct {
    ID       int64  `json:"id"`
    KeyValue string `json:"key_value"`
}
```

> **注意**：
> - `Project.KeyCount` 不使用 `omitempty`，确保值为 0 时仍然输出 `"key_count": 0`。
> - `ApiKey.IsEnabled` 是 Go `bool`，但 SQLite 中存储为 `INTEGER`（1/0）。`go-sqlite3` 驱动在 `Scan` 时可以自动将 `INTEGER` 转换为 `bool`，直接 `rows.Scan(&key.IsEnabled)` 即可。

### 3.3 脱敏函数

```go
func MaskKeyValue(key string) string {
    // key 格式: sk-<hex chars>
    // 保留: sk- + 前4位 + ****** + 后4位
    prefix := "sk-"
    if !strings.HasPrefix(key, prefix) || len(key) < len(prefix)+8 {
        return "sk-******"
    }
    body := key[len(prefix):]
    return prefix + body[:4] + "******" + body[len(body)-4:]
}
```

### 3.4 密钥生成

```go
import "crypto/rand"

func GenerateAPIKey() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return "sk-" + hex.EncodeToString(bytes), nil
}
```

## 4. 后端实现方案

### 4.1 入口 (`cmd/api/main.go`)

1. 解析配置（端口、DB 路径，默认 `:7080` / `dev.db`）
2. 调用 `db.Init()` 打开 SQLite 连接、执行 DDL、seed 数据
3. 注册路由（使用 Go 1.22+ `http.NewServeMux` 的方法匹配）
4. 包裹 CORS 中间件
5. 启动 `http.ListenAndServe`

### 4.2 路由注册

```go
mux := http.NewServeMux()

// Projects
mux.HandleFunc("GET /projects", h.ListProjects)
mux.HandleFunc("POST /projects", h.CreateProject)
mux.HandleFunc("GET /projects/{id}", h.GetProject)
mux.HandleFunc("PATCH /projects/{id}", h.UpdateProject)
mux.HandleFunc("DELETE /projects/{id}", h.DeleteProject)

// ApiKeys
mux.HandleFunc("GET /projects/{projectId}/keys", h.ListKeys)
mux.HandleFunc("POST /projects/{projectId}/keys", h.CreateKey)
mux.HandleFunc("PATCH /keys/{id}", h.UpdateKey)
mux.HandleFunc("DELETE /keys/{id}", h.DeleteKey)
mux.HandleFunc("GET /keys/{id}/reveal", h.RevealKey)
```

> Go 1.22+ 的 `ServeMux` 原生支持 `METHOD /path/{param}` 模式，无需第三方路由库。

### 4.3 Handler 要点

| Handler | 要点 |
|---------|------|
| `ListProjects` | 查询所有 project，LEFT JOIN 统计 key_count |
| `CreateProject` | 解析 JSON body，INSERT，返回创建的对象 |
| `GetProject` | 通过 path param 获取单个 project |
| `UpdateProject` | 部分更新（name/description/status），更新 updated_at |
| `DeleteProject` | DELETE，级联自动删除 keys（外键约束） |
| `ListKeys` | 查询 project_id 下的 keys，返回时使用 `MaskKeyValue()` 脱敏 |
| `CreateKey` | 调用 `GenerateAPIKey()` 生成密钥，INSERT，返回含明文的完整对象 |
| `UpdateKey` | 更新 `is_enabled` 字段 |
| `DeleteKey` | DELETE 单个密钥 |
| `RevealKey` | 返回单个密钥的 id + key_value 明文 |

### 4.4 CORS 中间件

```go
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### 4.5 Seed 数据

在 `db.Init()` 中，检查 `projects` 表是否为空，若空则插入预置数据：

- **Cloud Platform** (Active) — 4 个 key
- **Mobile App** (Active) — 3 个 key

每个 key 使用 `GenerateAPIKey()` 动态生成，保证每次 seed 密钥不同。

## 5. 前端实现方案

### 5.1 TypeScript 类型 (`app/types/index.ts`)

与 design.md 第 5 节保持一致，使用三个独立接口区分不同 API 场景：

```typescript
/** 项目 */
export interface Project {
  id: number
  name: string
  description: string
  status: 'Active' | 'Archived'
  key_count: number
  created_at: string
  updated_at: string
}

/** 创建/更新项目的请求体 */
export interface ProjectInput {
  name?: string
  description?: string
  status?: 'Active' | 'Archived'
}

/** 密钥（脱敏，列表接口返回） */
export interface ApiKeyMasked {
  id: number
  key_value_masked: string
  name: string
  is_enabled: boolean
  last_used_at: string | null
  project_id: number
  created_at: string
}

/** 密钥（完整明文，创建接口返回） */
export interface ApiKeyFull {
  id: number
  key_value: string
  name: string
  is_enabled: boolean
  last_used_at: string | null
  project_id: number
  created_at: string
}

/** 密钥明文（reveal 接口返回） */
export interface ApiKeyRevealed {
  id: number
  key_value: string
}

/** 通用错误响应 */
export interface ApiError {
  error: string
}
```

### 5.2 Composables

#### `useProjects.ts`

封装项目 CRUD 操作，使用 `useFetch` + `$fetch`：

- `fetchProjects()` — GET /projects
- `createProject(data)` — POST /projects
- `updateProject(id, data)` — PATCH /projects/:id
- `deleteProject(id)` — DELETE /projects/:id

#### `useApiKeys.ts`

封装密钥操作：

- `fetchKeys(projectId)` — GET /projects/:id/keys
- `createKey(projectId, name)` — POST /projects/:id/keys
- `toggleKey(id, isEnabled)` — PATCH /keys/:id
- `deleteKey(id)` — DELETE /keys/:id
- `revealKey(id)` — GET /keys/:id/reveal

#### `useToast.ts`

基于 provide/inject 模式（遵循前端规约）：

- `provideToastContext()` — 在 layout 中 provide
- `useToastContext()` — 在组件中 inject
- 支持 `success` / `error` / `info` 类型
- 自动 3 秒消失

### 5.3 页面

#### `pages/index.vue` — 项目列表

- 调用 `useProjects()` 获取数据
- Grid 布局展示 `ProjectCard`
- 顶部标题 + "新建项目" 按钮
- 空状态引导

#### `pages/projects/[id].vue` — 密钥管理

- 路由参数 `id` 获取当前项目
- 展示项目信息卡片（名称、描述、状态）
- `KeyTable` 组件展示密钥列表
- "生成密钥" 按钮触发 `KeyCreateModal`

### 5.4 组件设计

#### `ProjectCard.vue`

- 展示项目名称、描述、状态 badge、密钥数量
- 点击整体跳转到详情页
- 编辑/删除操作按钮

#### `KeyRow.vue`（核心交互组件）

```
┌────────────┬──────────────────┬──────────┬────────────┬────────────┬──────────────┐
│ Name       │ Key              │ Status   │ Created    │ Last Used  │ Actions      │
├────────────┼──────────────────┼──────────┼────────────┼────────────┼──────────────┤
│ Prod Key   │ sk-a1b2******p6  │ [Toggle] │ 2026-02-25 │ Never      │ 👁 📋 🗑    │
└────────────┴──────────────────┴──────────┴────────────┴────────────┴──────────────┘
```

状态管理：
- `isRevealed: Ref<boolean>` — 控制明文/脱敏显示
- `revealedValue: Ref<string | null>` — 缓存 reveal 结果
- 点击眼睛：若未 reveal 则调用 API，已 reveal 则切换显示状态
- 复制：若已 reveal 则直接复制缓存值，否则先 reveal 再复制

#### `ToggleSwitch.vue`

- Props: `modelValue: boolean`, `keyId: number`
- 乐观更新流程：
  1. 点击 → 立即翻转本地状态 → emit update
  2. 发起 PATCH 请求
  3. 失败 → 回滚状态 → toast 报错

#### `KeyCreateModal.vue`

- 创建成功后展示完整密钥
- 警告文案："此密钥仅显示一次，请妥善保存"
- 提供复制按钮
- 关闭弹窗后触发列表刷新

### 5.5 主题与样式

在 `tailwind.css` 中配置 DaisyUI 深色主题：

```css
@import 'tailwindcss';
@plugin "daisyui" {
  themes: business --default;
}
```

`business` 主题提供深灰/蓝色调专业风格，符合"安全感设计"要求。

### 5.6 布局 (`layouts/default.vue`)

```
┌───────────────────────────────────────────────┐
│  🔑 Dev Portal                    [主题切换]  │  ← Navbar
├───────────────────────────────────────────────┤
│                                               │
│               Page Content                    │  ← <slot />
│                                               │
└───────────────────────────────────────────────┘
│             Toast Container                   │  ← 固定定位右上角
```

## 6. 开发命令

```bash
# 安装依赖
pnpm install

# 初始化 Go 模块（首次）
cd apps/api && go mod init dev-portal/api && go mod tidy

# 启动开发（前后端并行）
pnpm dev                    # turbo run dev (并行启动 web + api)

# 单独启动
pnpm --filter web dev       # Nuxt dev server :3000
pnpm --filter api dev       # Go server :7080

# 构建
pnpm build                  # turbo run build

# Lint
pnpm lint                   # turbo run lint (前端 eslint + 后端 go vet)

# 前端 lint 自动修复（需在 apps/web/package.json 中定义 "lint:fix": "eslint . --fix"）
pnpm --filter web lint:fix
```

## 7. 开发任务分解

以下为建议的实现顺序，每个任务可独立验证：

### 阶段一：后端基础

1. **初始化 Go 项目** — `go.mod`、目录结构、SQLite 驱动依赖
2. **数据库层** — `db.go`（连接、DDL 建表）、`seed.go`（预置数据）
3. **Model + 工具函数** — `model.go`、`GenerateAPIKey()`、`MaskKeyValue()`
4. **Project handlers** — CRUD 五个接口
5. **ApiKey handlers** — CRUD + reveal 五个接口
6. **CORS 中间件** — 允许前端跨域
7. **主入口集成** — `main.go` 串联所有模块，`pnpm --filter api dev` 可启动

### 阶段二：前端基础

8. **类型定义** — `app/types/index.ts`
9. **布局与主题** — `default.vue`、DaisyUI `business` 主题配置
10. **Toast 系统** — `useToast.ts` + `AppToast.vue`
11. **项目列表页** — `pages/index.vue` + `ProjectCard.vue`
12. **项目表单** — `ProjectForm.vue`（创建/编辑 Modal）
13. **删除确认** — `DeleteConfirm.vue`

### 阶段三：密钥管理（核心功能）

14. **密钥管理页** — `pages/projects/[id].vue`
15. **密钥列表** — `KeyTable.vue` + `KeyRow.vue`
16. **脱敏与 Reveal** — 小眼睛交互，调用 reveal API
17. **一键复制** — 集成 `useClipboard`（VueUse）
18. **Toggle 开关** — `ToggleSwitch.vue`，乐观更新 + 回滚
19. **创建密钥弹窗** — `KeyCreateModal.vue`，一次性展示明文

### 阶段四：收尾

20. **错误处理完善** — 所有 API 调用的 loading/error 状态
21. **空状态** — 无项目、无密钥时的引导 UI
22. **自检验证** — 按 REQUIREMENTS.md 第 8 节逐项验收
