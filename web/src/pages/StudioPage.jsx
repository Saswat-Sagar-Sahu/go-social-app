import { useCallback, useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import PostImage from '../components/PostImage'
import { apiRequest } from '../lib/api'
import { clearToken, getUserIdFromToken } from '../lib/auth'

function toTagsInput(tags) {
  if (!Array.isArray(tags)) return ''
  return tags.join(', ')
}

function parseTags(input) {
  return input
    .split(',')
    .map((tag) => tag.trim())
    .filter(Boolean)
}

export default function StudioPage() {
  const navigate = useNavigate()
  const currentUserId = useMemo(() => getUserIdFromToken(), [])

  const [form, setForm] = useState({ title: '', content: '', image_url: '', tags: '' })
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const [posts, setPosts] = useState([])
  const [loadingPosts, setLoadingPosts] = useState(true)
  const [editById, setEditById] = useState({})
  const [postActionLoading, setPostActionLoading] = useState({})

  const handleUnauthorized = useCallback(
    (err) => {
      if (err?.status === 401) {
        clearToken()
        navigate('/login', { replace: true })
      }
    },
    [navigate],
  )

  const loadMyPosts = useCallback(async () => {
    setLoadingPosts(true)
    setError('')

    try {
      const data = await apiRequest('/v1/posts/me', { auth: true })
      setPosts(Array.isArray(data) ? data : [])
    } catch (err) {
      setError(err?.message || 'Failed to load your posts')
      handleUnauthorized(err)
    } finally {
      setLoadingPosts(false)
    }
  }, [handleUnauthorized])

  useEffect(() => {
    loadMyPosts()
  }, [loadMyPosts])

  function onChange(event) {
    const { name, value } = event.target
    setForm((prev) => ({ ...prev, [name]: value }))
  }

  async function createPost(event) {
    event.preventDefault()
    setCreating(true)
    setError('')
    setSuccess('')

    try {
      const created = await apiRequest('/v1/posts', {
        method: 'POST',
        auth: true,
        body: {
          title: form.title,
          content: form.content,
          image_url: form.image_url.trim(),
          tags: parseTags(form.tags),
        },
      })
      setPosts((prev) => [created, ...prev])
      setForm({ title: '', content: '', image_url: '', tags: '' })
      setSuccess('Post published in your stream.')
    } catch (err) {
      setError(err?.message || 'Failed to publish post')
      handleUnauthorized(err)
    } finally {
      setCreating(false)
    }
  }

  function startEdit(post) {
    setEditById((prev) => ({
      ...prev,
      [post.id]: {
        title: post.title,
        content: post.content,
        image_url: post.image_url || '',
        tags: toTagsInput(post.tags),
      },
    }))
  }

  async function saveEdit(postId) {
    const edit = editById[postId]
    if (!edit) return

    setPostActionLoading((prev) => ({ ...prev, [`save-${postId}`]: true }))
    setError('')

    try {
      const updated = await apiRequest(`/v1/posts/${postId}`, {
        method: 'PUT',
        auth: true,
        body: {
          title: edit.title,
          content: edit.content,
          image_url: edit.image_url.trim(),
          tags: parseTags(edit.tags),
        },
      })

      setPosts((prev) => prev.map((post) => (post.id === postId ? { ...post, ...updated } : post)))
      setEditById((prev) => ({ ...prev, [postId]: undefined }))
    } catch (err) {
      setError(err?.message || 'Failed to update post')
      handleUnauthorized(err)
    } finally {
      setPostActionLoading((prev) => ({ ...prev, [`save-${postId}`]: false }))
    }
  }

  async function deletePost(postId) {
    if (!window.confirm('Delete this post permanently?')) {
      return
    }

    let removedPost
    setPostActionLoading((prev) => ({ ...prev, [`delete-${postId}`]: true }))
    setError('')

    setPosts((prev) => {
      removedPost = prev.find((post) => post.id === postId)
      return prev.filter((post) => post.id !== postId)
    })

    try {
      await apiRequest(`/v1/posts/${postId}`, {
        method: 'DELETE',
        auth: true,
      })
    } catch (err) {
      if (removedPost) {
        setPosts((prev) => [removedPost, ...prev])
      }
      setError(err?.message || 'Failed to delete post')
      handleUnauthorized(err)
    } finally {
      setPostActionLoading((prev) => ({ ...prev, [`delete-${postId}`]: false }))
    }
  }

  return (
    <section>
      <div className="section-banner">
        <div className="section-head">
          <div>
            <h2>Creator Studio</h2>
            <p>Publish with visuals and manage your posts from one place.</p>
          </div>
        </div>
        {!loadingPosts && <p className="banner-meta">{posts.length} posts authored by you</p>}
      </div>

      <article className="hero-card">
        <h3>Publish a new post</h3>
        <p className="status-text">Tip: add an image URL for richer cards in the stream.</p>
        <form onSubmit={createPost} className="form-grid">
          <label htmlFor="studio-title">Title</label>
          <input id="studio-title" name="title" value={form.title} onChange={onChange} required />

          <label htmlFor="studio-content">Content</label>
          <textarea
            id="studio-content"
            name="content"
            rows={5}
            maxLength={1000}
            value={form.content}
            onChange={onChange}
            required
          />

          <label htmlFor="studio-image-url">Image URL (optional)</label>
          <input
            id="studio-image-url"
            name="image_url"
            type="url"
            value={form.image_url}
            onChange={onChange}
            placeholder="https://example.com/image.jpg"
          />

          <label htmlFor="studio-tags">Tags</label>
          <input
            id="studio-tags"
            name="tags"
            value={form.tags}
            onChange={onChange}
            placeholder="react, golang, community"
          />

          {error && <p className="form-error">{error}</p>}
          {success && <p className="form-success">{success}</p>}

          <button type="submit" disabled={creating}>
            {creating ? 'Publishing...' : 'Publish Post'}
          </button>
        </form>
      </article>

      <article className="panel-card">
          <div className="section-head">
            <div>
              <h3>Your Posts</h3>
              <p>Edit content, update image, or remove outdated posts.</p>
            </div>
          <button type="button" className="secondary-btn" onClick={loadMyPosts}>
            Refresh
          </button>
        </div>

        {loadingPosts && <p className="status-text">Loading your posts...</p>}
        {!loadingPosts && posts.length === 0 && <p className="status-text">You have not posted yet.</p>}

        <div className="post-grid">
          {posts.map((post) => {
            const canManage = currentUserId !== null && post.user_id === currentUserId
            const edit = editById[post.id]
            const isEditing = Boolean(edit)

            return (
              <article key={post.id} className="post-card">
                {!isEditing && (
                  <>
                    <h3>{post.title}</h3>
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
                  </>
                )}

                {isEditing && (
                  <div className="inline-edit">
                    <label htmlFor={`title-${post.id}`}>Title</label>
                    <input
                      id={`title-${post.id}`}
                      value={edit.title}
                      onChange={(event) =>
                        setEditById((prev) => ({
                          ...prev,
                          [post.id]: { ...edit, title: event.target.value },
                        }))
                      }
                    />
                    <label htmlFor={`content-${post.id}`}>Content</label>
                    <textarea
                      id={`content-${post.id}`}
                      rows={4}
                      maxLength={1000}
                      value={edit.content}
                      onChange={(event) =>
                        setEditById((prev) => ({
                          ...prev,
                          [post.id]: { ...edit, content: event.target.value },
                        }))
                      }
                    />
                    <label htmlFor={`image-url-${post.id}`}>Image URL</label>
                    <input
                      id={`image-url-${post.id}`}
                      type="url"
                      value={edit.image_url}
                      onChange={(event) =>
                        setEditById((prev) => ({
                          ...prev,
                          [post.id]: { ...edit, image_url: event.target.value },
                        }))
                      }
                    />
                    <label htmlFor={`tags-${post.id}`}>Tags</label>
                    <input
                      id={`tags-${post.id}`}
                      value={edit.tags}
                      onChange={(event) =>
                        setEditById((prev) => ({
                          ...prev,
                          [post.id]: { ...edit, tags: event.target.value },
                        }))
                      }
                    />
                  </div>
                )}

                <div className="inline-actions">
                  {canManage && !isEditing && (
                    <button type="button" className="secondary-btn" onClick={() => startEdit(post)}>
                      Edit
                    </button>
                  )}

                  {canManage && isEditing && (
                    <>
                      <button
                        type="button"
                        onClick={() => saveEdit(post.id)}
                        disabled={postActionLoading[`save-${post.id}`]}
                      >
                        {postActionLoading[`save-${post.id}`] ? 'Saving...' : 'Save'}
                      </button>
                      <button
                        type="button"
                        className="secondary-btn"
                        onClick={() => setEditById((prev) => ({ ...prev, [post.id]: undefined }))}
                      >
                        Cancel
                      </button>
                    </>
                  )}

                  {canManage && (
                    <button
                      type="button"
                      className="danger-btn"
                      onClick={() => deletePost(post.id)}
                      disabled={postActionLoading[`delete-${post.id}`]}
                    >
                      {postActionLoading[`delete-${post.id}`] ? 'Deleting...' : 'Delete'}
                    </button>
                  )}
                </div>
              </article>
            )
          })}
        </div>
      </article>
    </section>
  )
}
