package session

import (
	"time"
)

type Session struct {
	User      string        `json:"user,omitempty"`
	ExpiredAt time.Duration `json:"expired_at,omitempty"`
}
