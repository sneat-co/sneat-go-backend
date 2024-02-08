package models

type SplitMember interface { // This class is an abstraction for common parts of Bill & Group members
	GetID() string
	GetName() string
	GetShares() int
}
