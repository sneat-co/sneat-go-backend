package facade4invitus

//func stringHash(s string) (hash int32) {
//	if len(s) == 0 {
//		return
//	}
//	for _, char := range s {
//		hash = ((hash << 5) - hash) + char
//	}
//	return
//}

//func getPin(teamID, role, uid string) (pin int32) {
//	if pin = stringHash(teamID + "-" + role + "-" + uid); pin < 0 {
//		pin = -pin
//	}
//	return
//}

//func verifyPin(teamID string, pin int32, members []*briefs4memberus.MemberBrief) (role string, inviter *briefs4memberus.MemberBrief) {
//	for _, m := range members {
//		if m.UserID != "" {
//			if pin == getPin(teamID, briefs4memberus.TeamMemberRoleContributor, m.UserID) {
//				return briefs4memberus.TeamMemberRoleContributor, m
//			}
//			if pin == getPin(teamID, briefs4memberus.TeamMemberRoleSpectator, m.UserID) {
//				return briefs4memberus.TeamMemberRoleSpectator, m
//			}
//		}
//	}
//	return "", nil
//}
