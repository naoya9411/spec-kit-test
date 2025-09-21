# Implementation Plan: シンプルなTODOアプリ

**Branch**: `001-todo` | **Date**: 2025年9月21日 | **Spec**: [./spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-todo/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
4. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
5. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, or `GEMINI.md` for Gemini CLI).
6. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
7. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
8. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
シンプルなTODOアプリケーションで、ユーザーがタスクの作成、確認、完了、削除を行えるフルスタックWebアプリケーション。技術スタックはTypeScript/React/TailwindCSS + Go + PostgreSQL（またはMongoDB）を使用し、マイクロサービスアーキテクチャとDDDを採用。フロントエンド・バックエンド分離型のWebアプリケーションとして構築。

## Technical Context
**Language/Version**: TypeScript 5.x (frontend), Go 1.21+ (backend)
**Primary Dependencies**: React 18, TailwindCSS, shadcn/ui, Go Gin/Fiber, GORM/MongoDB Driver  
**Storage**: PostgreSQL 15+ または MongoDB 7.0+（SQLiteより堅牢でスケーラブル）  
**Testing**: 
  - Frontend: Vitest (unit), Playwright (E2E)
  - Backend: Testify, Go標準テスト
  - Integration: Testcontainers
**Target Platform**: Webブラウザー、Dockerコンテナ環境
**Project Type**: web (frontend + backend分離)  
**Performance Goals**: 100ms未満のAPI応答時間、60fpsのUI操作性  
**Constraints**: マイクロサービス間通信、コンテナ環境での動作保証、DDD原則遵守  
**Scale/Scope**: 1000ユーザー、10万TODOアイテム対応可能

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Simplicity**:
- Projects: 3 (frontend, backend-api, backend-todo-service) - マイクロサービス要件により必要
- Using framework directly? Yes (React直接使用、Go標準ライブラリ重視)
- Single data model? Yes (TODO entityを中心とした単一ドメイン)
- Avoiding patterns? Yes (過度な抽象化を避け、必要最小限のパターンのみ)

**Architecture**:
- EVERY feature as library? Yes (frontend/backend共にライブラリ化)
- Libraries listed: 
  - todo-domain (Go): ドメインロジック・エンティティ
  - todo-api (Go): REST API・インターフェース層
  - todo-ui (TypeScript): UIコンポーネント・サービス層
- CLI per library: Yes (開発用CLIコマンド提供)
- Library docs: llms.txt format planned? Yes

**Testing (NON-NEGOTIABLE)**:
- RED-GREEN-Refactor cycle enforced? Yes (厳格に遵守)
- Git commits show tests before implementation? Yes (テストファーストコミット)
- Order: Contract→Integration→E2E→Unit strictly followed? Yes (順守)
- Real dependencies used? Yes (TestcontainersでPostgreSQL/MongoDB使用)
- Integration tests for: new libraries, contract changes, shared schemas? Yes
- FORBIDDEN: Implementation before test, skipping RED phase

**Observability**:
- Structured logging included? Yes (JSON形式ログ出力)
- Frontend logs → backend? Yes (統一ログストリーム)
- Error context sufficient? Yes (エラートレーシング)

**Versioning**:
- Version number assigned? Yes (1.0.0 - MAJOR.MINOR.BUILD)
- BUILD increments on every change? Yes (各変更でビルド番号増分)
- Breaking changes handled? Yes (API互換性維持、移行計画策定)

## Project Structure

### Documentation (this feature)
```
specs/001-todo/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (frontend + backend detected)
backend/
├── cmd/
│   └── api/              # API サーバーエントリーポイント・DI
├── internal/
│   ├── domain/           # ドメイン層（中心）- エンティティ・値オブジェクト・リポジトリIF
│   ├── application/      # アプリケーション層 - ユースケース・オーケストレーション
│   └── infrastructure/   # インフラ層（外側）- リポジトリ実装・DB・外部API
├── pkg/                  # 共有ライブラリ
└── tests/
    ├── unit/             # ドメイン・アプリケーション層テスト
    ├── integration/      # インフラ層結合テスト
    └── e2e/              # API全体テスト

frontend/
├── src/
│   ├── features/             # Feature-based コンポーネント構成
│   │   └── todo/
│   │       ├── components/   # TODO機能のコンポーネント
│   │       │   ├── TodoList.tsx
│   │       │   ├── TodoItem.tsx
│   │       │   ├── TodoForm.tsx
│   │       │   └── index.ts
│   │       ├── hooks/        # カスタムフック（状態管理）
│   │       │   ├── useTodos.ts
│   │       │   └── useTodoMutation.ts
│   │       ├── services/     # API通信サービス
│   │       │   └── todoApi.ts
│   │       ├── types/        # 機能固有の型定義
│   │       │   └── todo.types.ts
│   │       ├── __tests__/    # 機能単位のテスト（近接性重視）
│   │       │   ├── components/
│   │       │   │   ├── TodoList.test.tsx
│   │       │   │   ├── TodoItem.test.tsx
│   │       │   │   └── TodoForm.test.tsx
│   │       │   ├── hooks/
│   │       │   │   ├── useTodos.test.ts
│   │       │   │   └── useTodoMutation.test.ts
│   │       │   └── services/
│   │       │       └── todoApi.test.ts
│   │       └── index.ts      # 機能エクスポート
│   ├── shared/               # 共通コンポーネント・ユーティリティ
│   │   ├── components/
│   │   │   ├── ui/           # shadcn/ui ベースコンポーネント
│   │   │   └── layout/
│   │   ├── hooks/
│   │   ├── services/
│   │   ├── types/
│   │   └── __tests__/        # 共通機能のテスト
│   │       ├── components/
│   │       ├── hooks/
│   │       └── services/
│   ├── app/                  # アプリケーション設定
│   │   ├── App.tsx
│   │   ├── router.tsx
│   │   ├── providers.tsx
│   │   └── __tests__/        # アプリケーション層テスト
│   │       ├── App.test.tsx
│   │       └── router.test.tsx
│   └── lib/                  # ユーティリティライブラリ
├── tests/
│   └── e2e/                  # Playwright E2Eテスト
│       └── todo.spec.ts
├── playwright.config.ts
└── vite.config.ts

docker/
├── backend.Dockerfile
├── frontend.Dockerfile
└── docker-compose.yml

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]
```

**Structure Decision**: Option 2 (Web application) - フロントエンドとバックエンドの分離

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - データベース選択（PostgreSQL vs MongoDB）
   - マイクロサービス通信方式（HTTP REST vs gRPC）
   - フロントエンドテストライブラリの最適化（Vitest vs Jest）
   - コンテナオーケストレーション方式（Docker Compose vs Kubernetes）
   - DDD実装パターン（Go言語でのベストプラクティス）

2. **Generate and dispatch research agents**:
   ```
   Task: "Research PostgreSQL vs MongoDB for TODO app scalability and DDD compatibility"
   Task: "Find best practices for Go microservices with DDD architecture"
   Task: "Evaluate Vitest vs Jest for TypeScript React testing in 2025"
   Task: "Research Docker multi-stage builds for Go + TypeScript optimization"
   Task: "Compare HTTP REST vs gRPC for microservice communication"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `/scripts/update-agent-context.sh [claude|gemini|copilot]` for your AI assistant
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each API endpoint → contract test task [P]
- Each domain entity → model creation task [P] 
- Each feature → frontend components task [P]
- Each layer (domain/application/infrastructure/interface) → implementation tasks
- Docker containerization → deployment tasks

**Task Categories**:
1. **Contract Tests** (4 tasks): GET, POST, PUT, DELETE endpoints
2. **Domain Layer** (3 tasks): TODO entity, value objects, domain services  
3. **Application Layer** (2 tasks): Use cases, application services
4. **Infrastructure Layer** (3 tasks): Repository implementation, database setup, migrations
5. **Interface Layer** (4 tasks): HTTP handlers, DTOs, middleware, error handling
6. **Todo Feature Frontend** (6 tasks): 
   - Todo types & API service + unit tests
   - Custom hooks (useTodos, useTodoMutation) + unit tests
   - Components (TodoList, TodoItem, TodoForm) + unit tests
   - Feature integration & routing
7. **Shared Frontend** (3 tasks): UI components, layout, providers + tests
8. **Integration Tests** (6 tasks): API integration, database integration, E2E scenarios
9. **Docker/Deployment** (3 tasks): Dockerfiles, docker-compose, environment configuration

**Ordering Strategy**:
- **TDD order**: Contract tests → Domain tests → Integration tests → Implementation
- **Dependency order**: 
  1. Domain models & tests (no dependencies)
  2. Application services & tests (depends on domain)
  3. Infrastructure & tests (depends on domain)
  4. Interface layer & tests (depends on all)
  5. Frontend shared components (no dependencies) [P]
  6. Frontend todo feature (depends on API) [P]
  7. E2E tests (depends on everything)
- **Mark [P] for parallel execution**: Domain/Infrastructure tasks, Frontend feature tasks

**Frontend Task Structure**:
```
## Task N: TODO Feature - Types & API Service [P]
- Create todo.types.ts with TypeScript interfaces
- Implement todoApi.ts service class
- Add error handling and request/response validation
- Create __tests__/services/todoApi.test.ts with unit tests

## Task N+1: TODO Feature - Custom Hooks [P]  
- Implement useTodos hook with React Query
- Implement useTodoMutation hook for CRUD operations
- Add optimistic updates and error handling
- Create __tests__/hooks/ unit tests for both hooks

## Task N+2: TODO Feature - Components [P]
- Create TodoForm component with validation
- Create TodoItem component with actions
- Create TodoList component with loading states  
- Create __tests__/components/ unit tests for all components
```

**Estimated Output**: 34-38 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| 3rd project (microservices) | DDDとマイクロサービス要求により必須 | モノリス構成では将来のスケーラビリティとドメイン分離が困難 |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*