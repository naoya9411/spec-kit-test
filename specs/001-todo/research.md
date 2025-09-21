# Research: シンプルなTODOアプリ技術選択

**作成日**: 2025年9月21日  
**対象機能**: シンプルなTODOアプリ  
**研究範囲**: データベース、マイクロサービス通信、テストライブラリ、コンテナ化

---

## データベース選択

### Decision: PostgreSQL
**選択理由**:
- 強力なACID特性とデータ整合性保証
- DDDにおけるドメインモデルの永続化に最適
- Go言語のORMライブラリ（GORM）との親和性
- JSON型サポートによる柔軟性
- SQLiteより本格的な本番運用に対応

**代替案検討**:
- **MongoDB**: NoSQLの柔軟性はあるが、TODOアプリの構造化データにはオーバースペック
- **SQLite**: 開発初期には適しているが、マイクロサービス環境では制約が多い

## コンテナオーケストレーション

### Decision: Docker Compose
**選択理由**:
- 開発環境での簡単なマルチコンテナ管理
- 本番環境への移行時にKubernetesに変更可能
- TODOアプリの規模に適している
- CI/CD統合が容易

**代替案検討**:
- **Kubernetes**: 高機能だが現段階では複雑すぎる
- **単一コンテナ**: マイクロサービス要件に適さない

### 実装方針
- Multi-stage Dockerfileでビルド最適化
- development/production環境の設定分離
- ヘルスチェックとロギング設定

---

## React アーキテクチャ統合戦略

### Feature-driven Development
- バックエンドのOnion Architecture（Domain中心）
- フロントエンドのFeature-based Components（機能中心）
- 両者がTODOドメインで統一された設計思想

### 利点
1. **一貫性**: バックエンドのドメイン境界とフロントエンドの機能境界が一致
2. **スケーラビリティ**: 新機能追加時に独立したフォルダー構造で開発可能
3. **保守性**: 機能変更時の影響範囲が明確
4. **テスト性**: 機能単位でのテスト戦略が立てやすい

---

## マイクロサービス通信方式

### Decision: HTTP REST
**選択理由**:
- シンプルで理解しやすい
- デバッグとモニタリングが容易
- TODOアプリの比較的シンプルな操作に適している
- フロントエンドとの統合が自然

**代替案検討**:
- **gRPC**: 高性能だが、TODOアプリには過剰。学習コストとデバッグ複雑性がデメリット

### 実装方針
- Go Gin frameworkでRESTful API実装
- OpenAPI 3.0仕様書による契約駆動開発
- JSON形式での統一したデータ交換

---

## フロントエンドアーキテクチャ

### Decision: Feature-based Components
**選択理由**:
- 機能単位でコンポーネントを組織化、スケーラビリティが高い
- ドメイン駆動設計（DDD）のアプローチと整合性がある
- コンポーネントの再利用性と保守性が向上
- TODOアプリの機能（作成、表示、編集、削除）に自然にマッピング

**代替案検討**:
- **Atomic Design**: デザインシステム重視だが、ビジネスロジックの組織化には不向き
- **Layer-based**: 技術的な層で分離するが、機能追加時に複数フォルダーを跨ぐ必要がある
- **Page-based**: シンプルだが、コンポーネント再利用が困難

### Feature-based 構成例
```
src/
├── features/
│   └── todo/
│       ├── components/          # TODO機能のコンポーネント
│       │   ├── TodoList.tsx
│       │   ├── TodoItem.tsx
│       │   ├── TodoForm.tsx
│       │   └── index.ts
│       ├── hooks/              # TODO機能のカスタムフック
│       │   ├── useTodos.ts
│       │   └── useTodoMutation.ts
│       ├── services/           # TODO機能のAPIサービス
│       │   └── todoApi.ts
│       ├── types/              # TODO機能の型定義
│       │   └── todo.types.ts
│       ├── __tests__/          # 機能単位のテスト（近接性重視）
│       │   ├── components/     # コンポーネントテスト
│       │   │   ├── TodoList.test.tsx
│       │   │   ├── TodoItem.test.tsx
│       │   │   └── TodoForm.test.tsx
│       │   ├── hooks/          # フックテスト
│       │   │   ├── useTodos.test.ts
│       │   │   └── useTodoMutation.test.ts
│       │   └── services/       # APIサービステスト
│       │       └── todoApi.test.ts
│       └── index.ts            # 機能エクスポート
├── shared/                     # 共通コンポーネント・ユーティリティ
│   ├── components/
│   │   ├── ui/                 # shadcn/ui ベースコンポーネント
│   │   └── layout/
│   ├── hooks/
│   ├── services/
│   ├── types/
│   └── __tests__/              # 共通機能のテスト
│       ├── components/
│       ├── hooks/
│       └── services/
└── app/                        # アプリケーション設定・ルーティング
    ├── App.tsx
    ├── router.tsx
    ├── providers.tsx
    └── __tests__/              # アプリケーション層のテスト
        ├── App.test.tsx
        └── router.test.tsx
```

---

## フロントエンドテストライブラリ

### Decision: Vitest + Feature-based Testing
**選択理由**:
- Viteベースで高速な実行速度
- TypeScript/ESModulesのネイティブサポート
- Jestとほぼ互換のAPI
- Feature-based構成と相性が良い
- 2025年時点での最新ベストプラクティス

**代替案検討**:
- **Jest**: 成熟しているが、ESModules対応やViteとの統合で設定が複雑
- **Testing Library**: 単体では不十分、Vitestとの組み合わせで使用

### テスト構成（近接性重視）
```
src/features/todo/
├── components/
│   ├── TodoList.tsx
│   ├── TodoItem.tsx
│   └── TodoForm.tsx
├── hooks/
│   ├── useTodos.ts
│   └── useTodoMutation.ts
├── services/
│   └── todoApi.ts
└── __tests__/                  # 機能内テスト（近接性重視）
    ├── components/
    │   ├── TodoList.test.tsx   # コンポーネントテスト
    │   ├── TodoItem.test.tsx
    │   └── TodoForm.test.tsx
    ├── hooks/
    │   ├── useTodos.test.ts    # フックテスト
    │   └── useTodoMutation.test.ts
    └── services/
        └── todoApi.test.ts     # APIサービステスト

tests/e2e/
└── todo.spec.ts                # E2Eテスト（機能横断）
```

### テストのメリット
- **近接性**: 実装とテストが同じfeatureフォルダー内に配置
- **発見しやすさ**: 機能変更時にテストファイルがすぐに見つかる
- **メンテナンス性**: 機能削除時にテストも一緒に削除される
- **独立性**: 各feature単位でテストスイートが完結

### 実装方針
- 各機能フォルダー内にテストファイルを配置
- Vitest + @testing-library/reactでコンポーネントテスト
- MSW (Mock Service Worker) でAPIモック
- Playwright でE2Eテスト（機能単位のシナリオ）

---

## Go言語でのDDD実装アーキテクチャ

### Decision: Onion Architecture + DDD
**選択理由**:
- 依存性の方向が内側（ドメイン）に向かう明確な構造
- Clean Architectureより概念的にシンプル
- Go言語のパッケージ依存性管理と自然に一致
- テストでのモック不要（インターフェース分離で実現）
- TODOアプリの規模に最適

**代替案検討**:
- **Clean Architecture**: 4層構造だが、TODOアプリには過剰。Enterprise Business Rulesレイヤーが不要
- **Hexagonal Architecture**: ポート・アダプターパターンは良いが、概念がOnion Architectureより複雑
- **Traditional Layered Architecture**: 上位層が下位層に依存するため、テスタビリティが低い

### Onion Architecture 層構成
```
┌─────────────────────────────────┐
│        Infrastructure          │ ← 外部システム（DB、HTTP）
├─────────────────────────────────┤
│        Application             │ ← ユースケース・オーケストレーション
├─────────────────────────────────┤
│          Domain                │ ← エンティティ・ビジネスロジック
└─────────────────────────────────┘
            ↑ 依存方向
```

### 実装方針
### アーキテクチャ比較表

| 観点 | Clean Architecture | Hexagonal Architecture | Onion Architecture | Traditional Layered |
|------|-------------------|------------------------|-------------------|-------------------|
| 層数 | 4層（過剰） | 可変（複雑） | 3層（適切） | 3-4層（上位→下位依存） |
| 依存方向 | 内向き | ポート・アダプター | 内向き（シンプル） | 外向き（問題） |
| Go親和性 | 良い | 普通 | **最良** | 良い |
| テスト性 | 良い | 良い | **最良** | 悪い |
| 学習コスト | 高い | 高い | **低い** | 低い |
| TODO規模適性 | 過剰 | 過剰 | **最適** | 不十分 |

**結論**: Onion ArchitectureがTODOアプリの規模と要求に最も適している

---

## テストライブラリ統合戦略

### Frontend Testing Stack
- **Unit**: Vitest + @testing-library/react
- **Integration**: Vitest + MSW
- **E2E**: Playwright

### Backend Testing Stack  
- **Unit**: Go標準testing + testify/assert
- **Integration**: Testcontainers + PostgreSQL
- **Contract**: OpenAPI schema validation

### テスト実行順序（TDD準拠）
1. Contract tests (API仕様検証)
2. Integration tests (DB結合検証)  
3. E2E tests (ユーザーシナリオ検証)
4. Unit tests (個別機能検証)

---

## 技術選択まとめ

| カテゴリ | 選択技術 | 理由 |
|---------|---------|------|
| Database | PostgreSQL | 堅牢性とDDD親和性 |
| Backend Framework | Go + Gin | パフォーマンスとシンプルさ |
| Frontend Framework | TypeScript + React + TailwindCSS | 型安全性とモダンUI |
| Testing | Vitest + Playwright | 最新ベストプラクティス |
| Container | Docker + Docker Compose | 開発効率と移行容易性 |
| Communication | HTTP REST | シンプルさとデバッグ性 |

全ての選択は、TODOアプリの要求仕様と憲法原則（シンプルさ、テスタビリティ、観測可能性）に基づいて行われました。
