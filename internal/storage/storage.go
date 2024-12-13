package storage

type Client interface {
	AddUser() error
	User() (int, error)
	IsAdmin() (bool, error)
}

type Storage struct {
	client Client
}

func New(client Client) Storage {
	return Storage{client: client}
}

func (Storage) AddUser(email, password string) error {
	return nil
}

func (Storage) User() (int, error) {
	return 0, nil
}

func (Storage) IsAdmin() (bool, error) {
	return false, nil
}
