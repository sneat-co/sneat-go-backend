package models4debtus

func AddOrUpdateDebtusContact(debtusSpace DebtusSpaceEntry, c DebtusSpaceContactEntry) (debtusContactBrief *DebtusContactBrief, changed bool) {
	if c.Data == nil {
		panic("c.DebtusSpaceContactDbo == nil")
	}
	debtusContactBrief = NewDebtusContactJson(c.Data.Status, c.Data.Balanced)
	debtusContactBrief.Transfers = c.Data.GetTransfersInfo()
	found := false
	for c1id, c1 := range debtusSpace.Data.Contacts {
		if c1id == c.ID {
			found = true
			if !c1.Equal(debtusContactBrief) {
				debtusSpace.Data.Contacts[c.ID] = debtusContactBrief
				changed = true
			}
			break
		}
	}
	if !found {
		debtusSpace.Data.Contacts[c.ID] = debtusContactBrief
		changed = true
	}
	return
}
