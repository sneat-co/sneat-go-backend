package facade4retrospectus

//func updateTeamWithUpcomingRetroUserCounts(
//	ctx context.Context,
//	tx dal.ReadwriteTransaction,
//	now time.Time,
//	uid, teamID string,
//	itemsByType dbretro.RetroItemsByType,
//) (err error) {
//	if uid == "" {
//		return errors.New("uid is a required parameter")
//	}
//	if teamID == "" {
//		return errors.New("teamID is a required parameter")
//	}
//	if now.IsZero() {
//		return errors.New("now is a required parameter")
//	}
//	teamKey := newTeamKey(teamID)
//	var team models4teamus.TeamContext
//	var teamRecord dal.Record
//	if team, err = txGetTeamByID(ctx, tx, teamID); err != nil {
//		return err
//	}
//	if !teamRecord.Exists() {
//		return dal.NewErrNotFoundByKey(teamKey, dal.ErrRecordNotFound)
//	}
//	teamUpdates := make([]dal.Update, 0, 1)
//	path := fmt.Sprintf("upcomingRetro.itemsByUserAndType.%v", uid)
//	if len(itemsByType) == 0 {
//		upcomingRetro := team.Data.UpcomingRetro
//		if upcomingRetro != nil && upcomingRetro.ItemsByUserAndType != nil {
//			teamUpdates = append(teamUpdates, dal.Update{Field: path, Value: dal.DeleteField})
//			delete(upcomingRetro.ItemsByUserAndType, uid)
//		}
//	} else {
//		currentUserCounts := make(map[string]int)
//		for itemType, items := range itemsByType {
//			if !dbretro.IsKnownItemType(itemType) {
//				return validation.NewErrBadRecordFieldValue("itemsByType", fmt.Sprintf("unknown value: %v", itemType))
//			}
//			if count := len(items); count > 0 {
//				currentUserCounts[itemType] = count
//				teamUpdates = append(teamUpdates, dal.Update{Field: fmt.Sprintf("%v.%v", path, itemType), Value: count})
//			}
//		}
//
//		if team.Data.UpcomingRetro != nil {
//			if existingCounts, teamHasUserCounts := team.Data.UpcomingRetro.ItemsByUserAndType[uid]; teamHasUserCounts {
//				for itemType := range existingCounts {
//					if _, exist := currentUserCounts[itemType]; !exist {
//						teamUpdates = append(teamUpdates, dal.Update{Field: fmt.Sprintf("%v.%v", path, itemType), Value: dal.DeleteField})
//					}
//				}
//			}
//		}
//		if team.Data.UpcomingRetro == nil {
//			team.Data.UpcomingRetro = &dbmodels.RetrospectiveCounts{}
//		}
//		if team.Data.UpcomingRetro.ItemsByUserAndType == nil {
//			team.Data.UpcomingRetro.ItemsByUserAndType = make(map[string]map[string]int, 1)
//		}
//		team.Data.UpcomingRetro.ItemsByUserAndType[uid] = currentUserCounts
//	}
//	if len(teamUpdates) > 0 {
//		if err = txUpdateTeam(ctx, tx, now, team, teamUpdates); err != nil {
//			return err
//		}
//	}
//	return err
//}
