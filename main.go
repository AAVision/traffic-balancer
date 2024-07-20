package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/AAVision/traffic-balancer/cmd"
	"gopkg.in/yaml.v2"
)

var (
	algorithm = ""
)

type InternalConnections struct {
	mu          sync.Mutex
	connections map[string]int
}

type Config struct {
	Algorithm string `yaml:"algorithm"`
	Port      int    `yaml:"port"`
	Servers   []struct {
		Host        string  `yaml:"host"`
		Weight      float64 `yaml:"weight"`
		Connections int     `yaml:"connections"`
	} `yaml:"servers"`
}

func (c *Config) getConf() *Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

var serverPool cmd.ServerPool
var internalConnections InternalConnections

func (i *InternalConnections) IncreaseConnection(key string) {
	i.mu.Lock()
	i.connections[key] += 1
	i.mu.Unlock()
}

func lb(w http.ResponseWriter, r *http.Request) {
	var peer *cmd.Backend
	attempts := cmd.GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	switch algorithm {
	case "least-time":
		peer = serverPool.GetLowestLatency()
	case "weighted-round-robin":
		peer = serverPool.GetHighestWeight()
	case "connection-per-second":
		peer = serverPool.GetNextPeer()
		internalConnections.IncreaseConnection(peer.URL.String())
	default:
		peer = serverPool.GetNextPeer()
	}

	if peer != nil {
		if internalConnections.connections[peer.URL.String()] > peer.Connections {
			peer.SetAlive(false)
			return
		}
		log.Printf("MESSAGE FORWARDED to server: %s\n", peer.URL)
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

func healthCheck() {
	t := time.NewTicker(time.Minute * 1)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

func resetConnections() {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			internalConnections.connections = serverPool.InitConnections()
		}
	}
}

func main() {

	var c Config
	c.getConf()
	algorithm = c.Algorithm

	if len(c.Servers) == 0 {
		log.Fatal("Please provide one or more backends to load balance inside the config.yml file!!!")
	}

	for _, server := range c.Servers {
		serverUrl, err := url.Parse(server.Host)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := cmd.GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), cmd.Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			serverPool.MarkBackendStatus(serverUrl, false)
			attempts := cmd.GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), cmd.Attempts, attempts+1)
			lb(writer, request.WithContext(ctx))
		}

		serverPool.AddBackend(
			&cmd.Backend{
				URL:          serverUrl,
				Alive:        true,
				ReverseProxy: proxy,
				Weight:       server.Weight,
				Connections:  server.Connections,
			})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	internalConnections.connections = serverPool.InitConnections()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", c.Port),
		Handler: http.HandlerFunc(lb),
	}

	go healthCheck()
	go resetConnections()

	log.Printf("Load Balancer started at :%d\n", c.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
