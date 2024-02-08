package models

import (
	"fmt"
	"time"

	"errors"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/strongo/decimal"
)

//go:generate ffjson $GOFILE

type TransferReturnJson struct {
	TransferID string
	Time       time.Time
	Amount     decimal.Decimal64p2 `json:",omitempty"` // TODO: For legacy records, consider removing later
}

type TransferReturns []TransferReturnJson

func (t *TransferData) GetReturns() (returns TransferReturns) {
	if t.ReturnsCount == 0 && t.ReturnsJson == "" {
		return
	}
	if t.returns != nil {
		if len(t.returns) != t.ReturnsCount {
			panic(fmt.Sprintf("len(t.returns) != t.ReturnsCount: %v != %v", len(t.returns), t.ReturnsCount))
		}
		returns = make(TransferReturns, t.ReturnsCount)
		copy(returns, t.returns)
		return
	}
	if err := ffjson.Unmarshal([]byte(t.ReturnsJson), &returns); err != nil {
		panic(err)
	}
	if len(returns) != t.ReturnsCount {
		panic(fmt.Sprintf("len(returns) != t.ReturnsCount: %v != %v", len(returns), t.ReturnsCount))
	}
	t.returns = make(TransferReturns, len(returns))
	copy(t.returns, returns)
	return
}

func (t *TransferData) AddReturn(returnTransfer TransferReturnJson) error {
	if returnTransfer.TransferID == "" {
		return errors.New("returnTransfer.TransferID == 0")
	}
	// if returnTransfer.TransferID == t.ID {
	// 	return fmt.Errorf("returnTransfer.TransferID == t.ID => %v", t.ID)
	// }
	if returnTransfer.Time.IsZero() {
		return fmt.Errorf("returnTransfer.Time.IsZero(), ID=%v", returnTransfer.TransferID)
	}
	if returnTransfer.Amount <= 0 {
		return fmt.Errorf("transfer return amount <= 0 (ID=%v, amount=%v)", returnTransfer.TransferID, returnTransfer.Amount)
	}

	returns := t.GetReturns()
	if len(returns) != t.ReturnsCount {
		return fmt.Errorf("transfer data integrity issue: len(returns) != t.ReturnsCount => %v != %v", len(returns), t.ReturnsCount)
	}
	var returnedValue decimal.Decimal64p2
	returnTransferIDs := make([]string, 1, len(returns)+1)
	returnTransferIDs[0] = returnTransfer.TransferID
	for _, r := range returns {
		for _, tID := range returnTransferIDs {
			if tID == r.TransferID {
				return fmt.Errorf("duplicate return transfer ID=%v", r.TransferID)
			}
		}
		returnTransferIDs = append(returnTransferIDs, r.TransferID)
		returnedValue += r.Amount
	}
	if returnedValue != t.AmountInCentsReturned {
		return fmt.Errorf("transfer data integrity issue: sum(returns.Amount) != t.AmountInCentsReturned => %v != %v", returnedValue, t.AmountInCentsReturned)
	}

	transferAmount := t.GetAmount()
	totalDue := transferAmount.Value + t.GetInterestValue(time.Now())
	returnedValue += returnTransfer.Amount
	if returnedValue > totalDue {
		return fmt.Errorf("wrong transfer returns: returnedValue > totalDue (%v > %v)", returnedValue, totalDue)
	}
	if returnedValue == totalDue {
		t.IsOutstanding = false
	}
	returns = append(returns, returnTransfer)
	returnsJson, err := ffjson.Marshal(returns) // We'll assign to t.ReturnsJson when all checks are OK
	if err != nil {
		return fmt.Errorf("failed to marshal transfer returns: %w", err)
	}
	t.ReturnsJson = string(returnsJson)
	t.AmountInCentsReturned = returnedValue
	t.ReturnsCount = len(returns)
	t.returns = returns
	return nil
}
