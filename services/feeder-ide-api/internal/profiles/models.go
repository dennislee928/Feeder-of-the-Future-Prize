package profiles

// Profile 代表一個 feeder profile（Rural/Suburban/Urban）
type Profile struct {
	Type           string        `json:"type"`
	Name           string        `json:"name"`
	Characteristics Characteristics `json:"characteristics"`
}

// Characteristics 描述 profile 的特徵
type Characteristics struct {
	LoadComposition      LoadComposition      `json:"load_composition"`
	TypicalFeederLength  float64              `json:"typical_feeder_length_km"`
	TypicalNodeCount     int                  `json:"typical_node_count"`
	TargetSAIDI          float64              `json:"target_saidi_minutes_per_year"`
	TargetSAIFI          float64              `json:"target_saifi_interruptions_per_year"`
	DERPenetrationRange  DERPenetrationRange  `json:"der_penetration_range"`
	EVPenetrationRange   EVPenetrationRange   `json:"ev_penetration_range"`
}

// LoadComposition 負載組成（比例總和應為 1.0）
type LoadComposition struct {
	Residential float64 `json:"residential"` // 住宅
	Commercial  float64 `json:"commercial"`   // 商業
	Industrial  float64 `json:"industrial"`   // 工業
}

// DERPenetrationRange 分散式能源滲透率範圍
type DERPenetrationRange struct {
	Min float64 `json:"min"` // 最小比例（0-1）
	Max float64 `json:"max"` // 最大比例（0-1）
}

// EVPenetrationRange EV 充電樁滲透率範圍
type EVPenetrationRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

