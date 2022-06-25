package dag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDag(t *testing.T) {
	p := New(context.Background(), DefaultConfig())
	assert.NotEqual(t, nil, p)
}
