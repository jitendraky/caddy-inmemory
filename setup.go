package inmemory

import (
	"fmt"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("inmemory", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {

	conf, err := parseConfigurationBlock(c)

	if err != nil {
		return err
	}

	fmt.Printf("Resolved config is: %+v\n", conf)

	cfg := httpserver.GetConfig(c)

	middleware := func(next httpserver.Handler) httpserver.Handler {
		return cacheHandler{Next: next, Config: conf}
	}

	cfg.AddMiddleware(middleware)

	return nil
}
