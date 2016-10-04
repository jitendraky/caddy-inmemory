package inmemory

import (
	"fmt"
	"testing"

	"github.com/mholt/caddy"
	"github.com/pkg/errors"
)

func TestParseConfiguration(t *testing.T) {
	conf := `inmemory {
  
    eviction_time 10
    max_cache_size 1024
    max_entry_size 1
    secret dsdasdas-dasdasd-asdas

    cache exact .*?
    cache match .*?

    ban match /
    ban exact /
    ban match index?
  }`

	c := caddy.NewTestController("http", conf)

	config, err := parseConfigurationBlock(c)

	if err != nil {
		t.Fatal(errors.Wrap(err, "Could not parse configuration"))
	}

	fmt.Println(config)
}
