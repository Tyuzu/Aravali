package reports

import (
	"net/http"
	"testing"
)

func TestGetMyAppealsHandlerExists(t *testing.T) {
	var _ http.HandlerFunc = GetMyAppeals(nil)
}
