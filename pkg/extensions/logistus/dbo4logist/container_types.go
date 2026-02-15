package dbo4logist

type ContainerType = string

const (
	ContainerType8ft          ContainerType = "8ft"
	ContainerType10ft         ContainerType = "10ft"
	ContainerType20ft         ContainerType = "20ft"
	ContainerType20ftHighCube ContainerType = "20ftHighCube" // High Cube
	ContainerType40ft         ContainerType = "40ft"
	ContainerType40ftHighCube ContainerType = "40ftHighCube" // High Cube
)

var ContainerTypes = []ContainerType{
	ContainerType8ft,
	ContainerType10ft,
	ContainerType20ft,
	ContainerType20ftHighCube,
	ContainerType40ft,
	ContainerType40ftHighCube,
}
