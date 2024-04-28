package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/JosephJoshua/remana-backend/internal/gensql"
	"github.com/JosephJoshua/remana-backend/internal/typemapper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repairOrderSlugProvider struct {
	queries *gensql.Queries
}

func newRepairOrderSlugProvider(db *pgxpool.Pool) repairOrderSlugProvider {
	return repairOrderSlugProvider{
		queries: gensql.New(db),
	}
}

func (r repairOrderSlugProvider) Generate(ctx context.Context, storeID uuid.UUID) (string, error) {
	// R123-45678-9012

	const (
		firstPartDigits  = 3
		secondPartDigits = 5
		thirdPartDigits  = 4

		maxRetries = 20
	)

	generate := func() (string, error) {
		getMax := func(digits int) *big.Int {
			return big.NewInt(9*int64(math.Pow10(digits-1)) - 1)
		}

		getMin := func(digits int) *big.Int {
			return big.NewInt(int64(math.Pow10(digits - 1)))
		}

		firstPart, err := rand.Int(rand.Reader, getMax(firstPartDigits))
		if err != nil {
			return "", err
		}

		firstPart.Add(firstPart, getMin(firstPartDigits))

		secondPart, err := rand.Int(rand.Reader, getMax(secondPartDigits))
		if err != nil {
			return "", err
		}

		secondPart.Add(secondPart, getMin(secondPartDigits))

		thirdPart, err := rand.Int(rand.Reader, getMax(thirdPartDigits))
		if err != nil {
			return "", err
		}

		thirdPart.Add(thirdPart, getMin(thirdPartDigits))

		return "R" + firstPart.String() + "-" + secondPart.String() + "-" + thirdPart.String(), nil
	}

	for range maxRetries {
		slug, err := generate()
		if err != nil {
			return "", fmt.Errorf("failed to generate repair order slug: %w", err)
		}

		_, err = r.queries.IsRepairOrderSlugTaken(ctx, gensql.IsRepairOrderSlugTakenParams{
			StoreID: typemapper.UUIDToPgtypeUUID(storeID),
			Slug:    slug,
		})

		if errors.Is(err, pgx.ErrNoRows) {
			return slug, nil
		}

		if err != nil {
			return "", fmt.Errorf("failed to check if repair order slug is taken: %w", err)
		}
	}

	return "", errors.New("max retries exceeded when generating repair order slug")
}
