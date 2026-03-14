import { forwardRef } from "react"

interface EmptyStateProps {
  message?: string
}

export const EmptyState = forwardRef<HTMLDivElement, EmptyStateProps>(
  ({ message = "No tasks yet" }, ref) => {
    return (
      <div ref={ref} className="text-center py-14 px-6 text-[#b0b7c3]">
        <div className="text-[2.6rem] mb-2.5 opacity-50">🐒</div>
        <div className="text-[0.88rem] font-medium">{message}</div>
      </div>
    )
  }
)

EmptyState.displayName = "EmptyState"