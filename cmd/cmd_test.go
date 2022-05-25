package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	var err error
	ctx := context.Background()

	_, err = initConfig(ctx, "invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig(ctx, "../test/config/invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)
}

func TestInitScheduler(t *testing.T) {
	ctx := context.Background()

	c, err := initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)

	_, err = initScheduler(ctx, c)
	assert.Equal(t, nil, err)
}

func TestInitRunner(t *testing.T) {
	ctx := context.Background()

	c, err := initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)

	_, err = initScheduler(ctx, c)
	assert.Equal(t, nil, err)
}
