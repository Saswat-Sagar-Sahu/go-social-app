import { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import AuthLayout from '../components/AuthLayout'
import { apiRequest } from '../lib/api'

export default function ActivatePage() {
  const location = useLocation()
  const navigate = useNavigate()
  const [token, setToken] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const registeredEmail = location.state?.email

  async function onSubmit(event) {
    event.preventDefault()
    setError('')
    setLoading(true)

    try {
      await apiRequest('/v1/users/activate', {
        method: 'POST',
        body: { token },
      })

      navigate('/login', {
        replace: true,
        state: { activated: true },
      })
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthLayout
      title="Activate account"
      subtitle={
        registeredEmail
          ? `Enter the activation code generated for ${registeredEmail}.`
          : 'Enter the activation code generated during registration.'
      }
      footer={
        <p>
          Back to <Link to="/login">Login</Link>
        </p>
      }
    >
      <form onSubmit={onSubmit} className="form-grid">
        <label htmlFor="token">Activation Code</label>
        <input
          id="token"
          name="token"
          value={token}
          onChange={(event) => setToken(event.target.value)}
          placeholder="Paste token"
          required
        />

        {error && <p className="form-error">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Activating...' : 'Activate'}
        </button>
      </form>
    </AuthLayout>
  )
}
