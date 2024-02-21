package jsondb

type Object interface {
    ID() string
}


type Client[T Object] interface {
    Write(data T) error
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
