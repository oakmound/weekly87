package run

type Persistable interface {
	ShouldPersist() bool
}
