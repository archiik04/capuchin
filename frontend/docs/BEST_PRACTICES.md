# Frontend Best Practices

This document outlines coding standards, patterns, and best practices for the Capuchin frontend application.

---

## Table of Contents

1. [Project Structure](#project-structure)
2. [Component Guidelines](#component-guidelines)
3. [Custom Hooks](#custom-hooks)
4. [TypeScript Guidelines](#typescript-guidelines)
5. [API Integration](#api-integration)
6. [State Management](#state-management)
7. [Styling with Tailwind](#styling-with-tailwind)
8. [Error Handling](#error-handling)
9. [Testing Recommendations](#testing-recommendations)
10. [Code Organization](#code-organization)

---

## Project Structure

### Directory Organization

```
src/
├── components/       # Reusable UI components
│   ├── ui/         # Generic, reusable primitives
│   ├── feature/    # Feature-specific components
│   └── layout/     # Layout components
├── hooks/           # Custom React hooks
├── lib/             # Utilities and libraries
├── pages/           # Route-level components
├── config/          # Configuration files
└── types/           # TypeScript type definitions
```

### When to Create a New Folder

| Folder | When to Use |
|--------|-------------|
| `ui/` | Generic components reusable across features (Button, Input, Modal) |
| `feature/` | Components specific to one feature domain |
| `hooks/` | Custom hooks encapsulating stateful logic |
| `lib/` | Shared utilities, API clients, constants |
| `pages/` | Components that represent routes |

---

## Component Guidelines

### Component Structure

```typescript
// 1. Imports
import { useState, useCallback } from "react"
import type { ComponentProps } from "@/types"

// 2. Types
interface MyComponentProps {
  title: string
  onAction: () => void
  variant?: "primary" | "secondary"
}

// 3. Component with forwardRef (if needed)
export const MyComponent = forwardRef<HTMLDivElement, MyComponentProps>(
  ({ title, onAction, variant = "primary", className = "" }, ref) => {
    // 4. Hooks first
    const [state, setState] = useState(false)

    // 5. Callbacks
    const handleClick = useCallback(() => {
      onAction()
    }, [onAction])

    // 6. Render
    return (
      <div ref={ref} className={className}>
        <h1>{title}</h1>
        <button onClick={handleClick}>Click</button>
      </div>
    )
  }
)

MyComponent.displayName = "MyComponent"
```

### Component Best Practices

1. **Use functional components** - Always prefer functional components with hooks over class components.

2. **Forward refs when needed** - Use `forwardRef` for components that need to expose DOM elements (buttons, inputs).

3. **Define displayName** - Always set `displayName` for forwardRef components and components used in DevTools.

4. **Extract types** - Define prop interfaces at the top of the file.

5. **Default prop values** - Use destructuring with defaults instead of defaultProps (deprecated in React).

```typescript
// ✅ Good
function Button({ variant = "primary", size = "md" }) { ... }

// ❌ Avoid (deprecated)
function Button({ variant, size }) { ... }
Button.defaultProps = { variant: "primary", size: "md" }
```

6. **Colocate related files** - Keep component, tests, and types together.

```
components/
├── Button/
│   ├── Button.tsx
│   ├── Button.test.tsx
│   └── index.ts (barrel export)
```

---

## Custom Hooks

### When to Create a Custom Hook

Create a custom hook when:
- You have stateful logic reused across components
- You want to extract and reuse component logic
- You need to share stateful behavior between components

### Hook Naming Convention

- **Always prefix with `use`** - This is required for React to recognize it as a hook.
- Use verb or verb+noun pattern: `useAuth`, `useTodos`, `useLocalStorage`.

### Hook Structure

```typescript
export function useTodos(token: string | null, isAuthed: boolean) {
  // 1. State
  const [todos, setTodos] = useState<Todo[]>([])
  const [loading, setLoading] = useState(true)

  // 2. Effects
  useEffect(() => {
    // Side effects
  }, [dependencies])

  // 3. Callbacks with useCallback
  const addTodo = useCallback(async (item: string) => {
    // Implementation
  }, [dependencies])

  // 4. Computed values
  const activeCount = todos.filter(t => !t.completed).length
  const completedCount = todos.filter(t => t.completed).length

  // 5. Return stable API
  return {
    todos,
    loading,
    activeCount,
    completedCount,
    addTodo,
    // ... other methods
  }
}
```

### Hook Best Practices

1. **Return a stable object** - Avoid returning new object references on every render when possible.

```typescript
// ✅ Good - return individual properties
return { todos, loading, addTodo }

// ⚠️ Caution - creates new object each render
return { data: { todos, loading, addTodo } }
```

2. **Use useCallback for returned functions** - Prevents unnecessary re-renders in consumers.

3. **Document dependencies clearly** - Add JSDoc comments explaining hook behavior and dependencies.

4. **Handle cleanup** - Return cleanup functions from useEffect when needed.

---

## TypeScript Guidelines

### Type Definitions

1. **Use explicit return types for functions**

```typescript
// ✅ Good - explicit return type
function getTodos(): Promise<Todo[]> {
  return fetchTodos()
}

// ⚠️ Avoid - inferred return type
function getTodos() {
  return fetchTodos()
}
```

2. **Use type aliases for object shapes**

```typescript
// ✅ Good
type User = {
  id: string
  email: string
  name?: string
}

// ⚠️ Avoid - using interfaces for simple shapes is fine too
interface User {
  id: string
  email: string
}
```

3. **Use discriminated unions for state**

```typescript
type LoadingState<T> =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; data: T }
  | { status: "error"; error: Error }
```

### Path Aliases

Always use path aliases (`@/`) for imports:

```typescript
// ✅ Good
import { Button } from "@/components/ui/Button"
import { useAuth } from "@/hooks/useAuth"

// ❌ Avoid - relative paths
import { Button } from "../../components/ui/Button"
```

### Type Safety

1. **Avoid `any`** - Use `unknown` or proper types instead.

```typescript
// ❌ Avoid
function process(data: any): any { ... }

// ✅ Good
function process(data: unknown): Todo { ... }

// Or use proper typing
function process(data: RawTodo): Todo { ... }
```

2. **Use type guards** - Narrow types with type guards.

```typescript
function isTodo(value: unknown): value is Todo {
  return (
    typeof value === "object" &&
    value !== null &&
    "id" in value &&
    "item" in value
  )
}
```

---

## API Integration

### Centralize API Calls

All API calls should go through the `lib/api.ts` or dedicated API modules:

```typescript
// lib/api.ts
const BASE_URL = env.apiUrl

const authHeaders = (token: string): HeadersInit => ({
  "Content-Type": "application/json",
  Authorization: `Bearer ${token}`,
})

export const todosApi = {
  getAll: async (token: string): Promise<Todo[]> => {
    const res = await fetch(`${BASE_URL}/user/todos`, {
      headers: authHeaders(token),
    })
    if (!res.ok) throw new Error(`Failed to fetch todos: ${res.status}`)
    return res.json()
  },
  // ... other methods
}
```

### Error Handling in API

1. **Throw descriptive errors**

```typescript
// ✅ Good
if (!res.ok) {
  throw new Error(`Failed to fetch todos: ${res.status}`)
}

// ❌ Avoid
if (!res.ok) {
  throw new Error("Error")
}
```

2. **Normalize API responses** - Handle API-specific field naming.

```typescript
const normalise = (raw: Record<string, unknown>): Todo => ({
  id: String(raw.id ?? raw.ID),
  item: String(raw.item ?? raw.Item ?? ""),
  completed: Boolean(raw.completed ?? raw.Completed ?? false),
})
```

### API Best Practices

1. **Use typed interfaces** - Define input/output types for each API method.

2. **Separate concerns** - Group related API endpoints (authApi, todosApi).

3. **Handle loading states** - Always handle async states in hooks.

4. **Implement retries** - Consider retry logic for failed requests.

---

## State Management

### Local State with useState

```typescript
// ✅ Simple state
const [count, setCount] = useState(0)

// ✅ Lazy initialization for expensive computation
const [data, setData] = useState(() => {
  return computeExpensiveValue(props.id)
})
```

### Derived State

Calculate derived state during render rather than in state:

```typescript
// ✅ Good - derived during render
const activeTodos = todos.filter(t => !t.completed)
const completedTodos = todos.filter(t => t.completed)
const total = todos.length

return { activeTodos, completedTodos, total }

// ❌ Avoid - unnecessary state
const [activeTodos, setActiveTodos] = useState<Todo[]>([])
```

### Optimistic Updates

Implement optimistic updates with rollback for better UX:

```typescript
const updateTodo = useCallback(async (id: string, item: string) => {
  // Snapshot current state
  const snapshot = todos

  // Optimistic update
  setTodos(prev => prev.map(t => t.id === id ? { ...t, item } : t))

  try {
    await todosApi.update(token!, id, item)
  } catch (error) {
    // Rollback on failure
    setTodos(snapshot)
    console.error(error)
  }
}, [todos, token])
```

---

## Styling with Tailwind

### Class Organization

Organize Tailwind classes in a consistent order:

```typescript
// 1. Layout (position, display)
// 2. Spacing (margin, padding)
// 3. Sizing (width, height)
// 4. Visual (colors, borders, shadows)
// 5. Typography (font, text)
// 6. States (hover, focus, disabled)
// 7. Effects (transitions, animations)

<button
  className={[
    // Layout
    "flex items-center justify-center",
    // Spacing
    "py-2 px-4",
    // Visual
    "bg-blue-500 rounded-lg",
    // Typography
    "text-white font-semibold",
    // States
    "hover:bg-blue-600 focus:ring-2",
    // Effects
    "transition-colors duration-150",
  ].join(" ")}
>
  Click me
</button>
```

### Component Variants

Use variant patterns for flexible components:

```typescript
const variantClasses: Record<Variant, string> = {
  primary: "bg-blue-500 text-white hover:bg-blue-600",
  secondary: "bg-gray-100 text-gray-700 hover:bg-gray-200",
  danger: "bg-red-500 text-white hover:bg-red-600",
}

const sizeClasses: Record<Size, string> = {
  sm: "text-sm px-3 py-1",
  md: "text-base px-4 py-2",
  lg: "text-lg px-6 py-3",
}

<button className={`${variantClasses[variant]} ${sizeClasses[size]}`}>
  {children}
</button>
```

### Best Practices

1. **Use arbitrary values sparingly** - Prefer utility classes over arbitrary values.

```typescript
// ✅ Good
<div className="grid grid-cols-3 gap-4">

// ⚠️ Caution
<div className="grid grid-cols-[1fr_2fr_1fr] gap-4">
```

2. **Extract repeated patterns** - Create component abstractions for repeated class combinations.

3. **Use `clamp()` for responsive typography**

```typescript
className="text-[clamp(1rem,5vw,2rem)]"
```

4. **Consistent spacing** - Use Tailwind's spacing scale consistently.

---

## Error Handling

### Try-Catch Blocks

Always wrap async operations in try-catch:

```typescript
const addTodo = useCallback(async (item: string) => {
  try {
    const created = await todosApi.create(token!, item)
    setTodos(prev => [...prev, created])
  } catch (error) {
    console.error("Failed to add todo:", error)
    // Optionally show error to user
  }
}, [token])
```

### Error Boundaries

Consider adding error boundaries for graceful degradation:

```typescript
class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean }
> {
  state = { hasError: false }

  static getDerivedStateFromError() {
    return { hasError: true }
  }

  render() {
    if (this.state.hasError) {
      return <div>Something went wrong.</div>
    }
    return this.props.children
  }
}
```

### User-Facing Errors

Show meaningful error messages to users:

```typescript
// ✅ Good - specific error message
if (!res.ok) throw new Error("Invalid email or password")

// ❌ Avoid - generic message
if (!res.ok) throw new Error("Error")
```

---

## Testing Recommendations

### Unit Tests for Hooks

```typescript
import { renderHook, act } from "@testing-library/react"
import { useCounter } from "./useCounter"

test("should increment counter", () => {
  const { result } = renderHook(() => useCounter())

  act(() => {
    result.current.increment()
  })

  expect(result.current.count).toBe(1)
})
```

### Component Tests

```typescript
import { render, screen, fireEvent } from "@testing-library/react"
import { Button } from "./Button"

test("should call onClick when clicked", () => {
  const handleClick = vi.fn()
  render(<Button onClick={handleClick}>Click me</Button>)

  fireEvent.click(screen.getByText("Click me"))

  expect(handleClick).toHaveBeenCalled()
})
```

### Testing Best Practices

1. **Test behavior, not implementation** - Focus on what the component does, not how.

2. **Test edge cases** - Empty states, loading states, error states.

3. **Mock external dependencies** - Mock API calls, localStorage, etc.

4. **Keep tests close to source** - Co-locate tests with components.

---

## Code Organization

### Import Order

Organize imports consistently:

```typescript
// 1. React imports
import { useState, useEffect } from "react"

// 2. External libraries
import { useNavigate } from "react-router"
import { format } from "date-fns"

// 3. Path aliases (internal)
import { Button } from "@/components/ui/Button"
import { useAuth } from "@/hooks/useAuth"

// 4. Relative imports (local)
import type { Todo } from "../types"

// 5. Types
import type { MyType } from "@/types"
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Components | PascalCase | `MyComponent`, `TodoItem` |
| Hooks | camelCase with `use` prefix | `useAuth`, `useTodos` |
| Types/Interfaces | PascalCase | `Todo`, `UserProps` |
| Constants | SCREAMING_SNAKE_CASE | `TOKEN_KEY`, `API_BASE_URL` |
| Variables/Functions | camelCase | `handleClick`, `formatDate` |

### File Naming

- **Components**: PascalCase (`Button.tsx`, `TodoItem.tsx`)
- **Hooks**: camelCase with `use` prefix (`useAuth.ts`, `useTodos.ts`)
- **Utils**: camelCase (`api.ts`, `helpers.ts`)
- **Types**: PascalCase, often in `index.ts` or `types.ts`

---

## Security Considerations

### Never Store Sensitive Data in LocalStorage

While this project uses localStorage for tokens, be aware of XSS vulnerabilities:

```typescript
// ⚠️ Caution - tokens in localStorage can be accessed by XSS
localStorage.setItem("token", token)
```

### Sanitize User Input

Always sanitize and validate user input:

```typescript
const sanitizeInput = (input: string): string => {
  return input.trim().slice(0, MAX_LENGTH)
}
```

---

## Performance Tips

### useMemo and useCallback

```typescript
// Memoize expensive computations
const sortedTodos = useMemo(() => {
  return todos.sort((a, b) => b.createdAt - a.createdAt)
}, [todos])

// Memoize callback references
const handleSubmit = useCallback((data: FormData) => {
  submitForm(data)
}, [submitForm])
```

### Avoid Unnecessary Re-renders

1. Keep component trees shallow
2. Use React.memo for pure components
3. Pass stable callback references
4. Consider code splitting with React.lazy

---

## Git Conventions

### Commit Messages

```
type(scope): description

Types: feat, fix, docs, style, refactor, test, chore
```

Example:
```
feat(todos): add todo filtering by status
fix(auth): handle token expiration gracefully
docs: update API integration guidelines
```

---

## Resources

- [React Documentation](https://react.dev)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [React Router](https://reactrouter.com)
- [ESLint Rules](https://eslint.org/rules)
</CONTENTS>