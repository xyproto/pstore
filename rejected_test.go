package pstore

import (
	"net/http/httptest"
	"testing"
)

func TestRejectedIsCaseInsensitive(t *testing.T) {
	perm := NewPermissions(&UserState{})
	perm.AddUserPath("/Private")
	for _, path := range []string{"/admin", "/ADMIN/secret.lua", "/Admin/x", "/DATA/x", "/Repo", "/private/x", "/PRIVATE/x"} {
		if !perm.Rejected(httptest.NewRecorder(), httptest.NewRequest("GET", path, nil)) {
			t.Errorf("%s should be rejected", path)
		}
	}
	perm.SetPublicPath([]string{"/Pub"})
	for _, path := range []string{"/pub/x", "/PUB/x"} {
		if perm.Rejected(httptest.NewRecorder(), httptest.NewRequest("GET", path, nil)) {
			t.Errorf("%s should not be rejected", path)
		}
	}
}
