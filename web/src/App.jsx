import { Navigate, Route, Routes } from 'react-router-dom'
import AppShell from './components/AppShell'
import ProtectedRoute from './components/ProtectedRoute'
import ActivatePage from './pages/ActivatePage'
import AtlasPage from './pages/AtlasPage'
import HorizonPage from './pages/HorizonPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import StudioPage from './pages/StudioPage'
import { isAuthenticated } from './lib/auth'

function HomeRedirect() {
  return <Navigate to={isAuthenticated() ? '/horizon' : '/register'} replace />
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomeRedirect />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/activate" element={<ActivatePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route
        element={
          <ProtectedRoute>
            <AppShell />
          </ProtectedRoute>
        }
      >
        <Route path="/horizon" element={<HorizonPage />} />
        <Route path="/studio" element={<StudioPage />} />
        <Route path="/atlas" element={<AtlasPage />} />
      </Route>
      <Route path="/home" element={<Navigate to="/horizon" replace />} />
      <Route path="/feed" element={<Navigate to="/horizon" replace />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
