package logger

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// TODO: sanitize passwords

type PgxLogTracer struct{}

func (t *PgxLogTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l := zerolog.Ctx(ctx)
	queryCorrelationID := uuid.New()

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("query_correlation_id", queryCorrelationID.String())
	})

	l.Debug().
		Str("sql", data.SQL).
		Interface("args", data.Args).
		Msg("pgx query start")

	return ctx
}

func (t *PgxLogTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	l := zerolog.Ctx(ctx)

	if data.Err != nil {
		l.Error().Err(data.Err).Msg("pgx query end")
	} else {
		l.Debug().Int64("rows_affected", data.CommandTag.RowsAffected()).Msg("pgx query end")
	}
}

func (t *PgxLogTracer) TraceBatchQuery(ctx context.Context, _ *pgx.Conn, data pgx.TraceBatchQueryData) {
	l := zerolog.Ctx(ctx).With().
		Str("sql", data.SQL).
		Interface("args", data.Args).
		Int("rows_affected", int(data.CommandTag.RowsAffected())).
		Logger()

	if data.Err != nil {
		l.Error().Err(data.Err).Msg("pgx batch query")
	} else {
		l.Debug().Msg("pgx batch query")
	}
}

func (t *PgxLogTracer) TraceCopyFromStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceCopyFromStartData,
) context.Context {
	l := zerolog.Ctx(ctx)
	queryCorrelationID := uuid.New()

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("query_correlation_id", queryCorrelationID.String())
	})

	l.Debug().
		Interface("column_names", data.ColumnNames).
		Str("table_name", data.TableName.Sanitize()).
		Msg("pgx copy from start")

	return ctx
}

func (t *PgxLogTracer) TraceCopyFromEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromEndData) {
	l := zerolog.Ctx(ctx)

	if data.Err != nil {
		l.Error().Err(data.Err).Msg("pgx copy from end")
	} else {
		l.Debug().Int64("rows_affected", data.CommandTag.RowsAffected()).Msg("pgx copy from end")
	}
}

func (t *PgxLogTracer) TracePrepareStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TracePrepareStartData,
) context.Context {
	l := zerolog.Ctx(ctx)
	queryCorrelationID := uuid.New()

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("query_correlation_id", queryCorrelationID.String())
	})

	l.Debug().
		Str("sql", data.SQL).
		Str("name", data.Name).
		Msg("pgx prepare start")

	return ctx
}

func (t *PgxLogTracer) TracePrepareEnd(ctx context.Context, _ *pgx.Conn, data pgx.TracePrepareEndData) {
	l := zerolog.Ctx(ctx)

	if data.Err != nil {
		l.Error().Err(data.Err).Msg("pgx prepare end")
	} else {
		l.Debug().Msg("pgx prepare end")
	}
}

func (t *PgxLogTracer) TraceConnectStart(ctx context.Context, _ pgx.TraceConnectStartData) context.Context {
	l := zerolog.Ctx(ctx)
	connectionCorrelationID := uuid.New()

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("connection_correlation_id", connectionCorrelationID.String())
	})

	l.Debug().Msg("pgx connection start")
	return ctx
}

func (t *PgxLogTracer) TraceConnectEnd(ctx context.Context, data pgx.TraceConnectEndData) {
	l := zerolog.Ctx(ctx)

	if data.Err != nil {
		l.Error().Err(data.Err).Msg("pgx connection end")
	} else {
		l.Debug().Msg("pgx connection end")
	}
}
