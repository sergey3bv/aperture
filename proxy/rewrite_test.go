package proxy

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrepareServicesRewriteValidation(t *testing.T) {
	testCases := []struct {
		name    string
		rewrite map[string]string
		wantErr string
	}{
		{
			name: "valid prefix absolute path",
			rewrite: map[string]string{
				"prefix": "/v1/api",
			},
		},
		{
			name: "reject relative prefix",
			rewrite: map[string]string{
				"prefix": "v1/api",
			},
			wantErr: "invalid prefix format",
		},
		{
			name: "reject prefix with scheme and host",
			rewrite: map[string]string{
				"prefix": "https://example.com/v1/api",
			},
			wantErr: "invalid prefix format",
		},
		{
			name: "reject unknown rewrite key",
			rewrite: map[string]string{
				"suffix": "/ignored",
			},
			wantErr: "unknown rewrite key",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			services := []*Service{{
				Name:       "test",
				Address:    "backend:8080",
				Protocol:   "https",
				Auth:       "off",
				HostRegexp: "^example\\.com$",
				PathRegexp: "^/.*$",
				Rewrite:    tc.rewrite,
			}}

			err := prepareServices(services)
			if tc.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.Contains(t, err.Error(), tc.wantErr)
		})
	}
}

func TestDirectorRewritePrefix(t *testing.T) {
	testCases := []struct {
		name         string
		prefix       string
		requestPath  string
		expectedPath string
	}{
		{
			name:         "prefix is prepended",
			prefix:       "/api",
			requestPath:  "/users",
			expectedPath: "/api/users",
		},
		{
			name:         "joinpath normalizes trailing slashes",
			prefix:       "/api/",
			requestPath:  "/users/",
			expectedPath: "/api/users/",
		},
		{
			name:         "joinpath normalizes encoded slash segment",
			prefix:       "/api",
			requestPath:  "/accounts/%2Fspecial",
			expectedPath: "/api/accounts/special",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			services := []*Service{{
				Name:       "test",
				Address:    "backend:8080",
				Protocol:   "https",
				Auth:       "off",
				HostRegexp: "^example\\.com$",
				PathRegexp: "^/.*$",
				Rewrite: map[string]string{
					"prefix": tc.prefix,
				},
			}}

			err := prepareServices(services)
			require.NoError(t, err)

			p := &Proxy{
				services: services,
			}

			req := httptest.NewRequest(
				"GET", "http://example.com"+tc.requestPath, nil,
			)
			p.director(req)

			require.Equal(t, "backend:8080", req.Host)
			require.Equal(t, "backend:8080", req.URL.Host)
			require.Equal(t, "https", req.URL.Scheme)
			require.Equal(t, tc.expectedPath, req.URL.Path)
		})
	}
}
