package facade4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/facade"
)

type orderWorker = func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *OrderWorkerParams) (err error)

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
	SpaceWorkerParams *dal4spaceus.SpaceWorkerParams
	Order             dbo4logist.Order
	Changed           OrderChanges
	// OrderUpdates     []update.Update
}

// RunOrderWorker executes an order worker with transaction
var RunOrderWorker = func(ctx facade.ContextWithUser, request dto4logist.OrderRequest, worker orderWorker) (err error) {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("invalid order request: %w", err)
	}
	return dal4spaceus.RunSpaceWorkerWithUserContext(ctx, request.SpaceID,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, teamWorkerParams *dal4spaceus.SpaceWorkerParams) (err error) {
			order := dbo4logist.NewOrder(teamWorkerParams.Space.ID, request.OrderID)
			params := OrderWorkerParams{
				SpaceWorkerParams: teamWorkerParams,
				Order:             order,
			}
			if err := tx.Get(ctx, order.Record); err != nil {
				return fmt.Errorf("failed to load order record: %w", err)
			}
			if err := order.Dto.Validate(); err != nil {
				return fmt.Errorf("loaded order is not valid (ContactID=%s): %w", order.ID, err)
			}
			if err := worker(ctx, tx, &params); err != nil {
				return fmt.Errorf("failed in order worker: %w", err)
			}
			var orderUpdates []update.Update
			if params.Changed.Status {
				orderUpdates = append(orderUpdates, update.ByFieldName("status", order.Dto.Status))
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
				return fmt.Errorf("order is not valid before preparing updates for DB (ContactID=%s): %w", order.ID, err)
			}
			if err := order.Dto.KeysField.Validate(); err != nil {
				return err
			}
			orderUpdates = append(orderUpdates, order.Dto.UpdatesWhenKeysChanged()...)

			order.Dto.UpdateDates()
			if err := order.Dto.DatesFields.Validate(); err != nil {
				return err
			}
			orderUpdates = append(orderUpdates, order.Dto.UpdatesWhenDatesChanged()...)

			order.Dto.MarkAsUpdated(params.SpaceWorkerParams.UserID())
			if err := order.Dto.Validate(); err != nil {
				return fmt.Errorf(
					"order is not valid before pushing updates to DB (ContactID=%s): %w",
					order.ID, err)
			}
			orderUpdates = append(orderUpdates, order.Dto.UpdatesWhenUpdatedFieldsChanged()...)
			if err := tx.Update(ctx, order.Key, orderUpdates); err != nil {
				return fmt.Errorf("failed to update order record: %w", err)
			}
			return nil
		})
}
