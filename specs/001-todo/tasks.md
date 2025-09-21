# Tasks: シンプルなTODOアプリ

**Input**: Design documents from `/specs/001-todo/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: TypeScript/React/TailwindCSS + Go + PostgreSQL
   → Structure: Web app (frontend + backend分離)
   → Architecture: Onion Architecture + Feature-based Components
2. Load design documents:
   → data-model.md: TODO entity, Onion layers → model tasks
   → contracts/api-spec.yaml: 4 endpoints (GET, POST, PUT, DELETE) → contract test tasks
   → quickstart.md: 4 scenarios → integration test tasks
3. Generate tasks by category:
   → Setup: Go/TypeScript projects, dependencies, Docker
   → Tests: contract tests, integration tests (TDD required)
   → Core: domain models, services, handlers, React components
   → Integration: DB, API layer, React Query
   → Polish: E2E tests, Docker deployment
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Tests before implementation (TDD)
   → Domain → Application → Infrastructure → Interface layers
5. Number tasks sequentially (T001-T040)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- File paths based on plan.md structure (backend/, frontend/)

## Phase 3.1: Project Setup
- [ ] T001 Create backend project structure with Go modules (go.mod, internal/, cmd/, pkg/)
- [ ] T002 Create frontend project structure with Vite + TypeScript + TailwindCSS
- [ ] T003 [P] Configure backend linting (golangci-lint) and formatting (gofmt)
- [ ] T004 [P] Configure frontend linting (ESLint) and formatting (Prettier)
- [ ] T005 Install backend dependencies: Gin, GORM, PostgreSQL driver, Testify, Testcontainers

## Phase 3.2: Contract Tests (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T006 [P] Contract test GET /api/v1/todos in backend/tests/contract/test_todos_get.go
- [ ] T007 [P] Contract test POST /api/v1/todos in backend/tests/contract/test_todos_post.go  
- [ ] T008 [P] Contract test PUT /api/v1/todos/{id} in backend/tests/contract/test_todos_put.go
- [ ] T009 [P] Contract test DELETE /api/v1/todos/{id} in backend/tests/contract/test_todos_delete.go

## Phase 3.3: Domain Layer (Backend - Core)
- [ ] T010 [P] TODO entity in backend/internal/domain/todo/entity.go
- [ ] T011 [P] Value objects (TodoID, Title) in backend/internal/domain/todo/value_objects.go
- [ ] T012 [P] Repository interface in backend/internal/domain/todo/repository.go
- [ ] T013 [P] Domain service interface in backend/internal/domain/todo/service.go

## Phase 3.4: Application Layer (Backend - Use Cases)
- [ ] T014 [P] Create TODO use case in backend/internal/application/todo/create_todo.go
- [ ] T015 [P] Get TODOs use case in backend/internal/application/todo/get_todos.go
- [ ] T016 [P] Update TODO use case in backend/internal/application/todo/update_todo.go
- [ ] T017 [P] Delete TODO use case in backend/internal/application/todo/delete_todo.go

## Phase 3.5: Infrastructure Layer (Backend - Database)
- [ ] T018 TODO repository implementation with GORM in backend/internal/infrastructure/persistence/todo_repository.go
- [ ] T019 Database configuration and migrations in backend/internal/infrastructure/database/config.go
- [ ] T020 Database connection setup in backend/cmd/api/main.go

## Phase 3.6: Interface Layer (Backend - API)
- [ ] T021 TODO DTO structs in backend/internal/interface/dto/todo_dto.go
- [ ] T022 GET /api/v1/todos handler in backend/internal/interface/handler/todo_handler.go
- [ ] T023 POST /api/v1/todos handler (same file as T022)
- [ ] T024 PUT /api/v1/todos/{id} handler (same file as T022)
- [ ] T025 DELETE /api/v1/todos/{id} handler (same file as T022)
- [ ] T026 Error handling middleware in backend/internal/interface/middleware/error_handler.go
- [ ] T027 CORS middleware in backend/internal/interface/middleware/cors.go
- [ ] T028 Router setup with Gin in backend/internal/interface/router/router.go

## Phase 3.7: Frontend Shared Components
- [ ] T029 [P] shadcn/ui setup and base components in frontend/src/shared/components/ui/
- [ ] T030 [P] Layout components in frontend/src/shared/components/layout/
- [ ] T031 [P] React Query provider in frontend/src/app/providers.tsx

## Phase 3.8: Frontend TODO Feature
- [ ] T032 [P] TODO types in frontend/src/features/todo/types/todo.types.ts
- [ ] T033 [P] TODO API service in frontend/src/features/todo/services/todoApi.ts
- [ ] T034 [P] useTodos hook with React Query in frontend/src/features/todo/hooks/useTodos.ts
- [ ] T035 [P] useTodoMutation hook in frontend/src/features/todo/hooks/useTodoMutation.ts
- [ ] T036 [P] TodoForm component in frontend/src/features/todo/components/TodoForm.tsx
- [ ] T037 [P] TodoItem component in frontend/src/features/todo/components/TodoItem.tsx
- [ ] T038 [P] TodoList component in frontend/src/features/todo/components/TodoList.tsx

## Phase 3.9: Integration Tests
- [ ] T039 [P] Backend integration test with Testcontainers in backend/tests/integration/todo_integration_test.go
- [ ] T040 [P] Frontend integration test with MSW in frontend/src/features/todo/__tests__/integration/todo-integration.test.tsx

## Phase 3.10: Unit Tests
- [ ] T041 [P] Domain entity unit tests in backend/tests/unit/domain/todo_test.go
- [ ] T042 [P] Use case unit tests in backend/tests/unit/application/todo_use_case_test.go
- [ ] T043 [P] Frontend component unit tests in frontend/src/features/todo/__tests__/components/
- [ ] T044 [P] Frontend hook unit tests in frontend/src/features/todo/__tests__/hooks/

## Phase 3.11: E2E Tests & Deployment
- [ ] T045 Playwright E2E test covering quickstart scenarios in frontend/tests/e2e/todo.spec.ts
- [ ] T046 Backend Dockerfile in docker/backend.Dockerfile
- [ ] T047 Frontend Dockerfile in docker/frontend.Dockerfile
- [ ] T048 Docker Compose configuration in docker-compose.yml

## Dependencies
```
Setup Phase (T001-T005) 
    ↓
Contract Tests (T006-T009) - Must fail before implementation
    ↓
Domain Layer (T010-T013)
    ↓  
Application Layer (T014-T017) - Depends on Domain
    ↓
Infrastructure (T018-T020) - Depends on Domain interfaces
    ↓
Interface Layer (T021-T028) - Depends on Application & Infrastructure
    ↓
Frontend Development (T029-T038) - Can run parallel with backend after API contracts
    ↓
Integration & Unit Tests (T039-T044)
    ↓
E2E & Deployment (T045-T048)
```

## Parallel Execution Examples

### Phase 3.2: Contract Tests (All parallel)
```bash
# Launch T006-T009 together:
Task: "Contract test GET /api/v1/todos in backend/tests/contract/test_todos_get.go"
Task: "Contract test POST /api/v1/todos in backend/tests/contract/test_todos_post.go"
Task: "Contract test PUT /api/v1/todos/{id} in backend/tests/contract/test_todos_put.go"
Task: "Contract test DELETE /api/v1/todos/{id} in backend/tests/contract/test_todos_delete.go"
```

### Phase 3.3: Domain Layer (All parallel)
```bash
# Launch T010-T013 together:
Task: "TODO entity in backend/internal/domain/todo/entity.go"
Task: "Value objects in backend/internal/domain/todo/value_objects.go"
Task: "Repository interface in backend/internal/domain/todo/repository.go"  
Task: "Domain service interface in backend/internal/domain/todo/service.go"
```

### Phase 3.8: Frontend Feature (Components parallel)
```bash
# Launch T036-T038 together:
Task: "TodoForm component in frontend/src/features/todo/components/TodoForm.tsx"
Task: "TodoItem component in frontend/src/features/todo/components/TodoItem.tsx"
Task: "TodoList component in frontend/src/features/todo/components/TodoList.tsx"
```

## Task Generation Rules Applied
1. **From Contracts (api-spec.yaml)**:
   - 4 endpoints → 4 contract test tasks [P] (T006-T009)
   - 4 endpoints → 4 handler implementation tasks (T022-T025, same file)

2. **From Data Model (data-model.md)**:
   - TODO entity → entity task [P] (T010)
   - Value objects → value objects task [P] (T011)
   - Repository pattern → interface + implementation (T012, T018)
   - Use cases → application layer tasks [P] (T014-T017)

3. **From Quickstart (quickstart.md)**:
   - 4 scenarios → integration test coverage (T039-T040)
   - Manual verification → E2E test (T045)

4. **Architecture (Onion + Feature-based)**:
   - Backend: Domain → Application → Infrastructure → Interface
   - Frontend: Feature-based structure with __tests__ proximity
   - TDD: Tests before implementation enforced

## Validation Checklist
- [x] All 4 contracts have corresponding tests (T006-T009)
- [x] TODO entity has model task (T010)
- [x] All tests come before implementation (T006-T009 before T022-T025)
- [x] Parallel tasks truly independent (different files)
- [x] Each task specifies exact file path
- [x] No [P] task modifies same file (handlers T022-T025 not marked [P])
- [x] TDD order enforced (Contract tests → Domain → Application → etc.)
- [x] Feature-based frontend structure with __tests__ directories
- [x] Integration tests use Testcontainers (backend) and MSW (frontend)

## Notes
- Contract tests (T006-T009) **MUST FAIL** before any implementation
- Backend handlers (T022-T025) in same file, so not marked [P]
- Frontend components follow Feature-based architecture
- Integration tests use real dependencies (PostgreSQL via Testcontainers)
- Follow TDD cycle: RED (failing tests) → GREEN (implementation) → REFACTOR
- Commit after each task completion
- Run quickstart.md scenarios for final validation
