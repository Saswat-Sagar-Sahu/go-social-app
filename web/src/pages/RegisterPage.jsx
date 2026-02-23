import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import AuthLayout from '../components/AuthLayout'
import { apiRequest } from '../lib/api'

export default function RegisterPage() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ username: '', email: '', password: '' })
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
      await apiRequest('/v1/users/register', {
        method: 'POST',
        body: form,
      })

      navigate('/activate', {
        replace: true,
        state: { email: form.email },
      })
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout
      title="Create account"
      subtitle="Sign up, then activate your account using the token from backend logs."
      footer={
        <p>
          Already have an account? <Link to="/login">Login</Link>
        </p>
      }
    >
      <form onSubmit={onSubmit} className="form-grid">
        <label htmlFor="username">Username</label>
        <input
          id="username"
          name="username"
          value={form.username}
          onChange={onChange}
          autoComplete="username"
          required
        />

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
          autoComplete="new-password"
          required
        />

        {error && <p className="form-error">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Creating...' : 'Register'}
        </button>
      </form>
    </AuthLayout>
  )
}
