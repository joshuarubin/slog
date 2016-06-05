package slog

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	{
		level := ParseLevel("info", WarnLevel)
		assert.Equal(t, InfoLevel, level)
	}

	{
		level := ParseLevel("warn", WarnLevel)
		assert.Equal(t, WarnLevel, level)
	}

	{
		level := ParseLevel("whatever", WarnLevel)
		assert.Equal(t, WarnLevel, level)
	}
}

func TestLevel_MarshalJSON(t *testing.T) {
	e := Entry{
		Level:   InfoLevel,
		Message: "hello",
		Fields:  Fields{},
	}

	expect := `{"fields":{},"level":"info","timestamp":"0001-01-01T00:00:00Z","message":"hello"}`

	b, err := json.Marshal(e)
	assert.NoError(t, err)
	assert.Equal(t, expect, string(b))
}
