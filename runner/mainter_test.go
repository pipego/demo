package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainter(t *testing.T) {
	m := MainterNew(context.Background(), MainterDefaultConfig())
	assert.NotEqual(t, nil, m)
}
