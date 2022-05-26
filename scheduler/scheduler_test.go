package scheduler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	s := New(context.Background(), DefaultConfig())
	assert.NotEqual(t, nil, s)
}
