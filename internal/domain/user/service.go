package user

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNameRequired     = errors.New("name is required")
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrUsernameTaken    = errors.New("username already exists")
	ErrEmailTaken       = errors.New("email already exists")
	ErrPhoneTaken       = errors.New("phone already exists")
	ErrWeakPassword     = errors.New("password must be at least 4 characters")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	req CreateUserRequestDto,
) (*UserResponseDto, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Name == "" {
		return nil, ErrNameRequired
	}

	if req.Username == "" {
		return nil, ErrUsernameRequired
	}

	if req.Email == "" {
		return nil, ErrEmailRequired
	}

	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(req.Password) < 4 {
		return nil, ErrWeakPassword
	}

	if !strings.Contains(req.Email, "@") {
		return nil, ErrInvalidEmail
	}

	// Check whether the username already exists.
	_, err := s.repo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, ErrUsernameTaken
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Check whether the email already exists.
	_, err = s.repo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailTaken
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Hash the password before storing it.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Balance:  0,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return &UserResponseDto{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Phone:     newUser.Phone,
		Balance:   formatTaka(newUser.Balance),
		CreatedAt: newUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func formatTaka(paisa int64) string {
	taka := paisa / 100
	remainingPaisa := paisa % 100

	return formatMoney(taka, remainingPaisa)
}

func formatMoney(taka, paisa int64) string {
	sign := ""
	if taka < 0 || paisa < 0 {
		sign = "-"
		taka = abs(taka)
		paisa = abs(paisa)
	}

	return sign + formatInt(taka) + "." + twoDigits(paisa)
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func twoDigits(n int64) string {
	if n < 10 {
		return "0" + formatInt(n)
	}
	return formatInt(n)
}

func formatInt(n int64) string {
	if n == 0 {
		return "0"
	}

	var result string
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
