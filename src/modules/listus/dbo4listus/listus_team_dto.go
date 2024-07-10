package dbo4listus

import (
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
)

type ListusSpaceDbo struct {
	with.CreatedFields
	Lists map[string]*ListBrief `json:"lists,omitempty" firestore:"lists,omitempty"`
	//ListGroups []*ListGroup          `json:"listGroups,omitempty" firestore:"listGroups,omitempty"`
}

func (v ListusSpaceDbo) Validate() error {
	if err := validateListBriefs(v.Lists); err != nil {
		return validation.NewErrBadRecordFieldValue("lists", err.Error())
	}
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	return nil
}
