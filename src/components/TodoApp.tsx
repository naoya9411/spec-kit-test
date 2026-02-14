import { useState } from 'react';
import type { Todo } from '@/types/todo';
import { Header } from './Header';
import { TodoForm } from './TodoForm';
import { TodoList } from './TodoList';
import { TodoEditModal } from './TodoEditModal';

export function TodoApp() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [editingTodo, setEditingTodo] = useState<Todo | null>(null);

  const addTodo = (title: string) => {
    const newTodo: Todo = {
      id: crypto.randomUUID(),
      title,
      completed: false,
      priority: 'medium',
      tags: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    setTodos([newTodo, ...todos]);
  };

  const toggleTodo = (id: string) => {
    setTodos(
      todos.map((todo) =>
        todo.id === id
          ? { ...todo, completed: !todo.completed, updatedAt: new Date().toISOString() }
          : todo
      )
    );
  };

  const deleteTodo = (id: string) => {
    setTodos(todos.filter((todo) => todo.id !== id));
  };

  const startEditTodo = (id: string) => {
    const todo = todos.find((t) => t.id === id);
    if (todo) {
      setEditingTodo(todo);
    }
  };

  const saveTodo = (id: string, title: string, priority: 'low' | 'medium' | 'high') => {
    setTodos(
      todos.map((todo) =>
        todo.id === id
          ? { ...todo, title, priority, updatedAt: new Date().toISOString() }
          : todo
      )
    );
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />
      <main className="container mx-auto px-4 py-8 max-w-2xl">
        <TodoForm onAddTodo={addTodo} />
        <TodoList
          todos={todos}
          onToggle={toggleTodo}
          onDelete={deleteTodo}
          onEdit={startEditTodo}
        />
      </main>
      <TodoEditModal
        todo={editingTodo}
        onClose={() => setEditingTodo(null)}
        onSave={saveTodo}
      />
    </div>
  );
}
