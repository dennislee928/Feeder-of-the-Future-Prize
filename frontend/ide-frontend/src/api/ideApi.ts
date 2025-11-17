import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface Topology {
  id: string
  name: string
  description?: string
  profile_type: 'rural' | 'suburban' | 'urban'
  nodes: Node[]
  lines: Line[]
  created_at: string
  updated_at: string
}

export interface Node {
  id: string
  type: string
  name: string
  position: {
    x: number
    y: number
  }
  properties?: Record<string, any>
}

export interface Line {
  id: string
  from_node_id: string
  to_node_id: string
  name?: string
  properties?: Record<string, any>
}

export interface CreateTopologyRequest {
  name: string
  description?: string
  profile_type: 'rural' | 'suburban' | 'urban'
  nodes: Node[]
  lines: Line[]
}

export interface Profile {
  type: string
  name: string
  characteristics: {
    load_composition: {
      residential: number
      commercial: number
      industrial: number
    }
    typical_feeder_length_km: number
    typical_node_count: number
    target_saidi_minutes_per_year: number
    target_saifi_interruptions_per_year: number
    der_penetration_range: {
      min: number
      max: number
    }
    ev_penetration_range: {
      min: number
      max: number
    }
  }
}

export const ideApi = {
  // Topology APIs
  async createTopology(data: CreateTopologyRequest): Promise<Topology> {
    const response = await apiClient.post<Topology>('/topologies', data)
    return response.data
  },

  async getTopology(id: string): Promise<Topology> {
    const response = await apiClient.get<Topology>(`/topologies/${id}`)
    return response.data
  },

  async updateTopology(id: string, data: Partial<CreateTopologyRequest>): Promise<Topology> {
    const response = await apiClient.put<Topology>(`/topologies/${id}`, data)
    return response.data
  },

  async deleteTopology(id: string): Promise<void> {
    await apiClient.delete(`/topologies/${id}`)
  },

  async listTopologies(): Promise<Topology[]> {
    const response = await apiClient.get<Topology[]>('/topologies')
    return response.data
  },

  // Profile APIs
  async listProfiles(): Promise<Profile[]> {
    const response = await apiClient.get<Profile[]>('/profiles')
    return response.data
  },

  async getProfile(type: string): Promise<Profile> {
    const response = await apiClient.get<Profile>(`/profiles/${type}`)
    return response.data
  },
}

