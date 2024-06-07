package models

type SplitMember interface { // This class is an abstraction for common parts of Bill & GroupEntry members
	GetID() string
	GetName() string
	GetShares() int
}
