package ookla

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/notifications"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/components/sqlite"
	"github.com/italypaleale/francis/host/local"
)

func TestSpeedtestWorkflow(t *testing.T) {
	tests := []struct {
		name                   string
		notifications          bool
		preNotificationErr     error
		postNotificationErr    error
		freshErr               error
		runnerErr              error
		updateErr              error
		complete               bool
		wantStatus             workflow.Status
		wantFreshCalls         int
		wantRunnerCalls        int
		wantUpdateCalls        int
		wantNotificationPhases []notifications.Phase
		wantResultError        string
		wantAttempts           map[string]int
		wantStepStatus         map[string]workflow.StepStatus
	}{
		{
			name:                   "stores a successful result and sends notifications",
			notifications:          true,
			wantStatus:             workflow.StatusCompleted,
			wantFreshCalls:         1,
			wantRunnerCalls:        1,
			wantUpdateCalls:        1,
			wantNotificationPhases: []notifications.Phase{notifications.PreRun, notifications.PostRun},
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  1,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 1,
			},
		},
		{
			name:                   "stores runner failures and sends post-run notifications",
			notifications:          true,
			runnerErr:              errors.New("speedtest failed"),
			wantStatus:             workflow.StatusCompleted,
			wantFreshCalls:         1,
			wantRunnerCalls:        1,
			wantUpdateCalls:        1,
			wantNotificationPhases: []notifications.Phase{notifications.PreRun, notifications.PostRun},
			wantResultError:        "speedtest failed",
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  1,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 1,
			},
		},
		{
			name: "fails when refreshing the result fails",
			freshErr: errors.Join(
				actor.ErrJobPermanentFailure,
				errors.New("database unavailable"),
			),
			wantStatus:     workflow.StatusFailed,
			wantFreshCalls: 1,
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  0,
				"run test":               1,
				"store result":           0,
				"post-run notifications": 0,
			},
		},
		{
			name: "fails when storing the result fails",
			updateErr: errors.Join(
				actor.ErrJobPermanentFailure,
				errors.New("database unavailable"),
			),
			wantStatus:      workflow.StatusFailed,
			wantFreshCalls:  1,
			wantRunnerCalls: 1,
			wantUpdateCalls: 1,
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  0,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 0,
			},
		},
		{
			name:               "continues when pre-run notifications fail",
			notifications:      true,
			preNotificationErr: errors.New("notification failed"),
			wantStatus:         workflow.StatusCompleted,
			wantFreshCalls:     1,
			wantRunnerCalls:    1,
			wantUpdateCalls:    1,
			wantNotificationPhases: []notifications.Phase{
				notifications.PreRun,
				notifications.PreRun,
				notifications.PreRun,
				notifications.PostRun,
			},
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  3,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 1,
			},
			wantStepStatus: map[string]workflow.StepStatus{
				"pre-run notifications": workflow.StepFailed,
			},
		},
		{
			name:                "continues when post-run notifications fail",
			notifications:       true,
			postNotificationErr: errors.New("notification failed"),
			wantStatus:          workflow.StatusCompleted,
			wantFreshCalls:      1,
			wantRunnerCalls:     1,
			wantUpdateCalls:     1,
			wantNotificationPhases: []notifications.Phase{
				notifications.PreRun,
				notifications.PostRun,
				notifications.PostRun,
				notifications.PostRun,
			},
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  1,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 3,
			},
			wantStepStatus: map[string]workflow.StepStatus{
				"post-run notifications": workflow.StepFailed,
			},
		},
		{
			name:            "skips notifications when none are configured",
			wantStatus:      workflow.StatusCompleted,
			wantFreshCalls:  1,
			wantRunnerCalls: 1,
			wantUpdateCalls: 1,
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  0,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 0,
			},
		},
		{
			name:                   "does not rerun an already complete result",
			notifications:          true,
			complete:               true,
			wantStatus:             workflow.StatusCompleted,
			wantFreshCalls:         1,
			wantUpdateCalls:        1,
			wantNotificationPhases: []notifications.Phase{notifications.PreRun, notifications.PostRun},
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  1,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &workflowRepoStub{freshErr: tt.freshErr, updateErr: tt.updateErr}
			runner := &workflowRunnerStub{err: tt.runnerErr}
			notifier := &workflowNotifierStub{
				shouldNotify: tt.notifications,
				errs: map[notifications.Phase]error{
					notifications.PreRun:  tt.preNotificationErr,
					notifications.PostRun: tt.postNotificationErr,
				},
			}

			wf, err := buildWorkflow(WorkflowOpts{
				Repo:     repo,
				Runner:   runner.run,
				Notifier: notifier,
			})
			if err != nil {
				t.Fatal(err)
			}
			host := startWorkflowHost(t, wf)
			service := wf.Service(host.Service())
			result := results.Result{
				Complete: tt.complete,
				ID:       uuid.MustOrdered(),
			}

			id, _, err := service.Start(t.Context(), RunWorkflow{ServerID: 42, Result: result})
			if err != nil {
				t.Fatal(err)
			}
			status := awaitWorkflowStatus(t, service, id)
			if status.Status != tt.wantStatus {
				t.Errorf(
					"workflow status = %q, want %q (cause: %s)",
					status.Status,
					tt.wantStatus,
					status.Cause,
				)
			}

			if got := repo.freshCallCount(); got != tt.wantFreshCalls {
				t.Errorf("Fresh calls = %d, want %d", got, tt.wantFreshCalls)
			}
			if got := runner.callCount(); got != tt.wantRunnerCalls {
				t.Errorf("runner calls = %d, want %d", got, tt.wantRunnerCalls)
			}
			if got := repo.updateCallCount(); got != tt.wantUpdateCalls {
				t.Errorf("Update calls = %d, want %d", got, tt.wantUpdateCalls)
			}
			if got := notifier.phases(); !slices.Equal(got, tt.wantNotificationPhases) {
				t.Errorf("notification phases = %v, want %v", got, tt.wantNotificationPhases)
			}
			for step, want := range tt.wantAttempts {
				if got := workflowStepAttempts(t, status, step); got != want {
					t.Errorf("%s attempts = %d, want %d", step, got, want)
				}
			}
			for step, want := range tt.wantStepStatus {
				if got := workflowStepStatus(t, status, step); got != want {
					t.Errorf("%s status = %q, want %q", step, got, want)
				}
			}

			if tt.wantRunnerCalls > 0 && runner.serverID != 42 {
				t.Errorf("runner server ID = %d, want 42", runner.serverID)
			}
			if tt.wantUpdateCalls > 0 {
				updated := repo.lastUpdatedResult()
				if updated.ID != result.ID || !updated.Complete {
					t.Errorf("stored result = %#v, want completed result %q", updated, result.ID)
				}
				if got := errorString(updated.Error); got != tt.wantResultError {
					t.Errorf("stored result error = %q, want %q", got, tt.wantResultError)
				}
			}
		})
	}
}

type workflowRepoStub struct {
	mu        sync.Mutex
	freshErr  error
	updateErr error
	freshes   int
	updates   []results.Result
}

func (r *workflowRepoStub) Fresh(_ context.Context, _ *results.Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.freshes++
	return r.freshErr
}

func (r *workflowRepoStub) Update(_ context.Context, result *results.Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.updates = append(r.updates, *result)
	return r.updateErr
}

func (r *workflowRepoStub) freshCallCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.freshes
}

func (r *workflowRepoStub) updateCallCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.updates)
}

func (r *workflowRepoStub) lastUpdatedResult() results.Result {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.updates[len(r.updates)-1]
}

type workflowRunnerStub struct {
	mu       sync.Mutex
	err      error
	calls    int
	serverID int
}

func (r *workflowRunnerStub) run(_ context.Context, serverID int, _ *results.Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	r.serverID = serverID
	return r.err
}

func (r *workflowRunnerStub) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
}

type workflowNotifierStub struct {
	mu           sync.Mutex
	shouldNotify bool
	errs         map[notifications.Phase]error
	calls        []notifications.Phase
}

func (n *workflowNotifierStub) ShouldNotify() bool {
	return n.shouldNotify
}

func (n *workflowNotifierStub) Notify(
	_ context.Context,
	phase notifications.Phase,
	_ results.Result,
) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.calls = append(n.calls, phase)
	return n.errs[phase]
}

func (n *workflowNotifierStub) phases() []notifications.Phase {
	n.mu.Lock()
	defer n.mu.Unlock()
	return slices.Clone(n.calls)
}

func startWorkflowHost(t *testing.T, wf *workflow.Workflow) *local.Host {
	t.Helper()

	address := freeUDPAddress(t)
	host, err := local.NewHost(
		local.WithAddress(address),
		local.WithSQLiteProvider(sqlite.SQLiteProviderOptions{
			ConnectionString: filepath.Join(t.TempDir(), "workflow.db"),
		}),
		local.WithRuntimePSKs([]byte("workflow-test-runtime-psk-0123456789")),
		local.WithAlarmsPollInterval(100*time.Millisecond),
		local.WithLogger(slog.New(slog.DiscardHandler)),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := host.RegisterBuiltInActor(wf); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- host.Run(ctx) }()
	select {
	case <-host.Ready():
	case <-time.After(10 * time.Second):
		cancel()
		t.Fatal("workflow host did not become ready")
	}
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-errCh:
			if err != nil && !errors.Is(err, context.Canceled) {
				t.Errorf("workflow host stopped with error: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Error("workflow host did not stop")
		}
	})
	return host
}

func freeUDPAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatal(err)
	}
	address := listener.LocalAddr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	return address
}

func awaitWorkflowStatus(
	t *testing.T,
	service *workflow.WorkflowService,
	id string,
) workflow.InstanceStatus {
	t.Helper()

	deadline := time.Now().Add(20 * time.Second)
	var last workflow.InstanceStatus
	for time.Now().Before(deadline) {
		status, err := service.GetStatus(t.Context(), id)
		if err == nil {
			last = status
			if status.Status == workflow.StatusCompleted || status.Status == workflow.StatusFailed {
				return status
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf(
		"workflow %q did not finish; last status %q on step %q (%s)",
		id,
		last.Status,
		last.CurrentStep,
		last.Cause,
	)
	return workflow.InstanceStatus{}
}

func workflowStepAttempts(t *testing.T, status workflow.InstanceStatus, name string) int {
	t.Helper()
	for _, step := range status.Steps {
		if step.Name == name {
			return step.Attempts
		}
	}
	t.Fatalf("workflow status has no %q step", name)
	return 0
}

func workflowStepStatus(t *testing.T, status workflow.InstanceStatus, name string) workflow.StepStatus {
	t.Helper()
	for _, step := range status.Steps {
		if step.Name == name {
			return step.Status
		}
	}
	t.Fatalf("workflow status has no %q step", name)
	return ""
}

func errorString(err *string) string {
	if err == nil {
		return ""
	}
	return *err
}
