package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	dbmodels2 "github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"strconv"
	"time"
)

// CreateOrder creates a new order
func CreateOrder(
	ctx context.Context,
	userContext facade.User,
	request dto4logist.CreateOrderRequest,
) (
	orderBrief *models4logist.OrderBrief, err error,
) {
	err = dal4teamus.RunTeamWorker(ctx, userContext, request.TeamID,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamWorkerParams) (err error) {
			orderBrief, err = createOrderTxWorker(ctx, tx, params, params.UserID, request)
			return err
		},
	)
	return orderBrief, err
}

func createOrderTxWorker(
	ctx context.Context, tx dal.ReadwriteTransaction, params *dal4teamus.TeamWorkerParams,
	userID string,
	request dto4logist.CreateOrderRequest,
) (orderBrief *models4logist.OrderBrief, err error) {
	logistTeam := models4logist.NewLogistTeamContext(request.TeamID)
	if err = tx.Get(ctx, logistTeam.Record); err != nil {
		if dal.IsNotFound(err) {
			err = nil // OK
		} else {
			return nil, fmt.Errorf("failed to load logistus team record: %w", err)
		}
	}

	if logistTeam.Dto.OrderCounters == nil {
		logistTeam.Dto.OrderCounters = make(map[string]models4logist.OrderCounter, 1)
	}
	const counterName = "all"
	counter := logistTeam.Dto.OrderCounters[counterName]
	counter.LastNumber++
	logistTeam.Dto.OrderCounters[counterName] = counter

	if err := logistTeam.Dto.Validate(); err != nil {
		return nil, fmt.Errorf("logistus team record is not valid: %w", err)
	}

	orderNumberPrefixed := logistTeam.Dto.OrderCounters["all"].Prefix + strconv.Itoa(counter.LastNumber)
	order := models4logist.NewOrder(params.Team.ID, orderNumberPrefixed)
	fillOrderDtoFromRequest(order.Dto, request, params, userID)

	if err := addContactsFromCounterparties(ctx, tx, params.Team.ID, order.Dto); err != nil {
		return nil, fmt.Errorf("failed to add contacts from counterparties: %w", err)
	}

	if logistTeam.Record.Exists() {
		logistTeamUpdates := []dal.Update{
			{Field: "orderCounters.all.lastNumber", Value: counter.LastNumber},
		}
		if err := tx.Update(ctx, logistTeam.Key, logistTeamUpdates); err != nil {
			return nil, fmt.Errorf("failed to update logistus team record: %w", err)
		}
	} else if err := tx.Insert(ctx, logistTeam.Record); err != nil {
		return nil, fmt.Errorf("failed to insert logistus team record: %w", err)
	}

	//order.Data.OrderIDs = []string{params.Team.ContactID + ":" + orderNumberPrefixed}

	order.Dto.UpdateKeys()
	order.Dto.UpdateDates()
	if err := order.Dto.Validate(); err != nil {
		return nil, fmt.Errorf("order record is not valid before insert: %w", err)
	}
	if err := tx.Insert(ctx, order.Record); err != nil {
		return nil, fmt.Errorf("failed to insert order record: %w", err)
	}

	orderBrief = &models4logist.OrderBrief{
		ID:        order.Key.ID.(string),
		OrderBase: order.Dto.OrderBase,
	}

	return orderBrief, err
}

func fillOrderDtoFromRequest(orderDto *models4logist.OrderDto, request dto4logist.CreateOrderRequest, params *dal4teamus.TeamWorkerParams, userID string) {
	orderDto.OrderBase = request.Order

	orderDto.Status = "active"
	orderDto.UserIDs = params.Team.Data.UserIDs
	orderDto.TeamID = params.Team.ID
	orderDto.TeamIDs = []string{params.Team.ID}
	modified := dbmodels2.Modified{
		By: userID,
		At: time.Now(),
	}
	orderDto.CreatedFields = with.CreatedFields{
		CreatedAtField: with.CreatedAtField{
			CreatedAt: modified.At,
		},
		CreatedByField: with.CreatedByField{
			CreatedBy: modified.By,
		},
	}
	orderDto.UpdatedFields = with.UpdatedFields{
		UpdatedAt: modified.At,
		UpdatedBy: modified.By,
	}

	if orderDto.Route != nil {
		if countryID := orderDto.Route.Origin.CountryID; slice.Index(orderDto.CountryIDs, countryID) < 0 {
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID)
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID+":origin")
		}
		if countryID := orderDto.Route.Destination.CountryID; slice.Index(orderDto.CountryIDs, countryID) < 0 {
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID)
			orderDto.CountryIDs = append(orderDto.CountryIDs, countryID+":destination")
		}
		for _, t := range orderDto.Route.TransitPoints {
			if countryID := t.CountryID; slice.Index(orderDto.CountryIDs, countryID) < 0 {
				orderDto.CountryIDs = append(orderDto.CountryIDs, countryID)
				orderDto.CountryIDs = append(orderDto.CountryIDs, countryID+":transit")
			}
		}
	}
	for size, count := range request.NumberOfContainers {
		for i := 1; i <= count; i++ {
			container := models4logist.OrderContainer{
				ID: fmt.Sprintf("%s%d", size, i),
				OrderContainerBase: models4logist.OrderContainerBase{
					Type: size,
				},
			}
			orderDto.Containers = append(orderDto.Containers, &container)
		}
	}
}

func addContactsFromCounterparties(ctx context.Context, tx dal.ReadTransaction, teamID string, order *models4logist.OrderDto) error {
	if len(order.Counterparties) == 0 {
		panic("at least 1 counterparty should be added to a new order")
	}
	contacts := make([]dal4contactus.ContactEntry, 0, len(order.Counterparties))
	records := make([]dal.Record, 0, len(order.Counterparties))
	contactIDs := make([]string, 0, len(order.Counterparties))
	for _, cp := range order.Counterparties {
		if slice.Index(contactIDs, cp.ContactID) < 0 {
			contactIDs = append(contactIDs, cp.ContactID)
			contact := dal4contactus.NewContactEntry(teamID, cp.ContactID)
			contacts = append(contacts, contact)
			records = append(records, contact.Record)
		}
	}
	if err := tx.GetMulti(ctx, records); err != nil {
		return fmt.Errorf("failed to load contacts: %w", err)
	}

	for _, contact := range contacts {
		_, orderContact := order.GetContactByID(contact.ID)
		if orderContact == nil {
			orderContact = &models4logist.OrderContact{
				ID:        contact.ID,
				Type:      contact.Data.Type,
				Title:     contact.Data.Title,
				ParentID:  contact.Data.ParentID,
				CountryID: contact.Data.CountryID,
			}
			//if contact.Data.Address != nil {
			//	orderContact.Address = *contact.Data.Address
			//}
			//if orderContact.Address.CountryID == "" {
			//	orderContact.Address.CountryID = contact.Data.CountryID
			//}
			order.Contacts = append(order.Contacts, orderContact)
		}
	}

	return nil
}
