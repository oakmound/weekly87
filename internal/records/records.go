package records

// Records serves as our safe file and all variables we track across
// multiple runs
type Records struct {
	SectionsCleared int64 `json:"sectionsCleared"`
	BaseSeed        int64 `json:"baseSeed"`
	// Todo: more
	FarthestGoneInSections int64 `json:"farthestGoneInSections"`
}

func (r *Records) Store() {
	// Todo
}
func (r *Records) Load() {
	// Todo
}
