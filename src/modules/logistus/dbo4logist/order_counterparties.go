package dbo4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/update"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

// WithCounterparties is a base type that defines Counterparties for OrderDbo
type WithCounterparties struct {
	Counterparties []*OrderCounterparty `json:"counterparties" firestore:"counterparties"`
}

var allowedDuplicateRoles = []string{
	CounterpartyRoleTrucker,
	CounterpartyRoleReceiver,
	CounterpartyRoleDispatcher,
	CounterpartyRoleReceivePoint,
	CounterpartyRoleDispatchPoint,
	CounterpartyRoleDropPoint,
	CounterpartyRolePickPoint,
}

// Validate validates counterparties
func (v WithCounterparties) Validate() error {
	hasRole := make([]string, 0, len(v.Counterparties))

	for i, counterparty := range v.Counterparties {
		if err := counterparty.Validate(); err != nil {
			err = fmt.Errorf("counterarty{role:%v,contactID:%v}: %w",
				counterparty.Role,
				counterparty.ContactID,
				err,
			)
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("counterparties[%v]", i), err.Error())
		}
		for j, c2 := range v.Counterparties {
			if j == i {
				continue
			}
			if c2.Role == counterparty.Role {
				if c2.ContactID == counterparty.ContactID {
					return validation.NewErrBadRecordFieldValue("counterparties", fmt.Sprintf("duplicate counterparty{role:%v,contactID:%v}", counterparty.Role, counterparty.ContactID))
				}
				if slice.Index(allowedDuplicateRoles, counterparty.Role) < 0 {
					if slice.Index(hasRole, counterparty.Role) >= 0 {
						return validation.NewErrBadRecordFieldValue("counterparties", fmt.Sprintf(
							"order has at least 2 counterparties with role [%s] at indexes %d & %d: %+v,%+v",
							counterparty.Role, j, i, counterparty, c2))
					}
					hasRole = append(hasRole, counterparty.Role)
				}
			}
			if counterparty.ContactID == c2.ContactID {
				if counterparty.Role == CounterpartyRolePortFrom && c2.Role == CounterpartyRolePortTo {
					if counterparty.ContactID == c2.ContactID {
						return validation.NewErrBadRecordFieldValue("counterparties", fmt.Sprintf(
							"same counterparty is set as portFrom and portTo at indexes %d & %d: %+v, %+v", i, j, counterparty, c2))
					}
				}
			}
		}
	}
	return nil
}

// GetCounterpartyByRole returns first order counterparty with the given role
func (v WithCounterparties) GetCounterpartyByRole(role CounterpartyRole) (i int, counterparty *OrderCounterparty) {
	for i, c := range v.Counterparties {
		if c.Role == role {
			return i, c
		}
	}
	return -1, nil
}

// GetCounterpartiesByRole returns all order counterparties with the given role
func (v WithCounterparties) GetCounterpartiesByRole(role CounterpartyRole) (counterparties []*OrderCounterparty) {
	counterparties = make([]*OrderCounterparty, 0, len(v.Counterparties))
	for _, c := range v.Counterparties {
		if c.Role == role {
			counterparties = append(counterparties, c)
		}
	}
	return counterparties
}

// GetCounterpartiesByContactID returns all order counterparties with the given contactID
func (v WithCounterparties) GetCounterpartiesByContactID(contactID string) (counterparties []*OrderCounterparty) {
	counterparties = make([]*OrderCounterparty, 0, len(v.Counterparties))
	for _, c := range v.Counterparties {
		if c.ContactID == contactID {
			counterparties = append(counterparties, c)
		}
	}
	return counterparties
}

// GetCounterpartyByRoleAndContactID returns first order counterparty with the given role and contactID
func (v WithCounterparties) GetCounterpartyByRoleAndContactID(role CounterpartyRole, contactID string) (i int, counterparty *OrderCounterparty) {
	for i, counterparty = range v.Counterparties {
		if counterparty.ContactID == contactID && counterparty.Role == role {
			return i, counterparty
		}
	}
	return -1, nil
}

// GetCounterpartyByContactID returns first order counterparty with the given contactID
func (v WithCounterparties) GetCounterpartyByContactID(contactID string) (i int, counterparty *OrderCounterparty) {
	for i, c := range v.Counterparties {
		if c.ContactID == contactID {
			return i, c
		}
	}
	return -1, nil
}

// Updates returns updates for WithCounterparties
func (v WithCounterparties) Updates() []update.Update {
	return []update.Update{update.ByFieldName("counterparties", v.Counterparties)}
}
