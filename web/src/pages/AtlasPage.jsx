import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { apiRequest } from '../lib/api'
import { clearToken, getUserIdFromToken } from '../lib/auth'

export default function AtlasPage() {
  const navigate = useNavigate()
  const currentUserId = useMemo(() => getUserIdFromToken(), [])

  const [query, setQuery] = useState('')
  const [activeQuery, setActiveQuery] = useState('')
  const [page, setPage] = useState(1)
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [total, setTotal] = useState(0)
  const [totalPages, setTotalPages] = useState(0)
  const [followState, setFollowState] = useState({})
  const [actionLoading, setActionLoading] = useState({})

  const handleUnauthorized = useCallback(
    (err) => {
      if (err?.status === 401) {
        clearToken()
        navigate('/login', { replace: true })
      }
    },
    [navigate],
  )

  const loadUsers = useCallback(async () => {
    setLoading(true)
    setError('')

    try {
      const params = new URLSearchParams({
        page: String(page),
        page_size: '10',
      })
      if (activeQuery) params.set('name', activeQuery)

      const data = await apiRequest(`/v1/users?${params.toString()}`, { auth: true })
      setUsers(Array.isArray(data?.data) ? data.data : [])
      setTotal(data?.total || 0)
      setTotalPages(data?.total_pages || 0)
    } catch (err) {
      setError(err?.message || 'Failed to load users')
      handleUnauthorized(err)
    } finally {
      setLoading(false)
    }
  }, [activeQuery, handleUnauthorized, page])

  useEffect(() => {
    loadUsers()
  }, [loadUsers])

  function onSearch(event) {
    event.preventDefault()
    setPage(1)
    setActiveQuery(query.trim())
  }

  async function follow(userId) {
    setActionLoading((prev) => ({ ...prev, [userId]: 'follow' }))
    try {
      await apiRequest(`/v1/users/${userId}/follow`, { method: 'POST', auth: true })
      setFollowState((prev) => ({ ...prev, [userId]: true }))
    } catch (err) {
      if (err?.status === 409) {
        setFollowState((prev) => ({ ...prev, [userId]: true }))
      } else {
        setError(err?.message || 'Follow failed')
      }
      handleUnauthorized(err)
    } finally {
      setActionLoading((prev) => ({ ...prev, [userId]: '' }))
    }
  }

  async function unfollow(userId) {
    setActionLoading((prev) => ({ ...prev, [userId]: 'unfollow' }))
    try {
      await apiRequest(`/v1/users/${userId}/unfollow`, { method: 'POST', auth: true })
      setFollowState((prev) => ({ ...prev, [userId]: false }))
    } catch (err) {
      if (err?.status === 404) {
        setFollowState((prev) => ({ ...prev, [userId]: false }))
      } else {
        setError(err?.message || 'Unfollow failed')
      }
      handleUnauthorized(err)
    } finally {
      setActionLoading((prev) => ({ ...prev, [userId]: '' }))
    }
  }

  return (
    <section>
      <div className="section-banner">
        <div className="section-head">
          <div>
            <h2>People Atlas</h2>
            <p>Find creators by name and tune your feed with follow controls.</p>
          </div>
        </div>
        {!loading && <p className="banner-meta">{users.length} users on this page</p>}
      </div>

      <article className="hero-card">
        <form className="search-row" onSubmit={onSearch}>
          <input
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search by username"
          />
          <button type="submit">Search</button>
        </form>

        <p className="status-text atlas-meta">
          Results: {total} | Page {page}
          {totalPages > 0 ? ` of ${totalPages}` : ''}
        </p>
      </article>

      {error && <p className="form-error">{error}</p>}
      {loading && <p className="status-text">Loading users...</p>}

      <div className="users-grid">
        {!loading &&
          users.map((user) => {
            const isSelf = currentUserId !== null && currentUserId === user.id
            const knownFollow = followState[user.id]

            return (
              <article key={user.id} className="user-card">
                <h3>{user.username}</h3>
                <p>{user.email}</p>
                <p className="status-text">User ID: {user.id}</p>
                <p className="status-text">{user.activated ? 'Activated' : 'Activation pending'}</p>
                {knownFollow === true && <p className="form-success">Following</p>}
                {knownFollow === false && <p className="status-text">Not following</p>}

                <div className="inline-actions">
                  <button
                    type="button"
                    onClick={() => follow(user.id)}
                    disabled={isSelf || actionLoading[user.id] === 'follow'}
                  >
                    {actionLoading[user.id] === 'follow' ? 'Following...' : 'Follow'}
                  </button>
                  <button
                    type="button"
                    className="secondary-btn"
                    onClick={() => unfollow(user.id)}
                    disabled={isSelf || actionLoading[user.id] === 'unfollow'}
                  >
                    {actionLoading[user.id] === 'unfollow' ? 'Unfollowing...' : 'Unfollow'}
                  </button>
                </div>
              </article>
            )
          })}
      </div>

      {!loading && users.length === 0 && (
        <p className="status-text empty-state">No users found. Try a shorter keyword or clear the search.</p>
      )}

      <div className="pagination-row">
        <button type="button" className="secondary-btn" onClick={() => setPage((p) => p - 1)} disabled={page <= 1}>
          Previous
        </button>
        <button
          type="button"
          className="secondary-btn"
          onClick={() => setPage((p) => p + 1)}
          disabled={totalPages > 0 ? page >= totalPages : users.length < 10}
        >
          Next
        </button>
      </div>
    </section>
  )
}
