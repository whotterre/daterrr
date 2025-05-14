package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDToPgType converts a uuid.UUID to a PostgreSQL-compatible pgtype.UUID
func UUIDToPgType(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes:  [16]byte(uuid),
		Valid: true,
	}
}
