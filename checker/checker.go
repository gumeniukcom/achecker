package checker

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const defaultTimeout = 30 * time.Second

type Check struct {
	client  *fasthttp.Client
	timeout time.Duration
}

func New() *Check {
	return &Check{
		client:  &fasthttp.Client{},
		timeout: defaultTimeout,
	}
}

func NewWithClient(client *fasthttp.Client) *Check {
	return &Check{
		client:  client,
		timeout: defaultTimeout,
	}
}

func (c *Check) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	log.Info().
		Float64("timeout", c.timeout.Seconds()).
		Msg("set new timeout")
}

func (c *Check) CheckDomain(ctx context.Context, domain string) (*CheckResult, error) {

	log.Debug().
		Str("domain", domain).
		Msg("start check domain")

	var req fasthttp.Request
	var resp fasthttp.Response
	req.SetRequestURI(domain)

	result := &CheckResult{
		domain: domain,
	}

	err := c.client.DoTimeout(&req, &resp, c.timeout)

	switch err {
	case fasthttp.ErrDialTimeout, fasthttp.ErrTimeout:
		log.Error().
			Err(err).
			Str("domain", domain).
			Msg("timeout error")
		result.error = ErrTimeout
	case fasthttp.ErrConnectionClosed:
		log.Error().
			Err(err).
			Str("domain", domain).
			Msg("the server closed connection before returning the first response byte")
		result.error = ErrConnectionClosed
	case fasthttp.ErrNoFreeConns:
		log.Error().
			Err(err).
			Str("domain", domain).
			Msg("internal error on request")
		return nil, err
	// TODO: research it
	/**
	case syscall.ECONNRESET:
		log.Error().
			Err(err).
			Str("domain", domain).
			Msg("connection reset by peer")
		result.error = ErrConnectionResetByPeer
	*/
	case nil:
		result.statusCode = resp.StatusCode()
	default:
		result.error = err
		if strings.Contains(err.Error(), "no such host") { // net.DNSError
			log.Error().
				Err(err).
				Str("domain", domain).
				Msg("some dns error")
			result.error = ErrHostNotFound
		} else {
			log.Error().
				Err(err).
				Str("domain", domain).
				Msg("unknown error on request")
			return nil, err
		}
	}

	log.Debug().
		Str("domain", result.domain).
		Int("status code", result.statusCode).
		Err(result.error).
		Msg("domain checked")

	return result, nil
}
