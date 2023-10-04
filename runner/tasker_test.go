package runner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasker(t *testing.T) {
	_t := TaskerNew(context.Background(), TaskerDefaultConfig())
	assert.NotEqual(t, nil, _t)
}
