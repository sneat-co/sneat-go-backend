package facade4companies

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// CreateCompanyRequest holds data for models2spotbuddies.CompanyDto creation
type CreateCompanyRequest struct {
	Type  string `json:"type"`
	Title string `json:"title"`
}

// CreateCompanyResponse DTO
type CreateCompanyResponse struct {
	ID string `json:"id"`
}

// CreateCompany creates a company // TODO: Obsolete?
func CreateCompany(ctx context.Context, request CreateCompanyRequest) (response CreateCompanyResponse, err error) {
	userID := facade.GetUserID(ctx)
	if userID == "" {
		return response, fmt.Errorf("context is missing user ContactID")
	}
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		key, err := dal.NewKeyWithOptions("facade4meetingus", dal.WithRandomStringID(dal.RandomLength(5)))
		if err != nil {
			return err
		}
		company := dal.NewRecordWithData(key, nil)
		if err = tx.Insert(ctx, company); err != nil {
			return fmt.Errorf("failed to create company record: %v", err)
		}
		userKey := dbo4userus.NewUserKey(userID)
		userData := make(map[string]interface{})
		userRecord := dal.NewRecordWithData(userKey, userData)
		if err = tx.Get(ctx, userRecord); err != nil {
			return err
		}

		return nil
	})
	return response, err
}
