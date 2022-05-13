package servo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HttpServer is a simple wrapper around the standard lib's HTTP server
// that adds the ability to start and stop the server asynchronously as
// well as collecting any exit errors synchronously.
//
// NOTE. The implementation is a bit too simplistic at the moment b/c it
// assumes
// (1) only one thread will ever call Start and Stop;
// (2) Start is called before Stop and they're called exactly once;
// (3) if called at all, ExitOutcome is called exactly once after Stop.
type HttpServer interface {
	// Add a route to the server to dispatch incoming requests for the
	// given path to the specified RouteHandler.
	Route(path string, handler RouteHandler)
	// Start the server in the calling thread if foreground is true or
	// asynchronously otherwise.
	Start(foreground bool)
	// Stop the server asynchronously.
	Stop()
	// Collect any server startup or shutdown errors, blocking the caller
	// until the server has exited. Only call this method after calling
	// Start followed by Stop.
	ExitOutcome() *ExitOutcome

	// TODO provide robust, thread-safe implementation so we can drop
	// assumptions (1), (2), and (3).
}

// RouteHandler serves a request for a given route.
type RouteHandler func(http.ResponseWriter, *http.Request)

// ExitOutcome holds any startup or shutdown errors collected after the
// HTTP server has exited.
type ExitOutcome struct {
	StartError error
	StopError  error
}

type hsrv struct {
	ctx                 context.Context
	startupOutcome      chan error
	signalStop          context.CancelFunc
	shutdownGracePeriod time.Duration
	shutdownOutcome     chan error
	svr                 *http.Server
}

// Create a new HttpServer to listen on the specified port and that will
// wait shutdownGracePeriod seconds for route handlers to complete on
// server shutdown.
func NewHttpServer(port uint16, shutdownGracePeriod uint8) HttpServer {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.NewServeMux(),
	}
	ctx, stop := context.WithCancel(context.Background())

	return &hsrv{
		ctx:                 ctx,
		startupOutcome:      make(chan error, 1),
		signalStop:          stop,
		shutdownGracePeriod: time.Duration(shutdownGracePeriod) * time.Second,
		shutdownOutcome:     make(chan error, 1),
		svr:                 server,
	}
}

func (s *hsrv) Route(path string, handler RouteHandler) {
	if mux, ok := s.svr.Handler.(*http.ServeMux); ok {
		mux.HandleFunc(path, handler)
	}
}

func (s *hsrv) serve() {
	s.startupOutcome <- s.svr.ListenAndServe()
	close(s.startupOutcome)
}

func (s *hsrv) shutdownHandler() {
	<-s.ctx.Done()

	ctx, signalShutdown := context.WithTimeout(s.ctx, s.shutdownGracePeriod)
	defer signalShutdown()

	s.shutdownOutcome <- s.svr.Shutdown(ctx)
	close(s.shutdownOutcome)
}

func (s *hsrv) Start(foreground bool) {
	go s.shutdownHandler()
	if foreground {
		s.serve()
	} else {
		go s.serve()
	}
}

func (s *hsrv) Stop() {
	s.signalStop()
}

func (s *hsrv) ExitOutcome() *ExitOutcome {
	outcome := &ExitOutcome{}
	outcome.StartError = <-s.startupOutcome
	outcome.StopError = <-s.shutdownOutcome

	return outcome
}
