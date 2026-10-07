package transaction

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0baydullah/FinSight/internal/domain/category"
	"gorm.io/gorm"
)

type Service struct {
	repo         Repository
	categoryRepo category.Repository
}

func NewService(
	repo Repository,
	categoryRepo category.Repository,
) *Service {
	return &Service{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

var (
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrTransactionAmountInvalid = errors.New("transaction amount is invalid")
	ErrTransactionTypeInvalid   = errors.New("transaction type is invalid")
	ErrTransactionCategory      = errors.New("invalid category")
	ErrUnauthorizedTransaction  = errors.New("unauthorized transaction")
)

const (
	TypeIncome  = "income"
	TypeExpense = "expense"
)

func (s *Service) Create(
	ctx context.Context,
	userID uint,
	req CreateTransactionRequestDto,
) (*TransactionResponseDto, error) {

	if req.CategoryID == 0 {
		return nil, ErrTransactionCategory
	}

	selectedCategory, err := s.categoryRepo.GetByID(
		ctx,
		req.CategoryID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionCategory
		}

		return nil, err
	}

	if selectedCategory.UserID != userID {
		return nil, ErrTransactionCategory
	}

	if req.Type != TypeIncome && req.Type != TypeExpense {
		return nil, ErrTransactionTypeInvalid
	}

	amount, err := parseAmount(req.Amount)
	if err != nil || amount <= 0 {
		return nil, ErrTransactionAmountInvalid
	}

	transactionDate, err := time.Parse(time.RFC3339, req.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction date: %w", err)
	}

	description := strings.TrimSpace(req.Description)

	transaction := &Transaction{
		UserID:          userID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          amount,
		TransactionDate: transactionDate,
	}

	if description != "" {
		transaction.Description = &description
	}

	balanceDelta := amount

	if req.Type == TypeExpense {
		balanceDelta = -amount
	}

	if err := s.repo.Create(ctx, transaction, balanceDelta); err != nil {
		return nil, err
	}

	return toTransactionResponse(transaction), nil
}

func (s *Service) GetAll(
	ctx context.Context,
	userID uint,
) ([]TransactionResponseDto, error) {

	transactions, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]TransactionResponseDto, 0, len(transactions))

	for _, transaction := range transactions {
		result = append(result, *toTransactionResponse(&transaction))
	}

	return result, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	userID uint,
	transactionID uint,
) (*TransactionResponseDto, error) {

	transaction, err := s.repo.GetByID(ctx, transactionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}

		return nil, err
	}

	if transaction.UserID != userID {
		return nil, ErrUnauthorizedTransaction
	}

	return toTransactionResponse(transaction), nil
}

func parseAmount(value string) (int64, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return 0, ErrTransactionAmountInvalid
	}

	parts := strings.Split(value, ".")

	if len(parts) > 2 {
		return 0, ErrTransactionAmountInvalid
	}

	whole, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || whole < 0 {
		return 0, ErrTransactionAmountInvalid
	}

	var fraction int64

	if len(parts) == 2 {
		fractionPart := parts[1]

		if len(fractionPart) > 2 {
			return 0, ErrTransactionAmountInvalid
		}

		if len(fractionPart) == 1 {
			fractionPart += "0"
		}

		if fractionPart != "" {
			fraction, err = strconv.ParseInt(fractionPart, 10, 64)
			if err != nil {
				return 0, ErrTransactionAmountInvalid
			}
		}
	}

	return whole*100 + fraction, nil
}

func (s *Service) Update(
	ctx context.Context,
	userID uint,
	transactionID uint,
	req UpdateTransactionRequestDto,
) (*TransactionResponseDto, error) {

	transaction, err := s.repo.GetByID(ctx, transactionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}

		return nil, err
	}

	if transaction.UserID != userID {
		return nil, ErrUnauthorizedTransaction
	}

	oldBalanceDelta := transactionBalanceDelta(
		transaction.Type,
		transaction.Amount,
	)

	if req.CategoryID != nil {
		if *req.CategoryID == 0 {
			return nil, ErrTransactionCategory
		}

		selectedCategory, err := s.categoryRepo.GetByID(
			ctx,
			*req.CategoryID,
		)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrTransactionCategory
			}

			return nil, err
		}

		if selectedCategory.UserID != userID {
			return nil, ErrTransactionCategory
		}

		transaction.CategoryID = *req.CategoryID
	}

	if req.Type != nil {
		if *req.Type != TypeIncome &&
			*req.Type != TypeExpense {
			return nil, ErrTransactionTypeInvalid
		}

		transaction.Type = *req.Type
	}

	if req.Amount != nil {
		amount, err := parseAmount(*req.Amount)

		if err != nil || amount <= 0 {
			return nil, ErrTransactionAmountInvalid
		}

		transaction.Amount = amount
	}

	if req.Description != nil {
		description := strings.TrimSpace(*req.Description)

		if description == "" {
			transaction.Description = nil
		} else {
			transaction.Description = &description
		}
	}

	if req.TransactionDate != nil {
		transactionDate, err := time.Parse(
			time.RFC3339,
			*req.TransactionDate,
		)

		if err != nil {
			return nil, err
		}

		transaction.TransactionDate = transactionDate
	}

	newBalanceDelta := transactionBalanceDelta(
		transaction.Type,
		transaction.Amount,
	)

	balanceDelta := newBalanceDelta - oldBalanceDelta

	if err := s.repo.Update(
		ctx,
		transaction,
		balanceDelta,
	); err != nil {
		return nil, err
	}

	return toTransactionResponse(transaction), nil
}

func (s *Service) Delete(
	ctx context.Context,
	userID uint,
	transactionID uint,
) error {

	transaction, err := s.repo.GetByID(ctx, transactionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTransactionNotFound
		}

		return err
	}

	if transaction.UserID != userID {
		return ErrUnauthorizedTransaction
	}

	oldBalanceDelta := transactionBalanceDelta(
		transaction.Type,
		transaction.Amount,
	)

	return s.repo.Delete(
		ctx,
		transaction,
		-oldBalanceDelta,
	)
}

func transactionBalanceDelta(
	transactionType string,
	amount int64,
) int64 {

	if transactionType == TypeExpense {
		return -amount
	}

	return amount
}

func toTransactionResponse(
	transaction *Transaction,
) *TransactionResponseDto {

	description := ""

	if transaction.Description != nil {
		description = *transaction.Description
	}

	return &TransactionResponseDto{
		ID:              transaction.ID,
		CategoryID:      transaction.CategoryID,
		Type:            transaction.Type,
		Amount:          formatAmount(transaction.Amount),
		Description:     description,
		TransactionDate: transaction.TransactionDate.Format(time.RFC3339),
		CreatedAt:       transaction.CreatedAt.Format(time.RFC3339),
	}
}

func formatAmount(amount int64) string {
	return fmt.Sprintf("%d.%02d", amount/100, amount%100)
}
