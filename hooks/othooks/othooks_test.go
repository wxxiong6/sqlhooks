package othooks

import (
	"context"
	"database/sql"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/require"
	"github.com/wxxiong6/sqlhooks"
	"go.opentelemetry.io/otel"
)

var (
	tracer *mocktracer.MockTracer
)

func init() {
	tracer = otel.Tracer("Test")
	driver := sqlhooks.Wrap(&sqlite3.SQLiteDriver{}, New(tracer))
	sql.Register("ot", driver)
}

func TestSpansAreRecorded(t *testing.T) {
	db, err := sql.Open("ot", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	ctx, span := tracer.Start(context.Background(), "sql")

	{
		rows, err := db.QueryContext(ctx, "SELECT 1+?", "1")
		require.NoError(t, err)
		rows.Close()
	}

	{
		rows, err := db.QueryContext(ctx, "SELECT 1+?", "1")
		require.NoError(t, err)
		rows.Close()
	}

	span.End()

	require.Len(t, span, 3)
}

func TestNoSpansAreRecorded(t *testing.T) {
	db, err := sql.Open("ot", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), "SELECT 1")
	require.NoError(t, err)
	rows.Close()

}
