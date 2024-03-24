package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithFields(t *testing.T) {
	ctx := context.Background()
	buf := bytes.Buffer{}

	{
		log := New().(logger)
		log.logrusLogger.Out = &buf
		ctx = NewContext(ctx, log)
	}

	log := FromContext(ctx)
	log = log.WithFields(map[string]any{
		"mock_key": "mock_value",
	})
	log.Infof("msg")

	type entryT struct {
		Msg     string `json:"msg"`
		MockKey string `json:"mock_key"`
	}
	var entry entryT

	require.NoError(t, json.NewDecoder(&buf).Decode(&entry))
	require.Equal(t, entry.Msg, "msg")
	require.Equal(t, entry.MockKey, "mock_value")
}
