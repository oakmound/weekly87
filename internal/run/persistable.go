package run

// Persistable defines characters that should stay on the map
// after their section has been passed
type Persistable interface {
	ShouldPersist() bool
}
