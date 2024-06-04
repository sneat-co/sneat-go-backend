package facade4retrospectus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/dbo4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
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
func AddRetroItem(ctx context.Context, userContext facade.User, request AddRetroItemRequest) (response AddRetroItemResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	if request.MeetingID == UpcomingRetrospectiveID {
		response, err = addRetroItemToUserRetro(ctx, userContext, request)
		if err != nil {
			err = fmt.Errorf("failed to add item to future retrospective: %w", err)
		}
	} else {
		response, err = addRetroItemToTeamRetro(ctx, userContext, request)
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

func addRetroItemToUserRetro(ctx context.Context, userContext facade.User, request AddRetroItemRequest) (response AddRetroItemResponse, err error) {
	uid := userContext.GetID()

	user := new(dbo4userus.UserDbo)
	userKey := dal.NewKeyWithID(dbo4userus.UsersCollection, uid)
	userRecord := dal.NewRecordWithData(userKey, user)

	db := facade.GetDatabase(ctx)
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, transaction dal.ReadwriteTransaction) error {
		started := time.Now()

		if err = facade4userus.TxGetUserByID(ctx, transaction, userRecord); err != nil {
			return err
		}

		userTeamInfo := user.GetUserTeamInfoByID(request.TeamID)
		if userTeamInfo == nil {
			return validation.NewErrBadRequestFieldValue("team", fmt.Sprintf("user does not belong to this team %v, uid=%v", request.TeamID, uid))
		}

		//if userTeamInfo.RetroItems == nil {
		//	userTeamInfo.RetroItems = make(dbretro.RetroItemsByType)
		//}
		//
		//items, existingType := userTeamInfo.RetroItems[request.Role]

		//if !existingType {
		//	items = make([]*dbretro.RetroItem, 0, 1)
		//}
		item := dbo4retrospectus.RetroItem{
			Title:   request.Title,
			Created: started,
		}
		//userTeamInfo.RetroItems[request.Role] = addItemWithUniqueID(&item, items)

		if err = user.Validate(); err != nil {
			return err
		}

		//if err := updateTeamWithUpcomingRetroUserCounts(ctx, transaction, started, uid, request.TeamID, userTeamInfo.RetroItems); err != nil {
		//	return fmt.Errorf("failed to update team record: %w", err)
		//}

		updates := []dal.Update{
			{
				Field: fmt.Sprintf("api4meetingus.%v.retroItems.%v", request.TeamID, request.Type),
				Value: dal.ArrayUnion(item),
			},
		}
		if err = facade4userus.TxUpdateUser(ctx, transaction, started, userKey, updates); err != nil {
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

func addRetroItemToTeamRetro(ctx context.Context, userContext facade.User, request AddRetroItemRequest) (response AddRetroItemResponse, err error) {
	uid := userContext.GetID()
	retrospectiveKey := getTeamRetroDocKey(request.TeamID, request.MeetingID)

	user := new(dbo4userus.UserDbo)
	userKey := dal.NewKeyWithID(dbo4userus.UsersCollection, uid)
	userRecord := dal.NewRecordWithData(userKey, user)

	db := facade.GetDatabase(ctx)
	if err = facade4userus.GetUserByID(ctx, db, userRecord); err != nil {
		err = fmt.Errorf("failed to get user record: %w", err)
		return
	}

	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, transaction dal.ReadwriteTransaction) error {
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
			if user.Names.FullName == "" {
				return fmt.Errorf("user[%v].Title is empty: %+v", uid, user)
			}

			if request.MeetingID != UpcomingRetrospectiveID {
				item.By = &dbo4retrospectus.RetroUser{
					UserID: uid,
					Title:  user.Names.FullName,
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
