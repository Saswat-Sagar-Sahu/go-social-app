import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import PostImage from '../components/PostImage'
import { apiRequest } from '../lib/api'
import { clearToken, getUserIdFromToken } from '../lib/auth'

const DEFAULT_VISIBLE_COMMENTS = 1

function formatDate(value) {
  const dt = new Date(value)
  return Number.isNaN(dt.getTime()) ? value : dt.toLocaleString()
}

function normalizeErrorMessage(err) {
  return err?.message || 'Something went wrong'
}

export default function HorizonPage() {
  const navigate = useNavigate()
  const currentUserId = useMemo(() => getUserIdFromToken(), [])

  const [posts, setPosts] = useState([])
  const [loadingPosts, setLoadingPosts] = useState(true)
  const [error, setError] = useState('')

  const [commentsByPost, setCommentsByPost] = useState({})
  const [openComments, setOpenComments] = useState({})
  const [loadingComments, setLoadingComments] = useState({})
  const [visibleCommentsCount, setVisibleCommentsCount] = useState({})
  const [commentDraft, setCommentDraft] = useState({})
  const [commentActionLoading, setCommentActionLoading] = useState({})
  const [editCommentById, setEditCommentById] = useState({})

  const handleUnauthorized = useCallback(
    (err) => {
      if (err?.status === 401) {
        clearToken()
        navigate('/login', { replace: true })
      }
    },
    [navigate],
  )

  const loadFeed = useCallback(async () => {
    setLoadingPosts(true)
    setError('')

    try {
      const data = await apiRequest('/v1/users/feed', { auth: true })
      setPosts(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(normalizeErrorMessage(err))
      handleUnauthorized(err)
    } finally {
      setLoadingPosts(false)
    }
  }, [handleUnauthorized])

  const loadCommentsForPost = useCallback(
    async (postId) => {
      setLoadingComments((prev) => ({ ...prev, [postId]: true }))
      setError('')

      try {
        const data = await apiRequest(`/v1/comments/post/${postId}`, { auth: true })
        const comments = Array.isArray(data) ? data : []
        setCommentsByPost((prev) => ({ ...prev, [postId]: comments }))
        setVisibleCommentsCount((prev) => ({ ...prev, [postId]: DEFAULT_VISIBLE_COMMENTS }))
      } catch (err) {
        if (err?.status === 404) {
          setCommentsByPost((prev) => ({ ...prev, [postId]: [] }))
          setVisibleCommentsCount((prev) => ({ ...prev, [postId]: DEFAULT_VISIBLE_COMMENTS }))
        } else {
          setError(normalizeErrorMessage(err))
          handleUnauthorized(err)
        }
      } finally {
        setLoadingComments((prev) => ({ ...prev, [postId]: false }))
      }
    },
    [handleUnauthorized],
  )

  useEffect(() => {
    loadFeed()
  }, [loadFeed])

  async function toggleComments(postId) {
    const nextOpen = !openComments[postId]
    setOpenComments((prev) => ({ ...prev, [postId]: nextOpen }))

    if (nextOpen && !commentsByPost[postId]) {
      await loadCommentsForPost(postId)
    }
  }

  async function createComment(postId) {
    const content = (commentDraft[postId] || '').trim()
    if (!content) return

    setCommentActionLoading((prev) => ({ ...prev, [`create-${postId}`]: true }))
    setError('')

    try {
      const created = await apiRequest('/v1/comments/', {
        method: 'POST',
        auth: true,
        body: { post_id: postId, content },
      })

      setCommentDraft((prev) => ({ ...prev, [postId]: '' }))
      setCommentsByPost((prev) => {
        const current = prev[postId] || []
        return { ...prev, [postId]: [created, ...current] }
      })
      setVisibleCommentsCount((prev) => ({
        ...prev,
        [postId]: Math.max(DEFAULT_VISIBLE_COMMENTS, prev[postId] || DEFAULT_VISIBLE_COMMENTS),
      }))
    } catch (err) {
      setError(normalizeErrorMessage(err))
      handleUnauthorized(err)
    } finally {
      setCommentActionLoading((prev) => ({ ...prev, [`create-${postId}`]: false }))
    }
  }

  async function updateComment(comment) {
    const nextContent = (editCommentById[comment.id] || '').trim()
    if (!nextContent) return

    setCommentActionLoading((prev) => ({ ...prev, [`update-${comment.id}`]: true }))
    setError('')

    try {
      const updated = await apiRequest(`/v1/comments/${comment.id}`, {
        method: 'PUT',
        auth: true,
        body: {
          content: nextContent,
          post_id: comment.post_id,
        },
      })

      setEditCommentById((prev) => ({ ...prev, [comment.id]: undefined }))
      setCommentsByPost((prev) => ({
        ...prev,
        [comment.post_id]: (prev[comment.post_id] || []).map((item) =>
          item.id === comment.id ? { ...item, ...updated } : item,
        ),
      }))
    } catch (err) {
      setError(normalizeErrorMessage(err))
      handleUnauthorized(err)
    } finally {
      setCommentActionLoading((prev) => ({ ...prev, [`update-${comment.id}`]: false }))
    }
  }

  async function deleteComment(comment) {
    setCommentActionLoading((prev) => ({ ...prev, [`delete-${comment.id}`]: true }))
    setError('')

    try {
      await apiRequest(`/v1/comments/${comment.id}`, {
        method: 'DELETE',
        auth: true,
      })

      setCommentsByPost((prev) => ({
        ...prev,
        [comment.post_id]: (prev[comment.post_id] || []).filter((item) => item.id !== comment.id),
      }))
    } catch (err) {
      setError(normalizeErrorMessage(err))
      handleUnauthorized(err)
    } finally {
      setCommentActionLoading((prev) => ({ ...prev, [`delete-${comment.id}`]: false }))
    }
  }

  return (
    <section>
      <div className="section-banner">
        <div className="section-head">
          <div>
            <h2>Horizon Stream</h2>
            <p>Stories and conversations from people you follow.</p>
          </div>
          <button type="button" className="secondary-btn" onClick={loadFeed}>
            Refresh Stream
          </button>
        </div>
        {!loadingPosts && <p className="banner-meta">{posts.length} posts in your current stream</p>}
      </div>

      {error && <p className="form-error">{error}</p>}
      {loadingPosts && <p className="status-text">Loading stream...</p>}
      {!loadingPosts && posts.length === 0 && <p className="status-text">No feed posts yet.</p>}

      <div className="post-grid">
        {posts.map((post) => {
          const comments = commentsByPost[post.id] || []
          const visibleCount = visibleCommentsCount[post.id] || DEFAULT_VISIBLE_COMMENTS
          const isOpen = Boolean(openComments[post.id])
          const visibleComments = comments.slice(0, visibleCount)
          const hasHiddenComments = comments.length > visibleComments.length

          return (
            <article key={post.id} className="post-card">
              <header className="post-header">
                <h3>{post.title}</h3>
                <small>{formatDate(post.created_at)}</small>
              </header>

              <PostImage src={post.image_url} alt={post.title} />
              <p className="post-content">{post.content}</p>

              {post.tags?.length > 0 && (
                <div className="tag-row">
                  {post.tags.map((tag) => (
                    <span key={`${post.id}-${tag}`} className="tag-pill">
                      #{tag}
                    </span>
                  ))}
                </div>
              )}

              <div className="card-meta">
                <span>Author ID: {post.user_id ?? 'n/a'}</span>
                <button type="button" className="link-btn" onClick={() => toggleComments(post.id)}>
                  {isOpen ? 'Hide comments' : 'View comments'}
                </button>
              </div>

              {isOpen && (
                <div className="comments-block">
                  {loadingComments[post.id] && <p className="status-text">Loading replies...</p>}

                  {!loadingComments[post.id] && comments.length === 0 && (
                    <p className="status-text">No replies yet. Start the conversation.</p>
                  )}

                  {visibleComments.map((comment) => {
                    const canModify = currentUserId !== null && comment.user_id === currentUserId
                    const editValue = editCommentById[comment.id]
                    const editing = typeof editValue === 'string'

                    return (
                      <div key={comment.id} className="comment-item">
                        <div className="comment-head">
                          <strong>User {comment.user_id}</strong>
                          <small>{formatDate(comment.created_at)}</small>
                        </div>

                        {!editing && <p>{comment.content}</p>}

                        {editing && (
                          <div className="inline-edit">
                            <textarea
                              rows={3}
                              value={editValue}
                              onChange={(event) =>
                                setEditCommentById((prev) => ({ ...prev, [comment.id]: event.target.value }))
                              }
                            />
                            <div className="inline-actions">
                              <button
                                type="button"
                                onClick={() => updateComment(comment)}
                                disabled={commentActionLoading[`update-${comment.id}`]}
                              >
                                {commentActionLoading[`update-${comment.id}`] ? 'Saving...' : 'Save'}
                              </button>
                              <button
                                type="button"
                                className="secondary-btn"
                                onClick={() => setEditCommentById((prev) => ({ ...prev, [comment.id]: undefined }))}
                              >
                                Cancel
                              </button>
                            </div>
                          </div>
                        )}

                        {canModify && !editing && (
                          <div className="inline-actions">
                            <button
                              type="button"
                              className="secondary-btn"
                              onClick={() => setEditCommentById((prev) => ({ ...prev, [comment.id]: comment.content }))}
                            >
                              Edit
                            </button>
                            <button
                              type="button"
                              className="danger-btn"
                              onClick={() => deleteComment(comment)}
                              disabled={commentActionLoading[`delete-${comment.id}`]}
                            >
                              {commentActionLoading[`delete-${comment.id}`] ? 'Deleting...' : 'Delete'}
                            </button>
                          </div>
                        )}
                      </div>
                    )
                  })}

                  {hasHiddenComments && (
                    <button
                      type="button"
                      className="link-btn"
                      onClick={() => setVisibleCommentsCount((prev) => ({ ...prev, [post.id]: comments.length }))}
                    >
                      View all {comments.length} replies
                    </button>
                  )}

                  {!hasHiddenComments && comments.length > DEFAULT_VISIBLE_COMMENTS && (
                    <button
                      type="button"
                      className="link-btn"
                      onClick={() =>
                        setVisibleCommentsCount((prev) => ({ ...prev, [post.id]: DEFAULT_VISIBLE_COMMENTS }))
                      }
                    >
                      Show fewer replies
                    </button>
                  )}

                  <div className="comment-create">
                    <textarea
                      rows={3}
                      placeholder="Write a reply..."
                      value={commentDraft[post.id] || ''}
                      onChange={(event) =>
                        setCommentDraft((prev) => ({
                          ...prev,
                          [post.id]: event.target.value,
                        }))
                      }
                    />
                    <button
                      type="button"
                      onClick={() => createComment(post.id)}
                      disabled={commentActionLoading[`create-${post.id}`]}
                    >
                      {commentActionLoading[`create-${post.id}`] ? 'Posting...' : 'Add Reply'}
                    </button>
                  </div>
                </div>
              )}
            </article>
          )
        })}
      </div>
    </section>
  )
}
