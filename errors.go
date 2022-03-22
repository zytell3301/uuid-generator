package uuid_generator

import "errors"

var (
	CheckerAlreadyStoppedError = errors.New("reader checker already stopped")
	CheckerAlreadyStartedError = errors.New("reader checker already started")
)
