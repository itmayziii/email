package send_test

import (
	"github.com/itmayziii/email/send"
	"testing"
)

func TestNewApp_IsNotNil(t *testing.T) {
	t.Parallel()
	app := send.NewApp()
	if app == nil {
		t.Errorf("app should not be nil")
	}
}
