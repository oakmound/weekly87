package enemies

type enemyType int

const (
	Hare   enemyType = iota
	Mantis enemyType = iota
	TypeLimit
)

var enemyTypeList = [TypeLimit]enemyType{
	Hare,
	Mantis,
}

func Init() {
	initHare()
	initMantis()
}
