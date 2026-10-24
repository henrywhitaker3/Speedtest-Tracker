package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/require"
)

func TestNotifierNotify(t *testing.T) {
	result := results.Result{BaseModel: database.BaseModel{ID: uuid.MustOrdered()}, Complete: true}
	failedResult := result
	failedMessage := "speedtest failed"
	failedResult.Error = &failedMessage

	tests := []struct {
		name          string
		phase         Phase
		result        results.Result
		notifications func(string) config.Notifications
		statuses      map[string]int
		wantPaths     []string
		wantErr       bool
	}{
		{
			name:   "sends a start ping before a speedtest",
			phase:  PreRun,
			result: result,
			notifications: func(url string) config.Notifications {
				return healthcheck(url, config.PreRun)
			},
			wantPaths: []string{"/start"},
		},
		{
			name:   "sends a successful healthcheck ping after a speedtest",
			phase:  PostRun,
			result: result,
			notifications: func(url string) config.Notifications {
				return healthcheck(url, config.PostRun)
			},
			wantPaths: []string{"/"},
		},
		{
			name:   "sends a failed healthcheck ping after a failed speedtest",
			phase:  PostRun,
			result: failedResult,
			notifications: func(url string) config.Notifications {
				return healthcheck(url, config.PostRun)
			},
			wantPaths: []string{"/fail"},
		},
		{
			name:   "only sends notifications configured for the current phase",
			phase:  PreRun,
			result: result,
			notifications: func(url string) config.Notifications {
				return config.Notifications{
					{
						Kind:   config.Healthcheck,
						Phases: []config.NotificationPhase{config.PreRun},
						URL:    url + "/pre",
					},
					{
						Kind:   config.Healthcheck,
						Phases: []config.NotificationPhase{config.PostRun},
						URL:    url + "/post",
					},
				}
			},
			wantPaths: []string{"/pre/start"},
		},
		{
			name:   "sends webhook notifications as JSON",
			phase:  PreRun,
			result: result,
			notifications: func(url string) config.Notifications {
				return config.Notifications{{
					Kind:   config.Webhook,
					Phases: []config.NotificationPhase{config.PreRun},
					URL:    url + "/webhook",
				}}
			},
			wantPaths: []string{"/webhook"},
		},
		{
			name:   "returns an error when a notification endpoint fails",
			phase:  PreRun,
			result: result,
			notifications: func(url string) config.Notifications {
				return healthcheck(url, config.PreRun)
			},
			statuses:  map[string]int{"/start": http.StatusInternalServerError},
			wantPaths: []string{"/start"},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newNotificationServer(t, tt.statuses)
			err := New(tt.notifications(server.URL())).Notify(t.Context(), tt.phase, tt.result)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			requests := server.requests()
			require.Equal(t, tt.wantPaths, requestPaths(requests))
			for _, request := range requests {
				require.Equal(t, http.MethodPost, request.method)
				require.Equal(t, "application/json", request.contentType)
				require.Equal(t, tt.result.ID, request.result.ID)
			}
		})
	}
}

func TestNotificationErrorFormatsAndRedactsDestinationURLs(t *testing.T) {
	firstURL := "https://example.com/first-private-token"
	secondURL := "https://example.com/second-private-token"
	err := &NotificationError{}
	err.push(config.Notification{URL: firstURL}, fmt.Errorf("request to %s failed", firstURL))
	err.push(config.Notification{URL: secondURL}, fmt.Errorf("request to %s failed", secondURL))

	require.Equal(
		t,
		"request to [REDACTED] failed, request to [REDACTED] failed",
		err.Error(),
	)
}

func healthcheck(url string, phases ...config.NotificationPhase) config.Notifications {
	return config.Notifications{{
		Kind:   config.Healthcheck,
		Phases: phases,
		URL:    url,
	}}
}

type notificationRequest struct {
	method      string
	path        string
	contentType string
	result      results.Result
}

type notificationServer struct {
	mu       sync.Mutex
	server   *httptest.Server
	statuses map[string]int
	received []notificationRequest
}

func newNotificationServer(t *testing.T, statuses map[string]int) *notificationServer {
	t.Helper()

	s := &notificationServer{statuses: statuses}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result results.Result
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.mu.Lock()
		s.received = append(s.received, notificationRequest{
			method:      r.Method,
			path:        r.URL.Path,
			contentType: r.Header.Get("Content-Type"),
			result:      result,
		})
		status := s.statuses[r.URL.Path]
		s.mu.Unlock()

		if status == 0 {
			status = http.StatusOK
		}
		w.WriteHeader(status)
	}))
	t.Cleanup(s.server.Close)
	return s
}

func (s *notificationServer) URL() string {
	return s.server.URL
}

func (s *notificationServer) requests() []notificationRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.received)
}

func requestPaths(requests []notificationRequest) []string {
	paths := make([]string, len(requests))
	for i, request := range requests {
		paths[i] = request.path
	}
	return paths
}
