package dbo4logist

import (
	"errors"
	"fmt"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// OrderBase is a base for order - we use it for storing orders in Firestore and in request for creating a new order
type OrderBase struct {
	// draft, quote_requested, quoting, reviewing_quote, active, canceled, completed, archived
	Status OrderStatus `json:"status" firestore:"status"`

	// import, export, internal
	Direction OrderDirection `json:"direction" firestore:"direction"`

	Route *OrderRoute `json:"route,omitempty" firestore:"route"`

	//Buyer       *OrderCounterparty `json:"buyer" firestore:"buyer"`
	//BuyerAgent  *OrderCounterparty `json:"buyerAgent" firestore:"buyerAgent"`
	//By *OrderCounterparty `json:"transporter" firestore:"transporter"`
	//Seller      *OrderCounterparty `json:"seller" firestore:"seller"`
	//SellerAgent *OrderCounterparty `json:"sellerAgent" firestore:"sellerAgent"`

	WithCounterparties

	with.FlagsField
}

type OrderStatus string

const (
	OrderStatusDraft          OrderStatus = "draft"
	OrderStatusQuoteRequested OrderStatus = "quote_requested"
	OrderStatusQuoting        OrderStatus = "quoting"
	OrderStatusReviewingQuote OrderStatus = "reviewing_quote"
	OrderStatusActive         OrderStatus = "active"
	OrderStatusCanceled       OrderStatus = "canceled"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusArchived       OrderStatus = "archived"
)

type OrderDirection string

const (
	OrderDirectionImport   = "import"
	OrderDirectionExport   = "export"
	OrderDirectionInternal = "internal"
)

// Validate validates OrderBase
func (v OrderBase) Validate() error {
	var errs []error
	switch v.Status {
	case OrderStatusDraft:
		break // known status
	case OrderStatusQuoteRequested, OrderStatusQuoting, OrderStatusReviewingQuote, OrderStatusActive, OrderStatusCanceled, OrderStatusCompleted, OrderStatusArchived:
		if len(v.Counterparties) == 0 {
			errs = append(errs, validation.NewErrRecordIsMissingRequiredField("counterparties"))
		}
	case "":
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("status"))
	default:
		errs = append(errs, validation.NewErrBadRecordFieldValue("status", fmt.Sprintf("unknown value: [%v]", v.Status)))
	}
	switch v.Direction {
	case OrderDirectionImport, OrderDirectionExport, OrderDirectionInternal:
		break // OK
	case "":
		errs = append(errs, validation.NewErrRecordIsMissingRequiredField("direction"))
	default:
		errs = append(errs, validation.NewErrBadRecordFieldValue("direction", fmt.Sprintf("unknonw value: [%v]", v.Direction)))
	}

	if len(errs) == 1 {
		return errs[0]
	} else if len(errs) > 1 {
		return fmt.Errorf("%d validation errors:\n%w", len(errs), errors.Join(errs...))
	}

	if v.Route == nil {
		//if v.Status != "draft" {
		//	return validation.NewErrRecordIsMissingRequiredField("route")
		//}
	} else if err := v.Route.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("route", err.Error())
	}
	if err := v.validateCounterparties(); err != nil {
		return err
	}
	if err := v.FlagsField.Validate(); err != nil {
		return err
	}
	return nil
}

func (v OrderBase) validateCounterparties() error {
	if err := v.WithCounterparties.Validate(); err != nil {
		return err
	}
	contactIDsWithRole := make([]string, len(v.Counterparties))

	for i, c := range v.Counterparties {
		if err := v.validateCounterparty(c, contactIDsWithRole); err != nil {
			return validation.NewErrBadRecordFieldValue(
				fmt.Sprintf("counterparties[%v]", i), err.Error())
		}
	}
	//if v.Buyer != nil {
	//	if err := v.Buyer.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("buyer", err.Error())
	//	}
	//}
	//if v.Seller != nil {
	//	if err := v.Seller.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("seller", err.Error())
	//	}
	//}
	//if v.By != nil {
	//	if err := v.By.Validate(); err != nil {
	//		return validation.NewErrBadRecordFieldValue("transporter", err.Error())
	//	}
	//}
	return nil
}

func (v OrderBase) validateCounterparty(cp *OrderCounterparty, contactIDsWithRole []string) error {
	if cp == nil {
		return nil
	}
	if err := cp.Validate(); err != nil {
		return err
	}
	if cp.Parent != nil {
		_, parent := v.GetCounterpartyByRoleAndContactID(cp.Parent.Role, cp.Parent.ContactID)
		if parent == nil {
			return validation.NewErrBadRecordFieldValue("parent", fmt.Sprintf("counterparty not found in order by (role=%s, contactID=%s)", cp.Parent.Role, cp.Parent.ContactID))
		}
	}
	id := fmt.Sprintf("%s:%s", cp.ContactID, cp.Role)
	if slice.Index(contactIDsWithRole, id) >= 0 {
		return fmt.Errorf("at least 2 counterparties have same role and contactID: %v-%v",
			cp.ContactID, cp.Role)
	}

	return nil
}

// OrderBrief is brief information about an order
type OrderBrief struct {
	ID string `json:"id" firestore:"id"`
	OrderBase
}

// Validate validates OrderBrief
func (v OrderBrief) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if err := v.OrderBase.Validate(); err != nil {
		return err
	}
	return nil
}
