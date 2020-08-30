package checker

import (
	"github.com/gumeniukcom/achecker/configs"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const defaultTimeout = 30 * time.Second

type Check struct {
	client    *fasthttp.Client
	timeout   time.Duration
	normalize bool
}

func New(cfg *configs.Config) *Check {
	return &Check{
		client:    &fasthttp.Client{},
		timeout:   time.Duration(cfg.Checker.TimeoutSecond) * time.Second,
		normalize: cfg.Checker.Normalize,
	}
}

func NewWithClient(cfg *configs.Config, client *fasthttp.Client) *Check {
	return &Check{
		client:    client,
		timeout:   time.Duration(cfg.Checker.TimeoutSecond) * time.Second,
		normalize: cfg.Checker.Normalize,
	}
}

// SetTimeout set new timeout
func (c *Check) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	log.Debug().
		Float64("timeout", c.timeout.Seconds()).
		Msg("set new timeout")
}

// CheckDomain check domain
func (c *Check) CheckDomain(domain string) (CheckResoluter, error) {
	log.Debug().
		Str("domain", domain).
		Msg("start check domain")

	if c.normalize {
		var err error
		domain, err = c.normalizeDomain(domain)
		if err != nil {
			log.Error().
				Err(err).
				Str("domain", domain).
				Msg("error on normalize domain")
			return nil, err
		}
	}

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

func (c *Check) normalizeDomain(domain string) (string, error) {
	info, err := url.Parse(domain)
	if err != nil {
		log.Debug().
			Err(err).
			Str("domain", domain).
			Msg("error on parse domain")
		return "", err
	}
	if info.Scheme != "http" && info.Scheme != "https" {
		info.Scheme = "http"
	}
	return info.String(), nil
}
