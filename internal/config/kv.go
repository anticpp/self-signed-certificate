package config

// A kv object wraps a key and a value.
// Being more concrete, the `key` type is always string.
type kv struct {
	key string
	val any
}
