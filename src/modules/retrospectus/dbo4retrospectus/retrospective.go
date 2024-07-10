package dbo4retrospectus

import (
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/dbo4meetingus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/reorder"
	"github.com/strongo/validation"
	"time"
)

// RetroUser record
type RetroUser struct {
	UserID string `json:"userId" firestore:"userId"`
	Title  string `json:"title" firestore:"title"`
}

// RetroItem record
type RetroItem struct {
	ID          string         `json:"ContactID" firestore:"ContactID"`           // TODO: ask StackOverflow: if "id" it gets lost
	Type        string         `json:"type,omitempty" firestore:"type,omitempty"` // required only for root level items
	Title       string         `json:"title,omitempty" firestore:"title,omitempty"`
	Created     time.Time      `json:"created" firestore:"created"`
	By          *RetroUser     `json:"by,omitempty" firestore:"id,omitempty"`
	Children    []*RetroItem   `json:"children,omitempty" firestore:"children,omitempty"`
	VotesByUser map[string]int `json:"votesByUser,omitempty" firestore:"votesByUser,omitempty"`
}

// Validate validates RetroItem record
func (v *RetroItem) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("ContactID")
	}
	if v.Title == "" {
		return validation.NewErrRecordIsMissingRequiredField("Title")
	}
	if v.Created.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("Created")
	}
	if v.By != nil {
		if v.By.Title == "" {
			return validation.NewErrRecordIsMissingRequiredField("By.Title")
		}
		if v.By.UserID == "" {
			return validation.NewErrRecordIsMissingRequiredField("By.UserID")
		}
	}
	if itemsCount := len(v.Children); itemsCount > 0 {
		ids := make([]string, itemsCount)
		for i, item := range v.Children {
			if err := item.Validate(); err != nil {
				return fmt.Errorf("invalid item at index %v: %w", i, err)
			}
			if item.ID == v.ID {
				return fmt.Errorf("items[%v] has same ContactID as its parent: '%v'", i, v.ID)
			}
			if prevIndex := reorder.IndexOf(ids, item.ID); prevIndex >= 0 {
				return fmt.Errorf("child item of '%v' at index %v has same id '%v' as item at index %v", v.ID, i, item.ID, prevIndex)
			}
			ids = append(ids, item.ID)
		}
	}
	return nil
}

// TreePosition params
type TreePosition struct {
	Parent string `json:"parent"`
	Index  int    `json:"index"`
}

// Validate validates TreePosition params
func (v *TreePosition) Validate() error {
	if v.Index < 0 {
		return validation.NewErrBadRecordFieldValue("index", "is negative")
	}
	return nil
}

// RetroItemsByType is a type alias
type RetroItemsByType = map[string][]*RetroItem

// ErrItemNotFound error to indicate an item not found by ContactID
var ErrItemNotFound = errors.New("item not found")

const (
	// RetroItemTypeGood "good"
	RetroItemTypeGood = "good"

	// RetroItemTypeBad "bad"
	RetroItemTypeBad = "bad"

	// RetroItemTypeEndorsement "endorsement"
	RetroItemTypeEndorsement = "endorsement"

	// RetroItemTypeIdea "idea"
	RetroItemTypeIdea = "idea"
)

// IsKnownItemType checks if value is a known item type
func IsKnownItemType(v string) bool {
	switch v {
	case RetroItemTypeGood, RetroItemTypeBad, RetroItemTypeEndorsement, RetroItemTypeIdea: // do not use kudos for endorsement
		return true
	default:
		return false
	}
}

const (
	// StageUpcoming "upcoming"
	StageUpcoming = "upcoming"

	// StageFeedback "feedback"
	StageFeedback = "feedback"

	// StageReview "review"
	StageReview = "review"
)

// DefaultMaxVotesPerUser 3
const DefaultMaxVotesPerUser = 3

// RetrospectiveSettings record
type RetrospectiveSettings struct {
	MaxVotesPerUser int `json:"maxVotesPerUser" firestore:"maxVotesPerUser"` // Defines maximum number of votes per each user
}

// Retrospective record
type Retrospective struct {
	dbo4meetingus.Meeting
	//RetroProfileID  string               `json:"retroProfileId" firestore:"retroProfileId"`
	Stage          string                `json:"stage" firestore:"stage"`                               // TODO: consider move to dbo4meetingus.Meeting?
	TimeLastAction *time.Time            `json:"timeAction,omitempty" firestore:"timeAction,omitempty"` // TODO: document why we need this, or consider to remove
	TimeStarts     *time.Time            `json:"timeStarts,omitempty" firestore:"timeStarts,omitempty"`
	TimeStarted    *time.Time            `json:"timeStarted,omitempty" firestore:"timeStarted,omitempty"`
	TimeFinished   *time.Time            `json:"timeFinished,omitempty" firestore:"timeFinished,omitempty"`
	ScheduledBy    *dbmodels.ByUser      `json:"scheduledBy,omitempty" firestore:"scheduledBy,omitempty"`
	StartedBy      *dbmodels.ByUser      `json:"startedBy,omitempty" firestore:"startedBy,omitempty"`
	Items          []*RetroItem          `json:"items,omitempty" firestore:"items,omitempty"`
	Settings       RetrospectiveSettings `json:"settings" firestore:"settings"`
	//
	CountsByMemberAndType map[string]map[string]int `json:"countsByMemberAndType,omitempty" firestore:"countsByUserAndType,omitempty"`
}

// BaseMeeting returns base api4meetingus info
func (v *Retrospective) BaseMeeting() *dbo4meetingus.Meeting {
	return &v.Meeting
}

// MoveRetroItem reorder item
func MoveRetroItem(items []*RetroItem, id string, from, to TreePosition) error {
	if id == "" {
		return errors.New("missing required parameter: id")
	}
	if err := from.Validate(); err != nil {
		return fmt.Errorf("bad parameter 'from': %w", err)
	}
	if err := to.Validate(); err != nil {
		return fmt.Errorf("bad parameter 'to': %w", err)
	}
	var fromParent, toParent *RetroItem
	var movingNode *RetroItemTreeNode
	nodesByID, err := GetMapOfRetroItemsByID(&RetroItem{Children: items}, make(map[string]*RetroItemTreeNode))
	if err != nil {
		return err
	}
	movingNode = nodesByID[id]
	if movingNode == nil {
		return fmt.Errorf("%w: id=%v", ErrItemNotFound, id)
	}
	fromNode := nodesByID[from.Parent]
	if fromNode == nil {
		return fmt.Errorf("%w: unknown From parent: %v", ErrItemNotFound, from.Parent)
	}
	fromParent = fromNode.item
	toNode := nodesByID[to.Parent]
	if toNode == nil {
		return fmt.Errorf("%w: unknown To parent: %v", ErrItemNotFound, to.Parent)
	}
	toParent = toNode.item

	if from.Parent == to.Parent {
		if fromParent.Children[from.Index].ID == id {
			moveRetroItemWithinSameParent(fromParent.Children, from.Index, to.Index)
			return nil
		}
		indexAtParent := indexOfRetroItem(fromParent.Children, id)
		if indexAtParent >= 0 {
			moveRetroItemWithinSameParent(fromParent.Children, indexAtParent, to.Index)
			return nil
		}
		// return nil?: TODO: Decide if we want to move node if it was moved to another parent
	}
	if indexAtToParent := indexOfRetroItem(toParent.Children, id); indexAtToParent >= 0 { // Item already belongs to target parent
		if len(toParent.Children) > indexAtToParent+1 { // In case some items were removed since requested
			to.Index = len(toParent.Children) - 1
		}
		moveRetroItemWithinSameParent(toParent.Children, indexAtToParent, to.Index)
		return nil
	}
	movingNode.parent.Children = removeRetroItem(movingNode.parent.Children, id)
	if to.Index >= len(toParent.Children) { // In case some items were removed since requested
		toParent.Children = append(toParent.Children, movingNode.item)
	} else {
		toParent.Children = insertRetroItem(toParent.Children, to.Index, movingNode.item)
	}
	return nil
}

// GetMapOfRetroItemsByID return map of items by ContactID
func (v *Retrospective) GetMapOfRetroItemsByID() (map[string]*RetroItemTreeNode, error) {
	return GetMapOfRetroItemsByID(&RetroItem{
		Children: v.Items,
	}, make(map[string]*RetroItemTreeNode))
}

func indexOfRetroItem(arr []*RetroItem, id string) int {
	for i, item := range arr {
		if item.ID == id {
			return i
		}
	}
	return -1
}

func insertRetroItem(arr []*RetroItem, index int, v *RetroItem) []*RetroItem {
	a := make([]*RetroItem, 0, len(arr)+1)
	a = append(a, arr[:index]...)
	a = append(a, v)
	return append(a, arr[index:]...)
}

func removeRetroItem(arr []*RetroItem, id string) []*RetroItem {
	//a := make([]*RetroItem, 0, len(arr))
	shift := 0
	for i, v := range arr {
		if v.ID == id {
			shift++
		} else if shift > 0 {
			arr[i-shift] = v
		}
	}
	return arr[:len(arr)-shift]
}

func moveRetroItemWithinSameParent(items []*RetroItem, from, to int) {
	if from == to {
		return
	}
	item := items[from]
	if from < to {
		for i := from; i < to; i++ {
			items[i] = items[i+1]
		}
	} else if to < from {
		for i := from; i > to; i-- {
			items[i] = items[i-1]
		}
	}
	items[to] = item
}

// RetroItemTreeNode record
type RetroItemTreeNode struct {
	item   *RetroItem
	parent *RetroItem
	index  int
}

// Index returns index of item
func (v RetroItemTreeNode) Index() int {
	return v.index
}

// Item return item belonging to the node
func (v RetroItemTreeNode) Item() *RetroItem {
	return v.item
}

// Parent returns parent note
func (v RetroItemTreeNode) Parent() *RetroItem {
	return v.parent
}

// GetUpdatePath return update path
func (v RetroItemTreeNode) GetUpdatePath(byID map[string]*RetroItemTreeNode) string {
	if v.parent == nil {
		return fmt.Sprintf("items.%v", v.index)
	}
	parent := byID[v.parent.ID]
	parentPath := parent.GetUpdatePath(byID)
	return fmt.Sprintf("%v.Children.%v", parentPath, v.index)
}

// GetMapOfRetroItemsByID returns map of items by ContactID
func GetMapOfRetroItemsByID(parent *RetroItem, byID map[string]*RetroItemTreeNode) (map[string]*RetroItemTreeNode, error) {
	for i, item := range parent.Children {
		if prevItem, ok := byID[item.ID]; ok {
			return byID, fmt.Errorf("child '%v' of '%v' at index %v has same id as child of '%v'", item.ID, parent.ID, i, prevItem.parent.ID)
		}
		byID[item.ID] = &RetroItemTreeNode{item: item, parent: parent, index: i}
		if byID, err := GetMapOfRetroItemsByID(item, byID); err != nil {
			return byID, fmt.Errorf("%v: %w", item.ID, err)
		}
	}
	return byID, nil
}

// Validate validates Retrospective record
func (v *Retrospective) Validate() error {
	if err := v.Meeting.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("meeting", err.Error())
	}
	switch v.Stage {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("stage")
	case StageUpcoming, StageFeedback, StageReview:
		break
	default:
		return validation.NewErrBadRecordFieldValue("stage", fmt.Sprintf("unknown value: %v", v.Stage))
	}
	if v.TimeLastAction == nil || v.TimeLastAction.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("TimeLastAction")
	}
	if v.TimeFinished != nil && !v.TimeFinished.IsZero() && (v.TimeStarted == nil || v.TimeStarted.IsZero()) {
		return validation.NewErrRecordIsMissingRequiredField("TimeStarted")
	}

	if v.Settings.MaxVotesPerUser == 0 {
		return validation.NewErrRecordIsMissingRequiredField("settings.maxVotesPerUser")
	}
	retroNodes, err := GetMapOfRetroItemsByID(&RetroItem{Children: v.Items}, make(map[string]*RetroItemTreeNode))
	if err != nil {
		return validation.NewErrBadRecordFieldValue("items", err.Error())
	}
	if v.Stage == "feedback" && len(v.Items) > 0 {
		return validation.NewErrBadRecordFieldValue("items", "retrospective can not hold items while in feedback stage")
	}
	votesByUser := make(map[string]int, len(v.UserIDs))
	for _, retroNode := range retroNodes {
		for userID, votes := range retroNode.item.VotesByUser {
			if !v.HasUserID(userID) {
				return fmt.Errorf("a vote by unknown user ContactID: %v", userID)
			}
			// TODO: Verify userID is in v.UserIDs
			if totalVotes := votesByUser[userID] + votes; totalVotes > v.Settings.MaxVotesPerUser {
				return fmt.Errorf("user '%v' exceeded maximum number of votes: %v", userID, v.Settings.MaxVotesPerUser)
			}
		}
	}

	for memberID, counts := range v.CountsByMemberAndType {
		//var member *dbo4meetingus.MeetingMemberBrief
		for id /*, m*/ := range v.Contacts {
			if dbmodels.SpaceItemID(id).ItemID() == memberID {
				//member = m
				panic("TODO: add space ID validation") // TODO: add team ID validation
			}
		}
		for itemType, count := range counts {
			if !IsKnownItemType(itemType) {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("countsByMemberAndType[%v]", memberID), fmt.Sprintf("unknown item type: %v", itemType))
			}
			if count < 0 {
				return validation.NewErrBadRecordFieldValue(fmt.Sprintf("countsByMemberAndType[%v][%v]", memberID, itemType), fmt.Sprintf("negative value: %v", count))
			}
		}
	}
	return nil
}

// RetroDurations record
type RetroDurations struct {
	Total int `json:"total" firestore:"total"`
}

// RetroProfile record
type RetroProfile struct {
	Title     string         `json:"title" firestore:"title"`
	SpaceID   string         `json:"spaceId" firestore:"spaceId"`
	Durations RetroDurations `json:"durations" firestore:"durations"`
}
