package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type orderWorker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) (err error)

type OrderChanges struct {
	Status          bool
	Contacts        bool
	Containers      bool
	ContainerPoints bool
	Counterparties  bool
	Segments        bool
	ShippingPoints  bool
}

func (v *OrderChanges) HasChanges() bool {
	return v.Contacts || v.Containers || v.ContainerPoints || v.Counterparties || v.Segments || v.ShippingPoints
}

func (v *OrderChanges) AddChanges(v2 OrderChanges) OrderChanges {
	v.Counterparties = v.Counterparties || v2.Counterparties
	v.Contacts = v.Contacts || v2.Contacts
	return *v
}

// OrderWorkerParams passes data to a order worker
type OrderWorkerParams struct {
	TeamWorkerParams *dal4teamus.TeamWorkerParams
	Order            models4logist.Order
	Changed          OrderChanges
	// OrderUpdates     []dal.Update
}

// RunOrderWorker executes an order worker with transaction
var RunOrderWorker = func(ctx context.Context, user facade.User, request dto4logist.OrderRequest, worker orderWorker) (err error) {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("invalid order request: %w", err)
	}
	return dal4teamus.RunTeamWorker(ctx, user, request.TeamRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *dal4teamus.TeamWorkerParams) (err error) {
		order := models4logist.NewOrder(teamWorkerParams.Team.ID, request.OrderID)
		params := OrderWorkerParams{
			TeamWorkerParams: teamWorkerParams,
			Order:            order,
		}
		if err := tx.Get(ctx, order.Record); err != nil {
			return fmt.Errorf("failed to load order record: %w", err)
		}
		if err := order.Dto.Validate(); err != nil {
			return fmt.Errorf("loaded order is not valid (ID=%s): %w", order.ID, err)
		}
		if err := worker(ctx, tx, &params); err != nil {
			return fmt.Errorf("failed in order worker: %w", err)
		}
		var orderUpdates []dal.Update
		if params.Changed.Status {
			orderUpdates = append(orderUpdates, dal.Update{Field: "status", Value: order.Dto.Status})
		}
		if params.Changed.ContainerPoints {
			orderUpdates = append(orderUpdates, order.Dto.WithContainerPoints.Updates()...)
		}
		if params.Changed.Counterparties {
			orderUpdates = append(orderUpdates, order.Dto.WithCounterparties.Updates()...)
		}
		if params.Changed.Contacts {
			orderUpdates = append(orderUpdates, order.Dto.WithOrderContacts.Updates()...)
		}
		if params.Changed.ShippingPoints {
			orderUpdates = append(orderUpdates, order.Dto.WithShippingPoints.Updates()...)
		}
		if params.Changed.Containers {
			orderUpdates = append(orderUpdates, order.Dto.WithOrderContainers.Updates()...)
		}
		if params.Changed.Segments {
			orderUpdates = append(orderUpdates, order.Dto.WithSegments.Updates()...)
		}
		if len(orderUpdates) == 0 {
			return nil
		}
		order.Dto.UpdateKeys()
		if err := order.Dto.Validate(); err != nil {
			return fmt.Errorf("order is not valid before preparing updates for DB (ID=%s): %w", order.ID, err)
		}
		if err := order.Dto.KeysField.Validate(); err != nil {
			return err
		}
		orderUpdates = append(orderUpdates, order.Dto.KeysField.UpdatesWhenKeysChanged()...)

		order.Dto.UpdateDates()
		if err := order.Dto.DatesFields.Validate(); err != nil {
			return err
		}
		orderUpdates = append(orderUpdates, order.Dto.DatesFields.UpdatesWhenDatesChanged()...)

		order.Dto.WithModified.MarkAsUpdated(params.TeamWorkerParams.UserID)
		if err := order.Dto.Validate(); err != nil {
			return fmt.Errorf(
				"order is not valid before pushing updates to DB (ID=%s): %w",
				order.ID, err)
		}
		orderUpdates = append(orderUpdates, order.Dto.UpdatedFields.UpdatesWhenUpdatedFieldsChanged()...)
		if err := tx.Update(ctx, order.Key, orderUpdates); err != nil {
			return fmt.Errorf("failed to update order record: %w", err)
		}
		return nil
	})
}
