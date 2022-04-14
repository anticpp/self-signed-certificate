package config

// All config source should be of Config type.
type Config interface {
	// Get value with key. Key is multiple level splited by '.', for example: 'key.alg'.
	// Return nil if key does not exist.
	Get(string) *Value

	// For printable
	String() string
}
