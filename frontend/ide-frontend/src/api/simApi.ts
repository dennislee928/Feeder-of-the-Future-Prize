import axios from 'axios'

// 在 Render 上，使用環境變數
const SIM_API_BASE_URL = import.meta.env.VITE_SIM_API_BASE_URL || 
  (import.meta.env.MODE === 'production' ? 'https://feeder-sim-engine.onrender.com' : 'http://localhost:8081')

const simApiClient = axios.create({
  baseURL: SIM_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface PowerflowResult {
  nodes: Array<{
    node_id: string
    voltage_pu: number
    voltage_kv: number
    voltage_deviation_percent: number
    status: 'normal' | 'warning' | 'critical'
  }>
  lines: Array<{
    line_id: string
    loading_percent: number
    status: 'normal' | 'warning' | 'critical'
  }>
  summary: {
    average_voltage_pu: number
    max_line_loading_percent: number
    total_nodes: number
    total_lines: number
    converged: boolean
  }
}

export interface ReliabilityResult {
  saidi: number
  saifi: number
  expected_faults_per_year: number
  average_repair_time_minutes: number
  node_risks: Array<{
    node_id: string
    risk_score: number
    risk_level: 'low' | 'medium' | 'high'
  }>
  line_risks: Array<{
    line_id: string
    risk_score: number
    risk_level: 'low' | 'medium' | 'high'
  }>
  summary: {
    total_nodes: number
    total_lines: number
    estimated_length_km: number
    profile_type: string
  }
}

export interface SimulationRequest {
  topology: {
    nodes: any[]
    lines: any[]
    profile_type?: string
  }
  parameters?: {
    fault_rate?: number
    repair_time?: number
  }
}

export const simApi = {
  async runPowerflow(request: SimulationRequest): Promise<PowerflowResult> {
    const response = await simApiClient.post<PowerflowResult>('/simulate/powerflow', request)
    return response.data
  },

  async runReliability(request: SimulationRequest): Promise<ReliabilityResult> {
    const response = await simApiClient.post<ReliabilityResult>('/simulate/reliability', request)
    return response.data
  },
}

