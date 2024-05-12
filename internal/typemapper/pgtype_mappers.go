package typemapper

import (
	"fmt"
	"time"

	"github.com/JosephJoshua/remana-backend/internal/optional"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func StringToPgtypeText(s string) pgtype.Text {
	return OptionalStringToPgtypeText(optional.Some(s))
}

func BoolToPgtypeBool(b bool) pgtype.Bool {
	return OptionalBoolToPgtypeBool(optional.Some(b))
}

func Int32ToPgtypeInt4(i int32) pgtype.Int4 {
	return OptionalInt32ToPgtypeInt4(optional.Some(i))
}

func UUIDToPgtypeUUID(id uuid.UUID) pgtype.UUID {
	return OptionalUUIDToPgtypeUUID(optional.Some(id))
}

func OptionalStringToPgtypeText(s optional.Optional[string]) pgtype.Text {
	return pgtype.Text{
		String: s.GetOrElse(""),
		Valid:  s.IsSet(),
	}
}

func OptionalBoolToPgtypeBool(b optional.Optional[bool]) pgtype.Bool {
	return pgtype.Bool{
		Bool:  b.GetOrElse(false),
		Valid: b.IsSet(),
	}
}

func OptionalInt32ToPgtypeInt4(i optional.Optional[int32]) pgtype.Int4 {
	return pgtype.Int4{
		Int32: i.GetOrElse(0),
		Valid: i.IsSet(),
	}
}

func OptionalUUIDToPgtypeUUID(id optional.Optional[uuid.UUID]) pgtype.UUID {
	var bytes [16]byte

	if id.IsSet() {
		bytes = id.MustGet()
	}

	return pgtype.UUID{
		Bytes: bytes,
		Valid: id.IsSet(),
	}
}

func TimeToPgtypeTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, InfinityModifier: pgtype.Finite, Valid: true}
}

func MustPgtypeUUIDToUUID(id pgtype.UUID) uuid.UUID {
	uuid, err := PgtypeUUIDToUUID(id)
	if err != nil {
		panic(fmt.Sprintf("failed to convert pgtype UUID to google UUID: %v", err))
	}

	return uuid
}

func PgtypeUUIDToUUID(id pgtype.UUID) (uuid.UUID, error) {
	return uuid.FromBytes(id.Bytes[:])
}

func MustPgtypeUUIDsToUUIDs(ids []pgtype.UUID) []uuid.UUID {
	uuids, err := PgtypeUUIDsToUUIDs(ids)
	if err != nil {
		panic(fmt.Sprintf("failed to convert pgtype UUIDs to google UUIDs: %v", err))
	}

	return uuids
}

func PgtypeUUIDsToUUIDs(ids []pgtype.UUID) ([]uuid.UUID, error) {
	googleUUIDs := make([]uuid.UUID, len(ids))

	for i, id := range ids {
		googleID, err := PgtypeUUIDToUUID(id)

		if err != nil {
			return nil, err
		}

		googleUUIDs[i] = googleID
	}

	return googleUUIDs, nil
}

func UUIDsToPgtypeUUIDs(ids []uuid.UUID) []pgtype.UUID {
	pgtypeUUIDs := make([]pgtype.UUID, len(ids))

	for i, id := range ids {
		pgtypeUUIDs[i] = UUIDToPgtypeUUID(id)
	}

	return pgtypeUUIDs
}
