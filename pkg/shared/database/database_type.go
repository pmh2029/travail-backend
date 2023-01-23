package database

type DBType int

const (
	MySQL DBType = iota
	PostgreSQL
)

func (dbType DBType) String() string {
	names := [...]string{
		"mysql",
		"postgres",
	}
	if dbType < MySQL || dbType > PostgreSQL {
		return "Unknown"
	}
	return names[dbType]
}
