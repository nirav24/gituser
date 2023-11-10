package user

import (
	"encoding/json"
	"os"
)

type Store struct {
	path       string
	localUsers []User
}

// NewStore ...
func NewStore(path string) (*Store, error) {
	s := &Store{
		path:       path,
		localUsers: make([]User, 0),
	}

	err := s.loadUsers()
	return s, err
}

// Close saves local changes to file
func (s *Store) Close() error {
	return s.saveLocalUsers()
}
func (s *Store) GetUsers() []User {
	return s.localUsers
}

func (s *Store) AddUser(user User) error {
	s.localUsers = append(s.localUsers, user)
	return nil
}

func (s *Store) loadUsers() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &s.localUsers)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) RemoveUser(user User) {
	newUsers := make([]User, 0)
	for _, localUser := range s.localUsers {
		if user != localUser {
			newUsers = append(newUsers, user)
		}
	}
	s.localUsers = newUsers
}

func (s *Store) saveLocalUsers() error {
	data, err := json.Marshal(s.localUsers)
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0533)
}
