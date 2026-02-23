import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { clearToken } from '../lib/auth'

const navItems = [
  { to: '/horizon', label: 'Horizon', caption: 'Live conversations' },
  { to: '/studio', label: 'Studio', caption: 'Create and manage' },
  { to: '/atlas', label: 'Atlas', caption: 'Discover people' },
]

export default function AppShell() {
  const navigate = useNavigate()
  const today = new Date().toLocaleDateString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  })

  function logout() {
    clearToken()
    navigate('/login', { replace: true })
  }

  return (
    <div className="app-bg">
      <div className="orb orb-one" />
      <div className="orb orb-two" />
      <div className="shell-wrap">
        <header className="topbar">
          <div className="topbar-copy">
            <p className="kicker">GO SOCIAL</p>
            <h1>Community Hub</h1>
            <p className="topbar-subtitle">Create, discover, and engage without switching context.</p>
          </div>
          <div className="topbar-actions">
            <span className="date-pill">{today}</span>
            <button type="button" className="secondary-btn" onClick={logout}>
              Logout
            </button>
          </div>
        </header>

        <nav className="nav-grid" aria-label="Primary">
          {navItems.map((item) => (
            <NavLink key={item.to} to={item.to} className={({ isActive }) => (isActive ? 'nav-card active' : 'nav-card')}>
              <span>{item.label}</span>
              <small>{item.caption}</small>
            </NavLink>
          ))}
        </nav>

        <main className="page-panel">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
