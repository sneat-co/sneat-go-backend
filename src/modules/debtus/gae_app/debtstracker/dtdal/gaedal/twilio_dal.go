package gaedal

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/gotwilio"
	"github.com/strongo/logus"
	"google.golang.org/appengine/v2"
	"strconv"
)

type TwilioDalGae struct {
}

func NewTwilioDalGae() TwilioDalGae {
	return TwilioDalGae{}
}

func (TwilioDalGae) GetLastTwilioSmsesForUser(ctx context.Context, tx dal.ReadSession, userID string, to string, limit int) (result []models4debtus.TwilioSms, err error) {
	q := dal.From(models4debtus.TwilioSmsKind).
		WhereField("UserID", dal.Equal, userID).
		OrderBy(dal.DescendingField("DtCreated"))

	if to != "" {
		q = q.WhereField("To", dal.Equal, to)
	}
	query := q.Limit(limit).SelectInto(models4debtus.NewTwilioSmsRecord)
	var records []dal.Record
	if records, err = tx.QueryAllRecords(ctx, query); err != nil {
		return
	}
	result = models4debtus.NewTwilioSmsFromRecords(records)
	return
}

func (TwilioDalGae) SaveTwilioSms(
	ctx context.Context,
	smsResponse *gotwilio.SmsResponse,
	transfer models4debtus.TransferEntry,
	phoneContact dto4contactus.PhoneContact,
	userID string,
	tgChatID int64,
	smsStatusMessageID int,
) (twilioSms models4debtus.TwilioSms, err error) {
	var twilioSmsEntity models4debtus.TwilioSmsData
	if err = facade.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
		user := dbo4userus.NewUserEntry(userID)
		twilioSms = models4debtus.NewTwilioSms(smsResponse.Sid, nil)
		counterparty := transfer.Data.Counterparty()
		//counterpartyDebtusContact := models4debtus.NewDebtusSpaceContactEntry(counterparty.SpaceID, counterparty.ContactID, nil)
		counterpartyContact := dal4contactus.NewContactEntry(counterparty.SpaceID, counterparty.ContactID)
		if err := tx.GetMulti(tctx, []dal.Record{user.Record, twilioSms.Record, transfer.Record, counterpartyContact.Record}); err != nil {
			var multiError appengine.MultiError
			if errors.As(err, &multiError) {
				if errors.Is(multiError[1], dal.ErrNoMoreRecords) {
					twilioSmsEntity = models4debtus.NewTwilioSmsFromSmsResponse(userID, smsResponse)
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
					var phoneFound bool
					for _, counterpartyPhone := range counterpartyContact.Data.Phones {
						if phoneFound = counterpartyPhone.Number == strconv.FormatInt(phoneContact.PhoneNumber, 10); phoneFound {
							break
						}
					}
					if !phoneFound {
						counterpartyContact.Data.Phones = append(counterpartyContact.Data.Phones, dbmodels.PersonPhone{
							Number:   strconv.FormatInt(phoneContact.PhoneNumber, 10),
							Verified: false,
						})
						recordsToPut = append(recordsToPut, counterpartyContact.Record)
					}
					if err = tx.SetMulti(tctx, recordsToPut); err != nil {
						logus.Errorf(ctx, "Failed to save Twilio SMS")
						return err
					}
					return err
				} else if multiError[1] == nil {
					logus.Warningf(ctx, "Twillio SMS already saved to DB (1)")
				}
			}
		} else {
			logus.Warningf(ctx, "Twillio SMS already saved to DB (2)")
		}
		return nil
	}); err != nil {
		err = fmt.Errorf("failed to save Twilio response to DB: %w", err)
		return
	}
	return
}
