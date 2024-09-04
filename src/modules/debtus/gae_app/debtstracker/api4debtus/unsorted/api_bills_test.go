package unsorted

import (
	"encoding/json"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus/dto4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtmocks"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"context"
	"github.com/strongo/decimal"
)

const (
	mockBillID    = "123"
	creatorUserID = "1"
)

func TestBillApiCreateBill(t *testing.T) {
	t.Skip("TODO: fix")
	c := context.Background()
	dtmocks.SetupMocks(c)

	if contact, err := dtdal.Contact.InsertContact(c, nil, &models4debtus.DebtusSpaceContactDbo{
		//UserID: creatorUserID,
		ContactDetails: dto4contactus.ContactDetails{
			//FirstName: "First",
		},
	}); err != nil {
		t.Fatal(err)
	} else if contact.ID != "1" {
		t.Fatalf("contact.ContactID: %v", contact.ID)
	}
	// if contact, err := dtdal.DebtusSpaceContactEntry.InsertContact(c, &models.DebtusSpaceContactDbo{
	// 	UserID: creatorUserID,
	// 	ContactDetails: models.ContactDetails{
	// 		FirstName: "Second",
	// 	},
	// }); err != nil {
	// 	t.Fatal(err)
	// } else if contact.ContactID != 2 {
	// 	t.Fatalf("contact.ContactID != 2: %v", contact.ContactID)
	// }
	// if contact, err := dtdal.DebtusSpaceContactEntry.InsertContact(c, &models.DebtusSpaceContactDbo{
	// 	UserID: creatorUserID,
	// 	ContactDetails: models.ContactDetails{
	// 		FirstName: "Third",
	// 	},
	// }); err != nil {
	// 	t.Fatal(err)
	// } else if contact.ContactID != 3 {
	// 	t.Fatalf("contact.ContactID != 3: %v", contact.ContactID)
	// }

	responseRecorder := httptest.NewRecorder()

	body := strings.NewReader("")
	request, err := http.NewRequest(http.MethodPost, "/api/bill-create", body)
	if err != nil {
		t.Fatal(err)
	}
	HandleCreateBill(c, responseRecorder, request, token4auth.AuthInfo{UserID: mockBillID})

	if responseRecorder.Code != http.StatusBadRequest {
		t.Error("Expected to return http.StatusBadRequest on empty request body")
		return
	}

	form := make(url.Values, 3)
	form.Add("name", "Test bill")
	form.Add("currency", "EUR")
	form.Add("amount", "0.10")
	form.Add("split", "percentage")
	form.Add("members", `
	[
		{"UserID":"1","Percent":34,"Amount":0.04},
		{"ContactID":"62","Percent":33,"Amount":0.03},
		{"ContactID":"63","Percent":33,"Amount":0.03}
	]`)

	//body = strings.NewReader("name=Test+bill&currency=EUR&amount=1.23")
	responseRecorder = httptest.NewRecorder()
	request = &http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/api/bill-create"}, PostForm: form}
	HandleCreateBill(c, responseRecorder, request, token4auth.AuthInfo{UserID: creatorUserID})

	if responseRecorder.Code != http.StatusOK {
		t.Errorf(`Expected to get http.StatusOK==200, got responseRecorder.Code=%v
--- Response body ---
%v
--- End of response body ---
Request data: %v`,
			responseRecorder.Code, responseRecorder.Body.String(), form)
		return
	}

	responseObject := make(map[string]dto4debtus.BillDto, 1)

	responseBody := responseRecorder.Body.Bytes()
	if err = json.Unmarshal(responseBody, &responseObject); err != nil {
		t.Errorf("Response(code=%v) body is not valid JSON: %v", responseRecorder.Code, string(responseBody))
		return
	}
	responseBill := responseObject["Bill"]
	if responseBill.ID == "" {
		t.Errorf("Response Bill.ContactID field is empty: %v", string(responseBody))
	}
	if responseBill.Name != "Test bill" {
		t.Errorf("Response Bill.ContactName field has unexpected value: %v\n%v", responseBill.Name, string(responseBody))
	}
	if responseBill.Amount.Currency != "EUR" {
		t.Errorf("Response Bill.AmountTotal.Currency field has unexpected value: %v\n%v", responseBill.Amount.Currency, string(responseBody))
	}
	if responseBill.Amount.Value != decimal.NewDecimal64p2FromFloat64(0.10) {
		t.Errorf("Response Bill.AmountTotal.Value field has unexpected value: %v\n%v", responseBill.Amount.Value, string(responseBody))
	}
}
