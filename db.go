package jsondb

type Client[T any] interface {
	Write(id string, data T) error
	Read(id string) (T, error)
	List() ([]T, error)
	Delete(id string) error
}

type NotFoundError struct {
	id string
}

func (n NotFoundError) Error() string {
	return n.id + " not found"
}
