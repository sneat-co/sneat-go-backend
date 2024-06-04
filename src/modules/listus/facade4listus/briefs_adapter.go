package facade4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

var briefsAdapter = func(listType dbo4listus.ListType, listID string) *dal4teamus.BriefsAdapter[*dbo4listus.ListusTeamDto] {
	getListGroupByListID := func(moduleTeam *dbo4listus.ListusTeamDto) *dbo4listus.ListGroup {
		//for _, lg := range moduleTeam.ListGroups {
		//	for _, l := range lg.Lists {
		//		if l.ID == listID {
		//			if lg.Lists == nil {
		//				lg.Lists = make([]*dbo4listus.ListBrief, 0)
		//			}
		//			return lg
		//		}
		//	}
		//}
		return nil
	}

	getListGroupByID := func(moduleTeam *dbo4listus.ListusTeamDto) *dbo4listus.ListGroup {
		panic("not implemented")
		//for _, lg := range moduleTeam.ListGroups {
		//	if lg.Type == listType {
		//		return lg
		//	}
		//}
		//lg := &dbo4listus.ListGroup{Type: listType}
		//moduleTeam.ListGroups = append(moduleTeam.ListGroups, lg)
		//return lg
	}

	var teamCache *dbo4listus.ListusTeamDto
	var listGroupCache *dbo4listus.ListGroup

	getListGroup := func(team *dbo4listus.ListusTeamDto) *dbo4listus.ListGroup {
		if team == nil {
			panic("team == nil")
		}
		if team == teamCache && listGroupCache != nil {
			return listGroupCache
		}
		teamCache = nil
		if listType != "" {
			listGroupCache = getListGroupByID(team)
		} else if listID != "" {
			listGroupCache = getListGroupByListID(team)
		} else {
			panic("Both parameter `listType` and `listID` are empty")
		}
		teamCache = team
		return listGroupCache
	}

	return &dal4teamus.BriefsAdapter[*dbo4listus.ListusTeamDto]{
		BriefsFieldName: "listGroups." + listType,
		BriefsValue: func(team *dbo4listus.ListusTeamDto) interface{} {
			lg := getListGroup(team)
			return lg.Lists
		},
		GetBriefsCount: func(team *dbo4listus.ListusTeamDto) int {
			lg := getListGroup(team)
			return len(lg.Lists)
		},
		GetBriefItemID: func(team *dbo4listus.ListusTeamDto, i int) (id string) {
			panic("not implemented")
			//lg := getListGroup(team)
			//return lg.Lists[i].ID
		},
		ShiftBriefs: func(team *dbo4listus.ListusTeamDto, from dal4teamus.SliceIndexes, to dal4teamus.SliceIndexes) {
			panic("not implemented")
			//lg := getListGroup(team)
			//copy(
			//	lg.Lists[to.Start:to.End],
			//	lg.Lists[from.Start:from.End],
			//)
		},
		TrimBriefs: func(team *dbo4listus.ListusTeamDto, count int) {
			panic("not implemented")
			//lg := getListGroup(team)
			//lg.Lists = lg.Lists[:count]
		},
	}
}
