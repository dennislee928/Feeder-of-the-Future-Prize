import axios from 'axios'

// 在 Render 上，使用環境變數
// 如果沒有設置環境變數，開發環境使用 localhost，生產環境使用空字符串（相對路徑）
const SIM_API_BASE_URL = import.meta.env.VITE_SIM_API_BASE_URL || 
  (import.meta.env.MODE === 'production' ? 'https://feeder-sim-engine.onrender.com' : 'http://localhost:8081')

const penetrationApiClient = axios.create({
  baseURL: SIM_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface AttackScenario {
  id: string
  name: string
  layer: number
  severity: 'low' | 'medium' | 'high' | 'critical'
  description: string
}

export interface AttackPath {
  from: string
  to: string
  layer: string
}

export interface PenetrationAttackResult {
  attack_id: string
  scenario: string
  scenario_name: string
  layer: number
  severity: 'low' | 'medium' | 'high' | 'critical'
  successful: boolean
  affected_nodes: string[]
  affected_lines: string[]
  attack_path: AttackPath[]
  impact: string
  recommendations: string[]
  timestamp: string
}

export interface PenetrationTestResult {
  attacks: PenetrationAttackResult[]
  summary: {
    total_attacks: number
    successful: number
    failed: number
    critical_vulnerabilities: number
    total_nodes: number
    total_lines: number
    affected_nodes_count: number
    affected_lines_count: number
  }
}

export interface PenetrationTestRequest {
  topology: {
    nodes: any[]
    lines: any[]
    profile_type?: string
  }
  attack_scenarios: string[]
  target_nodes?: string[]
}

export const penetrationApi = {
  async runPenetrationTest(request: PenetrationTestRequest): Promise<PenetrationTestResult> {
    const response = await penetrationApiClient.post<PenetrationTestResult>('/simulate/penetration', request)
    return response.data
  },

  async getAvailableScenarios(): Promise<AttackScenario[]> {
    const response = await penetrationApiClient.get<{ scenarios: AttackScenario[] }>('/simulate/penetration/scenarios')
    return response.data.scenarios
  },
}

