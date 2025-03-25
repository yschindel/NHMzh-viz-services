package server

import "errors"

// Common errors for message processing
var (
	ErrMissingProject   = errors.New("missing project field")
	ErrMissingFilename  = errors.New("missing filename field")
	ErrMissingTimestamp = errors.New("missing timestamp field")
	ErrEmptyData        = errors.New("empty data field")
	ErrProcessingFailed = errors.New("failed to process message")
	ErrWritingFailed    = errors.New("failed to write message to database")
)
