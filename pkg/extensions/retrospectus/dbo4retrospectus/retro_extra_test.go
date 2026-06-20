package dbo4retrospectus

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/stretchr/testify/assert"
)

func ptrTime(t time.Time) *time.Time { return &t }

func TestTreePosition_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       TreePosition
		wantErr bool
	}{
		{"valid_zero", TreePosition{Parent: "p", Index: 0}, false},
		{"valid_positive", TreePosition{Parent: "p", Index: 5}, false},
		{"negative_index", TreePosition{Parent: "p", Index: -1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsKnownItemType(t *testing.T) {
	tests := []struct {
		v    string
		want bool
	}{
		{RetroItemTypeGood, true},
		{RetroItemTypeBad, true},
		{RetroItemTypeEndorsement, true},
		{RetroItemTypeIdea, true},
		{"kudos", false},
		{"", false},
		{"unknown", false},
	}
	for _, tt := range tests {
		t.Run(tt.v, func(t *testing.T) {
			assert.Equal(t, tt.want, IsKnownItemType(tt.v))
		})
	}
}

func TestRetrospective_BaseMeeting(t *testing.T) {
	v := &Retrospective{}
	assert.Same(t, &v.Meeting, v.BaseMeeting())
}

func TestRetrospective_Validate(t *testing.T) {
	// A minimal valid retrospective: empty Meeting is valid, stage set, TimeLastAction set, MaxVotesPerUser set.
	base := func() *Retrospective {
		return &Retrospective{
			Stage:          StageUpcoming,
			TimeLastAction: ptrTime(time.Now()),
			Settings:       RetrospectiveSettings{MaxVotesPerUser: 3},
		}
	}

	t.Run("valid_minimal", func(t *testing.T) {
		assert.NoError(t, base().Validate())
	})

	t.Run("missing_stage", func(t *testing.T) {
		v := base()
		v.Stage = ""
		assert.Error(t, v.Validate())
	})

	t.Run("unknown_stage", func(t *testing.T) {
		v := base()
		v.Stage = "no-such-stage"
		assert.Error(t, v.Validate())
	})

	t.Run("missing_time_last_action", func(t *testing.T) {
		v := base()
		v.TimeLastAction = nil
		assert.Error(t, v.Validate())
	})

	t.Run("finished_without_started", func(t *testing.T) {
		v := base()
		v.TimeFinished = ptrTime(time.Now())
		assert.Error(t, v.Validate())
	})

	t.Run("finished_with_started_ok", func(t *testing.T) {
		v := base()
		v.TimeStarted = ptrTime(time.Now())
		v.TimeFinished = ptrTime(time.Now())
		assert.NoError(t, v.Validate())
	})

	t.Run("missing_max_votes", func(t *testing.T) {
		v := base()
		v.Settings.MaxVotesPerUser = 0
		assert.Error(t, v.Validate())
	})

	t.Run("items_in_feedback_stage", func(t *testing.T) {
		v := base()
		v.Stage = StageFeedback
		v.Items = []*RetroItem{{ID: "i1", Title: "Item 1", Created: time.Now()}}
		assert.Error(t, v.Validate())
	})

	t.Run("duplicate_item_ids", func(t *testing.T) {
		v := base()
		v.Items = []*RetroItem{
			{ID: "dup", Title: "A", Created: time.Now()},
			{ID: "dup", Title: "B", Created: time.Now()},
		}
		assert.Error(t, v.Validate())
	})

	t.Run("vote_by_unknown_user", func(t *testing.T) {
		v := base()
		v.Items = []*RetroItem{
			{ID: "i1", Title: "A", Created: time.Now(), VotesByUser: map[string]int{"ghost": 1}},
		}
		assert.Error(t, v.Validate())
	})

	t.Run("unknown_count_item_type", func(t *testing.T) {
		v := base()
		v.CountsByMemberAndType = map[string]map[string]int{
			"m1": {"bogus": 1},
		}
		assert.Error(t, v.Validate())
	})

	t.Run("negative_count", func(t *testing.T) {
		v := base()
		v.CountsByMemberAndType = map[string]map[string]int{
			"m1": {RetroItemTypeGood: -1},
		}
		assert.Error(t, v.Validate())
	})

	t.Run("valid_counts", func(t *testing.T) {
		v := base()
		v.CountsByMemberAndType = map[string]map[string]int{
			"m1": {RetroItemTypeGood: 2, RetroItemTypeBad: 0},
		}
		assert.NoError(t, v.Validate())
	})
}

func TestRetrospective_GetMapOfRetroItemsByID(t *testing.T) {
	v := &Retrospective{
		Items: []*RetroItem{
			{ID: "a", Children: []*RetroItem{{ID: "a1"}}},
			{ID: "b"},
		},
	}
	byID, err := v.GetMapOfRetroItemsByID()
	assert.NoError(t, err)
	assert.Len(t, byID, 3)
	assert.Contains(t, byID, "a")
	assert.Contains(t, byID, "a1")
	assert.Contains(t, byID, "b")
}

func TestGetMapOfRetroItemsByID_DuplicateError(t *testing.T) {
	parent := &RetroItem{Children: []*RetroItem{
		{ID: "x"},
		{ID: "x"},
	}}
	_, err := GetMapOfRetroItemsByID(parent, make(map[string]*RetroItemTreeNode))
	assert.Error(t, err)
}

func TestRetroItemTreeNode_Accessors(t *testing.T) {
	parent := &RetroItem{ID: "p"}
	item := &RetroItem{ID: "c"}
	node := RetroItemTreeNode{item: item, parent: parent, index: 4}
	assert.Equal(t, 4, node.Index())
	assert.Same(t, item, node.Item())
	assert.Same(t, parent, node.Parent())
}

func TestRetroItemTreeNode_GetUpdatePath(t *testing.T) {
	t.Run("root_level", func(t *testing.T) {
		node := RetroItemTreeNode{item: &RetroItem{ID: "a"}, parent: nil, index: 2}
		assert.Equal(t, "items.2", node.GetUpdatePath(nil))
	})
	t.Run("nested", func(t *testing.T) {
		root := &RetroItem{ID: "root"}
		child := &RetroItem{ID: "child"}
		byID := map[string]*RetroItemTreeNode{
			"root":  {item: root, parent: nil, index: 1},
			"child": {item: child, parent: root, index: 3},
		}
		assert.Equal(t, "items.1.Children.3", byID["child"].GetUpdatePath(byID))
	})
}

func TestMoveRetroItem_Errors(t *testing.T) {
	items := func() []*RetroItem {
		return []*RetroItem{{ID: "p", Children: []*RetroItem{{ID: "c1"}, {ID: "c2"}}}}
	}
	t.Run("missing_id", func(t *testing.T) {
		err := MoveRetroItem(items(), "", TreePosition{Parent: "p"}, TreePosition{Parent: "p"})
		assert.Error(t, err)
	})
	t.Run("bad_from", func(t *testing.T) {
		err := MoveRetroItem(items(), "c1", TreePosition{Parent: "p", Index: -1}, TreePosition{Parent: "p"})
		assert.Error(t, err)
	})
	t.Run("bad_to", func(t *testing.T) {
		err := MoveRetroItem(items(), "c1", TreePosition{Parent: "p"}, TreePosition{Parent: "p", Index: -1})
		assert.Error(t, err)
	})
	t.Run("unknown_id", func(t *testing.T) {
		err := MoveRetroItem(items(), "ghost", TreePosition{Parent: "p"}, TreePosition{Parent: "p"})
		assert.ErrorIs(t, err, ErrItemNotFound)
	})
	t.Run("unknown_from_parent", func(t *testing.T) {
		err := MoveRetroItem(items(), "c1", TreePosition{Parent: "ghost"}, TreePosition{Parent: "p"})
		assert.ErrorIs(t, err, ErrItemNotFound)
	})
	t.Run("unknown_to_parent", func(t *testing.T) {
		err := MoveRetroItem(items(), "c1", TreePosition{Parent: "p"}, TreePosition{Parent: "ghost"})
		assert.ErrorIs(t, err, ErrItemNotFound)
	})
}

func TestRetroSettings_Validate(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var v *RetroSettings
		assert.NoError(t, v.Validate())
	})
	t.Run("missing_max_votes", func(t *testing.T) {
		v := &RetroSettings{}
		assert.Error(t, v.Validate())
	})
	t.Run("valid", func(t *testing.T) {
		v := &RetroSettings{MaxVotesPerUser: 3}
		assert.NoError(t, v.Validate())
	})
}

func TestRetrospectiveCounts_Validate(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := &RetrospectiveCounts{}
		assert.Error(t, v.Validate())
	})
	t.Run("valid", func(t *testing.T) {
		v := &RetrospectiveCounts{ItemsByUserAndType: map[string]map[string]int{"u1": {"good": 1}}}
		assert.NoError(t, v.Validate())
	})
}

func TestRetroSpaceDbo_Validate(t *testing.T) {
	t.Run("missing_settings", func(t *testing.T) {
		v := &RetroSpaceDbo{}
		assert.Error(t, v.Validate())
	})
	t.Run("valid_minimal", func(t *testing.T) {
		v := &RetroSpaceDbo{RetroSettings: RetroSettings{MaxVotesPerUser: 3}}
		assert.NoError(t, v.Validate())
	})
	t.Run("invalid_upcoming", func(t *testing.T) {
		v := &RetroSpaceDbo{
			RetroSettings: RetroSettings{MaxVotesPerUser: 3},
			UpcomingRetro: &RetrospectiveCounts{}, // empty -> invalid
		}
		assert.Error(t, v.Validate())
	})
	t.Run("valid_upcoming", func(t *testing.T) {
		v := &RetroSpaceDbo{
			RetroSettings: RetroSettings{MaxVotesPerUser: 3},
			UpcomingRetro: &RetrospectiveCounts{ItemsByUserAndType: map[string]map[string]int{"u1": {"good": 1}}},
		}
		assert.NoError(t, v.Validate())
	})
}

func TestRetroSpaceDbo_ActiveRetro(t *testing.T) {
	t.Run("nil_active", func(t *testing.T) {
		v := &RetroSpaceDbo{}
		assert.Equal(t, dbo4spaceus.SpaceMeetingInfo{}, v.ActiveRetro())
	})
	t.Run("with_active", func(t *testing.T) {
		active := &dbo4spaceus.SpaceMeetingInfo{ID: "m1", Stage: "active"}
		v := &RetroSpaceDbo{Active: active}
		assert.Equal(t, *active, v.ActiveRetro())
	})
}
