package transaction

import "time"

type TransactionFilter struct {
    UserID      uint
    Type        string
    CategoryIDs []uint
    From        *time.Time
    To          *time.Time
    Search      string
    Page        int
    Limit       int
}

type TransactionPage struct {
    Items      []Transaction
    Total      int64
    Page       int
    Limit      int
    TotalPages int
}