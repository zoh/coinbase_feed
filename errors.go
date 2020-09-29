package feeder

import "errors"

var (
	EmptyRequest = errors.New("empty symbols request")

	UnexpectedResponseType = errors.New("unexpected response type")
)
