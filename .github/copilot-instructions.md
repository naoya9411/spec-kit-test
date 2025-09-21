```instructions
## General Rules

- Prefer to respond in Japanese when interacting with the user
- マークダウンファイルのドキュメント系やコードコメントは日本語で書くこと

## Current Project: TODOアプリ (Branch: 001-todo)

### Tech Stack
- **Language/Version**: TypeScript 5.x (frontend), Go 1.21+ (backend)
- **Primary Dependencies**: React 18, TailwindCSS, shadcn/ui, Go Gin, GORM
- **Storage**: PostgreSQL 15+
- **Testing**: Vitest (frontend unit), Playwright (E2E), Testify (Go), Testcontainers
- **Project Type**: web (frontend + backend分離)

### Architecture Patterns
- Domain-Driven Design (DDD)
- Onion Architecture (依存性が内向き、ドメイン中心)
- Feature-based Components (React機能ベース構成)
- Microservices with HTTP REST
- Test-Driven Development (RED-GREEN-Refactor)

### Current Phase
- Phase 1: Design complete (data-model.md, contracts/, quickstart.md ready)
- Next: Phase 2 task generation via /tasks command

### Key Files
- `/specs/001-todo/spec.md` - 機能仕様書
- `/specs/001-todo/plan.md` - 実装計画
- `/specs/001-todo/data-model.md` - データモデル設計
- `/specs/001-todo/contracts/api-spec.yaml` - OpenAPI仕様
- `/specs/001-todo/quickstart.md` - 動作確認手順

### Recent Changes
- 001-todo: Adopted Feature-based Components for React frontend architecture
- Onion Architecture for backend + Feature-based for frontend alignment
- API contracts and data models updated with feature-driven approach

Last updated: 2025-09-21
```
