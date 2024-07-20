package cmd

import (
	"fmt"
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func (s *ServerPool) InitConnections() map[string]int {
	connections := make(map[string]int, len(s.backends))

	for _, server := range s.backends {
		server.Alive = true
		connections[server.URL.String()] = 0
	}

	return connections
}

func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, server := range s.backends {
		if server.URL.String() == backendUrl.String() {
			server.SetAlive(alive)
			break
		}
	}
}

func (s *ServerPool) HealthCheck() {
	for _, server := range s.backends {
		// status := "online"
		alive, latency := getBackendStatus(server.URL)
		server.SetAlive(alive)
		server.SetLatency(latency)
		// if !alive {
		// 	status = "offline"
		// }
		// log.Printf("%s is [%s]\n", server.URL, status)

	}
}

func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next

	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsActive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

func (s *ServerPool) GetLowestLatency() *Backend {
	smallestServer := s.backends[0]
	for _, server := range s.backends[1:] {
		if server.Latency < smallestServer.Latency {
			smallestServer = server
		}
	}

	return smallestServer
}

func (s *ServerPool) GetHighestWeight() *Backend {
	highestServerWeight := s.backends[0]

	for _, server := range s.backends[1:] {
		fmt.Println(float64(server.Weight))
		if server.Weight > highestServerWeight.Weight {
			highestServerWeight = server
		}
	}
	return highestServerWeight
}
