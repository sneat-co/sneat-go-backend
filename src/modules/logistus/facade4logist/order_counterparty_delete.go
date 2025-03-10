package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"strings"
)

// DeleteOrderCounterparty deletes counterparty from an order
func DeleteOrderCounterparty(
	ctx facade.ContextWithUser,
	request dto4logist.DeleteOrderCounterpartyRequest,
) (err error) {
	err = RunOrderWorker(ctx, ctx.User(), request.OrderRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) (err error) {
			return deleteOrderCounterpartyTxWorker(params, request)
		},
	)
	return err
}

func deleteCounterpartyAndChildren(params *OrderWorkerParams, role dbo4logist.CounterpartyRole, contactID string) {
	if err := dbo4logist.ValidateOrderCounterpartyRoles("role", role); err != nil {
		panic(fmt.Errorf("invalid role: %w", err))
	}
	if strings.TrimSpace(contactID) == "" {
		panic("contactID is empty")
	}
	order := params.Order.Dto
	counterparties := make([]*dbo4logist.OrderCounterparty, 0, len(order.Counterparties))
	for _, cp := range order.Counterparties {
		if cp.ContactID == contactID && cp.Role == role {
			params.Changed.Counterparties = true
			remainingSPs, _ := core.Filter(order.ShippingPoints, func(sp *dbo4logist.OrderShippingPoint) bool {
				return sp.Counterparty.ContactID != contactID
			})
			if len(remainingSPs) != len(order.ShippingPoints) {
				order.ShippingPoints = remainingSPs
				params.Changed.ShippingPoints = true
			}
			order.ContainerPoints, _ = core.Filter(order.ContainerPoints, func(cp *dbo4logist.ContainerPoint) bool {
				for _, sp := range remainingSPs {
					if sp.ID == cp.ShippingPointID {
						return true
					}
				}
				params.Changed.ContainerPoints = true
				return false
			})
			continue
		}
		_, c := order.WithOrderContacts.GetContactByID(cp.ContactID)
		if cp.Parent != nil && cp.Parent.ContactID == c.ParentID && cp.Parent.Role == role {
			params.Changed.Counterparties = true
			continue
		}
		counterparties = append(counterparties, cp)
	}
	order.Counterparties = counterparties
	RemoveUnusedContacts(params)
}

func RemoveUnusedContacts(params *OrderWorkerParams) {
	contacts := make([]*dbo4logist.OrderContact, 0, len(params.Order.Dto.Contacts))
	for _, contact := range params.Order.Dto.Contacts {
		i, _ := params.Order.Dto.WithCounterparties.GetCounterpartyByContactID(contact.ID)
		if i < 0 {
			params.Changed.Contacts = true
			continue
		}
		contacts = append(contacts, contact)
	}
	params.Order.Dto.Contacts = contacts
}

func deleteOrderCounterpartyTxWorker(
	params *OrderWorkerParams,
	request dto4logist.DeleteOrderCounterpartyRequest,
) (err error) {
	order := params.Order
	deleteCounterpartyAndChildren(params, request.Role, request.ContactID)

	deletedShippingPointIDs := make([]string, 0, len(order.Dto.ShippingPoints))
	{ // Delete relevant shipping points
		shippingPoints := make([]*dbo4logist.OrderShippingPoint, 0, len(order.Dto.ShippingPoints))
		for _, sp := range order.Dto.ShippingPoints {
			if request.Role == dbo4logist.CounterpartyRoleDispatchPoint && sp.Location.ContactID == request.ContactID ||
				request.Role == dbo4logist.CounterpartyRoleDispatcher && sp.Counterparty.ContactID == request.ContactID {
				deletedShippingPointIDs = append(deletedShippingPointIDs, sp.ID)
				continue
			}
			shippingPoints = append(shippingPoints, sp)
		}
		if len(shippingPoints) != len(order.Dto.ShippingPoints) {
			order.Dto.ShippingPoints = shippingPoints
			params.Changed.ShippingPoints = true
		}
	}

	{ // Delete relevant container points
		containerPoints := make([]*dbo4logist.ContainerPoint, 0, len(order.Dto.ContainerPoints))
		for _, cp := range order.Dto.ContainerPoints {
			if slice.Index(deletedShippingPointIDs, cp.ShippingPointID) >= 0 {
				continue
			}
			containerPoints = append(containerPoints, cp)
		}
		if len(containerPoints) < len(order.Dto.ContainerPoints) {
			order.Dto.ContainerPoints = containerPoints
			params.Changed.ContainerPoints = true
		}
	}

	{ // Delete relevant segments
		segments := make([]*dbo4logist.ContainerSegment, 0, len(order.Dto.Segments))
		for _, segment := range order.Dto.Segments {
			if request.Role == dbo4logist.CounterpartyRoleTrucker &&
				segment.ByContactID == request.ContactID ||
				slice.Index(deletedShippingPointIDs, segment.From.ShippingPointID) >= 0 ||
				slice.Index(deletedShippingPointIDs, segment.To.ShippingPointID) >= 0 {
				continue
			}
			segments = append(segments, segment)
		}
		if len(segments) < len(order.Dto.Segments) {
			order.Dto.Segments = segments
			params.Changed.Segments = true
		}
	}

	//{ // Clear truckers
	//	for _, counterparty := range order.Data.Contacts {
	//		if counterparty.Role == dbo4logist.CounterpartyRoleTrucker {
	//			hasSegments := false
	//			for _, segment := range order.Data.Segments {
	//				if segment.By != nil && segment.By.ContactID == counterparty.ContactID {
	//					hasSegments = true
	//					break
	//				}
	//			}
	//			if !hasSegments {
	//
	//			}
	//		}
	//	}
	//}

	return nil
}
