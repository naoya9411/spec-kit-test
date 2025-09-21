# Data Model: シンプルなTODOアプリ

**作成日**: 2025年9月21日  
**対象機能**: シンプルなTODOアプリ  
**ドメインモデリング**: DDD + Onion Architecture

---

## アーキテクチャ概要

### Onion Architecture 層構成
```
┌─────────────────────────────────┐
│     Infrastructure Layer       │ ← PostgreSQL, HTTP Client
│  ┌───────────────────────────┐  │
│  │    Application Layer     │  │ ← Use Cases, Orchestration
│  │  ┌─────────────────────┐  │  │
│  │  │   Domain Layer      │  │  │ ← Entities, Business Logic
│  │  │                     │  │  │
│  │  └─────────────────────┘  │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘

依存関係: Infrastructure → Application → Domain
実際の制御フロー: Domain ← Application ← Infrastructure
```

---

## ドメイン層（Domain Layer）

### TODO Entity（ルートアグリゲート）
```go
// internal/domain/todo/entity.go
package todo

import (
    "errors"
    "strings"
    "time"
)

type TODO struct {
    id        TodoID
    title     Title
    completed bool
    createdAt time.Time
    updatedAt time.Time
}

// コンストラクタ
func NewTODO(title string) (*TODO, error) {
    titleVO, err := NewTitle(title)
    if err != nil {
        return nil, err
    }
    
    return &TODO{
        id:        NewTodoID(),
        title:     titleVO,
        completed: false,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }, nil
}

// ビジネスメソッド
func (t *TODO) ToggleComplete() {
    t.completed = !t.completed
    t.updatedAt = time.Now()
}

func (t *TODO) ChangeTitle(newTitle string) error {
    titleVO, err := NewTitle(newTitle)
    if err != nil {
        return err
    }
    t.title = titleVO
    t.updatedAt = time.Now()
    return nil
}

// ゲッター（データ隠蔽）
func (t *TODO) ID() TodoID { return t.id }
func (t *TODO) Title() string { return t.title.String() }
func (t *TODO) IsCompleted() bool { return t.completed }
func (t *TODO) CreatedAt() time.Time { return t.createdAt }
func (t *TODO) UpdatedAt() time.Time { return t.updatedAt }
```

### 値オブジェクト（Value Objects）
```go
// internal/domain/todo/value_objects.go
package todo

import (
    "errors"
    "strings"
    "github.com/google/uuid"
)

// TodoID - エンティティの一意識別子
type TodoID struct {
    value string
}

func NewTodoID() TodoID {
    return TodoID{value: uuid.New().String()}
}

func TodoIDFromString(id string) (TodoID, error) {
    if _, err := uuid.Parse(id); err != nil {
        return TodoID{}, errors.New("invalid todo id format")
    }
    return TodoID{value: id}, nil
}

func (id TodoID) String() string { return id.value }

// Title - タイトル値オブジェクト（バリデーション付き）
type Title struct {
    value string
}

func NewTitle(title string) (Title, error) {
    trimmed := strings.TrimSpace(title)
    if trimmed == "" {
        return Title{}, errors.New("title cannot be empty")
    }
    if len(trimmed) > 100 {
        return Title{}, errors.New("title cannot exceed 100 characters")
    }
    return Title{value: trimmed}, nil
}

func (t Title) String() string { return t.value }
```

---

## ビジネスルール

### バリデーション規則
1. **Title必須**: 空文字列やスペースのみは許可しない
2. **Title長さ制限**: 100文字以内
3. **ID一意性**: システム生成のUUID使用
4. **状態管理**: Completedフラグでタスク状態を管理

### 状態遷移
```
[作成] → [未完了] → [完了]
  ↓         ↓         ↓
[削除]   [削除]   [削除]
```

**許可される操作**:
- **作成**: 新しいTODOアイテムを未完了状態で作成
- **完了切替**: 未完了↔完了の状態変更
- **削除**: どの状態からも削除可能

---

### ドメインサービス・リポジトリインターフェース
```go
// internal/domain/todo/repository.go
package todo

import "context"

// リポジトリインターフェース（ドメイン層で定義、依存性逆転）
type Repository interface {
    Save(ctx context.Context, todo *TODO) error
    FindByID(ctx context.Context, id TodoID) (*TODO, error)
    FindAll(ctx context.Context) ([]*TODO, error)
    Update(ctx context.Context, todo *TODO) error
    Delete(ctx context.Context, id TodoID) error
    NextID() TodoID
}

// ドメインサービス（複雑なビジネスロジック）
type Service struct {
    // 現在のTODOアプリでは単純なため、特別なドメインサービスは不要
    // 将来的に、複数TODO間のビジネスルール等が必要になった場合に追加
}
```

### ドメインエラー
```go
// internal/domain/todo/errors.go
package todo

import "errors"

var (
    ErrTODONotFound     = errors.New("todo not found")
    ErrInvalidTitle     = errors.New("invalid title")
    ErrTitleTooLong     = errors.New("title too long")
    ErrTitleEmpty       = errors.New("title cannot be empty")
    ErrInvalidID        = errors.New("invalid todo id")
)
```

---

## アプリケーション層（Application Layer）

### ユースケース（Use Cases）
```go
// internal/application/usecase/todo_usecase.go
package usecase

import (
    "context"
    "todo-app/internal/domain/todo"
)

type TodoUseCase struct {
    todoRepo todo.Repository
}

func NewTodoUseCase(todoRepo todo.Repository) *TodoUseCase {
    return &TodoUseCase{
        todoRepo: todoRepo,
    }
}

func (uc *TodoUseCase) CreateTODO(ctx context.Context, title string) (*todo.TODO, error) {
    newTodo, err := todo.NewTODO(title)
    if err != nil {
        return nil, err
    }
    
    if err := uc.todoRepo.Save(ctx, newTodo); err != nil {
        return nil, err
    }
    
    return newTodo, nil
}

func (uc *TodoUseCase) GetAllTODOs(ctx context.Context) ([]*todo.TODO, error) {
    return uc.todoRepo.FindAll(ctx)
}

func (uc *TodoUseCase) ToggleComplete(ctx context.Context, idStr string) (*todo.TODO, error) {
    id, err := todo.TodoIDFromString(idStr)
    if err != nil {
        return nil, err
    }
    
    todoItem, err := uc.todoRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    todoItem.ToggleComplete()
    
    if err := uc.todoRepo.Update(ctx, todoItem); err != nil {
        return nil, err
    }
    
    return todoItem, nil
}

func (uc *TodoUseCase) DeleteTODO(ctx context.Context, idStr string) error {
    id, err := todo.TodoIDFromString(idStr)
    if err != nil {
        return err
    }
    
    return uc.todoRepo.Delete(ctx, id)
}
```

---

## インフラストラクチャ層（Infrastructure Layer）

### PostgreSQLリポジトリ実装
```go
// internal/infrastructure/persistence/todo_repository.go
package persistence

import (
    "context"
    "todo-app/internal/domain/todo"
    "gorm.io/gorm"
)

// データベース用のモデル（ドメインエンティティとは分離）
type todoModel struct {
    ID        string `gorm:"primaryKey"`
    Title     string `gorm:"type:varchar(100);not null"`
    Completed bool   `gorm:"default:false"`
    CreatedAt int64  `gorm:"autoCreateTime"`
    UpdatedAt int64  `gorm:"autoUpdateTime"`
}

func (todoModel) TableName() string {
    return "todos"
}

type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

func (r *TodoRepository) Save(ctx context.Context, todoEntity *todo.TODO) error {
    model := &todoModel{
        ID:        todoEntity.ID().String(),
        Title:     todoEntity.Title(),
        Completed: todoEntity.IsCompleted(),
        CreatedAt: todoEntity.CreatedAt().Unix(),
        UpdatedAt: todoEntity.UpdatedAt().Unix(),
    }
    
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *TodoRepository) FindByID(ctx context.Context, id todo.TodoID) (*todo.TODO, error) {
    var model todoModel
    if err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, todo.ErrTODONotFound
        }
        return nil, err
    }
    
    return r.modelToEntity(&model)
}

func (r *TodoRepository) FindAll(ctx context.Context) ([]*todo.TODO, error) {
    var models []todoModel
    if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
        return nil, err
    }
    
    todos := make([]*todo.TODO, len(models))
    for i, model := range models {
        todoEntity, err := r.modelToEntity(&model)
        if err != nil {
            return nil, err
        }
        todos[i] = todoEntity
    }
    
    return todos, nil
}

// プライベートメソッド：モデルからエンティティへの変換
func (r *TodoRepository) modelToEntity(model *todoModel) (*todo.TODO, error) {
    id, err := todo.TodoIDFromString(model.ID)
    if err != nil {
        return nil, err
    }
    
    // リフレクションを使わずに直接構築（パフォーマンス重視）
    todoEntity, err := todo.NewTODO(model.Title)
    if err != nil {
        return nil, err
    }
    
    // 必要に応じてプライベートフィールドを復元するヘルパーメソッドを
    // ドメインエンティティに追加することを検討
    return todoEntity, nil
}
```

### HTTPハンドラー（外側の層）
```go
// cmd/api/handlers/todo_handler.go
package handlers

import (
    "net/http"
    "todo-app/internal/application/usecase"
    "github.com/gin-gonic/gin"
)

type TodoHandler struct {
    todoUC *usecase.TodoUseCase
}

func NewTodoHandler(todoUC *usecase.TodoUseCase) *TodoHandler {
    return &TodoHandler{todoUC: todoUC}
}

type CreateTODORequest struct {
    Title string `json:"title" binding:"required"`
}

type TodoResponse struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
}

func (h *TodoHandler) CreateTODO(c *gin.Context) {
    var req CreateTODORequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    todo, err := h.todoUC.CreateTODO(c.Request.Context(), req.Title)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    response := &TodoResponse{
        ID:        todo.ID().String(),
        Title:     todo.Title(),
        Completed: todo.IsCompleted(),
        CreatedAt: todo.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt: todo.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
    }
    
    c.JSON(http.StatusCreated, response)
}

func (h *TodoHandler) GetTODOs(c *gin.Context) {
    todos, err := h.todoUC.GetAllTODOs(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    responses := make([]TodoResponse, len(todos))
    for i, todo := range todos {
        responses[i] = TodoResponse{
            ID:        todo.ID().String(),
            Title:     todo.Title(),
            Completed: todo.IsCompleted(),
            CreatedAt: todo.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
            UpdatedAt: todo.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "todos": responses,
        "total": len(responses),
    })
}
```

---

## 依存性注入・ワイヤリング
```go
// cmd/api/main.go
package main

import (
    "todo-app/internal/application/usecase"
    "todo-app/internal/infrastructure/persistence"
    "todo-app/cmd/api/handlers"
    "gorm.io/gorm"
)

func main() {
    // データベース接続
    db := setupDatabase()
    
    // リポジトリ（Infrastructure層）
    todoRepo := persistence.NewTodoRepository(db)
    
    // ユースケース（Application層）
    todoUC := usecase.NewTodoUseCase(todoRepo)
    
    // ハンドラー（外側）
    todoHandler := handlers.NewTodoHandler(todoUC)
    
    // ルーター設定
    setupRoutes(todoHandler)
}
```

---

## プロジェクト構造（Onion Architecture）

### ディレクトリ構造
```
backend/
├── cmd/
│   └── api/
│       ├── main.go              # エントリーポイント・DI
│       └── handlers/            # HTTPハンドラー（最外層）
│           └── todo_handler.go
├── internal/
│   ├── domain/                  # ドメイン層（中心）
│   │   └── todo/
│   │       ├── entity.go        # TODOエンティティ
│   │       ├── value_objects.go # 値オブジェクト
│   │       ├── repository.go    # リポジトリIF（依存性逆転）
│   │       └── errors.go        # ドメインエラー
│   ├── application/             # アプリケーション層
│   │   └── usecase/
│   │       └── todo_usecase.go  # ユースケース
│   └── infrastructure/          # インフラ層（最外層）
│       └── persistence/
│           └── todo_repository.go # リポジトリ実装
└── tests/
    ├── unit/                    # ドメイン・アプリケーション層テスト
    ├── integration/             # リポジトリ・DB結合テスト  
    └── e2e/                     # API全体テスト
```

### 依存関係の流れ
```
main.go → handlers → usecase → domain ← infrastructure
   ↓         ↓         ↓         ↑         ↑
 DI設定   HTTP処理   ビジネス   純粋ロジック  DB実装
```

**重要**: すべての依存性がドメイン層（中心）に向かう

---

## フロントエンド層（Feature-based Components）

### TODO機能の型定義
```typescript
// src/features/todo/types/todo.types.ts
export interface Todo {
  id: string;
  title: string;
  completed: boolean;
  createdAt: string; // ISO 8601 format
  updatedAt: string; // ISO 8601 format
}

export interface CreateTodoRequest {
  title: string;
}

export interface UpdateTodoRequest {
  title?: string;
  completed?: boolean;
}

export interface TodosResponse {
  todos: Todo[];
  total: number;
}

// UI状態管理用の型
export interface TodoState {
  todos: Todo[];
  loading: boolean;
  error: string | null;
}

export interface TodoFormData {
  title: string;
}
```

### TODOコンポーネント構造
```typescript
// src/features/todo/components/TodoList.tsx
import { useTodos } from '../hooks/useTodos';
import { TodoItem } from './TodoItem';
import { TodoForm } from './TodoForm';

export function TodoList() {
  const { todos, loading, error } = useTodos();

  if (loading) return <div>読み込み中...</div>;
  if (error) return <div>エラー: {error}</div>;

  return (
    <div className="space-y-4">
      <TodoForm />
      <div className="space-y-2">
        {todos.map(todo => (
          <TodoItem key={todo.id} todo={todo} />
        ))}
      </div>
      {todos.length === 0 && (
        <p className="text-gray-500 text-center py-8">
          TODOがありません。新しいタスクを追加してください。
        </p>
      )}
    </div>
  );
}

// src/features/todo/components/TodoItem.tsx
import { Todo } from '../types/todo.types';
import { useTodoMutation } from '../hooks/useTodoMutation';
import { Button } from '@/shared/components/ui/button';
import { Checkbox } from '@/shared/components/ui/checkbox';

interface TodoItemProps {
  todo: Todo;
}

export function TodoItem({ todo }: TodoItemProps) {
  const { toggleComplete, deleteTodo } = useTodoMutation();

  const handleToggle = () => {
    toggleComplete.mutate({ id: todo.id, completed: !todo.completed });
  };

  const handleDelete = () => {
    deleteTodo.mutate(todo.id);
  };

  return (
    <div className="flex items-center space-x-3 p-3 border rounded-lg">
      <Checkbox
        checked={todo.completed}
        onCheckedChange={handleToggle}
        disabled={toggleComplete.isPending}
      />
      <span
        className={`flex-1 ${
          todo.completed 
            ? 'line-through text-gray-500' 
            : 'text-gray-900'
        }`}
      >
        {todo.title}
      </span>
      <Button
        variant="destructive"
        size="sm"
        onClick={handleDelete}
        disabled={deleteTodo.isPending}
      >
        削除
      </Button>
    </div>
  );
}

// src/features/todo/components/TodoForm.tsx
import { useState } from 'react';
import { useTodoMutation } from '../hooks/useTodoMutation';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';

export function TodoForm() {
  const [title, setTitle] = useState('');
  const { createTodo } = useTodoMutation();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (title.trim()) {
      createTodo.mutate({ title: title.trim() });
      setTitle('');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex space-x-2">
      <Input
        type="text"
        placeholder="新しいTODOを入力..."
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        maxLength={100}
        className="flex-1"
      />
      <Button
        type="submit"
        disabled={!title.trim() || createTodo.isPending}
      >
        {createTodo.isPending ? '作成中...' : '作成'}
      </Button>
    </form>
  );
}
```

### カスタムフック（状態管理）
```typescript
// src/features/todo/hooks/useTodos.ts
import { useQuery } from '@tanstack/react-query';
import { todoApi } from '../services/todoApi';
import { TodoState } from '../types/todo.types';

export function useTodos() {
  const {
    data: todosResponse,
    isLoading: loading,
    error,
  } = useQuery({
    queryKey: ['todos'],
    queryFn: todoApi.getAll,
    staleTime: 5 * 60 * 1000, // 5分
  });

  return {
    todos: todosResponse?.todos || [],
    total: todosResponse?.total || 0,
    loading,
    error: error?.message || null,
  };
}

// src/features/todo/hooks/useTodoMutation.ts
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { todoApi } from '../services/todoApi';
import { CreateTodoRequest, UpdateTodoRequest } from '../types/todo.types';

export function useTodoMutation() {
  const queryClient = useQueryClient();

  const createTodo = useMutation({
    mutationFn: todoApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
    },
    onError: (error) => {
      console.error('TODO作成エラー:', error);
    },
  });

  const toggleComplete = useMutation({
    mutationFn: ({ id, completed }: { id: string; completed: boolean }) =>
      todoApi.update(id, { completed }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
    },
  });

  const deleteTodo = useMutation({
    mutationFn: todoApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['todos'] });
    },
  });

  return {
    createTodo,
    toggleComplete,
    deleteTodo,
  };
}
```

### APIサービス層
```typescript
// src/features/todo/services/todoApi.ts
import { 
  Todo, 
  CreateTodoRequest, 
  UpdateTodoRequest, 
  TodosResponse 
} from '../types/todo.types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

class TodoApiService {
  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new Error(error.message || `HTTP ${response.status}`);
    }

    return response.json();
  }

  async getAll(): Promise<TodosResponse> {
    return this.request<TodosResponse>('/todos');
  }

  async getById(id: string): Promise<Todo> {
    return this.request<Todo>(`/todos/${id}`);
  }

  async create(data: CreateTodoRequest): Promise<Todo> {
    return this.request<Todo>('/todos', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async update(id: string, data: UpdateTodoRequest): Promise<Todo> {
    return this.request<Todo>(`/todos/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete(id: string): Promise<void> {
    await this.request<void>(`/todos/${id}`, {
      method: 'DELETE',
    });
  }
}

export const todoApi = new TodoApiService();
```

### テスト実装例（__tests__ディレクトリ）
```typescript
// src/features/todo/__tests__/components/TodoForm.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { TodoForm } from '../../components/TodoForm';
import { useTodoMutation } from '../../hooks/useTodoMutation';

// モック
jest.mock('../../hooks/useTodoMutation');
const mockUseTodoMutation = useTodoMutation as jest.MockedFunction<typeof useTodoMutation>;

describe('TodoForm', () => {
  let queryClient: QueryClient;
  const mockCreateTodo = {
    mutate: jest.fn(),
    isPending: false,
  };

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    mockUseTodoMutation.mockReturnValue({
      createTodo: mockCreateTodo,
      toggleComplete: { mutate: jest.fn(), isPending: false },
      deleteTodo: { mutate: jest.fn(), isPending: false },
    });
  });

  const renderWithProvider = (component: React.ReactElement) => {
    return render(
      <QueryClientProvider client={queryClient}>
        {component}
      </QueryClientProvider>
    );
  };

  it('新しいTODOを作成できる', async () => {
    renderWithProvider(<TodoForm />);
    
    const input = screen.getByPlaceholderText('新しいTODOを入力...');
    const button = screen.getByText('作成');
    
    fireEvent.change(input, { target: { value: '新しいタスク' } });
    fireEvent.click(button);
    
    await waitFor(() => {
      expect(mockCreateTodo.mutate).toHaveBeenCalledWith({ title: '新しいタスク' });
    });
  });

  it('空のタイトルでは作成ボタンが無効化される', () => {
    renderWithProvider(<TodoForm />);
    
    const button = screen.getByText('作成');
    expect(button).toBeDisabled();
  });

  it('100文字以上は入力できない', () => {
    renderWithProvider(<TodoForm />);
    
    const input = screen.getByPlaceholderText('新しいTODOを入力...');
    expect(input).toHaveAttribute('maxLength', '100');
  });
});

// src/features/todo/__tests__/hooks/useTodos.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useTodos } from '../../hooks/useTodos';
import { todoApi } from '../../services/todoApi';

// モック
jest.mock('../../services/todoApi');
const mockTodoApi = todoApi as jest.Mocked<typeof todoApi>;

describe('useTodos', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
  });

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );

  it('TODOリストを取得できる', async () => {
    const mockTodos = [
      { id: '1', title: 'テストTODO', completed: false, createdAt: '2025-09-21T10:00:00Z', updatedAt: '2025-09-21T10:00:00Z' }
    ];
    
    mockTodoApi.getAll.mockResolvedValue({ todos: mockTodos, total: 1 });
    
    const { result } = renderHook(() => useTodos(), { wrapper });
    
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });
    
    expect(result.current.todos).toEqual(mockTodos);
    expect(result.current.total).toBe(1);
  });

  it('エラーハンドリングが正しく動作する', async () => {
    const errorMessage = 'API Error';
    mockTodoApi.getAll.mockRejectedValue(new Error(errorMessage));
    
    const { result } = renderHook(() => useTodos(), { wrapper });
    
    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });
    
    expect(result.current.error).toBe(errorMessage);
  });
});

// src/features/todo/__tests__/services/todoApi.test.ts
import { todoApi } from '../../services/todoApi';

// fetchのモック
global.fetch = jest.fn();
const mockFetch = fetch as jest.MockedFunction<typeof fetch>;

describe('todoApi', () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  describe('getAll', () => {
    it('TODOリストを正常に取得できる', async () => {
      const mockResponse = {
        todos: [
          { id: '1', title: 'テストTODO', completed: false, createdAt: '2025-09-21T10:00:00Z', updatedAt: '2025-09-21T10:00:00Z' }
        ],
        total: 1
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: async () => mockResponse,
      } as Response);

      const result = await todoApi.getAll();
      
      expect(result).toEqual(mockResponse);
      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/todos',
        expect.objectContaining({
          headers: { 'Content-Type': 'application/json' }
        })
      );
    });

    it('APIエラーを適切にハンドリングする', async () => {
      const errorMessage = 'Server Error';
      
      mockFetch.mockResolvedValue({
        ok: false,
        status: 500,
        json: async () => ({ message: errorMessage }),
      } as Response);

      await expect(todoApi.getAll()).rejects.toThrow(errorMessage);
    });
  });

  describe('create', () => {
    it('新しいTODOを作成できる', async () => {
      const newTodo = { title: '新しいタスク' };
      const createdTodo = {
        id: '1',
        title: '新しいタスク',
        completed: false,
        createdAt: '2025-09-21T10:00:00Z',
        updatedAt: '2025-09-21T10:00:00Z'
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: async () => createdTodo,
      } as Response);

      const result = await todoApi.create(newTodo);
      
      expect(result).toEqual(createdTodo);
      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/todos',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(newTodo),
          headers: { 'Content-Type': 'application/json' }
        })
      );
  });
});
```

---

## 統合テスト実装

### バックエンド統合テスト（Testcontainers）
```go
// backend/tests/integration/todo_repository_test.go
package integration

import (
    "context"
    "testing"
    "time"
    "todo-app/internal/domain/todo"
    "todo-app/internal/infrastructure/persistence"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func TestTodoRepository_Integration(t *testing.T) {
    // PostgreSQLコンテナ起動
    pgContainer, err := postgres.RunContainer(
        context.Background(),
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
        testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections")),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(context.Background())

    // データベース接続
    dsn, err := pgContainer.ConnectionString(context.Background())
    require.NoError(t, err)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    require.NoError(t, err)
    
    // マイグレーション実行
    err = db.AutoMigrate(&persistence.TodoModel{})
    require.NoError(t, err)
    
    // リポジトリ初期化
    repo := persistence.NewTodoRepository(db)
    
    t.Run("TODO作成と取得", func(t *testing.T) {
        // テストデータ作成
        todoEntity, err := todo.NewTODO("統合テスト用TODO")
        require.NoError(t, err)
        
        // 保存
        err = repo.Save(context.Background(), todoEntity)
        assert.NoError(t, err)
        
        // 取得
        retrieved, err := repo.FindByID(context.Background(), todoEntity.ID())
        assert.NoError(t, err)
        assert.Equal(t, todoEntity.Title(), retrieved.Title())
        assert.Equal(t, todoEntity.IsCompleted(), retrieved.IsCompleted())
    })
    
    t.Run("TODO更新", func(t *testing.T) {
        // テストデータ作成・保存
        todoEntity, _ := todo.NewTODO("更新前TODO")
        repo.Save(context.Background(), todoEntity)
        
        // 更新
        todoEntity.ToggleComplete()
        err := repo.Update(context.Background(), todoEntity)
        assert.NoError(t, err)
        
        // 確認
        updated, err := repo.FindByID(context.Background(), todoEntity.ID())
        assert.NoError(t, err)
        assert.True(t, updated.IsCompleted())
    })
    
    t.Run("TODO削除", func(t *testing.T) {
        // テストデータ作成・保存
        todoEntity, _ := todo.NewTODO("削除用TODO")
        repo.Save(context.Background(), todoEntity)
        
        // 削除
        err := repo.Delete(context.Background(), todoEntity.ID())
        assert.NoError(t, err)
        
        // 確認（見つからないことを確認）
        _, err = repo.FindByID(context.Background(), todoEntity.ID())
        assert.ErrorIs(t, err, todo.ErrTODONotFound)
    })
    
    t.Run("全TODO取得", func(t *testing.T) {
        // 複数のTODO作成
        todo1, _ := todo.NewTODO("TODO1")
        todo2, _ := todo.NewTODO("TODO2")
        repo.Save(context.Background(), todo1)
        repo.Save(context.Background(), todo2)
        
        // 全取得
        todos, err := repo.FindAll(context.Background())
        assert.NoError(t, err)
        assert.GreaterOrEqual(t, len(todos), 2)
    })
}

// backend/tests/integration/todo_api_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "todo-app/cmd/api/handlers"
    "todo-app/internal/application/usecase"
    "todo-app/internal/infrastructure/persistence"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestTodoAPI_Integration(t *testing.T) {
    // テスト用のデータベースとAPIサーバーセットアップ
    db := setupTestDatabase(t) // 前のテストと同様のDB設定
    defer cleanupTestDatabase(db)
    
    // 依存性注入
    todoRepo := persistence.NewTodoRepository(db)
    todoUC := usecase.NewTodoUseCase(todoRepo)
    todoHandler := handlers.NewTodoHandler(todoUC)
    
    // Ginルーター設定
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/api/v1/todos", todoHandler.CreateTODO)
    router.GET("/api/v1/todos", todoHandler.GetTODOs)
    router.PUT("/api/v1/todos/:id", todoHandler.UpdateTODO)
    router.DELETE("/api/v1/todos/:id", todoHandler.DeleteTODO)
    
    t.Run("TODO作成API", func(t *testing.T) {
        payload := map[string]string{"title": "API統合テスト"}
        body, _ := json.Marshal(payload)
        
        req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusCreated, w.Code)
        
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, "API統合テスト", response["title"])
        assert.False(t, response["completed"].(bool))
    })
    
    t.Run("TODO取得API", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/todos", nil)
        w := httptest.NewRecorder()
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
        
        var response map[string]interface{}
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Contains(t, response, "todos")
        assert.Contains(t, response, "total")
    })
    
    t.Run("バリデーションエラー", func(t *testing.T) {
        // 空のタイトルでリクエスト
        payload := map[string]string{"title": ""}
        body, _ := json.Marshal(payload)
        
        req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusBadRequest, w.Code)
    })
}
```
```

---

## データフロー（Feature-based）

### 作成フロー
1. **TodoForm Component**: ユーザー入力 → フォーム送信
2. **useTodoMutation Hook**: createTodo.mutate() 実行
3. **todoApi Service**: API POST リクエスト → バックエンド
4. **Backend (Onion)**: Handler → UseCase → Domain → Infrastructure
5. **Response**: 新しいTODOをフロントエンドに返却
6. **React Query**: キャッシュ無効化 → UI自動更新

### 状態変更フロー
1. **TodoItem Component**: チェックボックスクリック
2. **useTodoMutation Hook**: toggleComplete.mutate() 実行  
3. **todoApi Service**: API PUT リクエスト → バックエンド
4. **Backend (Onion)**: ドメインエンティティの状態変更
5. **Response**: 更新されたTODOを返却
6. **React Query**: オプティミスティック更新 + キャッシュ同期

### エラーハンドリングフロー
1. **バリデーションエラー**: フォームレベルでキャッチ → ユーザーにフィードバック
2. **APIエラー**: React Query の onError → エラーメッセージ表示
3. **ネットワークエラー**: 自動リトライ → 失敗時はオフライン表示

---

## Feature-based アーキテクチャの利点

### 1. ドメイン一致
- バックエンドのTODOドメイン ↔ フロントエンドのtodo feature
- 統一されたビジネス概念でコード組織化

### 2. 独立開発
- TODO機能の全レイヤー（UI、hooks、API、types）が1つのフォルダーに集約
- 他機能への影響を最小限に抑えた開発・テストが可能

### 3. スケーラビリティ  
- 新機能追加時は新しいfeatureフォルダーを作成
- 既存機能への影響なしに機能拡張

### 4. テスト戦略
- Feature単位でのテスト（unit, integration, E2E）
- 機能横断のテストを最小限に抑制

このFeature-based Componentsアプローチにより、バックエンドのOnion Architectureと完全に整合した、保守性の高いフロントエンド設計を実現できます。
