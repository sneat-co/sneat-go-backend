package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func cp(containerID, shippingPointID string) *ContainerPoint {
	return &ContainerPoint{ContainerID: containerID, ShippingPointID: shippingPointID}
}

func TestWithContainerPoints_Validate_panics(t *testing.T) {
	v := &WithContainerPoints{}
	assert.Panics(t, func() { _ = v.Validate() })
}

func TestWithContainerPoints_Updates(t *testing.T) {
	empty := (&WithContainerPoints{}).Updates()
	assert.Len(t, empty, 1)
	assert.Equal(t, "containerPoints", empty[0].FieldName())

	withData := (&WithContainerPoints{ContainerPoints: []*ContainerPoint{cp("c1", "sp1")}}).Updates()
	assert.Len(t, withData, 1)
	assert.Equal(t, "containerPoints", withData[0].FieldName())
}

func TestWithContainerPoints_GetContainerPoint(t *testing.T) {
	p := cp("c1", "sp1")
	v := &WithContainerPoints{ContainerPoints: []*ContainerPoint{p, cp("c2", "sp2")}}

	assert.Equal(t, p, v.GetContainerPoint("c1", "sp1"))
	assert.Nil(t, v.GetContainerPoint("c1", "spX"))
	assert.Nil(t, v.GetContainerPoint("cX", "sp1"))

	assert.Panics(t, func() { v.GetContainerPoint("", "sp1") })
	assert.Panics(t, func() { v.GetContainerPoint("c1", "") })
}

func TestWithContainerPoints_RemoveContainerPointsByShippingPointID(t *testing.T) {
	v := &WithContainerPoints{ContainerPoints: []*ContainerPoint{cp("c1", "sp1"), cp("c2", "sp2")}}
	remaining := v.RemoveContainerPointsByShippingPointID("sp1")
	assert.Len(t, remaining, 1)
	assert.Equal(t, "sp2", remaining[0].ShippingPointID)
}

func TestWithContainerPoints_RemoveContainerPointsByContainerID(t *testing.T) {
	v := &WithContainerPoints{ContainerPoints: []*ContainerPoint{cp("c1", "sp1"), cp("c2", "sp2")}}
	remaining := v.RemoveContainerPointsByContainerID("c1")
	assert.Len(t, remaining, 1)
	assert.Equal(t, "c2", remaining[0].ContainerID)
}
