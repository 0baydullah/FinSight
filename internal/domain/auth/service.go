package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/0baydullah/FinSight/internal/domain/user"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIdentifierRequired = errors.New("username or email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidCredentials = errors.New("invalid username/email or password")
	ErrJWTConfiguration   = errors.New("JWT configuration is missing")
)

type Service struct {
	userRepo  user.Repository
	jwtSecret []byte
	expiresIn time.Duration
}

func NewService(
	userRepo user.Repository,
	jwtSecret string,
	expiresIn time.Duration,
) *Service {
	return &Service{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		expiresIn: expiresIn,
	}
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *Service) Login(ctx context.Context, req LoginRequestDto) (*LoginResponseDto, error) {
	req.Identifier = strings.TrimSpace(req.Identifier)

	if req.Identifier == "" {
		return nil, ErrIdentifierRequired
	}

	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(s.jwtSecret) == 0 || s.expiresIn <= 0 {
		return nil, ErrJWTConfiguration
	}

	var (
		foundUser *user.User
		err       error
	)

	// An identifier containing @ is treated as an email.
	if strings.Contains(req.Identifier, "@") {
		foundUser, err = s.userRepo.GetByEmail(
			ctx,
			strings.ToLower(req.Identifier),
		)
	} else {
		foundUser, err = s.userRepo.GetByUsername(
			ctx,
			req.Identifier,
		)
	}

	if errors.Is(err, user.ErrUserNotFound) {
		return nil, ErrInvalidCredentials
	}

	if err != nil {
		return nil, err
	}

	// Compare the submitted password with the stored bcrypt hash.
	if err := bcrypt.CompareHashAndPassword(
		[]byte(foundUser.Password),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(s.expiresIn)

	claims := Claims{
		UserID:   foundUser.ID,
		Username: foundUser.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user:" + foundUser.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResponseDto{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.Format(time.RFC3339),
		User: LoginUserDto{
			ID:       foundUser.ID,
			Name:     foundUser.Name,
			Username: foundUser.Username,
			Email:    foundUser.Email,
		},
	}, nil
}
