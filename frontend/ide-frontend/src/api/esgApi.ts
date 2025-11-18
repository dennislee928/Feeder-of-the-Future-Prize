import axios from 'axios'

const SIM_API_BASE_URL = import.meta.env.VITE_SIM_API_BASE_URL || 'http://localhost:8081'

const esgApiClient = axios.create({
  baseURL: SIM_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface NodeEmission {
  node_id: string
  node_type: string
  power_kw: number
  energy_kwh: number
  emission_kg_co2: number
}

export interface ESGRecommendation {
  type: string
  priority: 'low' | 'medium' | 'high'
  title: string
  description: string
  estimated_reduction_ton?: number
}

export interface ESGCalculationResult {
  timestamp: string
  time_hours: number
  total_emissions_kg_co2: number
  total_emissions_ton_co2: number
  carbon_credits_ton: number
  carbon_credit_value_usd: number
  esg_score: number
  node_emissions: NodeEmission[]
  recommendations: ESGRecommendation[]
}

export interface ESGCalculationRequest {
  topology: {
    nodes: any[]
    lines: any[]
  }
  parameters?: {
    time_hours?: number
    ev_charging_hours?: number
    solar_generation_hours?: number
    battery_cycles?: number
  }
}

export const esgApi = {
  async calculateESG(request: ESGCalculationRequest): Promise<ESGCalculationResult> {
    const response = await esgApiClient.post<ESGCalculationResult>('/simulate/esg', request)
    return response.data
  },
}

