package characters

var (

	// CurrentParty is the party you are playing with
	CurrentParty = NewStartingParty()
)

// Party is a construct to maintain stats about a party
type Party struct {
	Jobs    []int
	MaxSize int
}

// NewStartingParty creates the basic intro party
func NewStartingParty() *Party {
	p := &Party{}
	p.Jobs = []int{JobSwordsman}
	p.MaxSize = 2
	return p
}
