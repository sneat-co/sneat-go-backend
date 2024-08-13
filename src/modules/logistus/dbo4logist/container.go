package dbo4logist

import (
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/random"
	"github.com/strongo/slice"
	"github.com/strongo/strongoapp/with"
	"github.com/strongo/validation"
	"strings"
)

// OrderContainerBase is a base type for OrderContainer
type OrderContainerBase struct {
	Type   ContainerType `json:"type" firestore:"type"`
	Number string        `json:"number,omitempty" firestore:"number,omitempty"`
	with.FlagsField
	with.TagsField
	Instructions string `json:"instructions,omitempty" firestore:"instructions,omitempty"`
	//TotalLoad   *FreightLoad `json:"totalLoad,omitempty" firestore:"totalLoad,omitempty"`
	//TotalUnload *FreightLoad `json:"totalLoad,omitempty" firestore:"totalLoad,omitempty"`
}

// String returns string representation of the container
func (v OrderContainerBase) String() string {
	return fmt.Sprintf("{ExtraType=%s,Number=%s}", v.Type, v.Number)
}

// Validate returns nil if valid, or error if not
func (v OrderContainerBase) Validate() error {
	switch v.Type {
	case "unknown":
	case "":
		return validation.NewErrRecordIsMissingRequiredField("type")
	default:
		if v.Type != strings.TrimSpace(v.Type) {
			return validation.NewErrBadRecordFieldValue("type", "leading or closing spaces")
		}
		if slice.Index(ContainerTypes, v.Type) < 0 {
			return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("unknown value: [%v]", v.Type))
		}
	}
	if v.Number != strings.TrimSpace(v.Number) {
		return validation.NewErrBadRecordFieldValue("numbers", "leading or closing spaces")
	}
	//if err := v.FreightLoad.Validate(); err != nil {
	//	return err // don't wrap the error here
	//}
	if err := v.FlagsField.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	return nil
}

// OrderContainer is a container for an order
type OrderContainer struct {
	ID string `json:"id" firestore:"id"`
	OrderContainerBase
}

// String returns string representation of the OrderContainer
func (v OrderContainer) String() string {
	return fmt.Sprintf("OrderContainer{ContactID=%s,ExtraType=%s,Number=%s}", v.ID, v.Type, v.Number)
}

func validateContainerID(field, id string) error {
	if strings.TrimSpace(id) == "" {
		return validation.NewErrRecordIsMissingRequiredField(field)
	}
	if err := validate.RecordID(id); err != nil {
		return validation.NewErrBadRecordFieldValue(field, err.Error())
	}
	return nil
}

// Validate returns nil if valid, or error if not
func (v OrderContainer) Validate() error {
	if err := validateContainerID("id", v.ID); err != nil {
		return err
	}
	if err := v.OrderContainerBase.Validate(); err != nil {
		return err
	}
	return nil
}

// WithOrderContainers is a type that has order containers
type WithOrderContainers struct {
	Containers []*OrderContainer `json:"containers,omitempty" firestore:"containers,omitempty"`
}

// RemoveContainer removes container by ContactID
func (v WithOrderContainers) RemoveContainer(id string) (containers []*OrderContainer, found bool) {
	i, _ := v.GetContainerByID(id)
	if i >= 0 {
		return append(v.Containers[:i], v.Containers[i+1:]...), true
	}
	return v.Containers, false
}

// GetContainerIDs returns IDs of containers
func (v WithOrderContainers) GetContainerIDs() (containerIDs []string) {
	containerIDs = make([]string, len(v.Containers))
	for i, c := range v.Containers {
		containerIDs[i] = c.ID
	}
	return containerIDs
}

// GenerateRandomContainerID generates random container ContactID that is not used in the list of containers
func (v WithOrderContainers) GenerateRandomContainerID() string {
	var attempt int
	for {
		id := random.ID(2)
		for _, c := range v.Containers {
			if c.ID == id {
				if attempt++; attempt == 100 {
					panic("too many attempts to generate random container ContactID")
				}
				continue
			}
		}
		return id
	}
}

// GetContainerByID returns container by ContactID
func (v WithOrderContainers) GetContainerByID(id string) (i int, container *OrderContainer) {
	for i, c := range v.Containers {
		if c.ID == id {
			return i, c
		}
	}
	return -1, nil
}

//func (v WithOrderContainers) UpdateContainerTotals(containerPoints []ContainerPoint) {
//	for _, c := range v.Containers {
//		c.NumberOfPallets = 0
//		c.GrossWeightKg = 0
//		c.VolumeM3 = 0
//	}
//	for _, cp := range containerPoints {
//		_, container := v.GetContainerByID(cp.ContainerID)
//		if !cp.ToLoad.IsEmpty() {
//			container.NumberOfPallets += cp.ToLoad.NumberOfPallets
//			container.GrossWeightKg += cp.ToLoad.GrossWeightKg
//			container.VolumeM3 += cp.ToLoad.VolumeM3
//		}
//		if cp.ToUnload.IsEmpty() {
//			container.NumberOfPallets -= cp.ToUnload.NumberOfPallets
//			container.GrossWeightKg -= cp.ToUnload.GrossWeightKg
//			container.VolumeM3 -= cp.ToUnload.VolumeM3
//		}
//	}
//}

func (v WithOrderContainers) validateOrder(order OrderDbo) error {
	if err := v.Validate(); err != nil {
		return err
	}
	for _, container := range v.Containers {
		var freightLoad FreightLoad
		for _, cp := range order.ContainerPoints {
			if cp.ContainerID != container.ID {
				continue
			}
			if cp.ToLoad != nil {
				freightLoad.NumberOfPallets += cp.ToLoad.NumberOfPallets
				freightLoad.GrossWeightKg += cp.ToLoad.GrossWeightKg
				freightLoad.VolumeM3 += cp.ToLoad.VolumeM3
			}
			if cp.ToUnload != nil {
				freightLoad.NumberOfPallets -= cp.ToUnload.NumberOfPallets
				freightLoad.GrossWeightKg -= cp.ToUnload.GrossWeightKg
				freightLoad.VolumeM3 -= cp.ToUnload.VolumeM3
			}
		}
		//validateContainerTotal := func(field string, expected, actual int) error {
		//	// TODO: fix this or remove
		//	if expected != actual {
		//		return validation.NewErrBadRecordFieldValue(fmt.Sprintf("containers[%d].%s", i, field),
		//			fmt.Sprintf("does not match sum of container points: expected %v, actual: %v", expected, actual))
		//	}
		//	return nil
		//}
		//if err := validateContainerTotal("numberOfPallets", freightLoad.NumberOfPallets, container.NumberOfPallets); err != nil {
		//	return err
		//}
		//if err := validateContainerTotal("grossWeightKg", freightLoad.GrossWeightKg, container.GrossWeightKg); err != nil {
		//	return err
		//}
		//if err := validateContainerTotal("volumeM3", freightLoad.VolumeM3, container.VolumeM3); err != nil {
		//	return err
		//}
	}
	return nil
}

// Validate returns nil if valid, or error if not
func (v WithOrderContainers) Validate() error {
	for i, c := range v.Containers {
		if err := c.Validate(); err != nil {
			return validation.NewErrBadRecordFieldValue(fmt.Sprintf("containers[%v]", i), err.Error())
		}
		for j, c2 := range v.Containers {
			if j == i {
				continue
			}
			if c.ID == c2.ID {
				return validation.NewErrBadRecordFieldValue("containers",
					fmt.Sprintf("containers[%v].id == containers[%v].id: %s", i, j, c.ID))
			}
			if c.Number != "" && c2.Number == c.Number {
				return validation.NewErrBadRecordFieldValue("containers",
					fmt.Sprintf("containers[%v,id=%s].number == containers[%v,id=%s].number: %s", i, c.ID, j, c2.ID, c.Number))
			}
		}
	}
	return nil
}

// Updates returns updates for the order containers
func (v WithOrderContainers) Updates() []dal.Update {
	if len(v.Containers) == 0 {
		return []dal.Update{
			{Field: "containers", Value: dal.DeleteField},
		}
	}
	return []dal.Update{
		{Field: "containers", Value: v.Containers},
	}
}
