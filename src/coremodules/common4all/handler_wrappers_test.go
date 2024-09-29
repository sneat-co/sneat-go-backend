package common4all

import (
	"net/http/httptest"
	"strings"
	"testing"

	"context"
	"errors"
)

func TestBadRequest(t *testing.T) {
	// Disable logging
	//testLogger := &logus.TestLogger{}
	//logus.AddLogger(testLogger)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	errMsg := "Test error #1"
	c := context.Background()
	BadRequestError(c, rr, errors.New(errMsg))
	rr.Flush()
	if !strings.Contains(rr.Body.String(), errMsg) {
		t.Error("Output does not contain error message")
	}
	//if len(testLogger.Messages) == 0 {
	//	t.Error("Not logged")
	//}
	//if len(testLogger.Messages) > 1 {
	//	t.Errorf("Logged too many times: %v", len(testLogger.Messages))
	//}
	//logMessage := testLogger.Messages[0]

	//if !strings.Contains(fmt.Sprintf(logMessage.Format, logMessage.Args...), errMsg) {
	//	t.Error("Log message does not contain error message")
	//}
}
