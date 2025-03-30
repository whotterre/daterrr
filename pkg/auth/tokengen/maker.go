package tokengen

import "time"

// Holds the interface for Paseto

type Maker interface {
	CreateToken(username string, role string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
	
}