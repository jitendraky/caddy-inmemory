package inmemory

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/mholt/caddy"
	"github.com/pkg/errors"
)

type config struct {
	EvictionTime time.Duration
	MaxCacheSize int64
	MaxEntrySize int64
	PurgeSecret  string
	Cache        []configCache
	Ban          []configBan
}

type matchType int

const (
	EXACT  matchType = iota
	REGEXP matchType = iota
)

func (m matchType) String() string {
	switch m {
	case EXACT:
		return "exact"
	case REGEXP:
		return "match"
	default:
		return fmt.Sprintf("unknown(%d)", m)
	}
}

func resolveMatchType(val string) (matchType, error) {
	switch val {
	case "exact":
		return EXACT, nil
	case "match":
		return REGEXP, nil

	default:
		return 0, errors.Errorf("Could not recognize match type: %s", val)
	}
}

type configCache struct {
	Type   matchType
	Value  string
	Regexp *regexp.Regexp
}

type configBan struct {
	Type   matchType
	Value  string
	Regexp *regexp.Regexp
}

var requiredNumberOfArgs map[string]int = map[string]int{
	"eviction_time":  1,
	"max_cache_size": 1,
	"max_entry_size": 1,
	"secret":         1,
	"cache":          2,
	"ban":            2,
}

func parseConfigurationBlock(c *caddy.Controller) (config, error) {
	conf := config{
		Cache: make([]configCache, 0),
		Ban:   make([]configBan, 0),
	}

	for c.Next() {
		for c.NextBlock() {
			key := c.Val()
			args := c.RemainingArgs()

			if numberOfArgs, ok := requiredNumberOfArgs[key]; !ok {
				return conf, errors.Errorf("Not recognized setting %s", key)
			} else if numberOfArgs != len(args) {
				return conf, errors.Errorf("Setting %s has invalid number of arguments: %d", key, numberOfArgs)
			}

			switch key {
			case "eviction_time":
				value, err := parseNumber(args[0])

				if err != nil {
					return conf, err
				}

				conf.EvictionTime = time.Duration(value) * time.Second

			case "max_cache_size":
				value, err := parseNumber(args[0])

				if err != nil {
					return conf, err
				}

				conf.MaxCacheSize = value

			case "max_entry_size":
				value, err := parseNumber(args[0])

				if err != nil {
					return conf, err
				}

				conf.MaxEntrySize = value

			case "secret":
				conf.PurgeSecret = args[0]

			case "cache":
				match, err := resolveMatchType(args[0])

				if err != nil {
					return conf, err
				}

				switch match {
				case EXACT:
					conf.Cache = append(conf.Cache, configCache{Type: EXACT, Value: args[1]})

				case REGEXP:
					regexp, err := regexp.Compile(args[1])

					if err != nil {
						return conf, errors.Wrap(err, "Could not parse regexp expression")
					}

					conf.Cache = append(conf.Cache, configCache{Type: REGEXP, Regexp: regexp})
				}

			case "ban":
				match, err := resolveMatchType(args[0])

				if err != nil {
					return conf, err
				}

				switch match {
				case EXACT:
					conf.Ban = append(conf.Ban, configBan{Type: EXACT, Value: args[1]})

				case REGEXP:
					regexp, err := regexp.Compile(args[1])

					if err != nil {
						return conf, errors.Wrap(err, "Could not parse regexp expression")
					}

					conf.Ban = append(conf.Ban, configBan{Type: REGEXP, Regexp: regexp})
				}
			}

		}
	}

	return conf, nil
}

func requireNumberOfArgs(args []string, value int) (err error) {
	if len(args) != value {
		err = fmt.Errorf("Required number of arguments is %d", value)
	}

	return
}

func parseNumber(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}

	if parsed, err := strconv.Atoi(value); err != nil {
		return 0, errors.Wrap(err, "Could not parse number")
	} else {
		return int64(parsed), nil
	}
}
