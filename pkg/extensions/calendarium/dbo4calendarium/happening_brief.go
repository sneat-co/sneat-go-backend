package dbo4calendarium

import "github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"

// HappeningBrief hold data that stored both in entity record and in a brief.
type HappeningBrief struct {
	HappeningBase
	dbo4linkage.WithRelated
}

func (v *HappeningBrief) Validate() error {
	if err := v.HappeningBase.Validate(); err != nil {
		return err
	}
	if err := v.WithRelated.Validate(); err != nil {
		return err
	}
	return nil
}
