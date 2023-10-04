package cmd

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pipego/cli/runner"
)

func TestInitConfig(t *testing.T) {
	ctx := context.Background()

	_, err := initConfig(ctx, "invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig(ctx, "../test/config/invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)
}

func TestLoadFile(t *testing.T) {
	_, err := loadFile("invalid.json")
	assert.NotEqual(t, nil, err)

	buf, err := loadFile("../test/data/runner.json")
	assert.Equal(t, nil, err)

	var data runner.Proto
	err = json.Unmarshal(buf, &data)
	assert.Equal(t, nil, err)
}

func TestInitDag(t *testing.T) {
	ctx := context.Background()

	c, err := initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)

	_, err = initDag(ctx, c)
	assert.Equal(t, nil, err)
}

func TestInitRunner(t *testing.T) {
	ctx := context.Background()

	c, err := initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)

	d, err := initDag(ctx, c)
	assert.Equal(t, nil, err)

	_, _, err = initRunner(ctx, c, "invalid.json", d)
	assert.NotEqual(t, nil, err)

	_, _, err = initRunner(ctx, c, "../test/data/runner.json", d)
	assert.Equal(t, nil, err)
}

func TestInitScheduler(t *testing.T) {
	ctx := context.Background()

	c, err := initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)

	_, err = initScheduler(ctx, c, "invalid.json")
	assert.NotEqual(t, nil, err)

	_, err = initScheduler(ctx, c, "../test/data/scheduler1.json")
	assert.Equal(t, nil, err)
}

func TestInitPipeline(t *testing.T) {
	ctx := context.Background()

	c, err := initConfig(ctx, "../test/config/config.yml")
	assert.Equal(t, nil, err)

	d, err := initDag(ctx, c)
	assert.Equal(t, nil, err)

	_t, _, err := initRunner(ctx, c, "../test/data/runner.json", d)
	assert.Equal(t, nil, err)

	s, err := initScheduler(ctx, c, "../test/data/scheduler1.json")
	assert.Equal(t, nil, err)

	_, err = initPipeline(ctx, c, _t, s)
	assert.Equal(t, nil, err)
}
