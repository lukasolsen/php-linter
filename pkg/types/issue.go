package types

import "github.com/codevault-llc/php-lint/internal/token"

type Issue struct {
	RuleName string
	Message  string
	Pos      token.Pos // Add position information
}