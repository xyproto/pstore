package pstore

import (
	"testing"

	"github.com/xyproto/pinterface"
)

func TestInterface(t *testing.T) {
	perm, err := NewWithConf(defaultConnectionString)
	if err != nil {
		t.Error(err)
	}
	// Check that the value qualifies for the interface
	var _ pinterface.IPermissions = perm
}
