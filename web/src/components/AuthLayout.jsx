export default function AuthLayout({ title, subtitle, children, footer }) {
  return (
    <div className="page-wrap">
      <div className="auth-card">
        <h1>{title}</h1>
        <p className="subtitle">{subtitle}</p>
        {children}
        {footer && <div className="form-footer">{footer}</div>}
      </div>
    </div>
  )
}
