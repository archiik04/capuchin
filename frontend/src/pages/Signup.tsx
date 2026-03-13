import { useState } from "react"
import { useNavigate } from "react-router"
import { ArrowRight } from "lucide-react"
import { Navbar } from "@/components/layout/Navbar"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { Alert } from "@/components/ui/Alert"
import { Divider } from "@/components/ui/Divider"
import { authApi } from "@/lib/api"

export default function Signup() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSignup = async () => {
    if (!email || !password) return
    setLoading(true)
    setError("")
    setSuccess("")
    try {
      await authApi.signup(email, password)
      setSuccess("Account created! Taking you to login...")
      setTimeout(() => navigate("/login"), 1500)
    } catch (err: unknown) {
      const anyErr = err as { response?: { data?: { error?: string } } }
      setError(anyErr.response?.data?.error ?? "Signup failed. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  const onKey = (e: React.KeyboardEvent) => e.key === "Enter" && handleSignup()

  return (
    <div className="grid-bg min-h-screen flex flex-col">
      <Navbar
        right={
          <button
            className="text-[0.83rem] font-semibold text-[#6b7280] bg-transparent border-none cursor-pointer transition-colors hover:text-[#111]"
            onClick={() => navigate("/")}
          >
            ← Back to home
          </button>
        }
      />

      <div className="flex-1 flex items-center justify-center p-10 px-6">
        <div className="bg-white border-[1.5px] border-[#e5e7eb] rounded-3xl p-12 px-11 w-full max-w-[420px] shadow-[0_8px_48px_rgba(0,0,0,0.06)] max-[480px]:p-9 max-[480px]:px-6">
          <h1 className="text-[1.85rem] font-black tracking-[-0.04em] text-[#111] mb-1.5 leading-[1.1]">
            Create your account
          </h1>

          <div className="flex flex-col gap-4 mt-8">
            <Input
              id="email"
              label="Email"
              type="email"
              placeholder="you@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onKeyDown={onKey}
            />
            <Input
              id="password"
              label="Password"
              type="password"
              placeholder="Enter password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              onKeyDown={onKey}
            />
          </div>

          <Button
            variant="primary"
            size="lg"
            fullWidth
            className="mt-6 rounded-xl"
            onClick={handleSignup}
            loading={loading}
          >
            Create account <ArrowRight size={16} />
          </Button>

          {error && <Alert type="error" message={error} />}
          {success && <Alert type="success" message={success} />}

          <Divider label="already have an account?" />

          <Button
            variant="secondary"
            fullWidth
            className="rounded-xl"
            onClick={() => navigate("/login")}
          >
            Sign in instead
          </Button>
        </div>
      </div>
    </div>
  )
}