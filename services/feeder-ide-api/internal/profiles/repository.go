package profiles

import "sync"

// Repository 定義 profile 儲存介面
type Repository interface {
	GetByType(profileType string) (*Profile, error)
	List() ([]*Profile, error)
}

// InMemoryRepository 記憶體實作（預設包含三種 profile）
type InMemoryRepository struct {
	mu       sync.RWMutex
	profiles map[string]*Profile
}

// NewInMemoryRepository 建立新的記憶體 repository（包含預設 profiles）
func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{
		profiles: make(map[string]*Profile),
	}

	// 初始化預設 profiles
	repo.profiles["rural"] = &Profile{
		Type: "rural",
		Name: "Rural Feeder",
		Characteristics: Characteristics{
			LoadComposition: LoadComposition{
				Residential: 0.70,
				Commercial:  0.20,
				Industrial:  0.10,
			},
			TypicalFeederLength: 50.0, // km
			TypicalNodeCount:    50,
			TargetSAIDI:          480.0, // minutes/year
			TargetSAIFI:         2.5,  // interruptions/year
			DERPenetrationRange: DERPenetrationRange{
				Min: 0.0,
				Max: 0.15,
			},
			EVPenetrationRange: EVPenetrationRange{
				Min: 0.0,
				Max: 0.10,
			},
		},
	}

	repo.profiles["suburban"] = &Profile{
		Type: "suburban",
		Name: "Suburban Feeder",
		Characteristics: Characteristics{
			LoadComposition: LoadComposition{
				Residential: 0.60,
				Commercial:  0.35,
				Industrial:  0.05,
			},
			TypicalFeederLength: 30.0,
			TypicalNodeCount:    100,
			TargetSAIDI:          120.0,
			TargetSAIFI:         1.2,
			DERPenetrationRange: DERPenetrationRange{
				Min: 0.10,
				Max: 0.40,
			},
			EVPenetrationRange: EVPenetrationRange{
				Min: 0.05,
				Max: 0.30,
			},
		},
	}

	repo.profiles["urban"] = &Profile{
		Type: "urban",
		Name: "Urban Feeder",
		Characteristics: Characteristics{
			LoadComposition: LoadComposition{
				Residential: 0.40,
				Commercial:  0.50,
				Industrial:  0.10,
			},
			TypicalFeederLength: 15.0,
			TypicalNodeCount:    200,
			TargetSAIDI:          60.0,
			TargetSAIFI:         0.8,
			DERPenetrationRange: DERPenetrationRange{
				Min: 0.20,
				Max: 0.60,
			},
			EVPenetrationRange: EVPenetrationRange{
				Min: 0.10,
				Max: 0.50,
			},
		},
	}

	return repo
}

func (r *InMemoryRepository) GetByType(profileType string) (*Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profile, exists := r.profiles[profileType]
	if !exists {
		return nil, ErrProfileNotFound
	}
	return profile, nil
}

func (r *InMemoryRepository) List() ([]*Profile, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	profiles := make([]*Profile, 0, len(r.profiles))
	for _, profile := range r.profiles {
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

