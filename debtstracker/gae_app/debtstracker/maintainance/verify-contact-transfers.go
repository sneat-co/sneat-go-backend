package maintainance

//type verifyContactTransfers struct {
//	contactsAsyncJob
//}
//
//func (m *verifyContactTransfers) Next(c context.Context, counters mapper.Counters, key *datastore.Key) error {
//	return m.startContactWorker(c, counters, key, m.processContact)
//}
//
//func (m *verifyContactTransfers) processContact(c context.Context, counters *asyncCounters, contact models.Contact) (err error) {
//	//log.Debugf(c, "processContact(contact.ID=%v)", contact.ID)
//	buf := new(bytes.Buffer)
//	now := time.Now()
//	hasError := false
//	var (
//		user          models.AppUser
//		warningsCount int
//		transfers     []models.Transfer
//	)
//	contactBalance := contact.Data.Balance()
//
//	defer func() {
//		if hasError || warningsCount > 0 {
//			var logFunc log.Func
//			if hasError {
//				logFunc = log.Errorf
//			} else {
//				logFunc = log.Warningf
//			}
//			var userName, contactName string
//			if user.Data != nil {
//				userName = user.Data.FullName()
//			}
//			if contact.Data == nil {
//				contactName = contact.Data.FullName()
//			}
//			logFunc(c,
//				fmt.Sprintf(
//					"Contact(id=%v, name=%v): has %v warning, %v transfers\n"+
//						"\tcontact.Balance: %v\n"+
//						"\tUser(id=%v, name=%v)",
//					contact.ID,
//					contactName,
//					warningsCount,
//					len(transfers),
//					litter.Sdump(contactBalance),
//					contact.Data.UserID,
//					userName,
//				)+buf.String(),
//			)
//		}
//	}()
//
//	//TODO: Load outstanding transfer just for the specific contact & specific direction
//	q := dal.From(models.TransferKind).
//		WhereField("BothCounterpartyIDs", dal.Equal, contact.ID).
//		OrderBy(dal.DescendingField("DtCreated")).
//		SelectInto(func() dal.Record {
//			return models.NewTransferWithIncompleteKey(nil).Record
//		})
//
//	var db dal.DB
//	if db, err = facade.GetDatabase(c); err != nil {
//		return err
//	}
//	var transferRecords []dal.Record
//	if transferRecords, err = db.QueryAllRecords(c, q); err != nil {
//		return err
//	}
//	transfers = make([]models.Transfer, len(transferRecords))
//	for i, r := range transferRecords {
//		transfers[i] = models.NewTransfer(r.Key().ID.(int), r.Data().(*models.TransferData))
//	}
//	models.ReverseTransfers(transfers)
//
//	transfersByID := make(map[int]models.Transfer, len(transfers))
//
//	if len(transfers) != contact.Data.CountOfTransfers {
//		_, _ = fmt.Fprintf(buf, "\tlen(transfers) != contact.CountOfTransfers: %v != %v\n", len(transfers), contact.Data.CountOfTransfers)
//		warningsCount += 1
//	}
//
//	if contact.Data.CounterpartyCounterpartyID != 0 || contact.Data.CounterpartyUserID != 0 { // Fixing names
//		for _, transfer := range transfers {
//			originalTransfer := transfer
//			*originalTransfer.Data = *transfer.Data
//			changed := false
//			self := transfer.Data.UserInfoByUserID(contact.Data.UserID)
//			counterparty := transfer.Data.CounterpartyInfoByUserID(contact.Data.UserID)
//
//			if contact.Data.CounterpartyUserID != 0 && counterparty.UserID == 0 {
//				counterparty.UserID = contact.Data.CounterpartyUserID
//				changed = true
//			}
//			if counterparty.UserName == "" && counterparty.UserID != 0 {
//				if user, err := facade.User.GetUserByID(c, db, counterparty.UserID); err != nil {
//					log.Errorf(c, err.Error())
//					return err
//				} else {
//					counterparty.UserName = user.Data.FullName()
//				}
//				changed = true
//			}
//
//			if contact.Data.CounterpartyCounterpartyID != 0 && self.ContactID == 0 {
//				self.ContactID = contact.Data.CounterpartyCounterpartyID
//				changed = true
//			}
//
//			if self.ContactID != 0 && self.ContactName == "" {
//				if counterpartyContact, err := facade.GetContactByID(c, nil, self.ContactID); err != nil {
//					log.Errorf(c, err.Error())
//					return err
//				} else {
//					self.ContactName = counterpartyContact.Data.FullName()
//				}
//				changed = true
//			}
//
//			if self.UserID != 0 && self.UserName == "" {
//				if user, err := facade.User.GetUserByID(c, nil, self.UserID); err != nil {
//					log.Errorf(c, err.Error())
//					return err
//				} else {
//					self.UserName = user.Data.FullName()
//				}
//				changed = true
//			}
//
//			if changed {
//				log.Warningf(c, "Fixing contact details for transfer %v: From:%v, To: %v\n\noriginal: %v\n\n new: %v", transfer.ID, litter.Sdump(transfer.Data.From()), litter.Sdump(transfer.Data.To()), litter.Sdump(originalTransfer), litter.Sdump(transfer))
//				err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
//					return facade.Transfers.SaveTransfer(c, tx, transfer)
//				})
//				if err != nil {
//					log.Errorf(c, fmt.Errorf("failed to save transfer: %w", err).Error())
//					return
//				}
//			}
//		}
//	}
//
//	loggedTransfers := make(map[int]bool, len(transfers))
//
//	logTransfer := func(transfer models.Transfer, padding int) {
//		if _, ok := loggedTransfers[transfer.ID]; !ok {
//			p := strings.Repeat("\t", padding)
//			fmt.Fprintf(buf, p+"Transfer: %v\n", transfer.ID)
//			fmt.Fprintf(buf, p+"\tCreated: %v\n", transfer.Data.DtCreated)
//			fmt.Fprintf(buf, p+"\tFrom(): userID=%v, contactID=%v\n", transfer.Data.From().UserID, transfer.Data.From().ContactID)
//			fmt.Fprintf(buf, p+"\t  To(): userID=%v, contactID=%v\n", transfer.Data.To().UserID, transfer.Data.To().ContactID)
//			fmt.Fprintf(buf, p+"\tAmount: %v\n", transfer.Data.GetAmount())
//			fmt.Fprintf(buf, p+"\tReturned: %v\n", transfer.Data.AmountInCentsReturned)
//			fmt.Fprintf(buf, p+"\tOutstanding: %v\n", transfer.Data.GetOutstandingValue(time.Now()))
//			fmt.Fprintf(buf, p+"\tIsReturn: %v\n", transfer.Data.IsReturn)
//			fmt.Fprintf(buf, p+"\tReturnToTransferIDs: %v\n", transfer.Data.ReturnToTransferIDs)
//			if transfer.Data.HasInterest() {
//				fmt.Fprintf(buf, p+"\tInterest: %v @ %v%%/%v_days, min=%v, grace=%v",
//					transfer.Data.InterestType, transfer.Data.InterestPercent, transfer.Data.InterestPeriod,
//					transfer.Data.InterestMinimumPeriod, transfer.Data.InterestGracePeriod)
//			}
//			loggedTransfers[transfer.ID] = true
//		}
//	}
//
//	transfersBalance := m.getTransfersBalance(transfers, contact.ID)
//
//	verifyReturnIDs := func() (valid bool) {
//		valid = true
//		counters.Lock()
//		for _, transfer := range transfersByID {
//			for i, returnToTransferID := range transfer.Data.ReturnToTransferIDs {
//				if _, ok := transfersByID[returnToTransferID]; ok {
//					counters.Increment("good_ReturnToTransferID", 1)
//				} else {
//					valid = false
//					fmt.Fprintf(buf, "\t\tReturnToTransferIDs[%d]: %v\n", i, returnToTransferID)
//					counters.Increment("wrong_ReturnToTransferID", 1)
//					warningsCount += 1
//				}
//			}
//		}
//		counters.Unlock()
//		return
//	}
//
//	var lastTransfer models.Transfer
//
//	if len(transfers) > 0 {
//		lastTransfer = transfers[len(transfers)-1]
//	}
//
//	var needsFixingContactOrUser bool
//
//	if valid, warnsCount := m.assertTotals(buf, counters, contact, transfersBalance); !valid {
//		needsFixingContactOrUser = true
//		warningsCount += warnsCount
//	} else {
//		warningsCount += warnsCount
//	}
//
//	outstandingIsValid, outstandingWarningsCount := m.verifyOutstanding(c, 1, buf, contactBalance, transfersBalance)
//	warningsCount += outstandingWarningsCount
//	if !outstandingIsValid {
//		//rollingBalance := make(money.Balance, len(transfersBalance)+1)
//		transfersByCurrency, transfersToSave := m.fixTransfers(c, now, buf, contact, transfers)
//
//		for currency, currencyTransfers := range transfersByCurrency {
//			_, _ = fmt.Fprintf(buf, "\tcurrency: %v - %d transfers\n", currency, len(currencyTransfers))
//		}
//
//		if valid, _ := m.verifyOutstanding(c, 2, buf, contactBalance, transfersBalance); !valid {
//			_, _ = fmt.Fprint(buf, "Outstandings are invalid after fix!\n")
//			needsFixingContactOrUser = true
//		} else if valid, _ = m.assertTotals(buf, counters, contact, transfersBalance); !valid {
//			_, _ = fmt.Fprint(buf, "Totals are invalid after fix!\n")
//		} else if valid = verifyReturnIDs(); !valid {
//			_, _ = fmt.Fprint(buf, "Return IDs are invalid after fix!\n")
//		} else {
//			_, _ = fmt.Fprintf(buf, "%v transfers to save!\n", len(transfersToSave))
//			recordsToSave := make([]dal.Record, 0, len(transfersToSave))
//			for id, transfer := range transfersToSave {
//				recordsToSave = append(recordsToSave, models.NewTransfer(id, transfer).Record)
//			}
//			//gaedb.LoggingEnabled = true
//			err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
//				return tx.SetMulti(c, recordsToSave)
//			})
//			if err != nil {
//				_, _ = fmt.Fprintf(buf, "ERROR: failed to save transfers: %v\n", err)
//				hasError = true
//				return
//			}
//			_, _ = fmt.Fprintf(buf, "SAVED %v transfers!\n", len(recordsToSave))
//		}
//	}
//
//	if outstandingIsValid {
//		if user, err = facade.User.GetUserByID(c, nil, contact.Data.UserID); err != nil {
//			log.Errorf(c, fmt.Errorf("contact{ID=%v}: user not found by ID: %w", contact.ID, err).Error())
//			return
//		}
//	}
//
//	if !outstandingIsValid || !contactBalance.Equal(user.Data.ContactByID(contact.ID).Balance()) || !contactBalance.Equal(transfersBalance) {
//		needsFixingContactOrUser = true
//	}
//
//	if !needsFixingContactOrUser && contact.Data.CounterpartyCounterpartyID != 0 {
//		var counterpartyContact models.Contact
//		if counterpartyContact, err = facade.GetContactByID(c, nil, contact.Data.CounterpartyCounterpartyID); err != nil {
//			return
//		}
//		fmt.Fprintf(buf, "contact.Balance(): %v\n", contact.Data.Balance())
//		fmt.Fprintf(buf, "counterpartyContact.Balance(): %v\n", contact.Data.Balance())
//		if !counterpartyContact.Data.GetTransfersInfo().Equal(contact.Data.GetTransfersInfo()) || !counterpartyContact.Data.Balance().Equal(transfersBalance.Reversed()) {
//			needsFixingContactOrUser = true
//		}
//	} else {
//		fmt.Fprintf(buf, "needsFixingContactOrUser: %v, contact.CounterpartyCounterpartyID: %v", needsFixingContactOrUser, contact.Data.CounterpartyCounterpartyID)
//	}
//
//	if needsFixingContactOrUser {
//		for _, transfer := range transfers {
//			logTransfer(transfer, 1)
//		}
//		if contact, user, err = m.fixContactAndUser(c, buf, counters, contact.ID, transfersBalance, len(transfers), lastTransfer); err != nil {
//			return
//		}
//	}
//
//	if warningsCount == 0 {
//		counters.Increment("good_contacts", 1)
//		//log.Infof(c, contactPrefix + "is OK, %v transfers", len(transfers))
//	} else {
//		counters.Lock()
//		counters.Increment("bad_contacts", 1)
//		counters.Increment("warnings", int64(warningsCount))
//		counters.Unlock()
//
//		_ = contact.Data.Balance()
//
//		//if len(contactBalance) == 0 {
//		//	contactBalance = nil
//		//}
//	}
//	return nil
//}
//
//func (m *verifyContactTransfers) assertTotals(buf *bytes.Buffer, counters *asyncCounters, contact models.Contact, transfersBalance money.Balance) (valid bool, warningsCount int) {
//	valid = true
//	contactBalance := contact.Data.Balance()
//	for currency, transfersTotal := range transfersBalance {
//		if contactTotal := contactBalance[currency]; contactTotal != transfersTotal {
//			valid = false
//			fmt.Fprintf(buf, "currency %v: transfersTotal != contactTotal: %v != %v\n", currency, transfersTotal, contactTotal)
//			warningsCount += 1
//		}
//		delete(contactBalance, currency)
//	}
//	for currency, contactTotal := range contactBalance {
//		if contactTotal == 0 {
//			counters.Increment("zero_balance", 1)
//			fmt.Fprintf(buf, "\t0 value for currency %v\n", currency)
//			warningsCount += 1
//		} else {
//			counters.Increment("no_transfers_for_non_zero_balance", 1)
//			fmt.Fprintf(buf, "\tno transfers found for %v=%v\n", currency, contactTotal)
//			warningsCount += 1
//		}
//	}
//	return
//}
//
//func (m *verifyContactTransfers) fixContactAndUser(c context.Context, buf *bytes.Buffer, counters *asyncCounters, contactID int64, transfersBalance money.Balance, transfersCount int, lastTransfer models.Transfer) (contact models.Contact, user models.AppUser, err error) {
//	var db dal.DB
//	if db, err = facade.GetDatabase(c); err != nil {
//		return
//	}
//	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if contact, user, err = m.fixContactAndUserWithinTransaction(c, tx, buf, counters, contactID, transfersBalance, transfersCount, lastTransfer); err != nil {
//			return
//		}
//		if contact.Data.CounterpartyCounterpartyID != 0 {
//			if _, _, err = m.fixContactAndUserWithinTransaction(c, tx, buf, counters, contact.Data.CounterpartyCounterpartyID, transfersBalance.Reversed(), transfersCount, lastTransfer); err != nil {
//				return
//			}
//		}
//		return
//	}); err != nil {
//		return
//	}
//	return
//}
//
//func (m *verifyContactTransfers) fixContactAndUserWithinTransaction(c context.Context, tx dal.ReadwriteTransaction, buf *bytes.Buffer, counters *asyncCounters, contactID int64, transfersBalance money.Balance, transfersCount int, lastTransfer models.Transfer) (contact models.Contact, user models.AppUser, err error) {
//	fmt.Fprintf(buf, "Fixing contact %v...\n", contactID)
//	if contact, err = facade.GetContactByID(c, tx, contactID); err != nil {
//		return
//	}
//	changed := false
//	if lastTransfer.Data != nil && lastTransfer.ID != 0 {
//		if contact.Data.LastTransferAt.Before(lastTransfer.Data.DtCreated) {
//			_, _ = fmt.Fprintf(buf, "\tcontact.LastTransferAt changed from %v to %v\n", contact.Data.LastTransferID, lastTransfer.Data.DtCreated)
//			contact.Data.LastTransferAt = lastTransfer.Data.DtCreated
//
//			if contact.Data.LastTransferID != int64(lastTransfer.ID) {
//				_, _ = fmt.Fprintf(buf, "\tcontact.LastTransferID changed from %v to %v\n", contact.Data.LastTransferID, lastTransfer.ID)
//				contact.Data.LastTransferID = int64(lastTransfer.ID)
//			}
//			changed = true
//		}
//	}
//	if contact.Data.CountOfTransfers < transfersCount {
//		_, _ = fmt.Fprintf(buf, "\tcontact.CountOfTransfers changed from %v to %v\n", contact.Data.CountOfTransfers, transfersCount)
//		contact.Data.CountOfTransfers = transfersCount
//		changed = true
//	}
//	if !contact.Data.Balance().Equal(transfersBalance) {
//		if err = contact.Data.SetBalance(transfersBalance); err != nil {
//			return
//		}
//		changed = true
//	}
//	if changed {
//		if err = facade.SaveContact(c, contact); err != nil {
//			return
//		}
//		//var user models.AppUser
//		if user, err = facade.User.GetUserByID(c, nil, contact.Data.UserID); err != nil {
//			return
//		}
//		userContacts := user.Data.Contacts()
//		userChanged := false
//		for i := range userContacts {
//			if userContacts[i].ID == contact.ID {
//				if !userContacts[i].Balance().Equal(transfersBalance) {
//					if err = userContacts[i].SetBalance(transfersBalance); err != nil {
//						return
//					}
//					user.Data.SetContacts(userContacts)
//					userChanged = true
//				}
//				userTransferInfo, contactTransferInfo := userContacts[i].Transfers, contact.Data.GetTransfersInfo()
//				if !userTransferInfo.Equal(contactTransferInfo) {
//					userContacts[i].Transfers = contactTransferInfo
//					userChanged = true
//				}
//				goto contactFound
//			}
//		}
//		// Contact not found
//		if _, changed := user.AddOrUpdateContact(contact); changed {
//			userChanged = true
//		}
//	contactFound:
//		userTotalBalance := user.Data.Balance()
//		if userContactsBalance := user.Data.TotalBalanceFromContacts(); !userContactsBalance.Equal(userTotalBalance) {
//			if err = user.Data.SetBalance(userContactsBalance); err != nil {
//				return
//			}
//			userChanged = true
//			fmt.Fprintf(buf, "user total balance update from contacts\nwas: %v\nnew: %v\n", userTotalBalance, userContactsBalance)
//		}
//		if userChanged {
//			if err = facade.User.SaveUser(c, tx, user); err != nil {
//				return
//			}
//		}
//	}
//	return
//}
//
//func (verifyContactTransfers) getTransfersBalance(transfers []models.Transfer, contactID int64) (totalBalance money.Balance) {
//	totalBalance = make(money.Balance)
//	for _, transfer := range transfers {
//		direction := transfer.Data.DirectionForContact(contactID)
//		switch direction {
//		case models.TransferDirectionUser2Counterparty:
//			totalBalance[transfer.Data.Currency] += transfer.Data.AmountInCents
//		case models.TransferDirectionCounterparty2User:
//			totalBalance[transfer.Data.Currency] -= transfer.Data.AmountInCents
//		default:
//			panic(fmt.Sprintf("transfer.DirectionForContact(%v): %v", contactID, direction))
//		}
//	}
//	for c, v := range totalBalance {
//		if v == 0 {
//			delete(totalBalance, c)
//		}
//	}
//	return
//}
//
//func (m *verifyContactTransfers) verifyOutstanding(c context.Context, iteration int, buf *bytes.Buffer, contactBalance money.Balance, transfersBalance money.Balance) (valid bool, warningsCount int) {
//	fmt.Fprintf(buf, "\tverifyOutstanding(iteration=%v):\n", iteration)
//	valid = true
//
//	for currency, contactTotal := range contactBalance {
//		if transfersTotal := transfersBalance[currency]; transfersTotal == contactTotal {
//			fmt.Fprintf(buf, "\t\tcurrency %v: contactBalance == transfersTotal: %v\n", currency, contactTotal)
//		} else {
//			valid = false
//			fmt.Fprintf(buf, "\t\tcurrency %v: contactBalance != transfersTotal: %v != %v\n", currency, contactTotal, transfersTotal)
//			warningsCount += 1
//		}
//		//delete(transfersOutstanding, currency)
//	}
//	fmt.Fprintf(buf, "\tverifyOutstanding(iteration=%v) => valid=%v\n", iteration, valid)
//
//	return
//}
//
//func (m *verifyContactTransfers) fixTransfers(c context.Context, now time.Time, buf *bytes.Buffer, contact models.Contact, transfers []models.Transfer) (
//	transfersByCurrency map[money.CurrencyCode][]models.Transfer,
//	transfersToSave map[int]*models.TransferData,
//) {
//	fmt.Fprintln(buf, "fixTransfers()")
//
//	transfersByCurrency = make(map[money.CurrencyCode][]models.Transfer)
//
//	transfersToSave = make(map[int]*models.TransferData)
//
//	for _, transfer := range transfers {
//		if transfer.Data.AmountInCentsReturned != 0 {
//			transfer.Data.AmountInCentsReturned = 0
//			transfersToSave[transfer.ID] = transfer.Data
//		}
//		if len(transfer.Data.ReturnToTransferIDs) != 0 {
//			transfer.Data.ReturnToTransferIDs = []int{}
//			transfersToSave[transfer.ID] = transfer.Data
//		}
//		amountToAssign := transfer.Data.GetAmount().Value
//		for _, previousTransfer := range transfersByCurrency[transfer.Data.Currency] {
//			if previousTransfer.Data.IsOutstanding && previousTransfer.Data.IsReverseDirection(transfer.Data) {
//				// previousTransfer.ReturnTransferIDs = append(previousTransfer.ReturnTransferIDs, transfer.ID)"
//				transfer.Data.ReturnToTransferIDs = append(transfer.Data.ReturnToTransferIDs, previousTransfer.ID)
//				transfersToSave[previousTransfer.ID] = previousTransfer.Data
//				if previousTransferOutstandingValue := previousTransfer.Data.GetOutstandingValue(now); amountToAssign <= previousTransferOutstandingValue {
//					previousTransfer.Data.AmountInCentsReturned += amountToAssign
//					amountToAssign = 0
//					//break
//				} else /* previousTransfer.AmountInCentsOutstanding < amountToAssign */ {
//					amountToAssign -= previousTransferOutstandingValue
//					previousTransfer.Data.AmountInCentsReturned += previousTransferOutstandingValue
//					previousTransfer.Data.IsOutstanding = false
//				}
//				panic("not implemented")
//			}
//		}
//		transfer.Data.IsReturn = len(transfer.Data.ReturnToTransferIDs) > 0
//		if transfer.Data.IsOutstanding = amountToAssign != 0; transfer.Data.IsOutstanding {
//			transfer.Data.AmountInCentsReturned = transfer.Data.AmountInCents - amountToAssign
//			transfersToSave[transfer.ID] = transfer.Data
//		}
//		transfersByCurrency[transfer.Data.Currency] = append(transfersByCurrency[transfer.Data.Currency], transfer)
//	}
//	return
//}
