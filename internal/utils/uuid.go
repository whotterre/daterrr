package utils

import "github.com/google/uuid"

import "github.com/jackc/pgx/v5/pgtype"

func UUIDToPgType(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: [16]byte(uuid)}
}

