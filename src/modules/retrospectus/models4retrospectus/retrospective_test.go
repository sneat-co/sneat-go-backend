package models4retrospectus

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRetroItem_Validate(t *testing.T) {
	t.Run("should_succeed", func(t *testing.T) {
		t.Run("without_children", func(t *testing.T) {
			v := RetroItem{
				ID:      "123",
				Title:   "ABC",
				Created: time.Now(),
			}
			if err := v.Validate(); err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
		t.Run("with_children", func(t *testing.T) {
			v := &RetroItem{
				ID:      "123",
				Title:   "ABC",
				Created: time.Now(),
				Children: []*RetroItem{
					{ID: "child1", Title: "Child 1", Created: time.Now()},
					{ID: "child2", Title: "Child 2", Created: time.Now()},
				},
			}
			if err := v.Validate(); err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	})

	t.Run("should_fail", func(t *testing.T) {
		t.Run("no_ID", func(t *testing.T) {
			v := &RetroItem{
				Title:   "ABC",
				Created: time.Now(),
			}
			if err := v.Validate(); err == nil {
				t.Error("Expected to get error")
				return
			}
		})

		t.Run("no_Title", func(t *testing.T) {
			v := &RetroItem{
				ID:      "123",
				Created: time.Now(),
			}
			if err := v.Validate(); err == nil {
				t.Error("Expected to get error")
				return
			}
		})

		t.Run("bad_child", func(t *testing.T) {
			v := &RetroItem{
				ID:      "123",
				Title:   "ABC",
				Created: time.Now(),
				Children: []*RetroItem{
					{ID: "child1", Title: "Child 1", Created: time.Now()},
					{ID: "child2", Title: "Child 2", Created: time.Now()},
				},
			}
			if err := v.Validate(); err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
		})
	})
}

func TestMoveRetroItem(t *testing.T) {

	var test = func(id string, fromParent string, fromIndex int, toParent string, toIndex int, expectedSource, expectedTarget string, getItems func() []*RetroItem) error {
		items := getItems()
		rootItemsCount := len(items)
		from := TreePosition{Parent: fromParent, Index: fromIndex}
		to := TreePosition{Parent: toParent, Index: toIndex}
		if err := MoveRetroItem(items, id, from, to); err != nil {
			t.Fatalf("failed to move item: %v", err)
		}
		if count := len(items); count != rootItemsCount {
			t.Fatalf("len(items): %v != %v", count, rootItemsCount)
		}
		var fromParentItem, toParentItem *RetroItem
		for _, item := range items {
			if item.ID == fromParent {
				fromParentItem = item
			}
			if item.ID == toParent {
				toParentItem = item
			}
		}
		if fromParentItem == nil {
			t.Fatalf("fromParentItem == nil")
		}
		if toParentItem == nil {
			t.Fatalf("toParentItem == nil")
		}
		if fromParentItem.ID != fromParent {
			t.Fatalf("fromParentItem.InviteID:%v != %v", fromParentItem.ID, fromParent)
		}
		if toParentItem.ID != toParent {
			t.Fatalf("toParentItem.InviteID:%v != %v", toParentItem.ID, toParent)
		}
		check := func(expected string, children []*RetroItem) error {
			ids := make([]string, 0, len(children))
			for _, child := range children {
				ids = append(ids, child.ID)
			}
			if actual := strings.Join(ids, ","); actual != expected {
				return fmt.Errorf("expected %v got %v", expected, actual)
			}
			return nil
		}
		if err := check(expectedTarget, toParentItem.Children); err != nil {
			return fmt.Errorf("invalid target: %w", err)
		}
		if fromParent != toParent {
			if err := check(expectedSource, fromParentItem.Children); err != nil {
				return fmt.Errorf("invalid source: %w", err)
			}
		}
		return nil
	}

	t.Run("move_to_different_parent", func(t *testing.T) {
		getItems2 := func() []*RetroItem {
			return []*RetroItem{
				{ID: "good", Children: []*RetroItem{
					{ID: "g1", Title: "Good #1"},
					{ID: "g2", Title: "Good #2"},
					{ID: "g3", Title: "Good #3"},
				}},
				{ID: "bad", Children: []*RetroItem{
					{ID: "b1", Title: "Bad #1"},
					{ID: "b2", Title: "Bad #2"},
					{ID: "b3", Title: "Bad #3"},
				}},
			}
		}
		if err := test("g1", "good", 0, "bad", 0, "g2,g3", "g1,b1,b2,b3", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g1", "good", 1, "bad", 1, "g2,g3", "b1,g1,b2,b3", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g1", "good", 1, "bad", 2, "g2,g3", "b1,b2,g1,b3", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g1", "good", 1, "bad", 3, "g2,g3", "b1,b2,b3,g1", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g1", "good", 1, "bad", 100, "g2,g3", "b1,b2,b3,g1", getItems2); err != nil {
			t.Error(err)
		}
	})

	t.Run("move_within_same_parent", func(t *testing.T) {
		getItems2 := func() []*RetroItem {
			return []*RetroItem{
				{ID: "goods", Children: []*RetroItem{
					{ID: "g1", Title: "Good #1"},
					{ID: "g2", Title: "Good #2"},
					{ID: "g3", Title: "Good #3"},
					{ID: "g4", Title: "Good #4"},
					{ID: "g5", Title: "Good #5"},
				}},
			}
		}
		const parent = "goods"

		if err := test("g1", parent, 0, parent, 1, "", "g2,g1,g3,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g1", parent, 0, parent, 4, "", "g2,g3,g4,g5,g1", getItems2); err != nil {
			t.Error(err)
		}

		if err := test("g3", parent, 2, parent, 2, "", "g1,g2,g3,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 2, parent, 0, "", "g3,g1,g2,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 2, parent, 1, "", "g1,g3,g2,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 2, parent, 3, "", "g1,g2,g4,g3,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 2, parent, 4, "", "g1,g2,g4,g5,g3", getItems2); err != nil {
			t.Error(err)
		}

		if err := test("g3", parent, 1, parent, 2, "", "g1,g2,g3,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 1, parent, 0, "", "g3,g1,g2,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 2, parent, 1, "", "g1,g3,g2,g4,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 2, parent, 3, "", "g1,g2,g4,g3,g5", getItems2); err != nil {
			t.Error(err)
		}
		if err := test("g3", parent, 1, parent, 4, "", "g1,g2,g4,g5,g3", getItems2); err != nil {
			t.Error(err)
		}
	})
}
