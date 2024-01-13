package facade4calendarium

//var briefsAdapter = func(listType models4listus.ListType, listID string) facade4teamus.BriefsAdapter {
//	getListGroupByListID := func(team *models4teamus.TeamDto) *models4listus.ListGroup {
//		for _, lg := range team.ListGroups {
//			for _, l := range lg.Lists {
//				if l.ContactID == listID {
//					if lg.Lists == nil {
//						lg.Lists = make([]*models4listus.ListBrief, 0)
//					}
//					return lg
//				}
//			}
//		}
//		return nil
//	}
//
//	getListGroupByID := func(team *models4teamus.TeamDto) *models4listus.ListGroup {
//		for _, lg := range team.ListGroups {
//			if lg.Role == listType {
//				return lg
//			}
//		}
//		lg := &models4listus.ListGroup{Role: listType}
//		team.ListGroups = append(team.ListGroups, lg)
//		return lg
//	}
//
//	var teamCache *models4teamus.TeamDto
//	var listGroupCache *models4listus.ListGroup
//
//	getListGroup := func(team *models4teamus.TeamDto) *models4listus.ListGroup {
//		if team == nil {
//			panic("team == nil")
//		}
//		if team == teamCache && listGroupCache != nil {
//			return listGroupCache
//		}
//		teamCache = nil
//		if listType != "" {
//			listGroupCache = getListGroupByID(team)
//		} else if listType != "" {
//			listGroupCache = getListGroupByListID(team)
//		} else {
//			panic("Both parameter `listType` and `listID` are empty")
//		}
//		teamCache = team
//		return listGroupCache
//	}
//
//	return facade4teamus.BriefsAdapter{
//		BriefsFieldName: "listGroups." + listType,
//		BriefsValue: func(team *models4teamus.TeamDto) interface{} {
//			lg := getListGroup(team)
//			return lg.Lists
//		},
//		GetBriefsCount: func(team *models4teamus.TeamDto) int {
//			lg := getListGroup(team)
//			return len(lg.Lists)
//		},
//		GetBriefItemID: func(team *models4teamus.TeamDto, i int) (id string) {
//			lg := getListGroup(team)
//			return lg.Lists[i].ContactID
//		},
//		ShiftBriefs: func(team *models4teamus.TeamDto, from facade4teamus.SliceIndexes, to facade4teamus.SliceIndexes) {
//			lg := getListGroup(team)
//			copy(
//				lg.Lists[to.Departs:to.Arrives],
//				lg.Lists[from.Departs:from.Arrives],
//			)
//		},
//		TrimBriefs: func(team *models4teamus.TeamDto, count int) {
//			lg := getListGroup(team)
//			lg.Lists = lg.Lists[:count]
//		},
//	}
//}
