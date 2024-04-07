package facade4logist

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/models4logist"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// AddSegments adds segments to an order
func AddSegments(ctx context.Context, user facade.User, request dto4logist.AddSegmentsRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return RunOrderWorker(ctx, user, request.OrderRequest, func(ctx context.Context, tx dal.ReadwriteTransaction, params *OrderWorkerParams) error {
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

func updateContainerWithAddedSegment(orderDto *models4logist.OrderDto, containerData dto4logist.SegmentContainerData) error {
	_, container := orderDto.GetContainerByID(containerData.ID)

	if container == nil {
		return fmt.Errorf("container not found in order by ID=%s", containerData.ID)
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

	segmentKey := models4logist.ContainerSegmentKey{
		ContainerID: containerData.ID,
		From: models4logist.SegmentEndpoint{
			SegmentCounterparty: request.From.Counterparty,
		},
		To: models4logist.SegmentEndpoint{
			SegmentCounterparty: request.To.Counterparty,
		},
	}
	segment := orderDto.WithSegments.GetSegmentByKey(segmentKey)
	if segment != nil {
		return segmentChanges, fmt.Errorf("segment already exists")
	}
	segment = &models4logist.ContainerSegment{
		ContainerSegmentKey: segmentKey,
	}
	if request.From.Date != "" || request.To.Date != "" {
		segment.Dates = &models4logist.SegmentDates{
			Departs: request.From.Date,
			Arrives: request.To.Date,
		}
	}
	orderDto.Segments = append(orderDto.Segments, segment)
	teamID := params.TeamWorkerParams.Team.ID

	if changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, teamID, orderDto, "from", request.From); err != nil {
		return segmentChanges, err
	} else {
		segmentChanges.AddChanges(changes)
	}
	if changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, teamID, orderDto, "to", request.To); err != nil {
		return segmentChanges.AddChanges(changes), err
	} else {
		segmentChanges.AddChanges(changes)
	}

	if request.By != nil {
		segment.ByContactID = request.By.Counterparty.ContactID
		if changes, err := addCounterpartyToOrderIfNeeded(ctx, tx, teamID, orderDto, "by", dto4logist.AddSegmentEndpoint{
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
	orderDto *models4logist.OrderDto,
	shippingPointID string,
	containerData dto4logist.SegmentContainerData,
	containerEndpoints models4logist.ContainerEndpoints,
) error {
	containerPoint := orderDto.GetContainerPoint(containerData.ID, shippingPointID)
	if containerPoint == nil {
		containerPoint = &models4logist.ContainerPoint{
			ContainerID:     containerData.ID,
			ShippingPointID: shippingPointID,
			ShippingPointBase: models4logist.ShippingPointBase{
				Status: "pending",
				FreightPoint: models4logist.FreightPoint{
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
		if containerEndpoints != (models4logist.ContainerEndpoints{}) {
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
			if containerEndpoints != (models4logist.ContainerEndpoints{}) {
				containerPoint.ContainerEndpoints = containerEndpoints
			}
		}
	}
	return nil
}

//func updateContainerLoadForSegment(params *OrderWorkerParams, containerData dto4logist.SegmentContainerData, segment *models4logist.ContainerSegment) error {
//	orderDto := params.Order.Data
//	containerPoint := orderDto.WithContainerPoints.GetContainerPoint(containerData.ID, segment.From.ShippingPointID)
//	if containerPoint == nil {
//		return fmt.Errorf("container point not found")
//	}
//	return nil
//}

func addOrUpdateShippingPoints(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	params *OrderWorkerParams,
	orderDto *models4logist.OrderDto,
	segment *models4logist.ContainerSegment,
	containerData dto4logist.SegmentContainerData,
) error {
	teamID := params.TeamWorkerParams.Team.ID
	fromShippingPoint, toShippingPoint, err := addShippingPointsToOrderIfNeeded(ctx, tx, teamID, orderDto, segment)
	if err != nil {
		return fmt.Errorf("failed to add shipping points to order: %w", err)
	}
	if fromShippingPoint != nil {
		fromContainerData := containerData
		fromContainerData.ToUnload = nil
		var containerEndpoints models4logist.ContainerEndpoints
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
		var containerDates models4logist.ContainerEndpoints
		if segment.Dates != nil {
			containerDates.Arrival.ScheduledDate = segment.Dates.Arrives
		}
		if err := addContainerPoint(orderDto, toShippingPoint.ID, toContainerData, containerDates); err != nil {
			return fmt.Errorf("failed to add 'to' container point: %w", err)
		}
	}
	//if containerData.ToLoad != nil {
	//	if fromShippingPoint.ToLoad == nil {
	//		fromShippingPoint.ToLoad = &models4logist.FreightLoad{}
	//	}
	//	fromShippingPoint.ToLoad.Add(containerData.ToLoad)
	//}
	//if containerData.ToUnload != nil {
	//	if toShippingPoint.ToUnload == nil {
	//		toShippingPoint.ToUnload = &models4logist.FreightLoad{}
	//	}
	//	toShippingPoint.ToUnload.Add(containerData.ToUnload)
	//}
	return nil
}

func addShippingPointsToOrderIfNeeded(
	ctx context.Context,
	tx dal.ReadTransaction,
	teamID string,
	orderDto *models4logist.OrderDto,
	segment *models4logist.ContainerSegment,
) (fromShippingPoint, toShippingPoint *models4logist.OrderShippingPoint, err error) {
	add := func(end string, point *models4logist.SegmentEndpoint) (
		shippingPoint *models4logist.OrderShippingPoint,
		err error,
	) {
		if shippingPoint, err = addShippingPointToOrderIfNeeded(ctx, tx, teamID, orderDto, end, point.SegmentCounterparty); err != nil {
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
	teamID string,
	orderDto *models4logist.OrderDto,
	end string,
	segmentCounterparty models4logist.SegmentCounterparty,
) (
	shippingPoint *models4logist.OrderShippingPoint,
	err error,
) {
	_, shippingPoint = orderDto.GetShippingPointByContactID(segmentCounterparty.ContactID)
	if shippingPoint != nil {
		return shippingPoint, nil
	}

	location := dal4contactus.NewContactEntry(teamID, segmentCounterparty.ContactID)
	if err := tx.Get(ctx, location.Record); err != nil {
		return nil, fmt.Errorf("failed to get location contact: %w", err)
	}

	if err := location.Data.Validate(); err != nil {
		return nil, fmt.Errorf("contact loaded from DB failed validation (ID=%s): %w", location.ID, err)
	}

	if location.Data.ParentID == "" {
		return shippingPoint, fmt.Errorf("segment counteparty with role=[%s] reference contact by ID=[%s] that has no reference to parent: %w",
			segmentCounterparty.Role, segmentCounterparty.ContactID,
			validation.NewErrRecordIsMissingRequiredField("parentContactID"))
	}

	parent := dal4contactus.NewContactEntry(teamID, location.Data.ParentID)
	if err := tx.Get(ctx, parent.Record); err != nil {
		return shippingPoint, fmt.Errorf("failed to get counterparty contact: %w", err)
	}
	if err := parent.Data.Validate(); err != nil {
		return shippingPoint, fmt.Errorf("parent contact with ID=[%v] loaded from DB failed validation: %w", parent.ID, err)
	}

	shippingPoint = &models4logist.OrderShippingPoint{
		ID: orderDto.NewOrderShippingPointID(),
		ShippingPointBase: models4logist.ShippingPointBase{
			Status: "pending",
		},
	}

	switch segmentCounterparty.Role {
	case models4logist.CounterpartyRoleDispatchPoint:
		shippingPoint.Tasks = []models4logist.ShippingPointTask{models4logist.ShippingPointTaskLoad}
	case models4logist.CounterpartyRoleReceivePoint:
		shippingPoint.Tasks = []models4logist.ShippingPointTask{models4logist.ShippingPointTaskUnload}
	case models4logist.CounterpartyRolePickPoint, models4logist.CounterpartyRoleDropPoint:
		switch end {
		case "from":
			shippingPoint.Tasks = []models4logist.ShippingPointTask{models4logist.ShippingPointTaskPick}
		case "to":
			shippingPoint.Tasks = []models4logist.ShippingPointTask{models4logist.ShippingPointTaskDrop}
		case "":
			panic("parameter `end` is required")
		default:
			panic(fmt.Sprintf("parameter `end` has unexpected value: [%s]", end))
		}
	}

	if parent.Data == nil {
		panic(fmt.Sprintf("parent.Data is nil: %+v", parent))
		//shippingPoint.Counterparty = models4logist.ShippingPointCounterparty{
		//	ContactID: location.ID,
		//	Title:     location.Data.Title,
		//}
	} else {
		shippingPoint.Counterparty = models4logist.ShippingPointCounterparty{
			ContactID: parent.ID,
			Title:     parent.Data.Title,
		}
		shippingPoint.Location = &models4logist.ShippingPointLocation{
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
	teamID string,
	order *models4logist.OrderDto,
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
	contact := dal4contactus.NewContactEntry(teamID, segmentCounterparty.ContactID)
	if err := tx.Get(ctx, contact.Record); err != nil {
		return changes, fmt.Errorf("failed to get %v contact: %w", segmentCounterparty.Role, err)
	}
	_, orderContact := order.GetContactByID(contact.ID)
	if orderContact == nil {
		orderContact = &models4logist.OrderContact{
			ID:        contact.ID,
			Type:      contact.Data.Type,
			ParentID:  contact.Data.ParentID,
			CountryID: contact.Data.CountryID,
			Title:     contact.Data.Title,
		}
		//if orderContact.Type == dbmodels.ContactTypeLocation {
		//	orderContact.Address = *contact.Data.Address
		//} else {
		//	orderContact.Address.CountryID = contact.Data.CountryID
		//}
		order.Contacts = append(order.Contacts, orderContact)
		changes.Contacts = true

		if contact.Data.ParentID != "" {
			_, parentOrderContact := order.GetContactByID(contact.Data.ParentID)
			if parentOrderContact == nil {
				parentContact := dal4contactus.NewContactEntry(teamID, contact.Data.ParentID)
				if err := tx.Get(ctx, parentContact.Record); err != nil {
					return changes, fmt.Errorf("failed to get parent contact: %w", err)
				}
				parentOrderContact = &models4logist.OrderContact{
					ID:        parentContact.ID,
					Type:      parentContact.Data.Type,
					ParentID:  parentContact.Data.ParentID,
					CountryID: contact.Data.CountryID,
					Title:     parentContact.Data.Title,
				}
				//if parentOrderContact.Type == dbmodels.ContactTypeLocation {
				//	parentOrderContact.Address = *parentContact.Data.Address
				//} else {
				//	parentOrderContact.Address.CountryID = parentContact.Data.CountryID
				//}
				order.Contacts = append(order.Contacts, parentOrderContact)
			}
		}
	}
	counterparty := models4logist.OrderCounterparty{
		ContactID: contact.ID,
		Role:      counterpartyRole,
		CountryID: contact.Data.CountryID,
		Title:     contact.Data.Title,
	}
	if counterparty.Role != models4logist.CounterpartyRoleDispatchPoint {
		counterparty.RefNumber = segmentEndpoint.RefNumber
	}
	order.Counterparties = append(order.Counterparties, &counterparty)
	changes.Counterparties = true
	if contact.Data.ParentID != "" {
		parent := dal4contactus.NewContactEntry(teamID, contact.Data.ParentID)
		if err := tx.Get(ctx, parent.Record); err != nil {
			return changes, fmt.Errorf("failed to get parent contact by ID=[%s]: %w", contact.Data.ParentID, err)
		}

		var parentCounterpartyRole models4logist.CounterpartyRole

		switch counterpartyRole {
		case models4logist.CounterpartyRoleDispatchPoint:
			parentCounterpartyRole = models4logist.CounterpartyRoleDispatcher
		case models4logist.CounterpartyRoleReceivePoint:
			parentCounterpartyRole = models4logist.CounterpartyRoleReceiver
		case models4logist.CounterpartyRolePickPoint:
			parentCounterpartyRole = models4logist.CounterpartyRolePortFrom
		case models4logist.CounterpartyRoleDropPoint:
			parentCounterpartyRole = models4logist.CounterpartyRolePortTo
		default:
			return changes, fmt.Errorf("counterparty with role=%s references a contact with ID=%s that unexpectedely has non empty parentContactID=%s",
				counterpartyRole, contact.ID, contact.Data.ParentID)
		}

		counterparty.Parent = &models4logist.CounterpartyParent{
			ContactID: contact.Data.ParentID,
			Role:      parentCounterpartyRole,
		}

		if _, parentCounterparty := order.GetCounterpartyByRoleAndContactID(parentCounterpartyRole, parent.ID); parentCounterparty == nil {
			parentCounterparty = &models4logist.OrderCounterparty{
				ContactID: parent.ID,
				Role:      parentCounterpartyRole,
				Title:     parent.Data.Title,
				CountryID: parent.Data.CountryID,
			}
			order.Counterparties = append(order.Counterparties, parentCounterparty)
		} else if segmentEndpoint.RefNumber != "" && counterparty.Role == models4logist.CounterpartyRoleDispatchPoint && parentCounterparty.RefNumber == "" {
			parentCounterparty.RefNumber = segmentEndpoint.RefNumber
		}
	}
	return changes, nil
}
