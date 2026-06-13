package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TxType string
type TxStatus string

const (
	TxDeposit     TxType = "deposit"
	TxWithdrawal  TxType = "withdrawal"
	TxTransferIn  TxType = "transfer_in"
	TxTransferOut TxType = "transfer_out"

	TxSuccess  TxStatus = "success"
	TxFailed   TxStatus = "failed"
	TxReversed TxStatus = "reversed"
)

type Transaction struct {
	ID           uuid.UUID
	AccountID    uuid.UUID
	Type         TxType
	Amount       decimal.Decimal
	BalanceAfter decimal.Decimal
	Reference    string
	RelatedTxID  *uuid.UUID
	Status       TxStatus
	CreatedAt    time.Time
}
