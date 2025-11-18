import React from 'react'
import ReactDOM from 'react-dom/client'
import './i18n/config'
import { AuthProvider } from './contexts/AuthContext'
import App from './App.tsx'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <AuthProvider>
      <App />
    </AuthProvider>
  </React.StrictMode>,
)

