package main

import (
	"context"
	"fmt"

	"github.com/pipego/cli/cmd"
)

func main() {
	if err := cmd.Run(context.Background()); err != nil {
		fmt.Println(err.Error())
	}
}
