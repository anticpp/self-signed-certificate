package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// CommandlineConfig is used to handle config from commandline arguments.
// All commandline arguments will be serialized to tree: map[string]any.
//
// Type interpretation:
// The default value type is string as default, for example, `-animal=dog`.
// If a number is specified, for example, `-age=25`, the value is interpreted as int value.
// It's the same on float value, for example, `-angel=13.5`, the value is interpreted as float value.
// You can expcitly quote the value, for example, `-age="25"`, the value is always interpreted as string value.
//
// Value type interpretations examples:
//   -animal=dog          string: "dog"
//   -age=25              int: 25
//   -age="25"            string: "25"           -- Quoted string
//   -age='25'            string: "25"           -- Quoted string
//   -angel=13.5          float: 13.5
//   -angel="13.5"        string: "13.5"         -- Quoted string
//   -heigh=13m           string: "13m"
//
// About prefix:
// You can specify a `prefix` to filter commandline arguments with the common prefix, or leaving `prefix` be empty.
// For example, `./prog -someprefix.key.alg=rsa` will specify the key 'key.rsa' with prefix 'someprefix'.
//
// Priority:
// The parsing process is from left to righ arguments, the latter parsed key/value will replace the previous ones.
// Thus, the key/value priority is from right to left.
// For example, `./prog -key.alg=dsa .. -key.alg=rsa`, the value for `key.alg` should be `rsa`.
type CommandlineConfig struct {
	m      map[string]any
	prefix string
	args   []string
}

// Create CommandlineConfig
func NewCommandlineConfig(args []string, prefix string) *CommandlineConfig {
	return &CommandlineConfig{
		m:      make(map[string]any),
		prefix: prefix,
		args:   args,
	}
}

// Printable.
func (c *CommandlineConfig) String() string {
	return fmt.Sprintf("CommandlineConfig: {%v}", c.m)
}

// Parse all arguments.
func (c *CommandlineConfig) Parse() error {
	for c.hasMore() {
		kv, err := c.parseNext()
		if err != nil {
			return errors.Wrap(err, "parseNext error")
		}
		c.feed(kv)
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
// `*kv` return not nil: if a new key/val is parsed success.
// `*kv` return nil:     1. error happens, 2) the next argument is filtered out by prefix.
// `error` indicates invalid arguments, the caller can stop or continue the next parse..
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

// Feed key/val to the map
func (c *CommandlineConfig) feed(kv *kv) {
	// Commandline value must be string
	ss, ok := kv.val.(string)
	if !ok {
		return
	}

	// Parse scalar value
	val := parseScalar(ss)

	cur := c.m
	names := strings.Split(kv.key, ".")
	for i, name := range names {
		// Leaf node, update value
		if i == len(names)-1 {
			cur[name] = val
			continue
		}

		// Tree node does not exists, create new node.
		m, ok := cur[name]
		if !ok {
			cur[name] = make(map[string]any)
		}
		m = cur[name]

		// If it's not a tree node(but leaf node), replace it with a tree node.
		_, ok = m.(map[string]any)
		if !ok {
			cur[name] = make(map[string]any)
		}
		m = cur[name]

		// Move next
		cur = m.(map[string]any)
	}
}

// Get value with key, see Config.Get().
func (c *CommandlineConfig) Get(key string) *Value {
	var cur any = c.m
	names := strings.Split(key, ".")
	for _, name := range names {
		switch m := cur.(type) {
		case map[any]any:
			cur = m[name]
		case map[string]any:
			cur = m[name]
		default:
			cur = nil
		}
	}
	if cur == nil {
		return nil
	}
	return &Value{cur}
}
