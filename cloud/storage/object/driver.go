package object

type Driver interface {
	Put(key string, data []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
}
