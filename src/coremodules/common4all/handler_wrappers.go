package common4all

import (
	"fmt"
	"github.com/strongo/logus"
	"net/http"
	"strings"

	"context"
)

func OptionsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		panic("Method != OPTIONS")
	}
	// Pre-flight request
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080":
	case "http://localhost:8100":
	case "https://local-app.sneat.ws":
	case "https://sneat.app":
	case "":
		BadRequestMessage(ctx, w, "Missing required request header: Origin")
		return
	default:
		if !(strings.HasPrefix(origin, "http://") && strings.HasSuffix(origin, ":8100")) {
			err := fmt.Errorf("unknown origin: %s", origin)
			logus.Debugf(ctx, err.Error())
			BadRequestError(ctx, w, err)
			return
		}
	}
	logus.Debugf(ctx, "Request 'Origin' header: %s", origin)
	responseHeader := w.Header()
	if accessControlRequestMethod := r.Header.Get("Access-Control-Request-Method"); !(accessControlRequestMethod == http.MethodGet || accessControlRequestMethod == http.MethodPost) {
		BadRequestMessage(ctx, w, "Requested method is unsupported: "+accessControlRequestMethod)
		return
	} else {
		responseHeader.Set("Access-Control-Allow-Methods", accessControlRequestMethod)
	}
	if accessControlRequestHeaders := r.Header.Get("Access-Control-Request-Headers"); accessControlRequestHeaders != "" {
		logus.Debugf(ctx, "Request Access-Control-Request-Headers: %v", accessControlRequestHeaders)
		responseHeader.Set("Access-Control-Allow-Headers", accessControlRequestHeaders)
	} else {
		logus.Debugf(ctx, "Request header 'Access-Control-Allow-Headers' is empty or missing")
		// TODO(security): Is it wrong to return 200 in this case?
	}
	responseHeader.Set("Access-Control-Allow-Origin", origin)
}

//func getOnly(handler dtdal.ContextHandler) func(w http.ResponseWriter, r *http.Request) {
//	return dtdal.HttpAppHost.HandleWithContext(OptionsHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//		if r.Method != http.MethodGet {
//			BadRequestMessage(c, w, "Expecting to get request method GET, got: "+r.Method)
//			return
//		}
//		hashedWriter := NewHashedResponseWriter(w)
//		handler(ctx, hashedWriter, r)
//		hashedWriter.setETagOrNotModifiedAndFlushBuffer(ctx, w, r)
//	}))
//}
//
//func postOnly(handler dtdal.ContextHandler) func(w http.ResponseWriter, r *http.Request) {
//	return dtdal.HttpAppHost.HandleWithContext(OptionsHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//		if r.Method != http.MethodPost {
//			BadRequestMessage(ctx, w, "Expecting to get request method POST, got: "+r.Method)
//			return
//		}
//		handler(ctx, w, r)
//	}))
//}

func BadRequestMessage(ctx context.Context, w http.ResponseWriter, m string) {
	logus.Infof(ctx, m)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(m))
}

func BadRequestError(ctx context.Context, w http.ResponseWriter, err error) {
	BadRequestMessage(ctx, w, err.Error())
}

func InternalError(ctx context.Context, w http.ResponseWriter, err error) { // TODO: deprecate!
	m := err.Error()
	logus.Errorf(ctx, m)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(m))
}

//type HashedResponseWriter struct {
//	status         int
//	hash           hash.Hash
//	buffer         *bytes.Buffer
//	responseWriter http.ResponseWriter
//}
//
//func NewHashedResponseWriter(w http.ResponseWriter) HashedResponseWriter {
//	return HashedResponseWriter{
//		hash:           md5.New(),
//		responseWriter: w,
//		buffer:         new(bytes.Buffer),
//	}
//}

//var _ http.ResponseWriter = (*HashedResponseWriter)(nil)
//
//func (hashedWriter *HashedResponseWriter) Header() http.Header {
//	return hashedWriter.responseWriter.Header()
//}
//
//func (hashedWriter *HashedResponseWriter) Write(b []byte) (int, error) {
//	if _, err := hashedWriter.hash.Write(b); err != nil {
//		return 0, fmt.Errorf("failed to write to hash: %w", err)
//	}
//	return hashedWriter.buffer.Write(b)
//}
//
//func (hashedWriter *HashedResponseWriter) WriteHeader(v int) {
//	hashedWriter.status = v
//	hashedWriter.responseWriter.WriteHeader(v)
//}
//
//func (hashedWriter *HashedResponseWriter) flush(w http.ResponseWriter) (int, error) {
//	i, err := w.Write(hashedWriter.buffer.Bytes())
//	if err != nil {
//		err = fmt.Errorf("failed to flush buffer to response writer: %w", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		i2, _ := w.Write([]byte(err.Error()))
//		i += i2
//	}
//	return i, err
//}
//
//func (hashedWriter *HashedResponseWriter) setETagOrNotModifiedAndFlushBuffer(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	if hashedWriter.status != 0 && hashedWriter.status != 200 {
//		if _, err := hashedWriter.flush(w); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		if hashedWriter.buffer.Len() > 0 {
//			logus.Debugf(c, "No ETag check/set as response.status=%d", hashedWriter.status)
//		}
//		return
//	}
//	eTag := fmt.Sprintf("%x", hashedWriter.hash)
//	if match := r.Header.Get("If-None-Match"); match == eTag {
//		logus.Debugf(c, "Setting response status to 304 - not modified")
//		w.WriteHeader(http.StatusNotModified)
//		logus.Debugf(c, "Response status set to 304 - not modified")
//	} else {
//		contentLen := hashedWriter.buffer.Len()
//		if (len(eTag) + 5) < contentLen {
//			w.Header().Set("ETag", eTag)
//			logus.Debugf(c, "ETag: "+eTag)
//		} else if contentLen > 0 {
//			logus.Debugf(c, "ETag is not set as contentLength:%d is smaller then ETag header", contentLen)
//		}
//		if _, err := hashedWriter.flush(w); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	}
//}
