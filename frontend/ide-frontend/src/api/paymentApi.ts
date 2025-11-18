import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8090/api/v1'

const paymentApiClient = axios.create({
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
paymentApiClient.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export interface CreateCheckoutRequest {
  tier: 'premium'
  provider: 'stripe' | 'paypal'
}

export interface CreateCheckoutResponse {
  checkout_url?: string
  session_id?: string
  subscription_id?: string
  status?: string
}

export interface Payment {
  id: string
  user_id: string
  subscription_id?: string
  amount: number
  currency: string
  payment_provider: 'stripe' | 'paypal' | 'usdt'
  payment_provider_id: string
  status: 'pending' | 'completed' | 'failed' | 'refunded'
  metadata?: Record<string, any>
  created_at: string
  updated_at: string
}

export interface PaymentHistoryResponse {
  payments: Payment[]
}

export const paymentApi = {
  // 創建付費 session
  async createCheckout(request: CreateCheckoutRequest): Promise<CreateCheckoutResponse> {
    const response = await paymentApiClient.post<CreateCheckoutResponse>('/payments/create-checkout', request)
    return response.data
  },

  // 取得付費歷史
  async getPaymentHistory(): Promise<PaymentHistoryResponse> {
    const response = await paymentApiClient.get<PaymentHistoryResponse>('/payments/history')
    return response.data
  },
}

