package usecase

import (
	atRepo "payslip-generation-system/internal/repository/attendance"
	apRepo "payslip-generation-system/internal/repository/attendanceperiod"
	otRepo "payslip-generation-system/internal/repository/overtime"
	payRepo "payslip-generation-system/internal/repository/payroll"
	rbRepo "payslip-generation-system/internal/repository/reimbursement"
	repoTx "payslip-generation-system/internal/repository/tx"
)

// NewForTest creates a blank usecase instance for unit tests.
func NewForTest() IUsecase { return &usecase{} }

// InjectForTest wires mock dependencies into a test instance created by NewForTest.
func InjectForTest(target IUsecase,
	ap apRepo.Repo,
	at atRepo.Repo,
	ot otRepo.Repo,
	rb rbRepo.Repo,
	pay payRepo.Repo,
	tx repoTx.TxManager,
) {
	if u, ok := target.(*usecase); ok {
		u.apRepo = ap
		u.atRepo = at
		u.otRepo = ot
		u.rbRepo = rb
		u.payrollRepo = pay
		u.txManager = tx
	}
}
