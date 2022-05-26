package pipeline

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	p := New(context.Background(), DefaultConfig())
	assert.NotEqual(t, nil, p)
}
