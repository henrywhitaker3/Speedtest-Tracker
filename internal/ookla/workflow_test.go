package ookla

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/notifications"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/testutils"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
)

func TestSpeedtestWorkflow(t *testing.T) {
	tests := []struct {
		name                   string
		notifications          bool
		preNotificationErr     error
		postNotificationErr    error
		pickServerErr          error
		freshErr               error
		runnerErr              error
		updateErr              error
		complete               bool
		wantStatus             workflow.Status
		wantPickServerCalls    int
		wantFreshCalls         int
		wantRunnerCalls        int
		wantUpdateCalls        int
		wantNotificationPhases []notifications.Phase
		wantResultError        string
		wantAttempts           map[string]int
		wantStepStatus         map[string]workflow.StepStatus
	}{
		{
			name:                "stores a successful result and sends notifications",
			notifications:       true,
			wantStatus:          workflow.StatusCompleted,
			wantPickServerCalls: 1,
			wantFreshCalls:      2,
			wantRunnerCalls:     1,
			wantUpdateCalls:     1,
			wantNotificationPhases: []notifications.Phase{
				notifications.PreRun,
				notifications.PostRun,
			},
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  1,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 1,
			},
		},
		{
			name:                "stores runner failures and sends post-run notifications",
			notifications:       true,
			runnerErr:           errors.New("speedtest failed"),
			wantStatus:          workflow.StatusCompleted,
			wantPickServerCalls: 1,
			wantFreshCalls:      2,
			wantRunnerCalls:     1,
			wantUpdateCalls:     1,
			wantNotificationPhases: []notifications.Phase{
				notifications.PreRun,
				notifications.PostRun,
			},
			wantResultError: "speedtest failed",
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
				"pick server":            1,
				"run test":               0,
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
			wantStatus:          workflow.StatusFailed,
			wantPickServerCalls: 1,
			wantFreshCalls:      2,
			wantRunnerCalls:     1,
			wantUpdateCalls:     1,
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  0,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 0,
			},
		},
		{
			name:                "continues when pre-run notifications fail",
			notifications:       true,
			preNotificationErr:  errors.New("notification failed"),
			wantStatus:          workflow.StatusCompleted,
			wantPickServerCalls: 1,
			wantFreshCalls:      2,
			wantRunnerCalls:     1,
			wantUpdateCalls:     1,
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
			wantPickServerCalls: 1,
			wantFreshCalls:      2,
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
			name:                "skips notifications when none are configured",
			wantStatus:          workflow.StatusCompleted,
			wantPickServerCalls: 1,
			wantFreshCalls:      2,
			wantRunnerCalls:     1,
			wantUpdateCalls:     1,
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  0,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 0,
			},
		},
		{
			name:                "does not rerun an already complete result",
			notifications:       true,
			complete:            true,
			wantStatus:          workflow.StatusCompleted,
			wantPickServerCalls: 0,
			wantFreshCalls:      2,
			wantUpdateCalls:     1,
			wantNotificationPhases: []notifications.Phase{
				notifications.PreRun,
				notifications.PostRun,
			},
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  1,
				"run test":               1,
				"store result":           1,
				"post-run notifications": 1,
			},
		},
		{
			name: "fails when selecting a server fails",
			pickServerErr: errors.Join(
				actor.ErrJobPermanentFailure,
				errors.New("no nearby server available"),
			),
			wantStatus:          workflow.StatusFailed,
			wantPickServerCalls: 1,
			wantFreshCalls:      1,
			wantUpdateCalls:     1,
			wantResultError:     "job failed permanently\nno nearby server available",
			wantAttempts: map[string]int{
				"notify":                 1,
				"pre-run notifications":  0,
				"pick server":            1,
				"run test":               0,
				"store result":           0,
				"post-run notifications": 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &workflowRepoStub{freshErr: tt.freshErr, updateErr: tt.updateErr}
			runner := &workflowRunnerStub{
				pickErr:        tt.pickServerErr,
				runErr:         tt.runnerErr,
				pickedServerID: 42,
			}
			notifier := &workflowNotifierStub{
				shouldNotify: tt.notifications,
				errs: map[notifications.Phase]error{
					notifications.PreRun:  tt.preNotificationErr,
					notifications.PostRun: tt.postNotificationErr,
				},
			}

			wf, err := buildWorkflow(WorkflowOpts{
				Repo:     repo,
				Runner:   runner,
				Notifier: notifier,
			})
			require.NoError(t, err)
			host := startWorkflowHost(t, wf)
			service := wf.Service(host.Service())
			result := results.Result{
				Complete: tt.complete,
				ID:       uuid.MustOrdered(),
			}

			id, _, err := service.Start(t.Context(), RunWorkflow{Result: result})
			require.NoError(t, err)
			status := awaitWorkflowStatus(t, service, id)
			assert.Equalf(t, tt.wantStatus, status.Status, "workflow cause: %s", status.Cause)
			assert.Equal(t, tt.wantFreshCalls, repo.freshCallCount())
			assert.Equal(t, tt.wantRunnerCalls, runner.callCount())
			assert.Equal(t, tt.wantPickServerCalls, runner.pickServerCallCount())
			assert.Equal(t, tt.wantUpdateCalls, repo.updateCallCount())
			assert.Equal(t, tt.wantNotificationPhases, notifier.phases())
			for step, want := range tt.wantAttempts {
				assert.Equalf(t, want, workflowStepAttempts(t, status, step), "%s attempts", step)
			}
			for step, want := range tt.wantStepStatus {
				assert.Equalf(t, want, workflowStepStatus(t, status, step), "%s status", step)
			}

			if tt.wantRunnerCalls > 0 {
				assert.Equal(t, 42, runner.runServerID)
			}
			if repo.updateCallCount() > 0 {
				updated := repo.lastUpdatedResult()
				assert.Equal(t, result.ID, updated.ID)
				assert.True(t, updated.Complete)
				assert.Equal(t, tt.wantResultError, errorString(updated.Error))
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
	mu             sync.Mutex
	pickErr        error
	runErr         error
	pickCalls      int
	calls          int
	pickedServerID int
	runServerID    int
}

func (r *workflowRunnerStub) PickServer(context.Context) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.pickCalls++
	return r.pickedServerID, r.pickErr
}

func (r *workflowRunnerStub) Run(_ context.Context, serverID int, _ *results.Result) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls++
	r.runServerID = serverID
	return r.runErr
}

func (r *workflowRunnerStub) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls
}

func (r *workflowRunnerStub) pickServerCallCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.pickCalls
}

type workflowNotifierStub struct {
	mu           sync.Mutex
	shouldNotify bool
	errs         map[notifications.Phase]error
	calls        []notifications.Phase
}

func (n *workflowNotifierStub) Invoke(
	_ context.Context,
	_ string,
	_ string,
	method string,
	data any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	switch method {
	case notifications.ShouldNotify:
		return workflowNotifierEnvelope{shouldNotify: n.shouldNotify}, nil
	case notifications.Notify:
		input, ok := data.(notifications.NotifyInput)
		if !ok {
			return nil, errors.New("unexpected notification input")
		}
		phase := input.Phase
		n.mu.Lock()
		n.calls = append(n.calls, phase)
		err := n.errs[phase]
		n.mu.Unlock()
		if err == nil {
			return workflowNotifierEnvelope{}, nil
		}
		return workflowNotifierEnvelope{
			notificationError: &notifications.NotificationError{Errors: []error{err}},
		}, nil
	default:
		return nil, errors.New("unexpected notifier method")
	}
}

func (n *workflowNotifierStub) phases() []notifications.Phase {
	n.mu.Lock()
	defer n.mu.Unlock()
	return slices.Clone(n.calls)
}

type workflowNotifierEnvelope struct {
	shouldNotify      bool
	notificationError *notifications.NotificationError
}

func (e workflowNotifierEnvelope) Decode(into any) error {
	switch value := into.(type) {
	case *bool:
		*value = e.shouldNotify
	case **notifications.NotificationError:
		*value = e.notificationError
	default:
		return errors.New("unexpected notification response target")
	}
	return nil
}

func startWorkflowHost(t *testing.T, wf *workflow.Workflow) *local.Host {
	t.Helper()

	host := testutils.Host(t, testutils.WithWorkflow(wf))
	return host
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
	require.FailNowf(t, "assertion failed",
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
	require.FailNowf(t, "assertion failed", "workflow status has no %q step", name)
	return 0
}

func workflowStepStatus(
	t *testing.T,
	status workflow.InstanceStatus,
	name string,
) workflow.StepStatus {
	t.Helper()
	for _, step := range status.Steps {
		if step.Name == name {
			return step.Status
		}
	}
	require.FailNowf(t, "assertion failed", "workflow status has no %q step", name)
	return ""
}

func errorString(err *string) string {
	if err == nil {
		return ""
	}
	return *err
}
