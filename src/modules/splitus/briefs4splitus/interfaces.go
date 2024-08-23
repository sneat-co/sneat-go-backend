package briefs4splitus

type SplitMember interface { // This class is an abstraction for anybot parts of Bill & GroupEntry members
	GetID() string
	GetName() string
	GetShares() int
}
