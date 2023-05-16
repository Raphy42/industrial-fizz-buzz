package metrics

import (
	"context"
	"sync"

	"github.com/Raphy42/industrial-fizz-buzz/core/errors"
	"github.com/Raphy42/industrial-fizz-buzz/core/generics"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

type registry struct {
	lock           sync.RWMutex
	requestBuckets map[string]map[string]uint
	requestChans   map[string]<-chan []byte
}

var (
	once           sync.Once
	globalRegistry *registry
)

func init() {
	once.Do(func() {
		globalRegistry = newRegistry()
	})
}

func newRegistry() *registry {
	return &registry{
		requestBuckets: make(map[string]map[string]uint),
		requestChans:   make(map[string]<-chan []byte),
	}
}

func (r *registry) top(routes ...string) (map[string]TopRequest, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if len(routes) == 0 {
		routes = generics.MapKeys(r.requestChans)
	}

	result := make(map[string]TopRequest)
	for _, route := range routes {
		var topRequest TopRequest
		maxHit := uint(0)

		// absence of requests is normal if the application has juste started and has no traffic
		if !generics.MapHas(r.requestBuckets, route) {
			// if we don't have any associated channel for this route, it means it wasn't registered
			// dispatch a 404
			if !generics.MapHas(r.requestChans, route) {
				return nil, errors.NotFound()
			}
			// else we are still waiting for the first request
			result[route] = TopRequest{
				Route: route,
				Bytes: nil,
				Hits:  0,
			}
			continue
		}

		// find the most common request
		for request, hit := range r.requestBuckets[route] {
			if hit > maxHit {
				maxHit = hit
				topRequest = TopRequest{
					Route: route,
					Bytes: []byte(request),
					Hits:  hit,
				}
			}
		}
		result[route] = topRequest
	}
	return result, nil
}

func (r *registry) newRequestCounter(path string) chan<- []byte {
	r.lock.Lock()
	defer r.lock.Unlock()

	requestChan := make(chan []byte)
	r.requestChans[path] = requestChan
	return requestChan
}

func (r *registry) incr(route, payload string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if !generics.MapHas(r.requestBuckets, route) {
		r.requestBuckets[route] = make(map[string]uint)
	}
	if !generics.MapHas(r.requestBuckets[route], payload) {
		r.requestBuckets[route][payload] = 0
	}
	r.requestBuckets[route][payload] += 1
}

type innerEvent struct {
	route   string
	payload string
}

func (r *registry) start(ctx context.Context) {
	log := logger.FromContext(ctx)

	r.lock.RLock()
	defer r.lock.RUnlock()

	// routes * (n * 16) buffer
	muxedChan := make(chan innerEvent, len(r.requestChans)*16)

	// fan-in incoming request bodies from handlers
	for route, requestChan := range r.requestChans {
		// clone
		requestChan := requestChan
		route := route
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case request := <-requestChan:
					select {
					case muxedChan <- innerEvent{route, string(request)}:
						continue
					case <-ctx.Done():
						log.Warn("context was cancelled before request metrics could be dispatched")
						return
					}
				}
			}
		}()
	}

	// update the inner state
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-muxedChan:
				r.incr(event.route, event.payload)
			}
		}
	}()
}
