package api4contactus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/models4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/httpmock"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/sneatauth"
	"github.com/strongo/strongoapp/with"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpAddMember(t *testing.T) {

	const teamID = "unit-test"
	request := dal4contactus.CreateMemberRequest{
		TeamRequest: dto4teamus.TeamRequest{
			TeamID: teamID,
		},
		WithRelated: models4linkage.WithRelated{
			Related: models4linkage.RelatedByModuleID{
				const4contactus.ModuleID: models4linkage.RelatedByCollectionID{
					const4contactus.ContactsCollection: []*models4linkage.RelatedItem{
						{
							Keys: []models4linkage.RelatedItemKey{
								{TeamID: "team1", ItemID: "c1"},
							},
							RelatedAs: map[models4linkage.RelationshipID]*models4linkage.Relationship{
								"spouse": {
									//CreatedField: with.CreatedField{
									//	Created: with.Created{
									//		By: "u1",
									//		At: "2020-01-01",
									//	},
									//},
								},
							},
						},
					},
				},
			},
		},
		CreatePersonRequest: dto4contactus.CreatePersonRequest{
			ContactBase: briefs4contactus.ContactBase{
				ContactBrief: briefs4contactus.ContactBrief{
					Type:     briefs4contactus.ContactTypePerson,
					Gender:   "unknown",
					Title:    "Some new members",
					AgeGroup: "unknown",
					RolesField: with.RolesField{
						Roles: []string{const4contactus.TeamMemberRoleContributor},
					},
				},
				Status: "active",
				//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
				//	CountryID: dbmodels.UnknownCountryID,
				//},
				Emails: []dbmodels.PersonEmail{
					{Type: "personal", Address: "someone@example.com"},
				},
			},
		},
	}
	request.CountryID = "IE"

	defer func() {
		apicore.GetAuthTokenFromHttpRequest = nil
	}()
	apicore.GetAuthTokenFromHttpRequest = func(r *http.Request) (token *sneatauth.Token, err error) {
		return &sneatauth.Token{UID: "TestUserID"}, nil
	}

	//t.Log(buffer.String())

	req := httpmock.NewPostJSONRequest("POST", "/v0/team/create_member", request)
	req.Host = "localhost"
	req.Header.Set("Origin", "http://localhost:3000")

	createMember = func(ctx context.Context, userCtx facade.User, request dal4contactus.CreateMemberRequest) (response dto4contactus.CreateContactResponse, err error) {
		if request.TeamID != teamID {
			t.Fatalf("Expected teamID=%v, got: %v", teamID, request.TeamID)
		}
		response.ID = "abc1"
		response.Data = &models4contactus.ContactDbo{
			ContactBase: briefs4contactus.ContactBase{
				ContactBrief: briefs4contactus.ContactBrief{
					Type:  briefs4contactus.ContactTypeCompany,
					Title: "Some company",
					OptionalCountryID: with.OptionalCountryID{
						CountryID: "IE",
					},
					RolesField: with.RolesField{
						Roles: []string{const4contactus.TeamMemberRoleContributor},
					},
				},
				Status: "active",
				//WithRequiredCountryID: dbmodels.WithRequiredCountryID{
				//	CountryID: const4contactus.UnknownCountryID,
				//},
			},
		}
		response.Data = &models4contactus.ContactDbo{
			ContactBase: response.Data.ContactBase,
		}
		return
	}

	const uid = "unit-test-user"
	apicore.GetAuthTokenFromHttpRequest = func(r *http.Request) (token *sneatauth.Token, err error) {
		return &sneatauth.Token{UID: uid}, nil
	}
	//sneatfb.NewFirebaseAuthToken = func(ctx context.Context, fbIDToken func() (string, error), authRequired bool) (*auth.Token, error) {
	//	return &auth.Token{UID: uid}, nil
	//}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httpPostCreateMember)
	handler.ServeHTTP(rr, req)

	responseBody := rr.Body.String()

	if expected := http.StatusCreated; rr.Code != expected {
		t.Fatalf(
			"unexpected status: got (%v) expects (%v): %v",
			rr.Code,
			expected,
			responseBody,
		)
	}
}
