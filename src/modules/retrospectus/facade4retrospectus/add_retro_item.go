package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	dal4userus2 "github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/random"
	"github.com/strongo/validation"
	"time"
)

// UpcomingRetrospectiveID "upcoming"
const UpcomingRetrospectiveID = "upcoming"

// RetroItemRequest request params
type RetroItemRequest struct {
	facade4meetingus.Request
	Type string `json:"type"`
	Item string `json:"item"`
}

// Validate validates request
func (v *RetroItemRequest) Validate() error {
	if err := v.Request.Validate(); err != nil {
		return err
	}
	if v.Type == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	return nil
}

// AddRetroItemRequest request params
type AddRetroItemRequest struct {
	RetroItemRequest
	Title string `json:"title"`
}

// AddRetroItemResponse response
type AddRetroItemResponse struct {
	ID          string    `json:"id"`
	TimeCreated time.Time `json:"timeCreated"`
}

// Validate validates response
func (v *AddRetroItemRequest) Validate() error {
	if err := v.RetroItemRequest.Validate(); err != nil {
		return err
	}
	if err := validate.RequestTitle(v.Title, "title"); err != nil {
		return err
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// AddRetroItem adds item to retrospective
func AddRetroItem(ctx context.Context, userCtx facade.UserContext, request AddRetroItemRequest) (response AddRetroItemResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	if request.MeetingID == UpcomingRetrospectiveID {
		response, err = addRetroItemToUserRetro(ctx, userCtx, request)
		if err != nil {
			err = fmt.Errorf("failed to add item to future retrospective: %w", err)
		}
	} else {
		response, err = addRetroItemToSpaceRetro(ctx, userCtx, request)
		if err != nil {
			err = fmt.Errorf("failed to add item to specific retrospective: %w", err)
		}
	}
	return
}

func addItemWithUniqueID(item *dbo4retrospectus.RetroItem, items []*dbo4retrospectus.RetroItem) []*dbo4retrospectus.RetroItem {
UniqueID:
	for {
		item.ID = random.ID(5)
		for i := 0; i < len(items); i++ {
			if items[i].ID == item.ID {
				continue UniqueID
			}
		}
		break
	}
	return append(items, item)
}

func addRetroItemToUserRetro(ctx context.Context, userCtx facade.UserContext, request AddRetroItemRequest) (response AddRetroItemResponse, err error) {
	uid := userCtx.GetUserID()

	user := dbo4userus.NewUserEntry(uid)

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		started := time.Now()

		if err = dal4userus2.GetUser(ctx, tx, user); err != nil {
			return err
		}

		userSpaceInfo := user.Data.GetUserSpaceInfoByID(request.SpaceID)
		if userSpaceInfo == nil {
			return validation.NewErrBadRequestFieldValue("space", fmt.Sprintf("user does not belong to this team %v, uid=%v", request.SpaceID, uid))
		}

		//if userSpaceInfo.RetroItems == nil {
		//	userSpaceInfo.RetroItems = make(dbretro.RetroItemsByType)
		//}
		//
		//items, existingType := userSpaceInfo.RetroItems[request.Role]

		//if !existingType {
		//	items = make([]*dbretro.RetroItem, 0, 1)
		//}
		item := dbo4retrospectus.RetroItem{
			Title:   request.Title,
			Created: started,
		}
		//userSpaceInfo.RetroItems[request.Role] = addItemWithUniqueID(&item, items)

		if err = user.Data.Validate(); err != nil {
			return err
		}

		//if err := updateTeamWithUpcomingRetroUserCounts(ctx, tx, started, uid, request.Space, userSpaceInfo.RetroItems); err != nil {
		//	return fmt.Errorf("failed to update team record: %w", err)
		//}

		updates := []dal.Update{
			{
				Field: fmt.Sprintf("api4meetingus.%v.retroItems.%v", request.SpaceID, request.Type),
				Value: dal.ArrayUnion(item),
			},
		}
		if err = dal4userus2.TxUpdateUser(ctx, tx, started, user.Key, updates); err != nil {
			return fmt.Errorf("failed to update retrospective record: %w", err)
		}
		response = AddRetroItemResponse{
			ID:          item.ID,
			TimeCreated: item.Created,
		}
		return err
	})
	//panic("not implemented yet")
	return
}

func addRetroItemToSpaceRetro(ctx context.Context, userCtx facade.UserContext, request AddRetroItemRequest) (response AddRetroItemResponse, err error) {
	uid := userCtx.GetUserID()
	retrospectiveKey := getSpaceRetroDocKey(request.SpaceID, request.MeetingID)

	user := dbo4userus.NewUserEntry(uid)

	if err = dal4userus2.GetUser(ctx, nil, user); err != nil {
		err = fmt.Errorf("failed to get user record: %w", err)
		return
	}

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, transaction dal.ReadwriteTransaction) error {
		now := time.Now()
		retrospective := new(dbo4retrospectus.Retrospective)
		retrospectiveRecord := dal.NewRecordWithData(retrospectiveKey, retrospective)
		var isNew bool

		if err = txGetRetrospective(ctx, transaction, retrospectiveRecord); err != nil {
			return err
		} else if !retrospectiveRecord.Exists() {
			isNew = true
			response.TimeCreated = now
			retrospective = new(dbo4retrospectus.Retrospective)
			retrospective.TimeLastAction = &response.TimeCreated
		}

		if retrospective.Items == nil {
			retrospective.Items = make([]*dbo4retrospectus.RetroItem, 0, 1)
		}

		item := dbo4retrospectus.RetroItem{
			Title:   request.Title,
			Type:    request.Type,
			Created: now,
		}

		// adds item to retrospective
		{
			if user.Data.Names.FullName == "" {
				return fmt.Errorf("user[%v].Title is empty: %+v", uid, user)
			}

			if request.MeetingID != UpcomingRetrospectiveID {
				item.By = &dbo4retrospectus.RetroUser{
					UserID: uid,
					Title:  user.Data.Names.FullName,
				}
				if item.By.Title == "" {
					item.By.Title = user.Data.Names.GetFullName()
				}
			}
			retrospective.Items = addItemWithUniqueID(&item, retrospective.Items)
		}

		if err = retrospective.Validate(); err != nil {
			return err
		}

		if isNew {
			if err = txCreateRetrospective(ctx, transaction, retrospectiveKey, retrospective); err != nil {
				return fmt.Errorf("failed to create retrospective record: %w", err)
			}
		} else {
			updates := []dal.Update{
				{
					Field: "lastAction",
					Value: response.TimeCreated,
				},
				{
					Field: "items",
					Value: dal.ArrayUnion(item),
				},
				{
					Field: fmt.Sprintf("countsByMemberAndType.%v.%v", uid, request.Type),
					Value: dal.Increment(1),
				},
			}
			if err = txUpdate(ctx, transaction, retrospectiveKey, updates); err != nil {
				return fmt.Errorf("failed to update retrospective record: %w", err)
			}
		}
		return err
	})
	response.ID = retrospectiveKey.ID.(string)
	return
}
