package models4debtus

//func TestRemoveDefaults(t *testing.T) {
//	properties := []datastore.Property{
//		{Name: "Int1", Value: int64(1)}, // Keep
//		{Name: "Int0", Value: int64(0)},
//		{Name: "Int-1", Value: int64(-1)},       // Keep
//		{Name: "Int10", Value: int64(10)},       // Keep
//		{Name: "Float1.1", Value: float64(1.2)}, // Keep
//		{Name: "Float0", Value: float64(0.0)},
//		{Name: "Obsolete", Value: 15},
//		{Name: "Time0", Value: time.Time{}},
//		{Name: "TimeNow", Value: time.Now()}, // Keep
//
//	}
//
//	filtered, _ := gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
//		"Int0":     gaedb.IsZeroInt,
//		"Int1":     gaedb.IsZeroInt,
//		"Int-1":    gaedb.IsZeroInt,
//		"Obsolete": gaedb.IsObsolete,
//		"Float1.1": gaedb.IsZeroFloat,
//		"Float0":   gaedb.IsZeroFloat,
//		"Time0":    gaedb.IsZeroTime,
//		"TimeNow":  gaedb.IsZeroTime,
//	})
//
//	if len(filtered) != 5 {
//		t.Errorf("Expected 5 properties, got %d: %v", len(filtered), filtered)
//	}
//}
