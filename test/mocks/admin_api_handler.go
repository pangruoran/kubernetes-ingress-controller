package mocks

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/versions"
)

// AdminAPIHandler is a mock implementation of the Admin API. It only implements the endpoints that are
// required for the tests.
type AdminAPIHandler struct {
	mux *http.ServeMux
	t   *testing.T

	// version is the version string returned by mocked Kong instance, default is set to versions.KICv3VersionCutoff (3.4.1).
	version string

	// ready is a flag that indicates whether the server should return a 200 OK or a 503 Service Unavailable.
	// It's set to true by default.
	ready bool

	// workspaceExists makes `/workspace/workspaces/:id` return 200 when true, or 404 otherwise.
	workspaceExists bool

	// workspaceWasCreated is set to true when a workspace `POST /workspaces` was called.
	workspaceWasCreated atomic.Bool

	// configurationHash specifies the configuration hash of mocked Kong instance
	// return in /status response.
	configurationHash string

	// config holds the previously received config via `POST /config`.
	// It is returned when `GET /config` requests are received.
	config []byte

	// configPostErrorBody contains the error body which will be returned when
	// responding to a `POST /config` request.
	configPostErrorBody []byte
}

type AdminAPIHandlerOpt func(h *AdminAPIHandler)

func WithConfigurationHash(hash string) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.configurationHash = hash
	}
}

func WithWorkspaceExists(exists bool) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.workspaceExists = exists
	}
}

func WithReady(ready bool) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.ready = ready
	}
}

// WithVersion sets the version string returned by mocked Kong instance.
// If version is empty, the default version is used.
func WithVersion(version string) AdminAPIHandlerOpt {
	if version == "" {
		version = versions.KICv3VersionCutoff.String()
	}
	return func(h *AdminAPIHandler) {
		h.version = version
	}
}

func WithConfigPostError(errorbody []byte) AdminAPIHandlerOpt {
	return func(h *AdminAPIHandler) {
		h.configPostErrorBody = errorbody
	}
}

func NewAdminAPIHandler(t *testing.T, opts ...AdminAPIHandlerOpt) *AdminAPIHandler {
	h := &AdminAPIHandler{
		version: versions.KICv3VersionCutoff.String(),
		t:       t,
		ready:   true,
	}

	for _, opt := range opts {
		opt(h)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_, _ = w.Write(formatDefaultDBLessRootResponse(h.version))
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if !h.ready {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				if h.configurationHash != "" {
					_, _ = w.Write(formatDBLessStatusResponseWithConfigurationHash(h.configurationHash))
				} else {
					_, _ = w.Write([]byte(defaultDBLessStatusResponseWithoutConfigurationHash))
				}
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/workspace/workspaces/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if h.workspaceExists {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/workspaces", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if !h.workspaceExists {
				h.workspaceWasCreated.Store(true)
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id": "workspace"}`))
			} else {
				t.Errorf("unexpected workspace creation")
			}
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL)
	})
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if h.config != nil {
				_, _ = w.Write(h.config)
			} else {
				_, _ = w.Write([]byte(fmt.Sprintf(`{"version": "%s"}`, h.version)))
			}

		case http.MethodPost:
			if h.configPostErrorBody != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write(h.configPostErrorBody)
			} else {
				w.WriteHeader(http.StatusNoContent)
				b, _ := io.ReadAll(r.Body)
				h.t.Logf("got config: %v", string(b))
				h.config = b
			}
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL)
		}
	})
	h.mux = mux
	return h
}

func (m *AdminAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// this path gets spammed by the readiness checker and shouldn't be of interest for logging
	if r.URL.Path != "/status" {
		m.t.Logf("AdminAPIHandler received request: %s %s", r.Method, r.URL)
	}
	m.mux.ServeHTTP(w, r)
}

func (m *AdminAPIHandler) WasWorkspaceCreated() bool {
	return m.workspaceWasCreated.Load()
}

func formatDefaultDBLessRootResponse(version string) []byte {
	const defaultDBLessRootResponse = `{
		"version": "%s",
		"configuration": {
			"database": "off",
			"router_flavor": "traditional",
			"role": "traditional",
			"proxy_listeners": [
				{
					"ipv6only=on": false,
					"ipv6only=off": false,
					"ssl": false,
					"so_keepalive=off": false,
					"listener": "0.0.0.0:8000",
					"bind": false,
					"port": 8000,
					"deferred": false,
					"so_keepalive=on": false,
					"http2": false,
					"proxy_protocol": false,
					"ip": "0.0.0.0",
					"reuseport": false
				}
			]
		}
	}`
	return []byte(fmt.Sprintf(defaultDBLessRootResponse, version))
}

func formatDBLessStatusResponseWithConfigurationHash(hash string) []byte {
	const defaultDBLessStatusResponseWithConfigurationHash = `{
		"configuration_hash": "%s",
		"memory": {
		  "workers_lua_vms": [
			{
			  "http_allocated_gc": "43.99 MiB",
			  "pid": 1260
			},
			{
			  "http_allocated_gc": "43.98 MiB",
			  "pid": 1261
			}
		  ],
		  "lua_shared_dicts": {
			"kong_secrets": {
			  "allocated_slabs": "0.04 MiB",
			  "capacity": "5.00 MiB"
			},
			"prometheus_metrics": {
			  "allocated_slabs": "0.04 MiB",
			  "capacity": "5.00 MiB"
			},
			"kong": {
			  "allocated_slabs": "0.04 MiB",
			  "capacity": "5.00 MiB"
			},
			"kong_locks": {
			  "allocated_slabs": "0.06 MiB",
			  "capacity": "8.00 MiB"
			},
			"kong_healthchecks": {
			  "allocated_slabs": "0.04 MiB",
			  "capacity": "5.00 MiB"
			},
			"kong_cluster_events": {
			  "allocated_slabs": "0.04 MiB",
			  "capacity": "5.00 MiB"
			},
			"kong_rate_limiting_counters": {
			  "allocated_slabs": "0.08 MiB",
			  "capacity": "12.00 MiB"
			},
			"kong_core_db_cache": {
			  "allocated_slabs": "0.76 MiB",
			  "capacity": "128.00 MiB"
			},
			"kong_core_db_cache_miss": {
			  "allocated_slabs": "0.08 MiB",
			  "capacity": "12.00 MiB"
			},
			"kong_db_cache": {
			  "allocated_slabs": "0.76 MiB",
			  "capacity": "128.00 MiB"
			},
			"kong_db_cache_miss": {
			  "allocated_slabs": "0.08 MiB",
			  "capacity": "12.00 MiB"
			}
		  }
		},
		"server": {
		  "connections_reading": 0,
		  "total_requests": 615,
		  "connections_writing": 3,
		  "connections_handled": 615,
		  "connections_waiting": 0,
		  "connections_accepted": 615,
		  "connections_active": 3
		}
	  }`
	return []byte(fmt.Sprintf(defaultDBLessStatusResponseWithConfigurationHash, hash))
}

const defaultDBLessStatusResponseWithoutConfigurationHash = `{
	"memory": {
	  "workers_lua_vms": [
		{
		  "http_allocated_gc": "43.99 MiB",
		  "pid": 1260
		},
		{
		  "http_allocated_gc": "43.98 MiB",
		  "pid": 1261
		}
	  ],
	  "lua_shared_dicts": {
		"kong_secrets": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"prometheus_metrics": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_locks": {
		  "allocated_slabs": "0.06 MiB",
		  "capacity": "8.00 MiB"
		},
		"kong_healthchecks": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_cluster_events": {
		  "allocated_slabs": "0.04 MiB",
		  "capacity": "5.00 MiB"
		},
		"kong_rate_limiting_counters": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		},
		"kong_core_db_cache": {
		  "allocated_slabs": "0.76 MiB",
		  "capacity": "128.00 MiB"
		},
		"kong_core_db_cache_miss": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		},
		"kong_db_cache": {
		  "allocated_slabs": "0.76 MiB",
		  "capacity": "128.00 MiB"
		},
		"kong_db_cache_miss": {
		  "allocated_slabs": "0.08 MiB",
		  "capacity": "12.00 MiB"
		}
	  }
	},
	"server": {
	  "connections_reading": 0,
	  "total_requests": 615,
	  "connections_writing": 3,
	  "connections_handled": 615,
	  "connections_waiting": 0,
	  "connections_accepted": 615,
	  "connections_active": 3
	}
}`
