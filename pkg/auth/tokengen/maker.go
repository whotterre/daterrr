package tokengen

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Holds the interface for Paseto

type Maker interface {
	CreateToken(userID pgtype.UUID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
