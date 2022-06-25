package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pipego/cli/cmd"
)

const (
	TIMEOUT = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	if err := cmd.Run(ctx); err != nil {
		fmt.Println(err.Error())
	}
}
