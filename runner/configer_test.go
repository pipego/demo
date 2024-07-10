package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiger(t *testing.T) {
	m := ConfigerNew(context.Background(), ConfigerDefaultConfig())
	assert.NotEqual(t, nil, m)
}
