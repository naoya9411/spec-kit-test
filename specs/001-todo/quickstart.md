# Quickstart Guide: シンプルなTODOアプリ

**作成日**: 2025年9月21日  
**対象**: 開発者・QA担当者・プロダクトオーナー  
**目的**: TODOアプリの基本機能を5分で確認

---

## 前提条件

### 必要なツール
- Docker & Docker Compose
- curl または Postman
- Webブラウザー

### 環境起動
```bash
# リポジトリクローン後
cd todo-app
docker-compose up -d

# 起動確認（30秒程度待機）
curl http://localhost:8080/api/v1/todos
```

---

## 基本操作シナリオ

### シナリオ1: TODO作成
**目的**: 新しいTODOアイテムを作成できることを確認

```bash
# API経由でTODO作成
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "食材を買いに行く"}'

# 期待されるレスポンス例
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "食材を買いに行く",
  "completed": false,
  "createdAt": "2025-09-21T10:00:00Z",
  "updatedAt": "2025-09-21T10:00:00Z"
}
```

**Webアプリでの確認**:
1. http://localhost:3000 を開く
2. 入力フィールドに「食材を買いに行く」を入力
3. 「作成」ボタンをクリック
4. リストに新しいTODOが表示されることを確認

### シナリオ2: TODOリスト確認
**目的**: 作成したTODOがリストに表示されることを確認

```bash
# TODOリスト取得
curl http://localhost:8080/api/v1/todos

# 期待されるレスポンス例
{
  "todos": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "食材を買いに行く",
      "completed": false,
      "createdAt": "2025-09-21T10:00:00Z",
      "updatedAt": "2025-09-21T10:00:00Z"
    }
  ],
  "total": 1
}
```

### シナリオ3: TODO完了
**目的**: TODOの完了状態を変更できることを確認

```bash
# TODO完了に変更
TODO_ID="550e8400-e29b-41d4-a716-446655440000"
curl -X PUT http://localhost:8080/api/v1/todos/${TODO_ID} \
  -H "Content-Type: application/json" \
  -d '{"completed": true}'

# 期待されるレスポンス例（completedがtrueに変更）
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "食材を買いに行く",
  "completed": true,
  "createdAt": "2025-09-21T10:00:00Z",
  "updatedAt": "2025-09-21T10:00:15Z"
}
```

**Webアプリでの確認**:
1. 作成したTODOアイテムのチェックボックスをクリック
2. アイテムが完了状態（取り消し線など）で表示されることを確認

### シナリオ4: TODO削除
**目的**: TODOを削除できることを確認

```bash
# TODO削除
curl -X DELETE http://localhost:8080/api/v1/todos/${TODO_ID}

# 期待されるレスポンス: 204 No Content

# 削除確認
curl http://localhost:8080/api/v1/todos
# 期待されるレスポンス: 空のリスト
{
  "todos": [],
  "total": 0
}
```

---

## エラーケース検証

### バリデーションエラー
```bash
# 空のタイトルでTODO作成（エラー）
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": ""}'

# 期待されるレスポンス: 400 Bad Request
{
  "message": "title cannot be empty",
  "code": "VALIDATION_ERROR"
}

# 長すぎるタイトル（101文字）
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "この文字列は100文字を超えるように作られていますこの文字列は100文字を超えるように作られていますこの文字列は100文字を超えるように作られています"}'

# 期待されるレスポンス: 400 Bad Request
{
  "message": "title cannot exceed 100 characters",
  "code": "VALIDATION_ERROR"
}
```

### 存在しないリソースアクセス
```bash
# 存在しないTODOを取得
curl http://localhost:8080/api/v1/todos/nonexistent-id

# 期待されるレスポンス: 404 Not Found
{
  "message": "todo not found",
  "code": "NOT_FOUND"
}
```

---

## パフォーマンス検証

### レスポンス時間測定
```bash
# API応答時間測定（100ms以内の目標）
time curl http://localhost:8080/api/v1/todos

# 複数TODO作成してパフォーマンス確認
for i in {1..10}; do
  curl -X POST http://localhost:8080/api/v1/todos \
    -H "Content-Type: application/json" \
    -d "{\"title\": \"テストTODO ${i}\"}"
done

# リスト取得のパフォーマンス確認
time curl http://localhost:8080/api/v1/todos
```

---

## 単体テスト・統合テスト実行

### フロントエンドテスト
```bash
# 全ての単体テスト実行（feature内の__tests__）
docker exec -it todo-frontend npm run test

# TODO機能のテストのみ実行
docker exec -it todo-frontend npm run test -- src/features/todo

# テストカバレッジ確認
docker exec -it todo-frontend npm run test:coverage

# 期待される結果：
✓ src/features/todo/__tests__/components/TodoForm.test.tsx
✓ src/features/todo/__tests__/components/TodoItem.test.tsx  
✓ src/features/todo/__tests__/components/TodoList.test.tsx
✓ src/features/todo/__tests__/hooks/useTodos.test.ts
✓ src/features/todo/__tests__/hooks/useTodoMutation.test.ts
✓ src/features/todo/__tests__/services/todoApi.test.ts
```

## E2Eテスト実行

### 自動テスト実行
```bash
# Playwrightテスト実行
docker exec -it todo-frontend npm run test:e2e

# 期待される結果：全テストPASS
✓ should create a new todo item
✓ should mark todo as completed
✓ should delete todo item
✓ should show validation error for empty title
```

---

## ログ確認

### アプリケーションログ
```bash
# バックエンドログ確認
docker logs todo-backend

# フロントエンドログ確認
docker logs todo-frontend

# データベースログ確認
docker logs todo-db
```

### 期待されるログ出力例
```json
{
  "level": "info",
  "time": "2025-09-21T10:00:00Z",
  "message": "TODO created",
  "todo_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_action": "create_todo"
}
```

---

## 環境クリーンアップ

### 開発環境停止
```bash
# コンテナ停止・削除
docker-compose down -v

# イメージ削除（必要に応じて）
docker rmi $(docker images -q todo-*)
```

---

## トラブルシューティング

### よくある問題

**ポート競合**:
```bash
# ポート使用状況確認
lsof -i :8080
lsof -i :3000

# 別ポートで起動
PORT=8081 docker-compose up -d
```

**データベース接続エラー**:
```bash
# PostgreSQL接続確認
docker exec -it todo-db psql -U todouser -d tododb -c "SELECT 1;"
```

**API応答がない**:
```bash
# ヘルスチェック確認
curl http://localhost:8080/health

# コンテナ状態確認
docker-compose ps
```

---

## 成功基準

このクイックスタートが成功であれば：
✅ 全4つの基本シナリオが正常に実行できる  
✅ エラーケースが期待通りのレスポンスを返す  
✅ API応答時間が100ms以内  
✅ E2Eテストが全てPASS  
✅ ログに適切な情報が出力される  

これらが確認できれば、TODOアプリの基本機能は正常に動作しています。
