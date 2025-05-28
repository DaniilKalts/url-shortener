package main

import (
	"fmt"

	"github.com/DaniilKalts/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TO-DO: init logger (slog)

	// TO-DO: init storage (sqlite)

	// TO-DO: init router (chi)

	// TO-DO: run server (http)
}
