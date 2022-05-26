package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	r := New(context.Background(), DefaultConfig())
	assert.NotEqual(t, nil, r)
}
