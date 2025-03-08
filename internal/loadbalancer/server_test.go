package loadbalancer_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
	"time"

	"github.com/onkarbanerjee/roundbalancer/internal/loadbalancer"
	"github.com/onkarbanerjee/roundbalancer/pkg/backends"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var servers []*http.Server

func TestStart(t *testing.T) {
	t.Run("server should be able to dispatch request in round robin manner to only health backends", func(t *testing.T) {
		setUpServers(t)

		backend1URL, err := url.Parse("http://localhost:8080/livez")
		assert.NoError(t, err)
		backend2URL, err := url.Parse("http://localhost:8081/livez")
		assert.NoError(t, err)
		backend3URL, err := url.Parse("http://localhost:8082/livez")
		assert.NoError(t, err)

		parse, err := url.Parse("http://localhost:8080")
		assert.NoError(t, err)
		backend1 := backends.NewBackend(
			"1",
			httputil.NewSingleHostReverseProxy(parse),
			backend1URL,
		)
		parse, err = url.Parse("http://localhost:8081")
		assert.NoError(t, err)
		backend2 := backends.NewBackend(
			"2",
			httputil.NewSingleHostReverseProxy(parse),
			backend2URL,
		)
		parse, err = url.Parse("http://localhost:8082")
		assert.NoError(t, err)
		backend3 := backends.NewBackend(
			"3",
			httputil.NewSingleHostReverseProxy(parse),
			backend3URL,
		)

		logger, err := zap.NewProduction()
		assert.NoError(t, err)

		go func() {
			assert.NoError(t, loadbalancer.Start(
				[]*backends.Backend{
					backend1,
					backend2,
					backend3},
				logger,
				time.Second,
				9090))
		}()

		time.Sleep(10 * time.Second)

		// all 3 servers in rotation
		post, err := http.Post("http://localhost:9090/echo", "application/json", nil)
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err := io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 1", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", nil)
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 2", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 3", string(all))

		servers[1].Shutdown(context.Background()) //nolint:errcheck

		time.Sleep(10 * time.Second)

		// server 2 out of rotation
		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 1", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 3", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 1", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 3", string(all))

		go func() {
			mux2 := http.NewServeMux()
			mux2.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 2")) }) //nolint:errcheck
			mux2.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
			server2 := &http.Server{
				Addr:    ":8081",
				Handler: mux2,
			}
			server2.ListenAndServe() //nolint:errcheck
		}()

		time.Sleep(10 * time.Second)

		// server 2 rejoins the rotation
		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 1", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 2", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 3", string(all))

		post, err = http.Post("http://localhost:9090/echo", "application/json", bytes.NewBuffer([]byte("hello world")))
		assert.NoError(t, err)
		assert.Equal(t, 200, post.StatusCode)
		all, err = io.ReadAll(post.Body)
		assert.NoError(t, err)
		assert.Equal(t, "I am from server 1", string(all))
	})

	t.Run("server should be able to dispatch request even when backends keep going down and up at the same time", func(t *testing.T) {
		setUpServers(t)

		backend1URL, err := url.Parse("http://localhost:8080/livez")
		assert.NoError(t, err)
		backend2URL, err := url.Parse("http://localhost:8081/livez")
		assert.NoError(t, err)
		backend3URL, err := url.Parse("http://localhost:8082/livez")
		assert.NoError(t, err)

		parse, err := url.Parse("http://localhost:8080")
		assert.NoError(t, err)
		backend1 := backends.NewBackend(
			"1",
			httputil.NewSingleHostReverseProxy(parse),
			backend1URL,
		)
		parse, err = url.Parse("http://localhost:8081")
		assert.NoError(t, err)
		backend2 := backends.NewBackend(
			"2",
			httputil.NewSingleHostReverseProxy(parse),
			backend2URL,
		)
		parse, err = url.Parse("http://localhost:8082")
		assert.NoError(t, err)
		backend3 := backends.NewBackend(
			"3",
			httputil.NewSingleHostReverseProxy(parse),
			backend3URL,
		)

		logger, err := zap.NewProduction()
		assert.NoError(t, err)

		go func() {
			assert.NoError(t, loadbalancer.Start(
				[]*backends.Backend{
					backend1,
					backend2,
					backend3},
				logger,
				time.Second,
				9090))
		}()

		time.Sleep(10 * time.Second)

		for range 100 {
			time.Sleep(time.Millisecond)
			go func() {
				post, err := http.Post("http://localhost:9090/echo", "application/json", nil)
				assert.NoError(t, err)
				assert.Contains(t, []int{http.StatusOK, http.StatusBadGateway}, post.StatusCode)
				all, err := io.ReadAll(post.Body)
				assert.NoError(t, err)
				if post.StatusCode == http.StatusOK {
					assert.Contains(t, []string{
						"I am from server 1",
						"I am from server 2",
						"I am from server 3",
					}, string(all))
				}

			}()
		}

		for _, server := range servers {
			server.Shutdown(context.Background()) //nolint:errcheck
		}
		setUpServers(t)

		<-time.After(10 * time.Second)
	})
}

func setUpServers(t *testing.T) {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 1")) }) //nolint:errcheck
	mux1.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
	server1 := &http.Server{
		Addr:    ":8080",
		Handler: mux1,
	}
	servers = append(servers, server1)
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 2")) }) //nolint:errcheck
	mux2.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
	server2 := &http.Server{
		Addr:    ":8081",
		Handler: mux2,
	}
	servers = append(servers, server2)
	mux3 := http.NewServeMux()
	mux3.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 3")) }) //nolint:errcheck
	mux3.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
	server3 := &http.Server{
		Addr:    ":8082",
		Handler: mux3,
	}
	servers = append(servers, server3)

	for _, server := range servers {
		go func() {
			err := server.ListenAndServe()
			if err != nil {
				t.Log(err.Error())
			}
		}()
	}
}
