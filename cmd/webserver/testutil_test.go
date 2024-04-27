//go:build e2e
// +build e2e

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
	"github.com/JosephJoshua/remana-backend/internal/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
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

	pool, err := testutil.StartDockerPool()
	require.NoError(t, err, "error starting docker pool")

	postgresResource, db, err := testutil.StartPostgresContainer(pool)
	require.NoError(t, err, "error deploying postgres container")

	err = testutil.MigratePostgres(context.Background(), db)
	require.NoError(t, err, "error migrating database")

	resources := []*dockertest.Resource{postgresResource}

	return db, func() error {
		return testutil.PurgeDockerResources(pool, resources)
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

func getEndpointURL(baseAddr string, endpoint string) string {
	url := url.URL{
		Scheme: "https",
		Host:   baseAddr,
		Path:   endpoint,
	}

	return url.String()
}
