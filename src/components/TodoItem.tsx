import type { Todo } from '@/types/todo';

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string) => void;
  onDelete: (id: string) => void;
  onEdit: (id: string) => void;
}

export function TodoItem({ todo, onToggle, onDelete, onEdit }: TodoItemProps) {
  const priorityColors = {
    low: 'bg-green-100 text-green-800',
    medium: 'bg-yellow-100 text-yellow-800',
    high: 'bg-red-100 text-red-800',
  };

  return (
    <div className="flex items-center gap-3 p-4 bg-white border border-gray-200 rounded-lg hover:shadow-md transition-shadow">
      <input
        type="checkbox"
        checked={todo.completed}
        onChange={() => onToggle(todo.id)}
        className="w-5 h-5 text-blue-600 rounded focus:ring-2 focus:ring-blue-500"
      />
      <div className="flex-1 min-w-0">
        <p
          className={`text-lg ${
            todo.completed ? 'line-through text-gray-500' : 'text-gray-900'
          }`}
        >
          {todo.title}
        </p>
        <div className="flex gap-2 mt-1">
          <span
            className={`px-2 py-1 text-xs font-semibold rounded ${priorityColors[todo.priority]}`}
          >
            {todo.priority === 'low' && '低'}
            {todo.priority === 'medium' && '中'}
            {todo.priority === 'high' && '高'}
          </span>
        </div>
      </div>
      <div className="flex gap-2">
        <button
          onClick={() => onEdit(todo.id)}
          className="px-3 py-1 text-sm text-blue-600 hover:bg-blue-50 rounded transition-colors"
        >
          編集
        </button>
        <button
          onClick={() => onDelete(todo.id)}
          className="px-3 py-1 text-sm text-red-600 hover:bg-red-50 rounded transition-colors"
        >
          削除
        </button>
      </div>
    </div>
  );
}
