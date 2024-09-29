package common4all

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"net/http"
)

func GetEnvironment(r *http.Request) string {
	switch r.Host {
	case "api.sneat.ws":
		return "prod"
	case "localhost", "local-api.sneat.ws":
		return strongoapp.LocalHostEnv
	default:
		panic("Unknown host: " + r.Host)
	}
}

func GetStrID(ctx context.Context, w http.ResponseWriter, r *http.Request, idParamName string) string {
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

func HasError(ctx context.Context, w http.ResponseWriter, err error, entity string, id string, notFoundStatus int) bool {
	switch {
	case err == nil:
		return false
	case dal.IsNotFound(err):
		if notFoundStatus == 0 {
			notFoundStatus = http.StatusNotFound
		}
		w.WriteHeader(notFoundStatus)
		m := fmt.Sprintf("Entity %s not found by id: %s", entity, id)
		logus.Infof(ctx, m)
		_, _ = w.Write([]byte(m))
	default:
		err = fmt.Errorf("failed to get entity %v by id=%v: %w", entity, id, err)
		logus.Errorf(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
	return true
}

// TODO - replace with generic sneat one
func JsonToResponse(ctx context.Context, w http.ResponseWriter, v interface{}) {
	header := w.Header()
	if buffer, err := ffjson.Marshal(v); err != nil {
		logus.Errorf(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		header.Add("Access-Control-Allow-Origin", "*")
		logus.Debugf(ctx, "w.Header(): %v", header)
		_, _ = w.Write([]byte(err.Error()))
	} else {
		MarkResponseAsJson(header)
		logus.Debugf(ctx, "w.Header(): %v", header)
		_, err := w.Write(buffer)
		ffjson.Pool(buffer)
		if err != nil {
			InternalError(ctx, w, err)
		}
	}
}

func ErrorAsJson(ctx context.Context, w http.ResponseWriter, status int, err error) {
	if status == 0 {
		panic("status == 0")
	}
	if status == http.StatusInternalServerError {
		logus.Errorf(ctx, "Error: %v", err.Error())
	} else {
		logus.Infof(ctx, "Error: %v", err.Error())
	}
	w.WriteHeader(status)
	JsonToResponse(ctx, w, map[string]string{"error": err.Error()})
}

func MarkResponseAsJson(header http.Header) {
	header.Add("Content-Type", "application/json")
	header.Add("Access-Control-Allow-Origin", "*")
}
