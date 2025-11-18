import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { authApi, User, Subscription, UserQuota } from '../api/authApi'

interface AuthContextType {
  user: User | null
  subscription: Subscription | null
  quota: UserQuota | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (provider: 'google' | 'github', code: string) => Promise<void>
  logout: () => void
  refreshUser: () => Promise<void>
  getAuthURL: (provider: 'google' | 'github') => Promise<string>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [subscription, setSubscription] = useState<Subscription | null>(null)
  const [quota, setQuota] = useState<UserQuota | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  // 檢查是否已登入並載入用戶資訊
  useEffect(() => {
    const loadUser = async () => {
      if (authApi.isAuthenticated()) {
        try {
          const data = await authApi.getMe()
          setUser(data.user)
          setSubscription(data.subscription || null)
          setQuota(data.quota || null)
        } catch (error) {
          // Token 可能已過期，清除
          authApi.logout()
        }
      }
      setIsLoading(false)
    }

    loadUser()
  }, [])

  const login = async (provider: 'google' | 'github', code: string) => {
    try {
      const response = await authApi.oauthCallback({ provider, code })
      setUser(response.user)
      // 重新載入完整資訊
      await refreshUser()
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const logout = () => {
    authApi.logout()
    setUser(null)
    setSubscription(null)
    setQuota(null)
  }

  const refreshUser = async () => {
    if (!authApi.isAuthenticated()) {
      return
    }

    try {
      const data = await authApi.getMe()
      setUser(data.user)
      setSubscription(data.subscription || null)
      setQuota(data.quota || null)
    } catch (error) {
      console.error('Failed to refresh user:', error)
      // Token 可能已過期，登出
      logout()
    }
  }

  const getAuthURL = async (provider: 'google' | 'github'): Promise<string> => {
    const response = await authApi.getAuthURL({ provider })
    return response.auth_url
  }

  const value: AuthContextType = {
    user,
    subscription,
    quota,
    isLoading,
    isAuthenticated: !!user,
    login,
    logout,
    refreshUser,
    getAuthURL,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

