package id

type Loader interface {
	// load disctinct ids
	Load()
	GetData() []interface{}
}
