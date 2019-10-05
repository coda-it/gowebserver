package store

type IStore interface {
	AddDataSource(string, interface{})
	GetDataSource(string) interface{}
}

type Store struct {
	dataSources map[string]interface{}
}

func New() Store {
	return Store{
		dataSources: make(map[string]interface{}),
	}
}

func (s Store) AddDataSource(name string, ds interface{}) {
	s.dataSources[name] = ds
}

func (s Store) GetDataSource(collection string) interface{} {
	return s.dataSources[collection]
}
