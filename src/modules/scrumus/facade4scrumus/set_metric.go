package facade4scrumus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/models4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

// SetMetric sets metric
func SetMetric(ctx context.Context, userContext facade.User, request SetMetricRequest) (response *SetMetricRequest, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	uid := userContext.GetID()
	err = runScrumWorker(ctx, userContext, request.Request,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params facade4meetingus.WorkerParams) (err error) {
			var teamMetric *models4teamus.TeamMetric
			for _, m := range params.Team.Data.Metrics {
				if m.ID == request.Metric {
					teamMetric = m
					break
				}
			}
			if teamMetric == nil {
				return fmt.Errorf("unknown metric: %s", request.Metric)
			}
			var scrumUpdates []dal.Update
			p := setMetricParams{
				uid:        uid,
				request:    request,
				scrum:      params.Meeting.Record.Data().(*models4scrumus.Scrum),
				teamMetric: teamMetric,
			}
			switch teamMetric.Mode {
			case "team":
				if scrumUpdates, err = setTeamMetric(p); err != nil {
					return
				}
			case "personal":
				if scrumUpdates, err = setPersonalMetric(p, params.TeamModuleEntry.Data); err != nil {
					return
				}
			default:
				err = fmt.Errorf("unknown metric mode: %s", teamMetric.Mode)
			}
			if len(scrumUpdates) > 0 {
				if err = tx.Update(ctx, params.Meeting.Key, scrumUpdates); err != nil {
					return
				}
			}
			return err
		})
	return
}

type setMetricParams struct {
	uid        string
	request    SetMetricRequest
	scrum      *models4scrumus.Scrum
	teamMetric *models4teamus.TeamMetric
}

func setPersonalMetric(p setMetricParams, contactusTeam *models4contactus.ContactusTeamDbo) (scrumUpdates []dal.Update, err error) {
	var status *models4scrumus.MemberStatus
	var teamMember *briefs4contactus.ContactBrief
	var teamMemberContactID string
	for contactID, contact := range contactusTeam.Contacts {
		if contactID == p.request.Member {
			teamMember = contact
			teamMemberContactID = contactID
			break
		}
	}
	if teamMember == nil {
		err = validation.NewErrBadRequestFieldValue("members", fmt.Sprintf("unknown members ContactID: %s", p.request.Member))
		return
	}
	for id, s := range p.scrum.Statuses {
		if id == p.request.Member {
			status = s
			goto UpdateMember
		}
	}
	status = &models4scrumus.MemberStatus{
		Member: models4scrumus.ScrumMember{
			ID:    teamMemberContactID,
			Title: teamMember.Title,
		},
		Metrics: make([]*models4scrumus.MetricRecord, 0, 1),
	}
UpdateMember:
	var changed bool
	changed, status.Metrics, scrumUpdates, err = setMetric(p, status.Metrics)
	if changed {
		scrumUpdates = append(scrumUpdates, dal.Update{
			Field: "statuses",
			Value: p.scrum.Statuses,
		})
	}
	return
}

func setTeamMetric(p setMetricParams) (scrumUpdates []dal.Update, err error) {
	var changed bool
	changed, p.scrum.TeamMetrics, scrumUpdates, err = setMetric(p, p.scrum.TeamMetrics)
	if changed {
		scrumUpdates = append(scrumUpdates, dal.Update{
			Field: "metrics",
			Value: p.scrum.Metrics,
		})
	}
	return
}

func setMetric(p setMetricParams, metrics []*models4scrumus.MetricRecord) (changed bool, updatedMetrics []*models4scrumus.MetricRecord, scrumUpdates []dal.Update, err error) {
	var metric *models4scrumus.MetricRecord
	isExistingRecord := true

	for _, m := range metrics {
		if m.ID == p.request.Metric {
			metric = m
			goto UpdateMetric
		}
	}
	metric = &models4scrumus.MetricRecord{
		ID:          p.request.Metric,
		UID:         p.uid,
		MetricValue: models4scrumus.MetricValue{},
	}
	isExistingRecord = false
UpdateMetric:
	switch p.teamMetric.Type {
	case "bool":
		if p.request.Bool == nil {
			err = validation.NewErrRecordIsMissingRequiredField("bool")
			return
		}
		if isExistingRecord && metric.Bool != nil && *metric.Bool == *p.request.Bool {
			return
		}
		metric.Bool = p.request.Bool
	case "int":
		if p.request.Int == nil {
			err = validation.NewErrRecordIsMissingRequiredField("int")
			return
		}
		if isExistingRecord && metric.Int != nil && *metric.Int == *p.request.Int {
			return
		}
		metric.Int = p.request.Int
	case "str":
		if p.request.Str == nil {
			err = validation.NewErrRecordIsMissingRequiredField("str")
			return
		}
		if isExistingRecord && metric.Str != nil && *metric.Str == *p.request.Str {
			return
		}
		metric.Str = p.request.Str
	}
	metrics = append(metrics, metric)
	changed = true
	if !isExistingRecord {
		p.scrum.Metrics = append(p.scrum.Metrics, p.teamMetric)
		scrumUpdates = []dal.Update{{
			Field: "teamMetrics",
			Value: metrics,
		}}
	}
	return
}
