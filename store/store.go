package store

// IStore - store inferface
type IStore interface {
	AddDataSource(string, interface{})
	GetDataSource(string) interface{}
}

// Store - store struct
type Store struct {
	dataSources map[string]interface{}
}

// New - store factory
func New() Store {
	return Store{
		dataSources: make(map[string]interface{}),
	}
}

// AddDataSource - adds data source to the store
func (s Store) AddDataSource(name string, ds interface{}) {
	s.dataSources[name] = ds
}

// GetDataSource - gets data source to the store
func (s Store) GetDataSource(collection string) interface{} {
	return s.dataSources[collection]
}
