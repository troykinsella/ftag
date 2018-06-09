package tagmap

type Store interface {
	Load() (*TM, error)
	Put(tm *TM) error
}
