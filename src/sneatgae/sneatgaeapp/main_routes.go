package sneatgaeapp

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-core/capturer"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/sneat-co/sneat-go-core/security"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func initHTTPRouter(globalOptions http.HandlerFunc) *httprouter.Router {
	router := httprouter.New()
	if router.HandleOPTIONS = globalOptions != nil; router.HandleOPTIONS {
		router.GlobalOPTIONS = globalOptions
	}
	return router
}

// globalOptionsHandler handles OPTIONS requests
func globalOptionsHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("globalOptionsHandler()", r.URL)
	accessControlRequestMethod := r.Header.Get("Access-Control-Request-Method")
	if accessControlRequestMethod == "" {
		w.WriteHeader(http.StatusBadRequest)
		const m = "Missing required request header: Access-Control-Request-Method"
		log.Printf("globalOptionsHandler(%s): bad request: %s\n", r.URL.String(), m)
		_, _ = fmt.Println(w)
		return
	}
	origin, isAllowedOrigin := allowedOrigin(r, w)
	if !isAllowedOrigin {
		return
	}
	// Set CORS headers BEFORE calling w.WriteHeader() or w.Write()
	responseHeader := w.Header()
	responseHeader.Set("Access-Control-Allow-Origin", origin)
	responseHeader.Set("Access-Control-Allow-Methods", accessControlRequestMethod)
	accessControlRequestHeaders := r.Header.Get("Access-Control-Request-Headers")
	if accessControlRequestHeaders != "" {
		responseHeader.Set("Access-Control-Allow-Headers", accessControlRequestHeaders)
	}
	//log.Println("globalOptionsHandler(): OK, response code = 204 - no content")
	w.WriteHeader(http.StatusOK) // Do not use http.StatusNoContent here, it will cause error in Chrome
}

func allowedOrigin(r *http.Request, w http.ResponseWriter) (string, bool) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = r.Header.Get("Referer")
	}
	if origin == "" {
		return "", true
	}
	if !security.IsSupportedOrigin(origin) {
		w.WriteHeader(http.StatusForbidden)
		m := "Unsupported origin: " + origin
		log.Printf("globalOptionsHandler(%s): %s\n", r.URL.String(), m)
		_, _ = fmt.Println(w, m)
		return origin, false
	}
	return origin, true
}

var ReportPanic = func(err any) {
}

type HandlerWrapper interface {
	Handle(handler http.Handler) http.Handler
}

type noOpHandlerWrapper struct{}

func (noOpHandlerWrapper) Handle(handler http.Handler) http.Handler {
	return handler
}

var errsReporter HandlerWrapper = noOpHandlerWrapper{}

func wrapHTTPHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	var handlerWrapper http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		if _, isAllowedOrigin := allowedOrigin(r, w); !isAllowedOrigin { // Check origin, is this  unnecessary?
			return
		}
		uri := r.URL.Path
		if r.URL.RawQuery != "" {
			uri += "?" + r.URL.RawQuery
		}
		//log.Println(r.Method, uri, "started")
		defer func(started time.Time) {
			log.Println(r.Method, uri, "completed in", time.Since(started))
		}(time.Now())
		handler(w, r)
	}
	errorsReporterHandler := errsReporter.Handle(handlerWrapper)
	panicHandler := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				stack := string(debug.Stack())
				handlePanic(w, r, err, stack)
				ReportPanic(err)
				fmt.Println("PANIC:", err, "\nSTACKTRACE from panic:\n"+stack)
				w.WriteHeader(http.StatusInternalServerError)
				httpserver.AccessControlAllowOrigin(w, r)
				_, _ = fmt.Fprint(w, "PANIC:", err, "\nSTACKTRACE from panic:\n"+stack)
			}
		}()
		errorsReporterHandler.ServeHTTP(w, r)
	}
	return panicHandler
}

func handlePanic(w http.ResponseWriter, r *http.Request, err interface{}, stack string) {
	w.WriteHeader(http.StatusInternalServerError)
	if !httpserver.AccessControlAllowOrigin(w, r) {
		return
	}
	if n, err := fmt.Fprintf(w, "panic: %v\n\n%v", err, stack); err != nil {
		_ = capturer.CaptureError(r.Context(),
			fmt.Errorf("failed to write panic error to response output after writing %v bytes: %w", n, err),
		)
	}
}
