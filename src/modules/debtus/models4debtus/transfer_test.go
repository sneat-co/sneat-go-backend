package models4debtus

//func TestTransfer_LoadSaver(t *testing.T) {
//	var saved []struct {
//		kind       string
//		properties []datastore.Property
//	}
//	//checkHasProperties = func(kind string, properties []datastore.Property) {
//	//	saved = append(saved, struct {
//	//		kind       string
//	//		properties []datastore.Property
//	//	}{kind, properties})
//	//	return
//	//}
//	currency := money.CurrencyCode_IRR
//	creator := TransferCounterpartyInfo{
//		UserID:      1,
//		ContactID:   2,
//		ContactName: "Test1",
//	}
//	counterparty := TransferCounterpartyInfo{
//		ContactName: "Creator 1",
//	}
//	transfer := NewTransferData(creator.UserID, false, money.NewAmount(currency, decimal.NewDecimal64p2FromFloat64(123.45)), &creator, &counterparty)
//	properties, err := transfer.Save()
//	if err != nil {
//		t.Error(err)
//	} else if len(saved) == 1 {
//		if saved[0].kind != TransfersCollection {
//			t.Errorf("saved[0].kind:'%v' != '%v'", saved[0].kind, TransfersCollection)
//		}
//		for _, p := range properties {
//			if p.Name == "AcknowledgeTime" {
//				t.Error("AcknowledgeTime should not be saved")
//			}
//		}
//	} else {
//		t.Errorf("len(saved):%v != 1", len(saved))
//	}
//
//	loadedTransfer := new(TransferData)
//	if err = loadedTransfer.Load(properties); err != nil {
//		t.Fatal(err)
//	}
//	if len(loadedTransfer.BothUserIDs) == 0 {
//		t.Error("len(loadedTransfer.BothUserIDs) == 0")
//	}
//}

//func TestTransferDump(t *testing.T) {
//	now := time.Now()
//	litter.Config.HidePrivateFields = true
//	t.Log("litter.Config.HidePrivateFields = true: ", litter.Sdump(now))
//	litter.Config.HidePrivateFields = false
//	t.Log("litter.Config.HidePrivateFields = false: ", litter.Options{HidePrivateFields: false}.Sdump(now))
//}
