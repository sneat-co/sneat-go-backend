package facade4calendarium

//var briefsAdapter = func(listType dbo4listus.ListType, listID string) facade4spaceus.BriefsAdapter {
//	getListGroupByListID := func(team *dbo4spaceus.TeamDto) *dbo4listus.ListGroup {
//		for _, lg := range team.ListGroups {
//			for _, l := range lg.Lists {
//				if l.ContactID == listID {
//					if lg.Lists == nil {
//						lg.Lists = make([]*dbo4listus.ListBrief, 0)
//					}
//					return lg
//				}
//			}
//		}
//		return nil
//	}
//
//	getListGroupByID := func(team *dbo4spaceus.TeamDto) *dbo4listus.ListGroup {
//		for _, lg := range team.ListGroups {
//			if lg.Role == listType {
//				return lg
//			}
//		}
//		lg := &dbo4listus.ListGroup{Role: listType}
//		team.ListGroups = append(team.ListGroups, lg)
//		return lg
//	}
//
//	var teamCache *dbo4spaceus.TeamDto
//	var listGroupCache *dbo4listus.ListGroup
//
//	getListGroup := func(team *dbo4spaceus.TeamDto) *dbo4listus.ListGroup {
//		if team == nil {
//			panic("space == nil")
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
//	return facade4spaceus.BriefsAdapter{
//		BriefsFieldName: "listGroups." + listType,
//		BriefsValue: func(team *dbo4spaceus.TeamDto) interface{} {
//			lg := getListGroup(team)
//			return lg.Lists
//		},
//		GetBriefsCount: func(team *dbo4spaceus.TeamDto) int {
//			lg := getListGroup(team)
//			return len(lg.Lists)
//		},
//		GetBriefItemID: func(team *dbo4spaceus.TeamDto, i int) (id string) {
//			lg := getListGroup(team)
//			return lg.Lists[i].ContactID
//		},
//		ShiftBriefs: func(team *dbo4spaceus.TeamDto, from facade4spaceus.SliceIndexes, to facade4spaceus.SliceIndexes) {
//			lg := getListGroup(team)
//			copy(
//				lg.Lists[to.Departs:to.Arrives],
//				lg.Lists[from.Departs:from.Arrives],
//			)
//		},
//		TrimBriefs: func(team *dbo4spaceus.TeamDto, count int) {
//			lg := getListGroup(team)
//			lg.Lists = lg.Lists[:count]
//		},
//	}
//}
