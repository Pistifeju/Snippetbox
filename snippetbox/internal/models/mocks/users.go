package mocks

import (
	"github.com/Pistieju/snippetbox/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa55word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {

	switch id {
	case 1:
		return &models.User{
			ID:      1,
			Name:    "Alice",
			Email:   "alice@example.com",
			Created: time.Date(2020, 01, 01, 12, 0, 0, 0, time.UTC),
		}, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	switch id {
	case 1:
		return nil
	default:
		return models.ErrNoRecord
	}
}
