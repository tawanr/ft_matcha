package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/tawanr/ft_matcha/internal/models"
	"github.com/tawanr/ft_matcha/internal/models/mocks"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: NewMockModels(),
	}
}

type testServer struct {
	*httptest.Server
}

func NewMockModels() *models.Models {
	return &models.Models{
		Users:    &mocks.UserModelMock{},
		Profiles: &mocks.ProfileModelMock{},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
