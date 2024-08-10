package models

type SplitMember interface { // This class is an abstraction for shared parts of Bill & GroupEntry members
	GetID() string
	GetName() string
	GetShares() int
}
