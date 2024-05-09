package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

// AddOrderShippingPoint adds shipping point to an order
func AddOrderShippingPoint(
	ctx context.Context,
	user facade.User,
	request dto4logist.AddOrderShippingPointRequest,
) (
	response dto4logist.OrderResponse,
	err error,
) {
	err = RunOrderWorker(ctx, user, request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) (err error) {
		_, err = addOrderShippingPointTx(ctx, tx, request, params)
		response.OrderDto = params.Order.Dto
		return
	})
	return response, err
}

func addOrderShippingPointTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	request dto4logist.AddOrderShippingPointRequest,
	params *OrderWorkerParams,
) (
	shippingPoint *models4logist.OrderShippingPoint,
	err error,
) {
	orderDto := params.Order.Dto
	if err = orderDto.Validate(); err != nil {
		return nil, fmt.Errorf("order record is not valid after loading from DB: %w", err)
	}

	locationContact := dal4contactus.NewContactEntry(request.TeamID, request.LocationContactID)
	if err := tx.Get(ctx, locationContact.Record); err != nil {
		return nil, fmt.Errorf("failed to get locationContact referenced by shipping point: %w", err)
	}
	locationContact.Data.UpdateRelatedIDs()
	if err := locationContact.Data.Validate(); err != nil {
		return nil, fmt.Errorf("locationContact record referenced by request.LocationContactID is not valid (ID=%s): %w", locationContact.ID, err)
	}
	if locationContact.Data.Type != "location" {
		return nil, fmt.Errorf("locationContact referenced by shipping point is not a location: %w", err)
	}
	if locationContact.Data.ParentID == "" {
		return nil, validation.NewErrBadRecordFieldValue("parentContactID", "locationContact referenced by shipping point has no parent contact ContactID")
	}
	counterpartyContact := dal4contactus.NewContactEntry(request.TeamID, locationContact.Data.ParentID)
	if err := tx.Get(ctx, counterpartyContact.Record); err != nil {
		return nil, fmt.Errorf("failed to get counterpartyContact referenced by location point: %w", err)
		//} else if !counterpartyContact.Record.Exists() {
		//	return nil, fmt.Errorf("counterpartyContact referenced by location point does not exist (id=%v): %w", locationContact.Data.ParentID, err)
	}
	counterpartyContact.Data.UpdateRelatedIDs()
	if err := counterpartyContact.Data.Validate(); err != nil {
		return nil, fmt.Errorf("counterpartyContact record referenced by location contact and loaded from DB is not valid (ID=%s): %w", counterpartyContact.ID, err)
	}

	for _, container := range request.Containers {
		for _, task := range container.Tasks {
			if slice.Index(request.Tasks, task) < 0 {
				request.Tasks = append(request.Tasks, task)
			}
		}
	}
	shippingPoint = &models4logist.OrderShippingPoint{
		ID: orderDto.NewOrderShippingPointID(),
		ShippingPointBase: models4logist.ShippingPointBase{
			Status: "pending",
			FreightPoint: models4logist.FreightPoint{
				Tasks: request.Tasks,
			},
		},
		Location: &models4logist.ShippingPointLocation{
			ContactID: request.LocationContactID,
			Title:     locationContact.Data.Title,
			Address:   locationContact.Data.Address,
		},
		Counterparty: models4logist.ShippingPointCounterparty{
			ContactID: counterpartyContact.ID,
			Title:     counterpartyContact.Data.Title,
		},
	}
	orderDto.ShippingPoints = append(orderDto.ShippingPoints, shippingPoint)
	params.Changed.ShippingPoints = true

	for _, task := range request.Tasks {
		var counterpartyRole, locationRole models4logist.CounterpartyRole
		switch task {
		case models4logist.ShippingPointTaskLoad:
			counterpartyRole = models4logist.CounterpartyRoleDispatcher
			locationRole = models4logist.CounterpartyRoleDispatchPoint
		case models4logist.ShippingPointTaskUnload:
			counterpartyRole = models4logist.CounterpartyRoleReceiver
			locationRole = models4logist.CounterpartyRoleReceivePoint
		}
		if _, locationCounterparty := orderDto.GetCounterpartyByRoleAndContactID(locationRole, locationContact.ID); locationCounterparty == nil {
			locationCounterparty = &models4logist.OrderCounterparty{
				Role:      locationRole,
				ContactID: locationContact.ID,
				Title:     locationContact.Data.Title,
				CountryID: locationContact.Data.CountryID,
				Parent: &models4logist.CounterpartyParent{
					ContactID: counterpartyContact.ID,
					Role:      counterpartyRole,
				},
			}
			_, locationOrderContact := orderDto.GetContactByID(locationContact.ID)
			if locationOrderContact == nil {
				locationOrderContact = &models4logist.OrderContact{
					ID:        locationContact.ID,
					Type:      locationContact.Data.Type,
					ParentID:  locationContact.Data.ParentID,
					CountryID: locationContact.Data.CountryID,
					Title:     locationContact.Data.Title,
				}
				//if locationContact.Data.Address != nil {
				//	locationOrderContact.Address = *locationContact.Data.Address
				//}
				orderDto.Contacts = append(orderDto.Contacts, locationOrderContact)
				params.Changed.Contacts = true
			}
			orderDto.Counterparties = append(orderDto.Counterparties, locationCounterparty)
			params.Changed.Counterparties = true
		}
		if _, counterparty := orderDto.GetCounterpartyByRoleAndContactID(counterpartyRole, counterpartyContact.ID); counterparty == nil {
			counterparty = &models4logist.OrderCounterparty{
				Role:      counterpartyRole,
				ContactID: counterpartyContact.ID,
				Title:     counterpartyContact.Data.Title,
				CountryID: counterpartyContact.Data.CountryID,
			}
			_, counterpartyOrderContact := orderDto.GetContactByID(counterpartyContact.ID)
			if counterpartyOrderContact == nil {
				counterpartyOrderContact = &models4logist.OrderContact{
					ID:        counterpartyContact.ID,
					Type:      counterpartyContact.Data.Type,
					ParentID:  counterpartyContact.Data.ParentID,
					CountryID: counterpartyContact.Data.CountryID,
					Title:     counterpartyContact.Data.Title,
					//Address: dbmodels.Address{
					//	CountryID: counterpartyContact.Data.CountryID,
					//},
				}
				orderDto.Contacts = append(orderDto.Contacts, counterpartyOrderContact)
				params.Changed.Contacts = true
			}
			orderDto.Counterparties = append(orderDto.Counterparties, counterparty)
			params.Changed.Counterparties = true
		}
	}

	for _, container := range request.Containers {
		_, container := orderDto.GetContainerByID(container.ID)
		if container == nil {
			return nil, fmt.Errorf("container with ContactID=[%s] not found", container.ID)
		}
		containerPoint := &models4logist.ContainerPoint{
			ContainerID:       container.ID,
			ShippingPointID:   shippingPoint.ID,
			ShippingPointBase: shippingPoint.ShippingPointBase,
		}
		orderDto.ContainerPoints = append(orderDto.ContainerPoints, containerPoint)
		params.Changed.ContainerPoints = true
	}

	return shippingPoint, nil
}
