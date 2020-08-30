package checker

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseRequest(t *testing.T) {

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", "applicaition/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer serv.Close()

	checker := New()

	assert.NotNil(t, checker, "they should be Check instance")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	res, err := checker.CheckDomain(ctx, serv.URL)

	assert.Nil(t, err, "not error on default test")
	assert.NotNil(t, res, "res should be not nil")
	assert.Equal(t, res.Domain(), serv.URL, "domains not equal")
	assert.Equal(t, res.StatusCode(), http.StatusOK, "domains not equal")
	assert.Nil(t, res.Error(), "error should be nil")
}

func TestBaseRequestTimeoutExceeded(t *testing.T) {
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		time.Sleep(4 * time.Second)
		w.Header().Set("Content-Type", "applicaition/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer serv.Close()

	checker := New()
	checker.SetTimeout(2 * time.Second)

	assert.NotNil(t, checker, "they should be Checke instance")
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	res, err := checker.CheckDomain(ctx, serv.URL)

	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, res.Error(), "empty error, but need real error")
	assert.EqualError(t, res.Error(), ErrTimeout.Error(), "timeout error not equal")

}

func TestRequestOKStatusNot200(t *testing.T) {
	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", "applicaition/json")
		w.WriteHeader(http.StatusNotFound)
	}))
	defer serv.Close()

	checker := New()

	assert.NotNil(t, checker, "they should be Check instance")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	res, err := checker.CheckDomain(ctx, serv.URL)

	assert.Nil(t, err, "not error on default test")
	assert.NotNil(t, res, "res should be not nil")
	assert.Equal(t, res.Domain(), serv.URL, "domains not equal")
	assert.Equal(t, res.StatusCode(), http.StatusNotFound, "domains not equal")
	assert.Nil(t, res.Error(), "error should be nil")
}

func TestRequestOKStatusConnectionClosed(t *testing.T) {

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("server panic")
	}))
	defer serv.Close()

	checker := New()

	assert.NotNil(t, checker, "they should be Check instance")
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	res, err := checker.CheckDomain(ctx, serv.URL)
	t.Logf("%#v", res)
	assert.Nil(t, err, "not error on default test")
	assert.NotNil(t, res, "res should be not nil")
	assert.Equal(t, res.Domain(), serv.URL, "domains not equal")
	assert.Equal(t, res.StatusCode(), 0, "domains not equal")
	assert.EqualError(t, res.Error(), ErrConnectionClosed.Error(), "error should be nil")
}
