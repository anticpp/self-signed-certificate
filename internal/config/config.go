package config

// Abstract interface of Config.
type Config interface {
	// Get value with key.
	// Key is a multiple segments splited by '.', for example: 'key.alg'.
	Get(string) *Value

	// Handy for debug printing.
	String() string
}
