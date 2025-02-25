package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dbo4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// AddSegments adds segments to an order
func AddSegments(ctx context.Context, userCtx facade.UserContext, request dto4logist.AddSegmentsRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return RunOrderWorker(ctx, userCtx, request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
		return addSegmentsTx(ctx, tx, params, request)
	})
}

func addSegmentsTx(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams, request dto4logist.AddSegmentsRequest) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}
	for i, container := range request.Containers {
		if segmentChanges, err := addSegment(ctx, tx, params, request, container); err != nil {
			return fmt.Errorf("failed to add container for segment # %d: %w", i, err)
		} else {
			params.Changed.AddChanges(segmentChanges)
		}
	}

	params.Changed.Counterparties = true
	params.Changed.ShippingPoints = true
	params.Changed.Containers = true
	params.Changed.ContainerPoints = true
	params.Changed.Segments = true
	return nil
}

func updateContainerWithAddedSegment(orderDto *dbo4logist.OrderDbo, containerData dto4logist.SegmentContainerData) error {
	_, container := orderDto.GetContainerByID(containerData.ID)

	if container == nil {
		return fmt.Errorf("container not found in order by containerData.ID=%s", containerData.ID)
	}

	//if !containerData.ToLoad.IsEmpty() {
	//	container.NumberOfPallets += containerData.ToLoad.NumberOfPallets
	//	container.GrossWeightKg += containerData.ToLoad.GrossWeightKg
	//	container.VolumeM3 += containerData.ToLoad.VolumeM3
	//}
	//if !containerData.ToUnload.IsEmpty() {
	//	container.NumberOfPallets -= containerData.ToUnload.NumberOfPallets
	//	container.GrossWeightKg -= containerData.ToUnload.GrossWeightKg
	//	container.VolumeM3 -= containerData.ToUnload.VolumeM3
	//}
	return nil
}

func addSegment(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams, request dto4logist.AddSegmentsRequest, containerData dto4logist.SegmentContainerData) (segmentChanges OrderChanges, err error) {

	orderDto := params.Order.Dto

	if err := updateContainerWithAddedSegment(orderDto, containerData); err != nil {
		return segmentChanges, fmt.Errorf("failed to update container with added segment: %w", err)
	}

	segmentKey := dbo4logist.ContainerSegmentKey{
		ContainerID: containerData.ID,
		From: dbo4logist.SegmentEndpoint{
			SegmentCounterparty: request.From.Counterparty,
		},
		To: dbo4logist.SegmentEndpoint{
			SegmentCounterparty: request.To.Counterparty,
		},
	}
	segment := orderDto.WithSegments.GetSegmentByKey(segmentKey)
	if segment != nil {
		return segmentChanges, fmt.Errorf("segment already exists")
	}
	segment = &dbo4logist.ContainerSegment{
		ContainerSegmentKey: segmentKey,
	}
	if request.From.Date != "" || request.To.Date != "" {
		segment.Dates = &dbo4logist.SegmentDates{
			Departs: request.From.Date,
			Arrives: request.To.Date,
		}
	}
	orderDto.Segments = append(orderDto.Segments, segment)
	spaceID := params.SpaceWorkerParams.Space.ID

	if changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, spaceID, orderDto, "from", request.From); err != nil {
		return segmentChanges, err
	} else {
		segmentChanges.AddChanges(changes)
	}
	if changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, spaceID, orderDto, "to", request.To); err != nil {
		return segmentChanges.AddChanges(changes), err
	} else {
		segmentChanges.AddChanges(changes)
	}

	if request.By != nil {
		segment.ByContactID = request.By.Counterparty.ContactID
		if changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, spaceID, orderDto, "by", dto4logist.AddSegmentEndpoint{
			AddSegmentParty: *request.By,
		}); err != nil {
			return segmentChanges, err
		} else {
			segmentChanges.AddChanges(changes)
		}
	}
	if err := addOrUpdateShippingPoints(ctx, tx, params, orderDto, segment, containerData); err != nil {
		return segmentChanges, fmt.Errorf("failed to add or update shipping points: %w", err)
	}
	{ // Double check that we did what expected
		if request.From.Date != "" && (segment.Dates == nil || segment.Dates.Departs == "") {
			return segmentChanges, fmt.Errorf("segment departs date is empty at the end of adding segment")
		}
		if request.To.Date != "" && (segment.Dates == nil || segment.Dates.Arrives == "") {
			return segmentChanges, fmt.Errorf("segment arrives date is empty at the end of adding segment")
		}
	}
	return segmentChanges, nil
}

func addContainerPoint(
	orderDto *dbo4logist.OrderDbo,
	shippingPointID string,
	containerData dto4logist.SegmentContainerData,
	containerEndpoints dbo4logist.ContainerEndpoints,
) error {
	containerPoint := orderDto.GetContainerPoint(containerData.ID, shippingPointID)
	if containerPoint == nil {
		containerPoint = &dbo4logist.ContainerPoint{
			ContainerID:     containerData.ID,
			ShippingPointID: shippingPointID,
			ShippingPointBase: dbo4logist.ShippingPointBase{
				Status: "pending",
				FreightPoint: dbo4logist.FreightPoint{
					Tasks: containerData.Tasks,
				},
			},
		}
		if !containerData.ToLoad.IsEmpty() {
			containerPoint.ToLoad = containerData.ToLoad
		}
		if !containerData.ToUnload.IsEmpty() {
			containerPoint.ToUnload = containerData.ToUnload
		}
		if containerEndpoints != (dbo4logist.ContainerEndpoints{}) {
			containerPoint.ContainerEndpoints = containerEndpoints
		}
		orderDto.ContainerPoints = append(orderDto.ContainerPoints, containerPoint)
	} else {
		containerPoint.Status = "pending"
		if !containerData.ToLoad.IsEmpty() {
			containerPoint.ToLoad = containerData.ToLoad
		}
		if !containerData.ToUnload.IsEmpty() {
			containerPoint.ToUnload = containerData.ToUnload
		}
		{
			if containerEndpoints != (dbo4logist.ContainerEndpoints{}) {
				containerPoint.ContainerEndpoints = containerEndpoints
			}
		}
	}
	return nil
}

//func updateContainerLoadForSegment(params *OrderWorkerParams, containerData dto4logist.SegmentContainerData, segment *dbo4logist.ContainerSegment) error {
//	orderDto := params.Order.Data
//	containerPoint := orderDto.WithContainerPoints.GetContainerPoint(containerData.ContactID, segment.From.ShippingPointID)
//	if containerPoint == nil {
//		return fmt.Errorf("container point not found")
//	}
//	return nil
//}

func addOrUpdateShippingPoints(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *OrderWorkerParams,
	orderDto *dbo4logist.OrderDbo,
	segment *dbo4logist.ContainerSegment,
	containerData dto4logist.SegmentContainerData,
) error {
	spaceID := params.SpaceWorkerParams.Space.ID
	fromShippingPoint, toShippingPoint, err := addShippingPointsToOrderIfNeeded(ctx, tx, spaceID, orderDto, segment)
	if err != nil {
		return fmt.Errorf("failed to add shipping points to order: %w", err)
	}
	if fromShippingPoint != nil {
		fromContainerData := containerData
		fromContainerData.ToUnload = nil
		var containerEndpoints dbo4logist.ContainerEndpoints
		if segment.Dates != nil {
			containerEndpoints.Departure.ScheduledDate = segment.Dates.Departs
		}
		if err := addContainerPoint(orderDto, fromShippingPoint.ID, containerData, containerEndpoints); err != nil {
			return fmt.Errorf("failed to add 'from' container point: %w", err)
		}
	}
	if toShippingPoint != nil {
		toContainerData := containerData
		toContainerData.ToLoad = nil
		var containerDates dbo4logist.ContainerEndpoints
		if segment.Dates != nil {
			containerDates.Arrival.ScheduledDate = segment.Dates.Arrives
		}
		if err := addContainerPoint(orderDto, toShippingPoint.ID, toContainerData, containerDates); err != nil {
			return fmt.Errorf("failed to add 'to' container point: %w", err)
		}
	}
	//if containerData.ToLoad != nil {
	//	if fromShippingPoint.ToLoad == nil {
	//		fromShippingPoint.ToLoad = &dbo4logist.FreightLoad{}
	//	}
	//	fromShippingPoint.ToLoad.Add(containerData.ToLoad)
	//}
	//if containerData.ToUnload != nil {
	//	if toShippingPoint.ToUnload == nil {
	//		toShippingPoint.ToUnload = &dbo4logist.FreightLoad{}
	//	}
	//	toShippingPoint.ToUnload.Add(containerData.ToUnload)
	//}
	return nil
}

func addShippingPointsToOrderIfNeeded(
	ctx context.Context,
	tx dal.ReadTransaction,
	spaceID coretypes.SpaceID,
	orderDto *dbo4logist.OrderDbo,
	segment *dbo4logist.ContainerSegment,
) (fromShippingPoint, toShippingPoint *dbo4logist.OrderShippingPoint, err error) {
	add := func(end string, point *dbo4logist.SegmentEndpoint) (
		shippingPoint *dbo4logist.OrderShippingPoint,
		err error,
	) {
		if shippingPoint, err = addShippingPointToOrderIfNeeded(ctx, tx, spaceID, orderDto, end, point.SegmentCounterparty); err != nil {
			return shippingPoint, fmt.Errorf("failed to add shipping point for segment endpoint '%s': %w", end, err)
		} else if shippingPoint != nil {
			point.ShippingPointID = shippingPoint.ID
		}
		return
	}
	if fromShippingPoint, err = add("from", &segment.From); err != nil {
		return
	}
	if toShippingPoint, err = add("to", &segment.To); err != nil {
		return
	}
	return
}

func addShippingPointToOrderIfNeeded(
	ctx context.Context,
	tx dal.ReadTransaction,
	spaceID coretypes.SpaceID,
	orderDto *dbo4logist.OrderDbo,
	end string,
	segmentCounterparty dbo4logist.SegmentCounterparty,
) (
	shippingPoint *dbo4logist.OrderShippingPoint,
	err error,
) {
	_, shippingPoint = orderDto.GetShippingPointByContactID(segmentCounterparty.ContactID)
	if shippingPoint != nil {
		return shippingPoint, nil
	}

	location := dal4contactus.NewContactEntry(spaceID, segmentCounterparty.ContactID)
	if err := tx.Get(ctx, location.Record); err != nil {
		return nil, fmt.Errorf("failed to get location contact: %w", err)
	}

	if err := location.Data.Validate(); err != nil {
		return nil, fmt.Errorf("contact loaded from DB failed validation (ContactID=%s): %w", location.ID, err)
	}

	if location.Data.ParentID == "" {
		return shippingPoint, fmt.Errorf("segment counteparty with role=[%s] reference contact by segmentCounterparty.ContactID=[%s] that has no reference to parent: %w",
			segmentCounterparty.Role, segmentCounterparty.ContactID,
			validation.NewErrRecordIsMissingRequiredField("parentContactID"))
	}

	parent := dal4contactus.NewContactEntry(spaceID, location.Data.ParentID)
	if err := tx.Get(ctx, parent.Record); err != nil {
		return shippingPoint, fmt.Errorf("failed to get counterparty contact: %w", err)
	}
	if err := parent.Data.Validate(); err != nil {
		return shippingPoint, fmt.Errorf("parent contact with ContactID=[%v] loaded from DB failed validation: %w", parent.ID, err)
	}

	shippingPoint = &dbo4logist.OrderShippingPoint{
		ID: orderDto.NewOrderShippingPointID(),
		ShippingPointBase: dbo4logist.ShippingPointBase{
			Status: "pending",
		},
	}

	switch segmentCounterparty.Role {
	case dbo4logist.CounterpartyRoleDispatchPoint:
		shippingPoint.Tasks = []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskLoad}
	case dbo4logist.CounterpartyRoleReceivePoint:
		shippingPoint.Tasks = []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskUnload}
	case dbo4logist.CounterpartyRolePickPoint, dbo4logist.CounterpartyRoleDropPoint:
		switch end {
		case "from":
			shippingPoint.Tasks = []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskPick}
		case "to":
			shippingPoint.Tasks = []dbo4logist.ShippingPointTask{dbo4logist.ShippingPointTaskDrop}
		case "":
			panic("parameter `end` is required")
		default:
			panic(fmt.Sprintf("parameter `end` has unexpected value: [%s]", end))
		}
	}

	if parent.Data == nil {
		panic(fmt.Sprintf("parent.Data is nil: %+v", parent))
		//shippingPoint.Counterparty = dbo4logist.ShippingPointCounterparty{
		//	ContactID: location.ContactID,
		//	Title:     location.Data.Title,
		//}
	} else {
		shippingPoint.Counterparty = dbo4logist.ShippingPointCounterparty{
			ContactID: parent.ID,
			Title:     parent.Data.Title,
		}
		shippingPoint.Location = &dbo4logist.ShippingPointLocation{
			ContactID: location.ID,
			Title:     location.Data.Title,
			Address:   location.Data.Address,
		}
	}
	orderDto.ShippingPoints = append(orderDto.ShippingPoints, shippingPoint)
	//panic(fmt.Sprintf("orderDto: %+v", orderDto))
	return shippingPoint, nil
}

func addCounterpartyToOrderIfNeeded(
	ctx context.Context,
	tx dal.ReadTransaction,
	spaceID coretypes.SpaceID,
	order *dbo4logist.OrderDbo,
	endpointType string, // Either "from", "to" or "by"
	segmentEndpoint dto4logist.AddSegmentEndpoint,
) (changes OrderChanges, err error) {
	switch endpointType {
	case "from", "to", "by": // OK
	case "":
		panic("parameter endpointType expected to be either 'from' or 'to', got empty string")
	default:
		panic("parameter endpointType expected to be either 'from' or 'to', got: " + endpointType)
	}
	segmentCounterparty := segmentEndpoint.Counterparty
	counterpartyRole := segmentCounterparty.Role
	if _, c := order.GetCounterpartyByRoleAndContactID(counterpartyRole, segmentCounterparty.ContactID); c != nil {
		return changes, nil
	}
	contact := dal4contactus.NewContactEntry(spaceID, segmentCounterparty.ContactID)
	if err := tx.Get(ctx, contact.Record); err != nil {
		return changes, fmt.Errorf("failed to get %v contact: %w", segmentCounterparty.Role, err)
	}
	_, orderContact := order.GetContactByID(contact.ID)
	if orderContact == nil {
		orderContact = &dbo4logist.OrderContact{
			ID:        contact.ID,
			Type:      contact.Data.Type,
			ParentID:  contact.Data.ParentID,
			CountryID: contact.Data.CountryID,
			Title:     contact.Data.Title,
		}
		//if orderContact.ExtraType == dbmodels.ContactTypeLocation {
		//	orderContact.Address = *contact.Data.Address
		//} else {
		//	orderContact.Address.CountryID = contact.Data.CountryID
		//}
		order.Contacts = append(order.Contacts, orderContact)
		changes.Contacts = true

		if contact.Data.ParentID != "" {
			_, parentOrderContact := order.GetContactByID(contact.Data.ParentID)
			if parentOrderContact == nil {
				parentContact := dal4contactus.NewContactEntry(spaceID, contact.Data.ParentID)
				if err := tx.Get(ctx, parentContact.Record); err != nil {
					return changes, fmt.Errorf("failed to get parent contact: %w", err)
				}
				parentOrderContact = &dbo4logist.OrderContact{
					ID:        parentContact.ID,
					Type:      parentContact.Data.Type,
					ParentID:  parentContact.Data.ParentID,
					CountryID: contact.Data.CountryID,
					Title:     parentContact.Data.Title,
				}
				//if parentOrderContact.ExtraType == dbmodels.ContactTypeLocation {
				//	parentOrderContact.Address = *parentContact.Data.Address
				//} else {
				//	parentOrderContact.Address.CountryID = parentContact.Data.CountryID
				//}
				order.Contacts = append(order.Contacts, parentOrderContact)
			}
		}
	}
	counterparty := dbo4logist.OrderCounterparty{
		ContactID: contact.ID,
		Role:      counterpartyRole,
		CountryID: contact.Data.CountryID,
		Title:     contact.Data.Title,
	}
	if counterparty.Role != dbo4logist.CounterpartyRoleDispatchPoint {
		counterparty.RefNumber = segmentEndpoint.RefNumber
	}
	order.Counterparties = append(order.Counterparties, &counterparty)
	changes.Counterparties = true
	if contact.Data.ParentID != "" {
		parent := dal4contactus.NewContactEntry(spaceID, contact.Data.ParentID)
		if err := tx.Get(ctx, parent.Record); err != nil {
			return changes, fmt.Errorf("failed to get parent contact by contact.Data.ParentID=[%s]: %w", contact.Data.ParentID, err)
		}

		var parentCounterpartyRole dbo4logist.CounterpartyRole

		switch counterpartyRole {
		case dbo4logist.CounterpartyRoleDispatchPoint:
			parentCounterpartyRole = dbo4logist.CounterpartyRoleDispatcher
		case dbo4logist.CounterpartyRoleReceivePoint:
			parentCounterpartyRole = dbo4logist.CounterpartyRoleReceiver
		case dbo4logist.CounterpartyRolePickPoint:
			parentCounterpartyRole = dbo4logist.CounterpartyRolePortFrom
		case dbo4logist.CounterpartyRoleDropPoint:
			parentCounterpartyRole = dbo4logist.CounterpartyRolePortTo
		default:
			return changes, fmt.Errorf("counterparty with role=%s references a contact with ContactID=%s that unexpectedely has non empty parentContactID=%s",
				counterpartyRole, contact.ID, contact.Data.ParentID)
		}

		counterparty.Parent = &dbo4logist.CounterpartyParent{
			ContactID: contact.Data.ParentID,
			Role:      parentCounterpartyRole,
		}

		if _, parentCounterparty := order.GetCounterpartyByRoleAndContactID(parentCounterpartyRole, parent.ID); parentCounterparty == nil {
			parentCounterparty = &dbo4logist.OrderCounterparty{
				ContactID: parent.ID,
				Role:      parentCounterpartyRole,
				Title:     parent.Data.Title,
				CountryID: parent.Data.CountryID,
			}
			order.Counterparties = append(order.Counterparties, parentCounterparty)
		} else if segmentEndpoint.RefNumber != "" && counterparty.Role == dbo4logist.CounterpartyRoleDispatchPoint && parentCounterparty.RefNumber == "" {
			parentCounterparty.RefNumber = segmentEndpoint.RefNumber
		}
	}
	return changes, nil
}
