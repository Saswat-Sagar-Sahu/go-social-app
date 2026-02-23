import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import AuthLayout from '../components/AuthLayout'
import { saveToken } from '../lib/auth'
import { apiRequest } from '../lib/api'

export default function LoginPage() {
  const navigate = useNavigate()
  const location = useLocation()
  const [form, setForm] = useState({ email: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function onChange(event) {
    const { name, value } = event.target
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  async function onSubmit(event) {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await apiRequest('/v1/users/login', {
        method: 'POST',
        body: form,
      })

      saveToken(response.token)
      navigate('/horizon', { replace: true })
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout
      title="Welcome back"
      subtitle="Sign in to open your home dashboard."
      footer={
        <p>
          New user? <Link to="/register">Register</Link>
        </p>
      }
    >
      {location.state?.activated && <p className="form-success">Account activated. You can login now.</p>}

      <form onSubmit={onSubmit} className="form-grid">
        <label htmlFor="email">Email</label>
        <input
          id="email"
          name="email"
          type="email"
          value={form.email}
          onChange={onChange}
          autoComplete="email"
          required
        />

        <label htmlFor="password">Password</label>
        <input
          id="password"
          name="password"
          type="password"
          value={form.password}
          onChange={onChange}
          autoComplete="current-password"
          required
        />

        {error && <p className="form-error">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Signing in...' : 'Login'}
        </button>
      </form>
    </AuthLayout>
  )
}
