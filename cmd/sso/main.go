package main

import (
	"fmt"

	"github.com/TemaStatham/sso/internal/config"
)

// go run ./cmd/sso/main.go --config="./config/config.yaml"

func main() {
	cfg := config.MustLoad()
	fmt.Print(cfg)
}
