// internal/usecase/test/mock_txmanager.go
package test

import (
	"context"
	repoTx "payslip-generation-system/internal/repository/tx"
)

type FakeTxManager struct{}

func (f FakeTxManager) Begin(ctx context.Context) (context.Context, error) { return ctx, nil }
func (f FakeTxManager) Commit(ctx context.Context) error                   { return nil }
func (f FakeTxManager) Rollback(ctx context.Context) error                 { return nil }

var _ repoTx.TxManager = (*FakeTxManager)(nil)
