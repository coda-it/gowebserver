package store

type IStore interface {
	Add(string, string, string)
	Get(string, string) string
}

type IDataSource interface {
	Add(string, string)
	Get(string) string
}

type Store struct {
	dataSources	map[string]IDataSource
}

func New() Store {
	return Store{}
}

func (s *Store) AddDataSource(name string, ds IDataSource) {
	s.dataSources[name] = ds
}

func (s *Store) Add(collection string, key string, value string) {
	s.dataSources[collection].Add(key, value)
}

func (s* Store) Get(collection string, key string) string {
	return s.dataSources[collection].Get(key)
}