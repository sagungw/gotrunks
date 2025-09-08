package log

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestHandlerHandle(t *testing.T) {
	writer := &bytes.Buffer{}
	InitLogger(WithLevel("info"), WithWriter(writer))

	slog.InfoContext(context.TODO(), "this is the message")

	if !strings.Contains(writer.String(), `"func":"github.com/sagungw/gotrunks/log.TestHandlerHandle"`) {
		t.FailNow()
	}

	if !strings.Contains(writer.String(), `"file":"/home/agung/Projects/gotrunks/log/log_test.go:15"`) {
		t.FailNow()
	}
}
