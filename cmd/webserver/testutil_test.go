package main_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	main "github.com/JosephJoshua/remana-backend/cmd/webserver"
	"github.com/JosephJoshua/remana-backend/internal/core"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SSL certs to test Secure cookies.
const (
	serverCertPEM = `-----BEGIN CERTIFICATE-----
MIICvDCCAkOgAwIBAgIUcz9tOXBUnUr3VFasqFH9scHRUFowCgYIKoZIzj0EAwIw
gaQxCzAJBgNVBAYTAlhYMQwwCgYDVQQIDANOL0ExDDAKBgNVBAcMA04vQTE2MDQG
A1UECgwtTmV2ZXIgdXNlIHRoaXMgY2VydGlmaWNhdGUgaW4gcHJvZHVjdGlvbiBJ
bmMuMUEwPwYDVQQDDDgxMjAuMC4wLjE6IE5ldmVyIHVzZSB0aGlzIGNlcnRpZmlj
YXRlIGluIHByb2R1Y3Rpb24gSW5jLjAgFw0yNDA0MTcwNTQ2NTRaGA8yMTI0MDMy
NDA1NDY1NFowgaQxCzAJBgNVBAYTAlhYMQwwCgYDVQQIDANOL0ExDDAKBgNVBAcM
A04vQTE2MDQGA1UECgwtTmV2ZXIgdXNlIHRoaXMgY2VydGlmaWNhdGUgaW4gcHJv
ZHVjdGlvbiBJbmMuMUEwPwYDVQQDDDgxMjAuMC4wLjE6IE5ldmVyIHVzZSB0aGlz
IGNlcnRpZmljYXRlIGluIHByb2R1Y3Rpb24gSW5jLjB2MBAGByqGSM49AgEGBSuB
BAAiA2IABEUtltwGMYZyNp9UtLiSV4T252TKUFCOp1YhyXMlFJj0myWcLqnNZOTo
TBDs+TLvLyEtUKA7W3QSUxncSoq8czKV1Mgw6/I1KIh1B32pllqcikqkjgva0i6V
21ZfUHjJOKMyMDAwDwYDVR0RBAgwBocEfwAAATAdBgNVHQ4EFgQUHHSplxoiiH3u
1b2WKLHEwADOYOwwCgYIKoZIzj0EAwIDZwAwZAIwMWwwqjePCwft6QtFlchjeZK1
ZMOmUAW9uanwotzcCgXf5PorAPTIBtuuu1vB4HY0AjAwQelTF1uejXNTojnxnROf
oW7LwS3MWiRy4sheLCoUz7hiVvpzXL83bAA85HtlhOw=
-----END CERTIFICATE-----`

	serverKeyPEM = `-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBBSo27sWxbMMca8qGy+Vd5cKHp+8nY1fWXBb7VlzCxb9EucNokL7oN
ShC2nnnrc2egBwYFK4EEACKhZANiAARFLZbcBjGGcjafVLS4kleE9udkylBQjqdW
IclzJRSY9JslnC6pzWTk6EwQ7Pky7y8hLVCgO1t0ElMZ3EqKvHMyldTIMOvyNSiI
dQd9qZZanIpKpI4L2tIuldtWX1B4yTg=
-----END EC PRIVATE KEY-----`
)

func setupTest(t testing.TB) (*pgxpool.Pool, func() error) {
	t.Helper()

	const (
		migrationMaxWait = 5 * time.Second
	)

	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "error creating docker pool")

	err = pool.Client.Ping()
	require.NoError(t, err, "error pinging docker server")

	postgresResource, db, err := deployPostgres(pool)
	require.NoError(t, err, "error deploying postgres container")

	migrationCtx, cancelMigration := context.WithTimeout(context.Background(), migrationMaxWait)
	defer cancelMigration()

	err = migrate(migrationCtx, db)
	require.NoError(t, err, "error migrating database")

	resources := []*dockertest.Resource{postgresResource}

	return db, func() error {
		for _, resource := range resources {
			if purgeErr := pool.Purge(resource); purgeErr != nil {
				return fmt.Errorf("error purging resource: %w", purgeErr)
			}
		}

		return nil
	}
}

func waitForReady(ctx context.Context, t testing.TB, client http.Client, baseAddr string, timeout time.Duration) {
	t.Helper()

	startTime := time.Now()

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, getEndpointURL(baseAddr, "/healthz"), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()

			if resp.StatusCode == http.StatusNoContent {
				return
			}
		}

		select {
		case <-ctx.Done():
			return

		default:
			if time.Since(startTime) > timeout {
				t.Fatal("timed out waiting for server to be ready")
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}
}

func createHTTPClient(t testing.TB) http.Client {
	t.Helper()

	roots := x509.NewCertPool()

	ok := roots.AppendCertsFromPEM([]byte(serverCertPEM))
	require.True(t, ok)

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	return http.Client{
		Timeout: 3 * time.Second,
		Jar:     jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: roots,
			},
		},
	}
}

func runServer(ctx context.Context, t testing.TB, db *pgxpool.Pool) string {
	t.Helper()

	addr, err := getFreeAddress()
	require.NoError(t, err)

	go func() {
		err = main.Run(ctx, db, addr, serverCertPEM, serverKeyPEM)
		assert.NoError(t, err)
	}()

	return addr
}

func getFreeAddress() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", fmt.Errorf("failed to resolve tcp addr: %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("failed to listen to tcp addr: %w", err)
	}

	defer l.Close()

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return "", fmt.Errorf("failed to cast to TCPAddr: %w", err)
	}

	return addr.String(), nil
}

func deployPostgres(pool *dockertest.Pool) (*dockertest.Resource, *pgxpool.Pool, error) {
	const (
		pgUsername = "username"
		pgPassword = "secretpassword"
		pgDBName   = "remana"

		pgContainerLifetimeSecs = 30
		pgContainerMaxWait      = 30 * time.Second
	)

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", pgUsername),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPassword),
			fmt.Sprintf("POSTGRES_DB=%s", pgDBName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	})

	if err != nil {
		return nil, nil, fmt.Errorf("error running postgres container: %w", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", pgUsername, pgPassword, hostAndPort, pgDBName)

	if err = resource.Expire(pgContainerLifetimeSecs); err != nil {
		return nil, nil, fmt.Errorf("error setting expiry date on postgres container: %w", err)
	}

	pool.MaxWait = pgContainerMaxWait

	var db *pgxpool.Pool

	if err = pool.Retry(func() error {
		db, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			return err
		}

		return db.Ping(context.Background())
	}); err != nil {
		return nil, nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return resource, db, nil
}

func migrate(ctx context.Context, db *pgxpool.Pool) error {
	rawDB := stdlib.OpenDBFromPool(db)
	return core.Migrate(ctx, rawDB, "postgres")
}

func getEndpointURL(baseAddr string, endpoint string) string {
	url := url.URL{
		Scheme: "https",
		Host:   baseAddr,
		Path:   endpoint,
	}

	return url.String()
}
