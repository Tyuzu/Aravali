// File: internal/reports/appeals_tracking_test.go

package reports

import (
	"net/http"
	"testing"
)

func TestGetMyAppealsHandlerExists(t *testing.T) {
	var _ http.HandlerFunc = GetMyAppeals(nil)
}
