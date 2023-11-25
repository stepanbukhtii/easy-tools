package api

import (
	"errors"
)

var (
	PathNotFound          = errors.New("path not found")
	NoAuthorizationHeader = errors.New("no authorization header")
	InvalidTokenSignature = errors.New("invalid token signature")
	FailedToParseToken    = errors.New("failed to parse token")
	InvalidToken          = errors.New("invalid token")
)
