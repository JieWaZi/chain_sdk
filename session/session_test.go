package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSessionFromReader(t *testing.T) {
	_, err := FromReader(nil)
	assert.NotNil(t, err)
	t.Logf("%s", err)
}

func TestCreateSessionFromBytes(t *testing.T) {
	_, err := FromRaw(nil)
	assert.NotNil(t, err)
	t.Logf("%s", err)
}
