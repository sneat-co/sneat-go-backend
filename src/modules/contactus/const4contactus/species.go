package const4contactus

import "github.com/strongo/slice"

type PetKind = string

const (
	PetKindDog     PetKind = "dog"
	PetKindCat     PetKind = "cat"
	PetKindHamster PetKind = "hamster"
	PetKindRabbit  PetKind = "rabbit"
	PetKindBird    PetKind = "bird"
	PetKindFish    PetKind = "fish"
	PetKindTurtle  PetKind = "turtle"
	PetKindSnake   PetKind = "snake"
	PetKindLizard  PetKind = "lizard"
	PetKindHorse   PetKind = "horse"
	PetKindPig     PetKind = "pig"
	PetKindCow     PetKind = "cow"
	PetKindRat     PetKind = "rat"
	PetKindMouse   PetKind = "mouse"
	PetKindGoat    PetKind = "goat"
	PetKindSheep   PetKind = "sheep"
	PetKindOther   PetKind = "other"
)

var KnownPetKinds = []PetKind{
	PetKindDog,
	PetKindCat,
	PetKindCow,
	PetKindFish,
	PetKindBird,
	PetKindHamster,
	PetKindLizard,
	PetKindGoat,
	PetKindHorse,
	PetKindPig,
	PetKindRat,
	PetKindMouse,
	PetKindRabbit,
	PetKindSheep,
	PetKindSnake,
	PetKindTurtle,
	PetKindOther,
}

func IsKnownPetPetKind(v PetKind) bool {
	return slice.Contains(KnownPetKinds, v)
}
