package topology

import "time"

// Topology 代表一個配電 feeder 拓樸
type Topology struct {
	ID          string    `json:"id"`
	UserID      *string   `json:"user_id,omitempty"` // 可選，無註冊用戶為 nil
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	ProfileType string    `json:"profile_type"` // rural, suburban, urban
	Nodes       []Node    `json:"nodes"`
	Lines       []Line    `json:"lines"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Node 代表拓樸中的節點（bus）
type Node struct {
	ID       string  `json:"id"`
	Type     string  `json:"type"` // bus, transformer, switch, ev_charger, der
	Name     string  `json:"name"`
	Position Position `json:"position"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Line 代表連接兩個節點的線路
type Line struct {
	ID         string  `json:"id"`
	FromNodeID string  `json:"from_node_id"`
	ToNodeID   string  `json:"to_node_id"`
	Name       string  `json:"name,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Position 代表節點在畫布上的位置
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Transformer 變壓器屬性
type TransformerProperties struct {
	RatedVoltageKV    float64 `json:"rated_voltage_kv"`
	RatedCapacityKVA float64 `json:"rated_capacity_kva"`
	PrimaryVoltage    float64 `json:"primary_voltage"`
	SecondaryVoltage  float64 `json:"secondary_voltage"`
}

// SwitchProperties 開關屬性
type SwitchProperties struct {
	Type        string `json:"type"` // sectionalizer, recloser, breaker
	IsClosed    bool   `json:"is_closed"`
	IsAutomated bool   `json:"is_automated"`
}

// EVChargerProperties EV 充電樁屬性
type EVChargerProperties struct {
	RatedPowerKW    float64 `json:"rated_power_kw"`
	MaxChargingRate float64 `json:"max_charging_rate"`
	IsControllable  bool    `json:"is_controllable"`
}

// DERProperties 分散式能源屬性
type DERProperties struct {
	Type         string  `json:"type"` // pv, battery, wind
	RatedPowerKW float64 `json:"rated_power_kw"`
	IsControllable bool  `json:"is_controllable"`
}

