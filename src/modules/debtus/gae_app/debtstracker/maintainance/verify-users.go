package maintainance

//type verifyUsers struct {
//	asyncMapper
//	entity *models.AppUserData
//}
//
//func (m *verifyUsers) Make() interface{} {
//	m.entity = new(models.AppUserData)
//	return m.entity
//}
//
//func (m *verifyUsers) Query(r *http.Request) (query *mapper.Query, err error) {
//	return applyIDAndUserFilters(r, "verifyUsers", models.AppUserKind, filterByIntID, "")
//}
//
//func (m *verifyUsers) Next(ctx context.Context, counters mapper.Counters, key *dal.Key) (err error) {
//	userEntity := *m.entity
//	user := models.NewAppUser(key.ContactID.(int64), &userEntity)
//	return m.startWorker(ctx, counters, func() Worker {
//		return func(counters *asyncCounters) error {
//			return m.processUser(c, user, counters)
//		}
//	})
//}
//
//func (m *verifyUsers) processUser(ctx context.Context, user models.AppUser, counters *asyncCounters) (err error) {
//	buf := new(bytes.Buffer)
//	if user, err = m.checkContactsExistsAndRecreateIfNeeded(ctx, buf, counters, user); err != nil {
//		return
//	}
//	if err = m.verifyUserBalanceAndContacts(ctx, buf, counters, user); err != nil {
//		return
//	}
//	if buf.Len() > 0 {
//		logus.Infof(ctx, buf.String())
//	}
//	return
//}
//
//func (m *verifyUsers) checkContactsExistsAndRecreateIfNeeded(ctx context.Context, buf *bytes.Buffer, counters *asyncCounters, user models.AppUser) (models.AppUser, error) {
//	userContacts := user.Data.Contacts()
//	userChanged := false
//	var err error
//	for i, userContact := range userContacts {
//		contactID := userContact.ContactID
//		var contact models.DebtusSpaceContactEntry
//		if contact, err = facade4debtus.GetContactByID(ctx, nil, contactID); err != nil {
//			if dal.IsNotFound(err) {
//				if err = m.createContact(ctx, buf, counters, user, userContact); err != nil {
//					logus.Errorf(c, "Failed to create contact %v", userContact.ContactID)
//					err = nil
//					continue
//				}
//			} else {
//				return user, err
//			}
//		}
//		if contact.Data.CounterpartyUserID != 0 && userContact.UserID != contact.Data.CounterpartyUserID {
//			if userContact.UserID == 0 {
//				userContact.UserID = contact.Data.CounterpartyUserID
//				userContacts[i] = userContact
//				userChanged = true
//			} else {
//				err = fmt.Errorf(
//					"data integrity issue for contact %v: userContact.UserID != contact.CounterpartyUserID: %v != %v",
//					contact.ContactID, userContact.UserID, contact.Data.CounterpartyUserID)
//				return user, err
//			}
//		}
//	}
//	if userChanged {
//		var db dal.DB
//		if db, err = facade4debtus.GetDatabase(ctx); err != nil {
//			return user, err
//		}
//		if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
//			if user, err = dal4userus.GetUserByID(ctx, tx, user.ContactID); err != nil {
//				return err
//			}
//			user.Data.SetContacts(userContacts)
//			if err = facade4debtus.User.SaveUserOBSOLETE(ctx, tx, user); err != nil {
//				return err
//			}
//			return nil
//		}); err != nil {
//			return user, err
//		}
//
//	}
//	return user, err
//}
//
//func (m *verifyUsers) createContact(ctx context.Context, buf *bytes.Buffer, counters *asyncCounters, user models.AppUser, userContact models.UserContactJson) (err error) {
//	var contact models.DebtusSpaceContactEntry
//	var db dal.DB
//	if db, err = facade4debtus.GetDatabase(c); err != nil {
//		return
//	}
//	if err = db.RunReadwriteTransaction(ctx, func(tctx context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if contact, err = facade4debtus.GetContactByID(tctx, nil, userContact.ContactID); err != nil {
//			if dal.IsNotFound(err) {
//				contact = models.NewContact(userContact.ContactID, &models.ContactData{
//					UserID:    user.ContactID,
//					DtCreated: time.Now(),
//					Status:    models.STATUS_ACTIVE,
//					ContactDetails: models.ContactDetails{
//						Nickname:       userContact.Name,
//						TelegramUserID: userContact.TgUserID,
//					},
//				})
//				if err = contact.Data.SetBalance(userContact.Balance()); err != nil {
//					return
//				}
//				if err = contact.Data.SetTransfersInfo(*contact.Data.GetTransfersInfo()); err != nil {
//					return
//				}
//				if err = facade4debtus.SaveContact(tctx, contact); err != nil {
//					return
//				}
//			}
//			return
//		}
//		return
//	}); err != nil {
//		return
//	} else {
//		logus.Warningf(ctx, "Recreated contact %v[%v] for user %v[%v]", contact.ContactID, contact.Data.FullName(), user.ContactID, user.Data.FullName())
//	}
//	return
//}
//
//func (m *verifyUsers) verifyUserBalanceAndContacts(ctx context.Context, buf *bytes.Buffer, counters *asyncCounters, user models.AppUser) (err error) {
//	if user.Data.BalanceCount > 0 {
//		balance := user.Data.Balance()
//		var fixedContactsBalances bool
//		if fixedContactsBalances, err = fixUserContactsBalances(m.entity); err != nil {
//			return err
//		} else if fixedContactsBalances || FixBalanceCurrencies(balance) {
//			var db dal.DB
//			if db, err = facade4debtus.GetDatabase(c); err != nil {
//				return err
//			}
//			if err = db.RunReadwriteTransaction(c, func(tctx context.Context, tx dal.ReadwriteTransaction) error {
//				if user, err = dal4userus.GetUserByID(tctx, tx, user.ContactID); err != nil {
//					return err
//				}
//				balance = m.entity.Balance()
//				if err != nil {
//					return err
//				}
//				changed := false
//				if FixBalanceCurrencies(balance) {
//					if err = m.entity.SetBalance(balance); err != nil {
//						return err
//					}
//					changed = true
//				}
//				if fixedContactsBalances, err = fixUserContactsBalances(m.entity); err != nil {
//					return err
//				} else if fixedContactsBalances {
//					changed = true
//				}
//				if changed {
//					if err = facade4debtus.User.SaveUserOBSOLETE(tctx, tx, user); err != nil {
//						return err
//					}
//					fmt.Fprintf(buf, "User fixed: %d ", user.ContactID)
//				}
//				return nil
//			}, nil); err != nil {
//				return err
//			}
//		}
//	}
//	return
//}
//
//func fixUserContactsBalances(u *models.AppUserData) (changed bool, err error) {
//	contacts := u.Contacts()
//	for i, contact := range contacts {
//		if balance := contact.Balance(); FixBalanceCurrencies(balance) {
//			balanceJsonBytes, err := ffjson.Marshal(balance)
//			if err != nil {
//				return changed, err
//			}
//			balanceJson := json.RawMessage(balanceJsonBytes)
//			contact.BalanceJson = &balanceJson
//			contacts[i] = contact
//			changed = true
//		}
//	}
//	if changed {
//		u.SetContacts(contacts)
//	}
//	return
//}
