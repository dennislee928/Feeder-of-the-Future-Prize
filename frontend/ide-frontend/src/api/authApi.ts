import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090/api/v1'

const authApiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 從 localStorage 取得 token
const getToken = () => {
  return localStorage.getItem('auth_token')
}

// 設置 token 到請求頭
authApiClient.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export interface OAuthCallbackRequest {
  provider: 'google' | 'github'
  code: string
  state?: string
}

export interface OAuthResponse {
  token: string
  user: User
  expires_in: number
}

export interface User {
  id: string
  email: string
  name?: string
  avatar_url?: string
  subscription_tier: 'demo' | 'free' | 'premium'
  subscription_status: 'active' | 'cancelled' | 'expired'
  subscription_expires_at?: string
  api_key?: string
  created_at: string
  updated_at: string
}

export interface Subscription {
  id: string
  user_id: string
  tier: 'free' | 'premium'
  status: 'active' | 'cancelled' | 'expired' | 'pending'
  payment_provider?: 'stripe' | 'paypal'
  payment_subscription_id?: string
  current_period_start?: string
  current_period_end?: string
  cancel_at_period_end: boolean
  created_at: string
  updated_at: string
}

export interface UserQuota {
  id: string
  user_id: string
  max_topologies: number
  used_topologies: number
  max_simulations_per_day: number
  used_simulations_today: number
  last_simulation_reset_date: string
  can_use_3d_rendering: boolean
  can_use_ai_prediction: boolean
  can_use_advanced_security: boolean
  can_access_api: boolean
  created_at: string
  updated_at: string
}

export interface GetMeResponse {
  user: User
  subscription?: Subscription
  quota?: UserQuota
}

export interface GetAuthURLRequest {
  provider: 'google' | 'github'
  state?: string
}

export interface GetAuthURLResponse {
  auth_url: string
  state: string
}

export interface RefreshTokenRequest {
  token: string
}

export interface RefreshTokenResponse {
  token: string
  expires_in: number
}

export const authApi = {
  // OAuth 登入
  async oauthCallback(request: OAuthCallbackRequest): Promise<OAuthResponse> {
    const response = await authApiClient.post<OAuthResponse>('/auth/oauth/callback', request)
    // 保存 token
    if (response.data.token) {
      localStorage.setItem('auth_token', response.data.token)
    }
    return response.data
  },

  // 取得 OAuth 授權 URL
  async getAuthURL(request: GetAuthURLRequest): Promise<GetAuthURLResponse> {
    const response = await authApiClient.post<GetAuthURLResponse>('/auth/oauth/url', request)
    return response.data
  },

  // 刷新 token
  async refreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    const response = await authApiClient.post<RefreshTokenResponse>('/auth/refresh', request)
    // 更新 token
    if (response.data.token) {
      localStorage.setItem('auth_token', response.data.token)
    }
    return response.data
  },

  // 取得當前用戶資訊
  async getMe(): Promise<GetMeResponse> {
    const response = await authApiClient.get<GetMeResponse>('/auth/me')
    return response.data
  },

  // 登出
  logout() {
    localStorage.removeItem('auth_token')
  },

  // 檢查是否已登入
  isAuthenticated(): boolean {
    return !!getToken()
  },

  // 取得 token
  getToken(): string | null {
    return getToken()
  },
}

