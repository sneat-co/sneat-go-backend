package dbo4userus

import (
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/strongo/strongoapp/person"
)

var _ botsfwmodels.AppUserAdapter = (*userBotsFwAdapter)(nil)

type userBotsFwAdapter struct {
	*UserDbo
}

func (v *userBotsFwAdapter) SetPreferredLocale(locale string) (err error) {
	_, err = v.WithPreferredLocale.SetPreferredLocale(locale)
	return err
}

func (v *UserDbo) BotsFwAdapter() botsfwmodels.AppUserAdapter {
	return &userBotsFwAdapter{UserDbo: v}
}

func (u *userBotsFwAdapter) SetNames(firstName, lastName, fullName string) error {
	if firstName == "" && lastName == "" && fullName == "" {
		return nil
	}
	if u.Names == nil {
		u.Names = new(person.NameFields)
	}
	if firstName != "" && u.Names.FirstName == "" {
		u.Names.FirstName = firstName
	}
	if lastName != "" && u.Names.LastName == "" {
		u.Names.LastName = lastName
	}
	if fullName != "" && u.Names.FullName == "" {
		u.Names.FullName = fullName
	}
	return nil
}
