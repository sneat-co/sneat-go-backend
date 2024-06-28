package sneatgaeapp

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sneat-co/sneat-go-core/capturer"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/sneat-co/sneat-go-core/security"
	"github.com/strongo/log"
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
		log.Infof(r.Context(), "globalOptionsHandler(%s): bad request: %s\n", r.URL.String(), m)
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
		log.Warningf(r.Context(), "globalOptionsHandler(%s): %s\n", r.URL.String(), m)
		_, _ = fmt.Println(w, m)
		return origin, false
	}
	return origin, true
}

var ReportPanic = func(err any) {
}

type HandlerWrapper = func(handler http.Handler) http.Handler

var noWrapper = func(handler http.Handler) http.Handler {
	return handler
}

func wrapHTTPHandler(handler http.HandlerFunc, wrapHandler HandlerWrapper) http.HandlerFunc {
	var wrappedHandlerFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		if _, isAllowedOrigin := allowedOrigin(r, w); !isAllowedOrigin { // Check origin, is this  unnecessary?
			return
		}
		defer func(started time.Time) { // needs to be inside handler wrapped by wrapHandler to keep context with logger
			log.Infof(r.Context(), "%s %s  in %v", r.Method, r.URL.String(), time.Since(started))
		}(time.Now())
		handler.ServeHTTP(w, r)
	}
	wrappedHandler := wrapHandler(wrappedHandlerFunc)
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
		wrappedHandler.ServeHTTP(w, r)
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
