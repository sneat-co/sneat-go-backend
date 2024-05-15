package facade4scrumus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/models4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/random"
	"github.com/strongo/validation"
	"strings"
)

// AddCommentRequest request
type AddCommentRequest struct {
	TaskRequest
	Message string `json:"message"`
}

// Validate validates request
func (v *AddCommentRequest) Validate() error {
	if strings.TrimSpace(v.Message) == "" {
		return validation.NewErrRecordIsMissingRequiredField("message")
	}
	return v.TaskRequest.Validate()
}

// AddComment adds comment
func AddComment(ctx context.Context, userContext facade.User, request AddCommentRequest) (comment *models4scrumus.Comment, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("facade4retrospectus bad request: %v", err)
		return
	}

	db := facade.GetDatabase(ctx)

	uid := userContext.GetID()

	userKey := models4userus.NewUserKey(uid)
	var user models4userus.UserDto
	userRecord := dal.NewRecordWithData(userKey, &user)
	if err = facade4userus.GetUserByID(ctx, db, userRecord); err != nil {
		return nil, err
	}

	err = runTaskWorker(ctx, userContext, request.TaskRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params taskWorkerParams) (err error) {
			if params.task == nil {
				return errors.New("task not found by ContactID: " + request.TaskRequest.Task)
			}
			comment = &models4scrumus.Comment{
				ID:      random.ID(1),
				Message: request.Message,
				By: &dbmodels.ByUser{
					UID:   uid,
					Title: user.Names.FullName,
				},
			}

		UniqueID:
			for _, c := range params.task.Comments {
				if c.ID == comment.ID {
					comment.ID = random.ID(len(comment.ID) + 1)
					goto UniqueID
				}
			}

			params.task.Comments = append(params.task.Comments, comment)
			return tx.Update(ctx, params.Meeting.Key, []dal.Update{
				{
					Field: fmt.Sprintf("statuses.%s.byType.%s", request.TaskRequest.ContactID, request.TaskRequest.Type),
					Value: params.tasks,
				},
			})
		})
	return comment, err
}
