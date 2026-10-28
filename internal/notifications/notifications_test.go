package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/components/sqlite"
	"github.com/italypaleale/francis/host/local"
	"github.com/stretchr/testify/require"
)

func TestNotifierNotify(t *testing.T) {
	result := results.Result{ID: uuid.MustOrdered(), Complete: true}
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
			notificationErrors := &NotificationError{}
			for _, notification := range tt.notifications(server.URL()) {
				if !notifiesAtPhase(notification, tt.phase) {
					continue
				}
				if err := notify(
					t.Context(),
					"https://speedtest.example",
					nil,
					notification,
					tt.phase,
					tt.result,
				); err != nil {
					notificationErrors.push(notification, err)
				}
			}

			if tt.wantErr {
				require.False(t, notificationErrors.Nil())
			} else {
				require.True(t, notificationErrors.Nil())
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

func TestItDoesntSendToSuccessfulChannelOnRetry(t *testing.T) {
	server := newNotificationServer(t, map[string]int{"/start": http.StatusOK})
	service := startNotifierHost(t, NotifierOpts{
		Notifications: config.Notifications{
			{
				Kind:   config.Healthcheck,
				URL:    server.URL(),
				Phases: []config.NotificationPhase{config.PreRun},
			},
		},
	})
	res := results.Result{
		ID: uuid.MustNew(),
	}
	require.Nil(t, invokeNotifier(t, service, PreRun, res))
	require.Len(t, server.received, 1)
	require.Nil(t, invokeNotifier(t, service, PreRun, res))
	require.Len(t, server.received, 1)
}

func TestItDoesSendToErroredChannelOnRetry(t *testing.T) {
	server := newNotificationServer(t, map[string]int{"/start": http.StatusServiceUnavailable})
	service := startNotifierHost(t, NotifierOpts{
		Notifications: config.Notifications{
			{
				Kind:   config.Healthcheck,
				URL:    server.URL(),
				Phases: []config.NotificationPhase{config.PreRun},
			},
		},
	})
	res := results.Result{
		ID: uuid.MustNew(),
	}
	require.NotNil(t, invokeNotifier(t, service, PreRun, res))
	require.Len(t, server.received, 1)
	require.NotNil(t, invokeNotifier(t, service, PreRun, res))
	require.Len(t, server.received, 2)
}

func TestNotifierRetriesOnlyFailedChannels(t *testing.T) {
	server := newNotificationServer(t, map[string]int{
		"/successful/start": http.StatusOK,
		"/failed/start":     http.StatusServiceUnavailable,
	})
	service := startNotifierHost(t, NotifierOpts{
		Notifications: config.Notifications{
			{
				Kind:   config.Healthcheck,
				URL:    server.URL() + "/successful",
				Phases: []config.NotificationPhase{config.PreRun},
			},
			{
				Kind:   config.Healthcheck,
				URL:    server.URL() + "/failed",
				Phases: []config.NotificationPhase{config.PreRun},
			},
		},
	})
	result := results.Result{ID: uuid.MustNew()}

	require.NotNil(t, invokeNotifier(t, service, PreRun, result))
	require.NotNil(t, invokeNotifier(t, service, PreRun, result))

	require.Equal(
		t,
		[]string{"/successful/start", "/failed/start", "/failed/start"},
		requestPaths(server.requests()),
	)
}

func startNotifierHost(t *testing.T, opts NotifierOpts) *local.Host {
	t.Helper()

	address := freeNotifierUDPAddress(t)
	host, err := local.NewHost(
		local.WithAddress(address),
		local.WithSQLiteProvider(sqlite.SQLiteProviderOptions{
			ConnectionString: filepath.Join(t.TempDir(), "notifications.db"),
		}),
		local.WithRuntimePSKs([]byte("notifications-test-runtime-psk-0123456789")),
		local.WithLogger(slog.New(slog.DiscardHandler)),
	)
	require.NoError(t, err)
	opts.Host = host
	require.NoError(t, RegisterActors(opts))

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- host.Run(ctx) }()
	select {
	case <-host.Ready():
	case <-time.After(10 * time.Second):
		cancel()
		require.FailNow(t, "assertion failed", "notifications host did not become ready")
	}
	t.Cleanup(func() {
		cancel()
		select {
		case <-errCh:
		case <-time.After(10 * time.Second):
			t.Error("notifications host did not stop")
		}
	})

	return host
}

func invokeNotifier(
	t *testing.T,
	host *local.Host,
	phase Phase,
	result results.Result,
) *NotificationError {
	t.Helper()

	response, err := host.Service().Invoke(
		t.Context(),
		Actor,
		result.ID.String(),
		Notify,
		NotifyInput{Result: result, Phase: phase},
	)
	require.NoError(t, err)

	var notificationErrors *NotificationError
	require.NoError(t, response.Decode(&notificationErrors))
	return notificationErrors
}

func freeNotifierUDPAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	require.NoError(t, err)
	address := listener.LocalAddr().String()
	require.NoError(t, listener.Close())
	return address
}

func healthcheck(url string, phases ...config.NotificationPhase) config.Notifications {
	return config.Notifications{{
		Kind:   config.Healthcheck,
		Phases: phases,
		URL:    url,
	}}
}

func notificationPhase(phase Phase) config.NotificationPhase {
	if phase == PreRun {
		return config.PreRun
	}
	return config.PostRun
}

type telegramSenderStub struct {
	params []*bot.SendMessageParams
	err    error
}

func (s *telegramSenderStub) SendMessage(
	_ context.Context,
	params *bot.SendMessageParams,
) (*models.Message, error) {
	copy := *params
	s.params = append(s.params, &copy)
	return &models.Message{}, s.err
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
