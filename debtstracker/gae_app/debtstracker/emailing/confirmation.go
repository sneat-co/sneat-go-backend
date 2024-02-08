package emailing

import (
	"fmt"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func CreateConfirmationEmailAndQueueForSending(c context.Context, user models.AppUser, userEmail models.UserEmail) error {
	emailEntity := &models.EmailData{
		From:    "Alex @ DebtsTracker.io <alex@debtstracker.io>",
		To:      userEmail.ID,
		Subject: "Please confirm your account at DebtsTracker.io",
		BodyText: fmt.Sprintf(`%v, we are thrilled to have you on board!

To keep your account secure please confirm your email by clicking this link:

  >> https://debtstracker.io/confirm?email=%v&pin=%v

If you have any questions or issue please drop me an email to alex@debtstracker.io
--
Alex
Creator of https://DebtsTracker.io

We are social:
  FB page - https://www.facebook.com/debtstracker
  Twitter - https://twitter.com/debtstracker
`, user.Data.FullName(), userEmail.ID, userEmail.Data.ConfirmationPin()),
	}
	_, err := CreateEmailRecordAndQueueForSending(c, emailEntity)
	return err
}
