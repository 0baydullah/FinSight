package category

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

var (
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryNameTaken    = errors.New("category name already exists")
	ErrUnauthorizedCategory = errors.New("unauthorized category access")
)

func (s *Service) Create(
	ctx context.Context,
	userID uint,
	req CreateCategoryRequestDto,
) (*CategoryResponseDto, error) {

	name := strings.TrimSpace(req.Name)

	if name == "" {
		return nil, ErrCategoryNameRequired
	}

	categories, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, category := range categories {
		if strings.EqualFold(category.Name, name) {
			return nil, ErrCategoryNameTaken
		}
	}

	category := &Category{
		UserID:   userID,
		Name:     name,
		IsSystem: false,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return toCategoryResponse(category), nil
}

func (s *Service) GetAll(
	ctx context.Context,
	userID uint,
) ([]CategoryResponseDto, error) {

	categories, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryResponseDto, 0, len(categories))

	for _, category := range categories {
		result = append(result, *toCategoryResponse(&category))
	}

	return result, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	userID uint,
	categoryID uint,
) (*CategoryResponseDto, error) {

	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	if category.UserID != userID {
		return nil, ErrUnauthorizedCategory
	}

	return toCategoryResponse(category), nil
}

func (s *Service) Update(
	ctx context.Context,
	userID uint,
	categoryID uint,
	req UpdateCategoryRequestDto,
) (*CategoryResponseDto, error) {

	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}

		return nil, err
	}

	if category.UserID != userID {
		return nil, ErrUnauthorizedCategory
	}

	name := strings.TrimSpace(req.Name)

	if name == "" {
		return nil, ErrCategoryNameRequired
	}

	if !strings.EqualFold(category.Name, name) {
		categories, err := s.repo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		for _, existing := range categories {
			if existing.ID != category.ID &&
				strings.EqualFold(existing.Name, name) {
				return nil, ErrCategoryNameTaken
			}
		}
	}

	category.Name = name

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return toCategoryResponse(category), nil
}

func (s *Service) Delete(
	ctx context.Context,
	userID uint,
	categoryID uint,
) error {

	category, err := s.repo.GetByID(ctx, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}

		return err
	}

	if category.UserID != userID {
		return ErrUnauthorizedCategory
	}

	return s.repo.Delete(ctx, category)
}

func toCategoryResponse(category *Category) *CategoryResponseDto {
	return &CategoryResponseDto{
		ID:        category.ID,
		Name:      category.Name,
		IsSystem:  category.IsSystem,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
	}
}
