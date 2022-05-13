package servo

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func buildHttpServer(port int) HttpServer {
	if port > 0 {
		return NewHttpServer(uint16(port), 1)
	}

	server := NewHttpServer(0, 1)
	serverImpl := server.(*hsrv)
	serverImpl.svr.Addr = "non:sense, will:bomb out!"
	return server
}

func isStartError(err error) bool {
	// NOTE. Race conditions.
	// Since we call Start and Stop one after another, these two methods
	// will race to do their job. So either outcome is possible:
	// 1. Stop shuts down the server before or while Start is still busy.
	//    In this case we expect Start to exit with an ErrServerClosed.
	// 2. Start completes before Stop. In this case we expect an address
	//    error since the URL we specified is non-sense.
	//
	//    error(*net.OpError) *{
	// 	    Op: "listen", Net: "tcp", Source: net.Addr nil, Addr: net.Addr nil,
	//      Err: error(*net.AddrError) *{
	//	      Err: "too many colons in address",
	//        Addr: "non:sense, will:bomb out!"
	//      }
	//    }

	if err == http.ErrServerClosed {
		return true
	}
	if opErr, ok := err.(*net.OpError); ok {
		if _, ok := (opErr.Err).(*net.AddrError); ok {
			return true
		}
	}
	return false
}

const howzit = "howzit!\n"

func sayHowzit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", howzit)
}

func getBody(url string) (string, error) {
	if resp, err := http.Get(url); err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
		if body, err := io.ReadAll(resp.Body); err != nil {
			return "", err
		} else {
			return string(body), nil
		}
	}
}

func runFailedHttpServerStart(t *testing.T, foreground bool) {
	target := buildHttpServer(0)
	target.Start(foreground)
	target.Stop()
	outcome := target.ExitOutcome()

	if !isStartError(outcome.StartError) {
		t.Errorf("want: start error; got: %v", outcome)
	}
	if outcome.StopError != nil {
		t.Errorf("want: nil shutdown as server didn't start; got: %v", outcome)
	}
}

func TestFailedBackgroundHttpServerStart(t *testing.T) {
	runFailedHttpServerStart(t, false)
}

func TestFailedForegroundHttpServerStart(t *testing.T) {
	runFailedHttpServerStart(t, true)
}

func TestRequestRouting(t *testing.T) {
	target := buildHttpServer(8282)
	target.Route("/", sayHowzit)

	target.Start(false)
	defer target.Stop()

	var got string
	var err error
	for k := 0; k < 10; k++ {
		time.Sleep(500 * time.Millisecond) // cater for server startup time
		if got, err = getBody("http://localhost:8282/"); err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("want: %s; got: %v", howzit, err)
	}
	if got != howzit {
		t.Errorf("want: %s; got: %v", howzit, got)
	}
}
