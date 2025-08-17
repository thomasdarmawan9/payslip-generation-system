package model

import "time"

// Satu run payroll per AttendancePeriod (unik)
type PayrollRun struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	PeriodID  uint      `gorm:"uniqueIndex;not null"` // unique => setiap period hanya 1x run
	RunAt     time.Time `gorm:"type:timestamp;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now()"`
	// Optional: status, metadata, etc.
}

func (PayrollRun) TableName() string { return "payroll_runs" }

// Snapshot per karyawan (agar perubahan data setelah run tidak mengubah payslip)
type PayrollItem struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement"`
	PayrollRunID       uint      `gorm:"index;not null"`
	UserID             uint      `gorm:"index;not null"`
	SnapshotSalary     float64   `gorm:"type:numeric(12,2);not null"` // gaji bulanan saat run
	WorkingDays        int       `gorm:"not null"`                    // hari kerja (weekday) dalam period
	AttendanceDays     int       `gorm:"not null"`                    // jumlah hadir
	WorkingHours       int       `gorm:"not null"`                    // WorkingDays * 8
	AttendanceHours    int       `gorm:"not null"`                    // AttendanceDays * 8
	OvertimeHours      float64   `gorm:"type:numeric(6,2);not null"`  // total jam lembur
	BasePay            float64   `gorm:"type:numeric(14,2);not null"` // prorate
	OvertimePay        float64   `gorm:"type:numeric(14,2);not null"` // 2x hourly * hours
	ReimbursementTotal float64   `gorm:"type:numeric(14,2);not null"`
	GrandTotal         float64   `gorm:"type:numeric(14,2);not null"`
	CreatedAt          time.Time `gorm:"type:timestamp;default:now()"`
	UpdatedAt          time.Time `gorm:"type:timestamp;default:now()"`
}

func (PayrollItem) TableName() string { return "payroll_items" }
