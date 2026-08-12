package url

import (
	"errors"
	"fmt"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-playground/validator/v10"
)

type Storage interface {
	SaveURL(urlToSave, alias string) (int64, error)
	GetURL(alias string) (string, error)
	UpdateURL(alias, newURL string) (int64, error)
	DeleteURL(alias string) (int64, error)
}

type Service struct {
	storage     Storage
	aliasLength int
}

func NewService(storage Storage, aliasLength int) *Service {
	return &Service{
		storage:     storage,
		aliasLength: aliasLength,
	}
}

var (
	ErrInvalidURL  = errors.New("invalid url")
	ErrURLExists   = errors.New("url already exists")
	ErrURLNotFound = errors.New("url not found")
)

func (s *Service) SaveURL(urlToSave, alias string) (string, error) {
	if err := validator.New().Var(urlToSave, "required,url"); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if alias == "" {
		alias = random.NewRandomString(s.aliasLength)
	}
	_, err := s.storage.SaveURL(urlToSave, alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLExists) {
			return "", ErrURLExists
		}
		return "", fmt.Errorf("failed to save url: %w", err)
	}
	return alias, nil
}

func (s *Service) GetURL(alias string) (string, error) {
	url, err := s.storage.GetURL(alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return "", ErrURLNotFound
		}

		return "", fmt.Errorf("failed to get url: %w", err)
	}

	return url, nil
}

func (s *Service) UpdateURL(alias, newURL string) (int64, error) {
	if err := validator.New().Var(newURL, "required,url"); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	rows, err := s.storage.UpdateURL(alias, newURL)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return 0, ErrURLNotFound
		}
		return 0, fmt.Errorf("failed to update url: %w", err)
	}
	return rows, nil
}

func (s *Service) DeleteURL(alias string) (int64, error) {
	rows, err := s.storage.DeleteURL(alias)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			return 0, ErrURLNotFound
		}
		return 0, fmt.Errorf("failed to delete url: %w", err)
	}
	return rows, nil
}
