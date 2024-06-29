package gaedal

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/gotwilio"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
)

type TwilioDalGae struct {
}

func NewTwilioDalGae() TwilioDalGae {
	return TwilioDalGae{}
}

func (TwilioDalGae) GetLastTwilioSmsesForUser(c context.Context, tx dal.ReadSession, userID string, to string, limit int) (result []models.TwilioSms, err error) {
	q := dal.From(models.TwilioSmsKind).
		WhereField("UserID", dal.Equal, userID).
		OrderBy(dal.DescendingField("DtCreated"))

	if to != "" {
		q = q.WhereField("To", dal.Equal, to)
	}
	query := q.Limit(limit).SelectInto(models.NewTwilioSmsRecord)
	var records []dal.Record
	if records, err = tx.QueryAllRecords(c, query); err != nil {
		return
	}
	result = models.NewTwilioSmsFromRecords(records)
	return
}

func (TwilioDalGae) SaveTwilioSms(
	c context.Context,
	smsResponse *gotwilio.SmsResponse,
	transfer models.TransferEntry,
	phoneContact models.PhoneContact,
	userID string,
	tgChatID int64,
	smsStatusMessageID int,
) (twilioSms models.TwilioSms, err error) {
	var twilioSmsEntity models.TwilioSmsData
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	if err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		user := models.NewAppUser(userID, nil)
		twilioSms = models.NewTwilioSms(smsResponse.Sid, nil)
		counterparty := models.NewDebtusContact(transfer.Data.Counterparty().ContactID, nil)
		if err := tx.GetMulti(tc, []dal.Record{user.Record, twilioSms.Record, transfer.Record, counterparty.Record}); err != nil {
			var multiError appengine.MultiError
			if errors.As(err, &multiError) {
				if errors.Is(multiError[1], dal.ErrNoMoreRecords) {
					twilioSmsEntity = models.NewTwilioSmsFromSmsResponse(userID, smsResponse)
					twilioSmsEntity.CreatorTgChatID = tgChatID
					twilioSmsEntity.CreatorTgSmsStatusMessageID = smsStatusMessageID

					user.Data.SmsCount += 1
					transfer.Data.SmsCount += 1

					user.Data.SmsCost += float64(twilioSmsEntity.Price)
					transfer.Data.SmsCost += float64(twilioSmsEntity.Price)

					recordsToPut := []dal.Record{
						user.Record,
						twilioSms.Record,
						transfer.Record,
					}
					if counterparty.Data.PhoneContact.PhoneNumber != phoneContact.PhoneNumber {
						counterparty.Data.PhoneContact = phoneContact
						recordsToPut = append(recordsToPut, counterparty.Record)
					}
					if err = tx.SetMulti(tc, recordsToPut); err != nil {
						logus.Errorf(c, "Failed to save Twilio SMS")
						return err
					}
					return err
				} else if multiError[1] == nil {
					logus.Warningf(c, "Twillio SMS already saved to DB (1)")
				}
			}
		} else {
			logus.Warningf(c, "Twillio SMS already saved to DB (2)")
		}
		return nil
	}); err != nil {
		err = fmt.Errorf("failed to save Twilio response to DB: %w", err)
		return
	}
	return
}
