package api

import (
	"fmt"
	"net/http"
	"strings"

	"context"
	"github.com/strongo/log"
)

func optionsHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
	if r.Method != "OPTIONS" {
		panic("Method != OPTIONS")
	}
	// Pre-flight request
	origin := r.Header.Get("Origin")
	switch origin {
	case "http://localhost:8080":
	case "http://localhost:8100":
	case "https://debtstracker.local":
	case "https://debtstracker.io":
	case "":
		BadRequestMessage(c, w, "Missing required request header: Origin")
		return
	default:
		if !(strings.HasPrefix(origin, "http://") && strings.HasSuffix(origin, ":8100")) {
			err := fmt.Errorf("Unknown origin: %v", origin)
			log.Debugf(c, err.Error())
			BadRequestError(c, w, err)
			return
		}
	}
	log.Debugf(c, "Request 'Origin' header: %v", origin)
	responseHeader := w.Header()
	if accessControlRequestMethod := r.Header.Get("Access-Control-Request-Method"); !(accessControlRequestMethod == "GET" || accessControlRequestMethod == "POST") {
		BadRequestMessage(c, w, "Requested method is unsupported: "+accessControlRequestMethod)
		return
	} else {
		responseHeader.Set("Access-Control-Allow-Methods", accessControlRequestMethod)
	}
	if accessControlRequestHeaders := r.Header.Get("Access-Control-Request-Headers"); accessControlRequestHeaders != "" {
		log.Debugf(c, "Request Access-Control-Request-Headers: %v", accessControlRequestHeaders)
		responseHeader.Set("Access-Control-Allow-Headers", accessControlRequestHeaders)
	} else {
		log.Debugf(c, "Request header 'Access-Control-Allow-Headers' is empty or missing")
		// TODO(security): Is it wrong to return 200 in this case?
	}
	responseHeader.Set("Access-Control-Allow-Origin", origin)
}

//func getOnly(handler dtdal.ContextHandler) func(w http.ResponseWriter, r *http.Request) {
//	return dtdal.HttpAppHost.HandleWithContext(optionsHandler(func(c context.Context, w http.ResponseWriter, r *http.Request) {
//		if r.Method != "GET" {
//			BadRequestMessage(c, w, "Expecting to get request method GET, got: "+r.Method)
//			return
//		}
//		hashedWriter := NewHashedResponseWriter(w)
//		handler(c, hashedWriter, r)
//		hashedWriter.setETagOrNotModifiedAndFlushBuffer(c, w, r)
//	}))
//}
//
//func postOnly(handler dtdal.ContextHandler) func(w http.ResponseWriter, r *http.Request) {
//	return dtdal.HttpAppHost.HandleWithContext(optionsHandler(func(c context.Context, w http.ResponseWriter, r *http.Request) {
//		if r.Method != "POST" {
//			BadRequestMessage(c, w, "Expecting to get request method POST, got: "+r.Method)
//			return
//		}
//		handler(c, w, r)
//	}))
//}

func BadRequestMessage(c context.Context, w http.ResponseWriter, m string) {
	log.Infof(c, m)
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(m))
}

func BadRequestError(c context.Context, w http.ResponseWriter, err error) {
	BadRequestMessage(c, w, err.Error())
}

func InternalError(c context.Context, w http.ResponseWriter, err error) {
	m := err.Error()
	log.Errorf(c, m)
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
//func (hashedWriter *HashedResponseWriter) setETagOrNotModifiedAndFlushBuffer(c context.Context, w http.ResponseWriter, r *http.Request) {
//	if hashedWriter.status != 0 && hashedWriter.status != 200 {
//		if _, err := hashedWriter.flush(w); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		if hashedWriter.buffer.Len() > 0 {
//			log.Debugf(c, "No ETag check/set as response.status=%d", hashedWriter.status)
//		}
//		return
//	}
//	eTag := fmt.Sprintf("%x", hashedWriter.hash)
//	if match := r.Header.Get("If-None-Match"); match == eTag {
//		log.Debugf(c, "Setting response status to 304 - not modified")
//		w.WriteHeader(http.StatusNotModified)
//		log.Debugf(c, "Response status set to 304 - not modified")
//	} else {
//		contentLen := hashedWriter.buffer.Len()
//		if (len(eTag) + 5) < contentLen {
//			w.Header().Set("ETag", eTag)
//			log.Debugf(c, "ETag: "+eTag)
//		} else if contentLen > 0 {
//			log.Debugf(c, "ETag is not set as contentLength:%d is smaller then ETag header", contentLen)
//		}
//		if _, err := hashedWriter.flush(w); err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//	}
//}
