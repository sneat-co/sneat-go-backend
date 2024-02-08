package api

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/strongoapp"
	"net/http"
	//"encoding/json"
	"fmt"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/log"
)

func getEnvironment(r *http.Request) string {
	switch r.Host {
	case "debtstracker.io":
		return "prod"
	case "debtstracker-dev1.appspot.com":
		return "dev"
	case "debtstracker.local":
		return strongoapp.LocalHostEnv
	case "localhost":
		return strongoapp.LocalHostEnv
	default:
		panic("Unknown host: " + r.Host)
	}
}

func getStrID(c context.Context, w http.ResponseWriter, r *http.Request, idParamName string) string {
	q := r.URL.Query()
	if idParamName == "" {
		panic("idParamName is not specified")
	}
	idParamVal := q.Get(idParamName)
	if idParamVal == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Missing parameter: '" + idParamName + "'"))
		return ""
	}
	return idParamVal
}

func hasError(c context.Context, w http.ResponseWriter, err error, entity string, id string, notFoundStatus int) bool {
	switch {
	case err == nil:
		return false
	case dal.IsNotFound(err):
		if notFoundStatus == 0 {
			notFoundStatus = http.StatusNotFound
		}
		w.WriteHeader(notFoundStatus)
		m := fmt.Sprintf("Entity %s not found by id: %s", entity, id)
		log.Infof(c, m)
		_, _ = w.Write([]byte(m))
	default:
		err = fmt.Errorf("failed to get entity %v by id=%v: %w", entity, id, err)
		log.Errorf(c, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	return true
}

func jsonToResponse(c context.Context, w http.ResponseWriter, v interface{}) {
	header := w.Header()
	if buffer, err := ffjson.Marshal(v); err != nil {
		log.Errorf(c, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		header.Add("Access-Control-Allow-Origin", "*")
		log.Debugf(c, "w.Header(): %v", header)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		markResponseAsJson(header)
		log.Debugf(c, "w.Header(): %v", header)
		_, err := w.Write(buffer)
		ffjson.Pool(buffer)
		if err != nil {
			InternalError(c, w, err)
		}
	}
}

func ErrorAsJson(c context.Context, w http.ResponseWriter, status int, err error) {
	if status == 0 {
		panic("status == 0")
	}
	if status == http.StatusInternalServerError {
		log.Errorf(c, "Error: %v", err.Error())
	} else {
		log.Infof(c, "Error: %v", err.Error())
	}
	w.WriteHeader(status)
	jsonToResponse(c, w, map[string]string{"error": err.Error()})
}

func markResponseAsJson(header http.Header) {
	header.Add("Content-Type", "application/json")
	header.Add("Access-Control-Allow-Origin", "*")
}
