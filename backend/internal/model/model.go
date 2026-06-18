package model

import "time"

type SysUser struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	RealName  string    `gorm:"size:50;not null" json:"realName"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Role      int       `gorm:"not null;default:1" json:"role"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (SysUser) TableName() string { return "sys_user" }

type WorkTeam struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamName  string    `gorm:"size:100;not null" json:"teamName"`
	TeamCode  string    `gorm:"uniqueIndex;size:50;not null" json:"teamCode"`
	LeaderID  *uint64   `gorm:"" json:"leaderId"`
	Status    int       `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (WorkTeam) TableName() string { return "work_team" }

type TeamMember struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID    uint64    `gorm:"not null;uniqueIndex:uk_team_user" json:"teamId"`
	UserID    uint64    `gorm:"not null;uniqueIndex:uk_team_user" json:"userId"`
	JoinDate  string    `gorm:"type:date;not null" json:"joinDate"`
	LeaveDate *string   `gorm:"type:date" json:"leaveDate"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (TeamMember) TableName() string { return "team_member" }

type Product struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductCode string    `gorm:"uniqueIndex;size:50;not null" json:"productCode"`
	ProductName string    `gorm:"size:100;not null" json:"productName"`
	Spec        string    `gorm:"size:200" json:"spec"`
	Status      int       `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (Product) TableName() string { return "product" }

type ProcessStep struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProcessCode string    `gorm:"uniqueIndex;size:50;not null" json:"processCode"`
	ProcessName string    `gorm:"size:100;not null" json:"processName"`
	ProductID   uint64    `gorm:"not null;index" json:"productId"`
	Difficulty  int       `gorm:"not null;default:1" json:"difficulty"`
	Description string    `gorm:"size:500" json:"description"`
	IsShared    int       `gorm:"not null;default:0" json:"isShared"`
	Status      int       `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ProcessStep) TableName() string { return "process_step" }

type ProcessPrice struct {
	ID            uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProcessID     uint64      `gorm:"not null;uniqueIndex:uk_process_version" json:"processId"`
	VersionNo     int         `gorm:"not null;uniqueIndex:uk_process_version" json:"versionNo"`
	GradeLevel    string      `gorm:"size:20;not null;default:'STD';uniqueIndex:uk_process_version" json:"gradeLevel"`
	UnitPrice     float64     `gorm:"type:decimal(12,4);not null" json:"unitPrice"`
	EffectiveDate string      `gorm:"type:date;not null" json:"effectiveDate"`
	ExpiryDate    *string     `gorm:"type:date" json:"expiryDate"`
	Remark        string      `gorm:"size:500" json:"remark"`
	CreatedBy     uint64      `gorm:"not null" json:"createdBy"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	Process       *ProcessStep `gorm:"foreignKey:ProcessID" json:"process,omitempty"`
}

func (ProcessPrice) TableName() string { return "process_price" }

type ProductionReport struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReportNo     string    `gorm:"uniqueIndex;size:32;not null" json:"reportNo"`
	WorkerID     uint64    `gorm:"not null;index:idx_worker_date" json:"workerId"`
	TeamID       *uint64   `gorm:"index:idx_team_date" json:"teamId"`
	ProcessID    uint64    `gorm:"not null;index" json:"processId"`
	PriceID      uint64    `gorm:"not null" json:"priceId"`
	UnitPrice    float64   `gorm:"type:decimal(12,4);not null" json:"unitPrice"`
	GradeLevel   string    `gorm:"size:20;not null;default:'STD'" json:"gradeLevel"`
	ReportDate   string    `gorm:"type:date;not null;index:idx_worker_date" json:"reportDate"`
	QtyGood      int       `gorm:"not null;default:0" json:"qtyGood"`
	QtyDefect    int       `gorm:"not null;default:0" json:"qtyDefect"`
	QtyTotal     int       `gorm:"not null;default:0" json:"qtyTotal"`
	UnitDefect   float64   `gorm:"type:decimal(12,4);not null;default:0.5" json:"unitDefect"`
	GrossAmount  float64   `gorm:"type:decimal(14,2);not null;default:0" json:"grossAmount"`
	DefectAmount float64   `gorm:"type:decimal(14,2);not null;default:0" json:"defectAmount"`
	NetAmount    float64   `gorm:"type:decimal(14,2);not null;default:0" json:"netAmount"`
	WorkOrderNo  string    `gorm:"size:50" json:"workOrderNo"`
	Remark       string    `gorm:"size:500" json:"remark"`
	Status       int       `gorm:"not null;default:1" json:"status"`
	CreatedBy    uint64    `gorm:"not null" json:"createdBy"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Worker  *SysUser      `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
	Process *ProcessStep  `gorm:"foreignKey:ProcessID" json:"process,omitempty"`
	Price   *ProcessPrice `gorm:"foreignKey:PriceID" json:"price,omitempty"`
}

func (ProductionReport) TableName() string { return "production_report" }

type TeamWageAllocation struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ReportID       uint64    `gorm:"uniqueIndex;not null" json:"reportId"`
	TeamID         uint64    `gorm:"not null" json:"teamId"`
	TotalAmount    float64   `gorm:"type:decimal(14,2);not null" json:"totalAmount"`
	AllocationRule int       `gorm:"not null;default:1" json:"allocationRule"`
	CreatedBy      uint64    `gorm:"not null" json:"createdBy"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`

	Items []TeamWageAllocationItem `gorm:"foreignKey:AllocationID" json:"items,omitempty"`
}

func (TeamWageAllocation) TableName() string { return "team_wage_allocation" }

type TeamWageAllocationItem struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	AllocationID  uint64    `gorm:"not null;index" json:"allocationId"`
	WorkerID      uint64    `gorm:"not null;index" json:"workerId"`
	WeightRatio   float64   `gorm:"type:decimal(8,4);not null" json:"weightRatio"`
	WorkHours     *float64  `gorm:"type:decimal(8,2)" json:"workHours"`
	AllocatedAmt  float64   `gorm:"type:decimal(14,2);not null" json:"allocatedAmt"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`

	Worker *SysUser `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
}

func (TeamWageAllocationItem) TableName() string { return "team_wage_allocation_item" }

type WorkerWageDetail struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkerID      uint64    `gorm:"not null;index:idx_worker_date" json:"workerId"`
	ReportID      *uint64   `gorm:"" json:"reportId"`
	AllocationID  *uint64   `gorm:"" json:"allocationId"`
	DetailType    int       `gorm:"not null" json:"detailType"`
	ProcessID     *uint64   `gorm:"" json:"processId"`
	WageDate      string    `gorm:"type:date;not null;index:idx_worker_date" json:"wageDate"`
	QtyGood       int       `gorm:"not null;default:0" json:"qtyGood"`
	QtyDefect     int       `gorm:"not null;default:0" json:"qtyDefect"`
	UnitPrice     float64   `gorm:"type:decimal(12,4);not null;default:0" json:"unitPrice"`
	Amount        float64   `gorm:"type:decimal(14,2);not null" json:"amount"`
	Remark        string    `gorm:"size:500" json:"remark"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (WorkerWageDetail) TableName() string { return "worker_wage_detail" }

type WorkerWageSummary struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkerID       uint64     `gorm:"not null;uniqueIndex:uk_worker_month" json:"workerId"`
	SummaryMonth   string     `gorm:"size:7;not null;uniqueIndex:uk_worker_month" json:"summaryMonth"`
	TotalQtyGood   int        `gorm:"not null;default:0" json:"totalQtyGood"`
	TotalQtyDefect int        `gorm:"not null;default:0" json:"totalQtyDefect"`
	GrossAmount    float64    `gorm:"type:decimal(14,2);not null;default:0" json:"grossAmount"`
	DefectAmount   float64    `gorm:"type:decimal(14,2);not null;default:0" json:"defectAmount"`
	AllocationAmt  float64    `gorm:"type:decimal(14,2);not null;default:0" json:"allocationAmt"`
	AdjustAmount   float64    `gorm:"type:decimal(14,2);not null;default:0" json:"adjustAmount"`
	NetAmount      float64    `gorm:"type:decimal(14,2);not null;default:0" json:"netAmount"`
	CalcStatus     int        `gorm:"not null;default:0" json:"calcStatus"`
	LastCalcTime   *time.Time `gorm:"" json:"lastCalcTime"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	Worker *SysUser `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
}

func (WorkerWageSummary) TableName() string { return "worker_wage_summary" }
