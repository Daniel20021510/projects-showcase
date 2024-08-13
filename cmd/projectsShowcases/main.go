package main

import (
	"fmt"
	"projectsShowcase/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
