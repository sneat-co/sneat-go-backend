package facade4logist

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
	"strings"
	"time"
)

// SetContainerEndpointFields sets dates for a container point
func SetContainerEndpointFields(ctx context.Context, user facade.User, request dto4logist.SetContainerEndpointFieldsRequest) error {
	return RunOrderWorker(ctx, user, request.OrderRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
			return txSetContainerEndpointFields(ctx, tx, params, request)
		},
	)
}

func txSetContainerEndpointFields(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *OrderWorkerParams,
	request dto4logist.SetContainerEndpointFieldsRequest,
) error {
	orderDto := params.Order.Dto
	containerPoint := orderDto.GetContainerPoint(request.ContainerID, request.ShippingPointID)
	if containerPoint == nil {
		return errors.New("container point not found by containerID & shippingPointID")
	}
	var scheduledDatesDiff time.Duration
	if containerPoint.Arrival != nil && containerPoint.Departure != nil {
		if containerPoint.Arrival.ScheduledDate != containerPoint.Departure.ScheduledDate &&
			containerPoint.Arrival.ScheduledDate != "" && containerPoint.Departure.ScheduledDate != "" {
			arrives, err := time.Parse(time.DateOnly, containerPoint.Arrival.ScheduledDate)
			if err != nil {
				return validation.NewErrBadRecordFieldValue("arrival.scheduledDate", err.Error())
			}
			departs, err := time.Parse(time.DateOnly, containerPoint.Departure.ScheduledDate)
			if err != nil {
				return validation.NewErrBadRecordFieldValue("departure.scheduledDate", err.Error())
			}
			scheduledDatesDiff = departs.Sub(arrives)
		}
	}
	var endpoint *models4logist.ContainerEndpoint
	switch request.Side {
	case models4logist.EndpointSideArrival:
		if containerPoint.Arrival == nil {
			containerPoint.Arrival = &models4logist.ContainerEndpoint{}
		}
		endpoint = containerPoint.Arrival
	case models4logist.EndpointSideDeparture:
		if containerPoint.Departure == nil {
			containerPoint.Departure = &models4logist.ContainerEndpoint{}
		}
		endpoint = containerPoint.Departure
	case "":
		return validation.NewErrRecordIsMissingRequiredField("side")
	default:
		return validation.NewErrBadRequestFieldValue("side", "unknown side: "+request.Side)
	}
	if request.ByContactID != nil {
		byContactID := *request.ByContactID
		if byContactID != endpoint.ByContactID {
			_, orderContact := orderDto.WithOrderContacts.GetContactByID(byContactID)
			if orderContact == nil {
				byContact, err := facade4contactus.GetContactByID(ctx, tx, params.TeamWorkerParams.Team.ID, byContactID)
				if err != nil {
					return fmt.Errorf("failed to load 'by' contact: %w", err)
				}
				orderContact = &models4logist.OrderContact{
					ID:   byContactID,
					Type: byContact.Data.Type,
					//CountryID: byContact.Data.CountryID,
					ParentID: byContact.Data.ParentID,
					Title:    byContact.Data.Title,
				}
				if orderContact.CountryID == "" && byContact.Data.Address != nil {
					orderContact.CountryID = byContact.Data.Address.CountryID
				}
				orderDto.Contacts = append(orderDto.Contacts, orderContact)
				params.Changed.Contacts = true
			}
			const roleTrucker = models4logist.CounterpartyRoleTrucker
			_, truckerCounterparty := orderDto.WithCounterparties.GetCounterpartyByRoleAndContactID(roleTrucker, byContactID)
			if truckerCounterparty == nil {
				truckerCounterparty = &models4logist.OrderCounterparty{
					Role:      roleTrucker,
					ContactID: byContactID,
					CountryID: orderContact.CountryID,
					Title:     orderContact.Title,
				}
				orderDto.Counterparties = append(orderDto.Counterparties, truckerCounterparty)
				params.Changed.Counterparties = true
			}
			endpoint.ByContactID = byContactID
		}
	}
	for name, value := range request.Dates {
		switch strings.TrimSpace(name) {
		case "scheduledDate":
			endpoint.ScheduledDate = value
		case "actualDate":
			endpoint.ActualDate = value
		default:
			return validation.NewErrBadRequestFieldValue(request.Side+".name", "unknown name: "+name)
		}
	}
	for name, value := range request.Times {
		switch strings.TrimSpace(name) {
		case "scheduledTime":
			endpoint.ScheduledTime = value
		case "actualTime":
			endpoint.ActualTime = value
		default:
			return validation.NewErrBadRequestFieldValue(request.Side+".name", "unknown name: "+name)
		}
	}
	if endpoint.IsEmpty() {
		endpoint = nil
		switch request.Side {
		case models4logist.EndpointSideArrival:
			containerPoint.Arrival = endpoint
		case models4logist.EndpointSideDeparture:
			containerPoint.Departure = endpoint
		}
	}
	if request.Side == models4logist.EndpointSideArrival {
		if request.ByContactID != nil && containerPoint.Arrival.ByContactID != "" && (containerPoint.Departure == nil || containerPoint.Departure.ByContactID == "") {
			if containerPoint.Departure == nil {
				containerPoint.Departure = &models4logist.ContainerEndpoint{}
			}
			containerPoint.Departure.ByContactID = containerPoint.Arrival.ByContactID
		}
		if request.Dates["scheduledDate"] != "" && (containerPoint.Departure == nil || (containerPoint.Departure.ScheduledDate == "" || containerPoint.Departure.ScheduledDate < containerPoint.Arrival.ScheduledDate)) {
			if containerPoint.Departure == nil {
				containerPoint.Departure = &models4logist.ContainerEndpoint{}
			}
			// Ignore error as it was validated before
			arrives, _ := time.Parse(time.DateOnly, containerPoint.Arrival.ScheduledDate)
			containerPoint.Departure.ScheduledDate = arrives.Add(scheduledDatesDiff).Format(time.DateOnly)
		}
	}
	params.Changed.ContainerPoints = true
	return nil
}
