import { useState, useEffect } from "react"
import { signup } from "../api/auth"
import { useNavigate } from "react-router-dom"

export default function Signup() {

  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")

  const navigate = useNavigate()

  //CLEAR OLD TOKEN WHEN ENTERING SIGNUP
  useEffect(() => {
    localStorage.removeItem("token")
  }, [])

  const handleSignup = async () => {
    try {
      setError("")
      setSuccess("")

      await signup(email, password)

      setSuccess("Account created! Redirecting to login...")

      setTimeout(() => {
        navigate("/login")
      }, 1500)

    } catch (err: any) {

      console.error(err)

      if (err.response?.data?.error) {
        setError(err.response.data.error)
      } else {
        setError("Signup failed")
      }
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="w-full max-w-md bg-gray-800 rounded-xl shadow-2xl border border-gray-700">

        <div className="p-6 border-b border-gray-700">
          <h2 className="text-3xl font-bold text-white">Signup</h2>
          <p className="text-gray-400 text-sm mt-1">Create your account</p>
        </div>

        <div className="p-6 space-y-4">

          <input
            placeholder="Email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            className="w-full px-4 py-3 bg-gray-900 text-white rounded-lg border border-gray-700"
          />

          <input
            placeholder="Password"
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            className="w-full px-4 py-3 bg-gray-900 text-white rounded-lg border border-gray-700"
          />

          <button
            onClick={handleSignup}
            className="w-full py-3 bg-green-600 hover:bg-green-500 text-white font-medium rounded-lg"
          >
            Create Account
          </button>

          {error && (
            <p className="text-red-400 text-sm">{error}</p>
          )}

          {success && (
            <p className="text-green-400 text-sm">{success}</p>
          )}

        </div>
      </div>
    </div>
  )
}
