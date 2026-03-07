package handler

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
)

type IHandler interface {
	Start() error
	Stop()
}

type Registry struct {
	handlers  map[string]IHandler
	wg        *sync.WaitGroup
	errorChan chan handlerError
}

type handlerError struct {
	err error
	id  string
}

func NewRegistry() *Registry {
	return &Registry{
		handlers:  make(map[string]IHandler),
		wg:        new(sync.WaitGroup),
		errorChan: make(chan handlerError),
	}
}

func (r *Registry) StartAll() {
	for k, s := range r.handlers {
		r.wg.Add(1)
		r.run(k, s)
	}
	r.wait()
}

func (r *Registry) Register(id string, s IHandler) {
	r.handlers[id] = s
}

func (r *Registry) StopAll() {
	for k, s := range r.handlers {
		slog.Info("Stopping handler", "id", k)
		s.Stop()
		slog.Info("Handler stopped", "id", k)
	}
	r.wg.Wait()
	slog.Info("All handler stopped")
}

func (r *Registry) wait() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	slog.Info("waiting for signal")

	select {
	case <-signalCh:
		slog.Info("interrupted")
	case err := <-r.errorChan:
		slog.Error("fatal error for service:", "id", err.id)
		slog.Error(err.err.Error())
	}
}

func (r *Registry) run(k string, s IHandler) {
	go func() {
		defer r.wg.Done()
		slog.Info("Starting handler", "id", k)
		err := s.Start()
		if err != nil {
			r.errorChan <- handlerError{
				id:  k,
				err: err,
			}
		}
		slog.Info("Handler started", "id", k)
	}()
}
