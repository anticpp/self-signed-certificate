package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// CommandlineConfig is used to handle config from commandline arguments.
// All commandline arguments will be serialize to an m-way tree: map[any]any.
type CommandlineConfig struct {
	m      map[any]any
	prefix string
	args   []string
}

// Create YamlConfig with args.
// You can specify a `prefix` to filter commandline arguments with the common prefix, or leaving `prefix` be empty.
// For example, `./prog -someprefix.key.alg=rsa` will specify the key 'key.rsa'.
func NewCommandlineConfig(args []string, prefix string) *CommandlineConfig {
	return &CommandlineConfig{
		m:      make(map[any]any),
		prefix: prefix,
		args:   args,
	}
}

// Implementation of Config.String()
func (c *CommandlineConfig) String() string {
	return fmt.Sprintf("CommandlineConfig: {%v}", c.m)
}

// Start parse commandline arguments
func (c *CommandlineConfig) Parse() error {
	for c.hasMore() {
		kv, err := c.parseNext()
		if err != nil {
			return errors.Wrap(err, "parseNext error")
		}
		fmt.Printf("{key: %v, val: %v}\n", kv.key, kv.val)
	}
	return nil
}

func (c *CommandlineConfig) hasMore() bool {
	if len(c.args) == 0 {
		return false
	}
	return true
}

func (c *CommandlineConfig) popNextArg() (string, bool) {
	if len(c.args) == 0 {
		return "", false
	}
	arg := c.args[0]
	c.args = c.args[1:]
	return arg, true
}

// Parse next key/val from `args`.
// Return `*kv` if a new key/val is parsed success. `*kv` may return nil if error happens, or the next argument is filtered out by prefix.
// `error` indicates invalid arguments, the caller can stop or continue the process.
func (c *CommandlineConfig) parseNext() (*kv, error) {
	s0, ok := c.popNextArg()
	if !ok {
		return nil, nil
	}

	// Strip '-' or '--'
	pos := 1
	if len(s0) < 2 || s0[0] != '-' {
		return nil, errors.New(fmt.Sprintf("argument invalid: \"%v\"", s0))
	}
	if s0[1] == '-' {
		pos++
	}
	s0 = s0[pos:]

	// If '-key=val', get key/val from current argument split by '='.
	// If '-key val', get value from the next argument.
	var key, val string

	posEqual := strings.Index(s0, "=")
	if posEqual != -1 {
		key = s0[0:posEqual]
		val = s0[posEqual+1:]
	} else {
		key = s0
		val, ok = c.popNextArg()
		if !ok {
			return nil, errors.New(fmt.Sprintf("no value specified for key: \"%v\"", key))
		}
	}

	// Filter key prefix
	if len(c.prefix) > 0 {
		has := strings.HasPrefix(key, fmt.Sprintf("%v.", c.prefix))
		if !has {
			return nil, nil
		}

		key = key[len(c.prefix)+1:]
	}

	return &kv{key, val}, nil
}
