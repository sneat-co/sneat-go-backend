package api4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"net/http"
)

// httpPostGetSpaceJoinInfo is an API endpoint that return team info for user willing to join
func httpPostGetSpaceJoinInfo(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequest(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		httpserver.HandleError(ctx, err, "VerifyRequest", w, r)
		return
	}
	q := r.URL.Query()
	request := facade4invitus.JoinInfoRequest{
		InviteID: q.Get("id"),
	}
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
		return
	}
	var response facade4invitus.JoinInfoResponse
	if response, err = facade4invitus.GetSpaceJoinInfo(ctx, request); err != nil {
		return
		//} else {
		//	for i, m := range response.Space.Members {
		//		response.Space.Members[i] = &briefs4memberus.MemberBrief{
		//			MemberBase: briefs4memberus.MemberBase{
		//				ContactBaseWithUserID: dbmodels.ContactBaseWithUserID{
		//					Title:  m.Title,
		//					Roles:  m.Roles,
		//					Avatar: m.Avatar,
		//					Gender: m.Gender,
		//				},
		//			},
		//		}
		//	}
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
