package db

type Surreal interface {
	Save(i interface{}) error
}

type SurrealDB struct {
}

func (s *SurrealDB) Save(i interface{}) error {
	return nil
}
