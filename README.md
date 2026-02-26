# Dev Portal

开发者 API 密钥管理控制台，模拟 AWS / OpenAI 风格的密钥管理后台。支持项目管理、密钥生成与脱敏展示、启用/禁用切换等核心功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 包管理 | pnpm + Turborepo (Monorepo) |
| 后端 | Go 1.22+ / net/http / SQLite (`modernc.org/sqlite`) |
| 前端 | Nuxt 4 / TypeScript / VueUse / Tailwind CSS / DaisyUI |

## 快速开始

```bash
# 安装依赖
pnpm install

# 启动开发环境（前后端并行）
pnpm dev

# 或分别启动
pnpm --filter api dev   # Go API → http://localhost:7080
pnpm --filter web dev   # Nuxt  → http://localhost:3000
```

首次启动后端时会自动创建 `dev.db` 并写入预置数据（2 个项目 + 7 个密钥）。

## 项目结构

```
dev-portal/
├── apps/
│   ├── api/                  # Go 后端
│   │   ├── cmd/api/          # main.go 入口
│   │   └── internal/         # handler / db / model / middleware
│   └── web/                  # Nuxt 4 前端
│       └── app/
│           ├── pages/        # index.vue, projects/[id].vue
│           ├── components/   # ProjectCard, KeyTable, KeyRow, ...
│           ├── composables/  # useProjects, useApiKeys, useToast
│           └── types/        # TypeScript 接口定义
├── turbo.json
├── pnpm-workspace.yaml
└── package.json
```

## API 概览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/projects` | 项目列表（含密钥数量） |
| POST | `/projects` | 创建项目 |
| GET | `/projects/:id` | 项目详情 |
| PATCH | `/projects/:id` | 更新项目 |
| DELETE | `/projects/:id` | 删除项目（级联删除密钥） |
| GET | `/projects/:id/keys` | 密钥列表（脱敏） |
| POST | `/projects/:id/keys` | 创建密钥（返回明文，仅一次） |
| PATCH | `/keys/:id` | 切换密钥启用状态 |
| DELETE | `/keys/:id` | 删除密钥 |
| GET | `/keys/:id/reveal` | 获取密钥明文 |

## 安全设计

- **后端脱敏**：列表接口仅返回 `key_value_masked`（如 `sk-a1b2******a1b2`），不含明文
- **按需 Reveal**：明文通过独立的 `/keys/:id/reveal` 接口单独获取
- **一次性展示**：创建密钥时弹窗警告"此密钥仅显示一次，请妥善保存"
- **加密生成**：使用 `crypto/rand` 生成 256-bit 随机密钥（`sk-` + 64 位 hex）

## 常用命令

```bash
pnpm dev          # 并行启动前后端开发服务
pnpm build        # 构建前后端
pnpm lint         # ESLint (前端) + go vet (后端)
```

## 相关文档

- [需求文档](./REQUIREMENTS.md)
- [设计文档](./design.md)
- [开发文档](./DEVELOPMENT.md)
- [任务清单](./TASKS.md)
