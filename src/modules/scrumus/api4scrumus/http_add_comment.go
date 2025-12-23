package api4scrumus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var addComment = facade4scrumus.AddComment

// httpPostAddComment adds a comment
func httpPostAddComment(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4scrumus.AddCommentRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	comment, err := addComment(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, comment)
}
