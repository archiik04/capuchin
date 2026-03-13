import { type InputHTMLAttributes, forwardRef } from "react"

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(({ label, id, className = "", ...props }, ref) => {
  return (
    <div className="flex flex-col gap-2">
      {label && (
        <label htmlFor={id} className="text-[0.72rem] font-bold text-[#374151] uppercase tracking-[0.08em]">
          {label}
        </label>
      )}
      <input
        ref={ref}
        id={id}
        className={[
          "w-full py-3 px-4 border-[1.5px] border-[#e5e7eb] rounded-xl text-[0.9rem] text-[#111]",
          "outline-none bg-[#fafafa] transition-[border-color,box-shadow] duration-[180ms]",
          "focus:border-[#111] focus:shadow-[0_0_0_3px_rgba(0,0,0,0.06)] focus:bg-white",
          "placeholder:text-[#d1d5db]",
          className,
        ].join(" ")}
        {...props}
      />
    </div>
  )
})

Input.displayName = "Input"