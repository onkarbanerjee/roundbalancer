package loadbalancer_test

import (
	"bytes"
	"context"
	"fmt"
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
	t.Run("run", func(t *testing.T) {

		setUpServers()

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

		//client := &http.Client{
		//	Transport: &http.Transport{},
		//}
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

		servers[1].Shutdown(context.Background())

		time.Sleep(10 * time.Second)

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
			mux2.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 2")) })
			mux2.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
			server2 := &http.Server{
				Addr:    ":8081",
				Handler: mux2,
			}
			server2.ListenAndServe()
		}()

		time.Sleep(10 * time.Second)

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
}

func setUpServers() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 1")) })
	mux1.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
	server1 := &http.Server{
		Addr:    ":8080",
		Handler: mux1,
	}
	servers = append(servers, server1)
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 2")) })
	mux2.HandleFunc("/livez", func(w http.ResponseWriter, r *http.Request) {})
	server2 := &http.Server{
		Addr:    ":8081",
		Handler: mux2,
	}
	servers = append(servers, server2)
	mux3 := http.NewServeMux()
	mux3.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("I am from server 3")) })
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
				fmt.Println(err.Error())
			}
		}()
	}
}
