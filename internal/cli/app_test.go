package cli

import (
	"bytes"
	"testing"
)

func TestAppRunRejectsUnknownCommand(t *testing.T) {
	App := NewApp(bytes.NewBuffer(nil), &bytes.Buffer{}, "test.db")

	if err := App.Run([]string {"wat"}); err == nil {
		t.Fatalf("Command does not exist.")
	}
}
