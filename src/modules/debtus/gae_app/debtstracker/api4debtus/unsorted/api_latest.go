package unsorted

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-core-modules/auth/token4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"net/http"
)

func HandleAdminLatestUsers(ctx context.Context, w http.ResponseWriter, _ *http.Request, _ token4auth.AuthInfo) {
	common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, errors.New("not implemented yet"))
	//users, err := dtdal.Admin.LatestUsers(ctx)
	//if err != nil {
	//	api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	//	return
	//}
	//var buffer bytes.Buffer
	//buffer.WriteString("[")
	//lastIndex := len(users) - 1
	//var wg sync.WaitGroup
	//records := make([]*dto4debtus.Record, len(users))
	//
	//for i, user := range users {
	//	records[i] = &dto4debtus.Record{
	//		Id:                     user.ID,
	//		Name:                   user.Data.Names.GetFullName(),
	//		Transfers:              user.Data.CountOfTransfers,
	//		CountOfReceiptsCreated: user.Data.CountOfReceiptsCreated,
	//		TelegramUserIDs:        user.Data.GetTelegramUserIDs(),
	//		Balance:                user.Data.Balance,
	//	}
	//	userCounterpartiesIDs := user.Data.ContactIDs()
	//	if len(userCounterpartiesIDs) > 0 {
	//		wg.Add(1)
	//		go func(i int, userCounterpartiesIDs []string) {
	//			counterparties, err := facade4debtus.GetDebtusSpaceContactsByIDs(ctx, nil, userCounterpartiesIDs)
	//			if err != nil {
	//				logus.Errorf(ctx, fmt.Errorf("failed to get counterparties by ids=%+v: %w", userCounterpartiesIDs, err).Error())
	//				wg.Done()
	//				return
	//			}
	//			record := records[i]
	//			for j, counterparty := range counterparties {
	//				counterpartyDto := dto4debtus.CounterpartyDto{
	//					ID:     userCounterpartiesIDs[j],
	//					UserID: counterparty.Data.CounterpartyUserID,
	//					Name:   counterparty.Data.FullName(),
	//				}
	//				if counterparty.Data.BalanceJson != "" {
	//					balance := json.RawMessage(counterparty.Data.BalanceJson)
	//					counterpartyDto.Balance = &balance
	//				}
	//				record.Counterparties = append(record.Counterparties, counterpartyDto)
	//			}
	//			logus.Debugf(ctx, "Contacts goroutine completed.")
	//			wg.Done()
	//		}(i, userCounterpartiesIDs)
	//	}
	//	if user.Data.InvitedByUserID != "" {
	//		wg.Add(1)
	//		go func(i int, userID string) {
	//			inviter, err := dal4userus.GetUserByID(ctx, nil, userID)
	//			if err != nil {
	//				logus.Errorf(ctx, fmt.Errorf("failed to get user by id=%v: %w", userID, err).Error())
	//				return
	//			}
	//			records[i].InvitedByUser = &struct {
	//				Id   string
	//				Name string
	//			}{
	//				userID,
	//				inviter.Data.FullName(),
	//			}
	//			logus.Debugf(ctx, "User goroutine completed.")
	//			wg.Done()
	//		}(i, user.Data.InvitedByUserID)
	//	}
	//}
	//
	//wg.Wait()
	//
	//for i, record := range records {
	//	if userBytes, err := json.Marshal(record); err != nil {
	//		logus.Errorf(ctx, err.Error())
	//		w.WriteHeader(http.StatusInternalServerError)
	//		_, _ = w.Write([]byte(err.Error()))
	//		return
	//	} else {
	//		buffer.Write(userBytes)
	//	}
	//	if i < lastIndex {
	//		buffer.Write([]byte(","))
	//	}
	//}
	//
	//buffer.WriteString("]")
	//header := w.Header()
	//header.Add("Content-Type", "application/json")
	//header.Add("Access-Control-Allow-Origin", "*")
	//if _, err = w.Write(buffer.Bytes()); err != nil {
	//	logus.Errorf(ctx, err.Error())
	//	w.WriteHeader(http.StatusInternalServerError)
	//	_, _ = w.Write([]byte(err.Error()))
	//}
}
