import { useState, useEffect } from "react"
import { login } from "../api/auth"
import { useNavigate } from "react-router-dom"

export default function Login() {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")

  const navigate = useNavigate()

  // Auto redirect 
  useEffect(() => {
    const token = localStorage.getItem("token")
    if (token) navigate("/todos")
  }, [])

  const handleLogin = async () => {
    try {
      setError("")

      const data = await login(email, password)

      // Save JWT
      localStorage.setItem("token", data.token)

      // Redirect to Todos page
      navigate("/todos")

    } catch (err: any) {
      console.error(err)
      setError("Invalid email or password")
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-gray-800 rounded-xl shadow-2xl border border-gray-700">

        {/* Header */}
        <div className="p-6 border-b border-gray-700">
          <h2 className="text-3xl font-bold text-white">Login</h2>
          <p className="text-gray-400 text-sm mt-1">Welcome back</p>
        </div>

        {/* Form */}
        <div className="p-6 space-y-4">

          <input
            placeholder="Email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            className="w-full px-4 py-3 bg-gray-900 text-white rounded-lg border border-gray-700 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none placeholder-gray-500"
          />

          <input
            placeholder="Password"
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            className="w-full px-4 py-3 bg-gray-900 text-white rounded-lg border border-gray-700 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none placeholder-gray-500"
          />

          <button
            onClick={handleLogin}
            className="w-full py-3 bg-blue-600 hover:bg-blue-500 text-white font-medium rounded-lg transition"
          >
            Login
          </button>

          {error && (
            <p className="text-red-400 text-sm">{error}</p>
          )}

        </div>
      </div>
    </div>
  )
}
