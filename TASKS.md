# Dev Portal — 开发任务清单

> 按阶段顺序执行。每个任务有明确的验收条件（AC），完成后勾选 `[x]`。

---

## 阶段一：后端基础

### 1.1 初始化 Go 项目

- [x] 在 `apps/api/` 下执行 `go mod init dev-portal/api`
- [x] 添加 SQLite 驱动依赖（使用纯 Go 驱动 `modernc.org/sqlite`，无需 CGO）
- [x] 创建目录结构：`cmd/api/`、`internal/db/`、`internal/handler/`、`internal/middleware/`、`internal/model/`
- [x] `go mod tidy` 通过

**AC**：`apps/api/go.mod` 存在，包含 module 和 sqlite3 依赖；目录结构就绪。

### 1.2 Model + 工具函数

- [x] 实现 `internal/model/model.go`：定义 `Project`、`ApiKey`、`ApiKeyRevealed` 结构体（按 design.md 第 6 节）
- [x] 实现 `GenerateAPIKey() (string, error)`：使用 `crypto/rand` 生成 `sk-<64位hex>`
- [x] 实现 `MaskKeyValue(key string) string`：脱敏规则 `sk-` + 前4位 + `******` + 后4位

**AC**：结构体 JSON tag 与 design.md 一致；`GenerateAPIKey()` 每次返回不同的 67 字符密钥；`MaskKeyValue("sk-a1b2c3d4e5f6...")` 返回 `"sk-a1b2******f6...后4位"`。

### 1.3 数据库层

- [x] 实现 `internal/db/db.go`：`Init(dbPath string) (*sql.DB, error)`
  - 打开 SQLite 连接
  - 执行 `PRAGMA foreign_keys = ON`
  - 执行 DDL 建表（按 design.md 第 2 节）
- [x] 实现 `internal/db/seed.go`：`Seed(db *sql.DB) error`
  - 检查 `projects` 表是否为空
  - 为空时插入 2 个项目 + 7 个密钥（按 design.md 第 7 节）
  - 密钥通过 `GenerateAPIKey()` 动态生成

**AC**：首次启动后 `dev.db` 存在，包含 2 个项目和 7 个密钥；再次启动不重复插入；`PRAGMA foreign_keys` 已开启（可通过 `PRAGMA foreign_keys` 查询确认返回 1）。

### 1.4 通用响应 + CORS 中间件

- [x] 实现 `internal/handler/response.go`：`writeJSON(w, status, v)` 和 `writeError(w, status, msg)`
- [x] 实现 `internal/middleware/cors.go`：允许 `http://localhost:3000` 跨域，支持 GET/POST/PATCH/DELETE/OPTIONS

**AC**：`writeError(w, 400, "test")` 输出 `{"error":"test"}` + 400 状态码；OPTIONS 请求返回 204 + 正确的 CORS 头。

### 1.5 Project Handlers

- [x] 实现 `internal/handler/project.go`，包含以下 5 个 handler：
  - `ListProjects` — GET /projects，LEFT JOIN 统计 key_count
  - `CreateProject` — POST /projects，校验 name 必填 + 唯一
  - `GetProject` — GET /projects/{id}，含 key_count
  - `UpdateProject` — PATCH /projects/{id}，部分更新，校验 status 合法值
  - `DeleteProject` — DELETE /projects/{id}，返回 204

**AC**：使用 curl 或 HTTP 客户端验证全部 5 个接口的正常响应和错误响应均符合 design.md 第 3.1–3.5 节。

### 1.6 ApiKey Handlers

- [x] 实现 `internal/handler/apikey.go`，包含以下 5 个 handler：
  - `ListKeys` — GET /projects/{projectId}/keys，返回脱敏列表（`key_value_masked`，无 `key_value`）
  - `CreateKey` — POST /projects/{projectId}/keys，生成密钥，返回含明文的完整对象
  - `UpdateKey` — PATCH /keys/{id}，更新 `is_enabled`
  - `DeleteKey` — DELETE /keys/{id}，返回 204
  - `RevealKey` — GET /keys/{id}/reveal，返回 `{id, key_value}`

**AC**：`ListKeys` 响应不含 `key_value` 字段；`CreateKey` 响应含 `key_value`；`RevealKey` 仅返回 id 和 key_value；所有响应格式与 design.md 第 3.6–3.10 节一致。

### 1.7 主入口集成

- [x] 实现 `cmd/api/main.go`：
  - 读取环境变量 `PORT`（默认 `7080`）和 `DB_PATH`（默认 `dev.db`）
  - 调用 `db.Init()` + `db.Seed()`
  - 使用 Go 1.22+ `http.NewServeMux` 注册所有路由（按 DEVELOPMENT.md 第 4.2 节）
  - 包裹 CORS 中间件
  - 启动 `http.ListenAndServe`
- [x] 确认 `pnpm --filter api dev` 可正常启动服务，监听 7080 端口

**AC**：`pnpm --filter api dev` 启动成功，`curl http://localhost:7080/projects` 返回含 2 个项目的 JSON 数组。

---

## 阶段二：前端基础

### 2.1 TypeScript 类型定义

- [x] 创建 `apps/web/app/types/index.ts`，定义所有接口类型（按 design.md 第 5 节）：
  - `Project`、`ProjectInput`
  - `ApiKeyMasked`、`ApiKeyFull`、`ApiKeyRevealed`
  - `ApiError`

**AC**：类型文件存在且通过 TypeScript 编译；字段与 design.md 一致。

### 2.2 DaisyUI 主题 + 全局布局

- [x] 修改 `apps/web/app/tailwind.css`，配置 DaisyUI `business` 主题为默认
- [x] 创建 `apps/web/app/layouts/default.vue`：
  - 顶部 Navbar：左侧 Logo + "Dev Portal" 标题，可点击回首页
  - 中间内容区 `<slot />`，限制最大宽度居中
- [x] 确认 `app.vue` 使用 `<NuxtLayout>` + `<NuxtPage />`

**AC**：`pnpm --filter web dev` 启动后访问 `localhost:3000`，页面显示深色主题的 Navbar；整体色调为深灰/蓝调。

### 2.3 Toast 通知系统

- [x] 创建 `apps/web/app/composables/useToast.ts`：
  - 按前端规约的 provide/inject 模式封装
  - `provideToastContext()`：layout 中调用
  - `useToastContext()`：组件中注入
  - 支持 `success(msg)` / `error(msg)` / `info(msg)`
  - 消息自动 3 秒后消失
- [x] 创建 `apps/web/app/components/AppToast.vue`：
  - 固定定位右上角
  - 不同类型使用不同颜色（success=绿、error=红、info=蓝）
  - 进入/退出动画

**AC**：在任意组件中调用 `useToastContext().success('测试')` 后，页面右上角出现绿色 toast 并 3 秒后自动消失。

### 2.4 项目列表页

- [x] 创建 `apps/web/app/composables/useProjects.ts`：
  - `fetchProjects()` — `useFetch` GET /projects
  - `createProject(data: ProjectInput)` — `$fetch` POST /projects
  - `updateProject(id, data: ProjectInput)` — `$fetch` PATCH /projects/:id
  - `deleteProject(id)` — `$fetch` DELETE /projects/:id
- [x] 创建 `apps/web/app/components/ProjectCard.vue`：
  - 展示项目名称、描述（截断）、状态 badge（Active=绿、Archived=灰）、密钥数量
  - 点击卡片跳转到 `/projects/:id`
  - 编辑、删除操作按钮（右上角或底部）
- [x] 实现 `apps/web/app/pages/index.vue`：
  - 标题 "Projects" + "New Project" 按钮
  - Grid 布局展示 `ProjectCard` 列表
  - Loading 骨架屏 / Error 提示 / 空状态引导

**AC**：首页加载后显示 seed 的 2 个项目卡片；每个卡片展示名称、描述、状态和密钥数量；点击卡片跳转到详情页。

### 2.5 项目创建/编辑表单

- [x] 创建 `apps/web/app/components/ProjectForm.vue`：
  - Modal 弹窗形式
  - 输入字段：名称（必填）、描述（可选）
  - 编辑模式：回填现有数据，可修改状态（Active/Archived 下拉）
  - 提交后调用 API，成功关闭弹窗并刷新列表，失败 toast 报错

**AC**：点击 "New Project" 弹出表单；填写并提交后列表新增一个项目；编辑已有项目可修改名称/描述/状态。

### 2.6 删除确认弹窗

- [x] 创建 `apps/web/app/components/DeleteConfirm.vue`：
  - 通用确认弹窗，接受 `title`、`message`、`onConfirm` props
  - 项目删除时提示"将同时删除该项目下的所有 API 密钥"
  - 确认后调用 API 删除，成功后刷新列表并 toast 成功提示

**AC**：点击项目卡片上的删除按钮，弹出确认弹窗显示警告文案；确认后项目从列表消失。

---

## 阶段三：密钥管理（核心功能）

### 3.1 密钥管理页

- [x] 创建 `apps/web/app/composables/useApiKeys.ts`：
  - `fetchKeys(projectId)` — `useFetch` GET /projects/:id/keys
  - `createKey(projectId, name)` — `$fetch` POST /projects/:id/keys
  - `toggleKey(id, isEnabled)` — `$fetch` PATCH /keys/:id
  - `deleteKey(id)` — `$fetch` DELETE /keys/:id
  - `revealKey(id)` — `$fetch` GET /keys/:id/reveal
- [x] 创建 `apps/web/app/pages/projects/[id].vue`：
  - 获取路由参数 `id`，加载项目详情 + 密钥列表
  - 展示项目信息卡片（名称、描述、状态、编辑按钮）
  - "Generate New Key" 按钮
  - 密钥列表区域（`KeyTable` 组件）
  - 返回项目列表的导航链接

**AC**：从项目卡片点击跳转后，页面展示项目信息和密钥列表；列表中密钥以脱敏形式显示。

### 3.2 密钥列表 + 行组件

- [x] 创建 `apps/web/app/components/KeyTable.vue`：
  - 表头：密钥名称 / 密钥值 / 状态 / 创建时间 / 最后使用 / 操作
  - 空状态：无密钥时显示引导文案和 "Generate Key" 按钮
- [x] 创建 `apps/web/app/components/KeyRow.vue`：
  - 密钥名称
  - 脱敏密钥值（等宽字体 `font-mono`）
  - Toggle 开关
  - 创建时间（格式化）
  - 最后使用时间（空值显示 "Never"）
  - 操作列：查看（👁）、复制（📋）、删除（🗑）

**AC**：密钥列表渲染 seed 数据的全部密钥；每行展示脱敏值如 `sk-a1b2******a1b2`；禁用密钥的 Toggle 处于关闭状态。

### 3.3 小眼睛 Reveal + 一键复制

- [x] 在 `KeyRow.vue` 中实现"小眼睛"交互：
  - 维护 `isRevealed` + `revealedValue` 状态
  - 首次点击：调用 `GET /keys/:id/reveal`，获取明文后缓存，切换显示
  - 再次点击：直接切换显示（不重复请求），恢复脱敏
  - 加载中显示 spinner
- [x] 实现一键复制功能：
  - 使用 VueUse 的 `useClipboard`
  - 点击复制按钮：若已 reveal 则直接复制缓存值；否则先调用 reveal 再复制
  - 复制成功后 toast 提示 "已复制到剪贴板"

**AC**：点击小眼睛，密钥值从 `sk-a1b2******a1b2` 变为完整明文；再次点击恢复脱敏；点击复制按钮后剪贴板中为完整明文。

### 3.4 Toggle 开关（乐观更新）

- [x] 创建 `apps/web/app/components/ToggleSwitch.vue`（或直接在 KeyRow 中实现）：
  - 使用 DaisyUI toggle 组件
  - 点击后**立即翻转 UI 状态**
  - 同时发起 `PATCH /keys/:id { is_enabled: newValue }`
  - 成功：保持新状态
  - 失败：**回滚到原状态** + toast 错误提示
  - 请求中禁止重复点击（loading 态）

**AC**：点击 Toggle 后视觉立即切换；网络请求成功则保持；模拟网络错误时 Toggle 回滚到原状态并显示错误 toast。

### 3.5 创建密钥弹窗

- [x] 创建 `apps/web/app/components/KeyCreateModal.vue`：
  - 输入密钥名称（必填）
  - 提交后调用 `POST /projects/:id/keys`
  - 成功后进入"结果展示"状态：
    - 显示完整密钥明文（等宽字体，背景高亮）
    - 警告文案："请立即复制此密钥。关闭弹窗后将无法再次查看完整密钥。"
    - 复制按钮
  - 关闭弹窗后刷新密钥列表

**AC**：创建密钥后弹窗展示完整明文 + 警告文案；复制按钮可用；关闭弹窗后列表新增一行脱敏密钥。

---

## 阶段四：收尾与验收

### 4.1 错误处理与加载状态

- [x] 所有页面的 `useFetch` 调用处理 `{ data, error, status }` 三态
- [x] 加载中显示骨架屏或 spinner
- [x] API 错误通过 toast 展示具体错误消息
- [x] 密钥删除加二次确认弹窗

**AC**：后端停机时，前端显示错误提示而非空白页；所有 loading 过渡自然无闪烁。

### 4.2 空状态设计

- [x] 项目列表为空时：显示引导文案 + "Create Your First Project" 按钮
- [x] 密钥列表为空时：显示引导文案 + "Generate Your First Key" 按钮

**AC**：删除所有项目后首页显示空状态引导；进入无密钥的项目详情页显示空状态引导。

### 4.3 自检验收

按 REQUIREMENTS.md 第 8 节逐项验证：

- [x] **安全意识**：`GET /projects/:id/keys` 响应中无 `key_value` 字段，仅有 `key_value_masked`
- [x] **交互体验**：Toggle 点击后立即更新视觉状态 + 触发 PATCH 请求；失败回滚
- [x] **级联关系**：删除项目后，直接查询 SQLite 确认对应 api_keys 已被删除
- [x] **密钥生成**：检查 `key_value` 格式为 `sk-` + 64 位 hex，长度 67 字符
- [x] **一键复制**：复制到剪贴板的是完整明文
- [x] **创建提醒**：创建密钥后弹窗展示 "仅显示一次" 警告
- [x] **CORS 配置**：前端 `localhost:3000` 可正常调用后端 `localhost:7080` 的所有接口
- [x] **错误处理**：API 返回 `{"error": "..."}` 格式；前端展示错误 toast

### 4.4 代码质量

- [x] `pnpm lint` 全部通过（前端 ESLint + 后端 go vet）
- [x] 无 TypeScript 类型错误
- [x] 无控制台 warning/error
