package faqs

import "testing"

func TestIsValidEntityTypeAcceptsEvent(t *testing.T) {
	if !isValidEntityType("event") {
		t.Fatal("expected event to be treated as a valid entity type")
	}
}
