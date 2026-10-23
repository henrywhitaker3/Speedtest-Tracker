package auth

import (
	"reflect"
	"testing"
)

func TestApplyIdentityClaimsReplacesRefreshedGroups(t *testing.T) {
	response := OIDCResponse{
		Email:  "old@example.com",
		Name:   "Old Name",
		Groups: []string{"speedtest-admins", "users"},
	}

	applyIdentityClaims(&response, oidcClaims{
		Email:  "new@example.com",
		Name:   "New Name",
		Groups: []string{"users"},
	})

	if response.Email != "new@example.com" || response.Name != "New Name" {
		t.Errorf("identity = email %q name %q, want refreshed claims", response.Email, response.Name)
	}
	if !reflect.DeepEqual(response.Groups, []string{"users"}) {
		t.Errorf("groups = %#v, want refreshed groups without admin access", response.Groups)
	}
}
