package notifications

import (
	"net/http"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/testutils"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/require"
)

func TestItOnlyRetriesNotificationsWhenTheyHaveFailed(t *testing.T) {
	result := results.Result{
		ID: uuid.MustNew(),
	}

	failing := newNotificationServer(t, map[string]int{
		"/": http.StatusServiceUnavailable,
	})
	passing := newNotificationServer(t, map[string]int{
		"/": http.StatusOK,
	})
	host := testutils.Host(t, testutils.WithActor(
		Actor, newNotifier(NotifierOpts{
			Notifications: config.Notifications{
				{
					Kind:   config.Webhook,
					URL:    passing.URL(),
					Phases: []config.NotificationPhase{config.PreRun},
				},
				{
					Kind:   config.Webhook,
					URL:    failing.URL(),
					Phases: []config.NotificationPhase{config.PreRun},
				},
			},
		}),
	))

	nErrs, err := host.Service().Invoke(
		t.Context(),
		Actor,
		result.ID.String(),
		Notify,
		NotifyInput{
			Result: result,
			Phase:  PreRun,
		},
	)
	require.Nil(t, err)
	require.NotNil(t, nErrs)
	var outErr *NotificationError
	require.Nil(t, nErrs.Decode(&outErr))
	require.Len(t, outErr.Errors, 1)

	// Run it again, we should only get additional requests
	// to the failing webhook server
	_, err = host.Service().Invoke(
		t.Context(),
		Actor,
		result.ID.String(),
		Notify,
		NotifyInput{
			Result: result,
			Phase:  PreRun,
		},
	)
	require.Nil(t, err)

	require.Len(t, passing.received, 1)
	require.Len(t, failing.received, 2)
}

func TestItDoesntCallANotificationWhenItIsNotSetupInThatPhase(t *testing.T) {
	result := results.Result{
		ID: uuid.MustNew(),
	}

	pre := newNotificationServer(t, map[string]int{
		"/": http.StatusOK,
	})
	post := newNotificationServer(t, map[string]int{
		"/": http.StatusOK,
	})
	host := testutils.Host(t, testutils.WithActor(
		Actor, newNotifier(NotifierOpts{
			Notifications: config.Notifications{
				{
					Kind:   config.Webhook,
					URL:    pre.URL(),
					Phases: []config.NotificationPhase{config.PreRun},
				},
				{
					Kind:   config.Webhook,
					URL:    post.URL(),
					Phases: []config.NotificationPhase{config.PostRun},
				},
			},
		}),
	))

	_, err := host.Service().Invoke(
		t.Context(),
		Actor,
		result.ID.String(),
		Notify,
		NotifyInput{
			Result: result,
			Phase:  PreRun,
		},
	)
	require.Nil(t, err)
	require.Len(t, pre.received, 1)
	require.Len(t, post.received, 0)
	_, err = host.Service().Invoke(
		t.Context(),
		Actor,
		result.ID.String(),
		Notify,
		NotifyInput{
			Result: result,
			Phase:  PostRun,
		},
	)
	require.Nil(t, err)
	require.Len(t, pre.received, 1)
	require.Len(t, post.received, 1)
}
