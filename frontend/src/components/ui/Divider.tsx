interface DividerProps {
  label?: string
}

export function Divider({ label }: DividerProps) {
  return (
    <div className="flex items-center gap-3 my-6">
      <div className="flex-1 h-px bg-[#f3f4f6]" />
      {label && <span className="text-[#d1d5db] font-semibold text-[0.75rem]">{label}</span>}
      <div className="flex-1 h-px bg-[#f3f4f6]" />
    </div>
  )
}