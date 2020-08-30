package checker

import (
	"time"
)

type Checker interface {
	SetTimeout(timeout time.Duration)
	CheckDomain(domain string) (CheckResoluter, error)
}

type CheckResoluter interface {
	Domain() string
	StatusCode() int
	Error() error
}
