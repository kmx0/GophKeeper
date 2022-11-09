package secret

import "time"

type Secret struct {
	Owner    string        `json:"owner,omitempty"`
	UpdateAt time.Duration `json:"update_at,omitempty"`
	// Expired  bool
	Value interface{}
}

func (s *Secret) Get(owner string) (secret *Secret, error){
	if s.Owner != owner: return nil, fmt.Errorf("permission denied for %", owner)

	return s.Value, nil
}
