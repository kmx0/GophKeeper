package secret

import "time"

type CreateSecretDTO struct {
	Owner string `json:"owner"`
}

type UpdateSecretDTO struct {
	Owner    string        `json:"owner"`
	UpdateAt time.Duration `json:"update_at"`
	Value    interface{}   `json:"value"`
	// Expired  bool
}
