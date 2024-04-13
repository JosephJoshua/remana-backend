package repository

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func pgtypeUUIDToGoogleUUID(id pgtype.UUID) (uuid.UUID, error) {
	return uuid.FromBytes(id.Bytes[:])
}
