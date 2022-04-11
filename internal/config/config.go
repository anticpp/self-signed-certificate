package config

// Abstract interface of Config.
type Config interface {
	// Get value with key.
	// Key is a multiple segments splited by '.', for example: 'key.alg'.
	// The `bool` return false if key doesn't exist.
	Get(string) (Value, bool)

	// Handy for debug printing.
	String() string
}
