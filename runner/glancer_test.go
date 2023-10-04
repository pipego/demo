package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlancer(t *testing.T) {
	g := GlancerNew(context.Background(), GlancerDefaultConfig())
	assert.NotEqual(t, nil, g)
}
