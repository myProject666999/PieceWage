package model

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token    string  `json:"token"`
	UserID   uint64  `json:"userId"`
	Username string  `json:"username"`
	RealName string  `json:"realName"`
	Role     int     `json:"role"`
}

type CreateProcessPriceReq struct {
	ProcessID     uint64  `json:"processId" binding:"required"`
	GradeLevel    string  `json:"gradeLevel" binding:"required"`
	UnitPrice     float64 `json:"unitPrice" binding:"required,gt=0"`
	EffectiveDate string  `json:"effectiveDate" binding:"required"`
	Remark        string  `json:"remark"`
}

type CreateReportReq struct {
	WorkerID     uint64  `json:"workerId" binding:"required"`
	TeamID       *uint64 `json:"teamId"`
	ProcessID    uint64  `json:"processId" binding:"required"`
	GradeLevel   string  `json:"gradeLevel"`
	ReportDate   string  `json:"reportDate" binding:"required"`
	QtyGood      int     `json:"qtyGood" binding:"required,gte=0"`
	QtyDefect    int     `json:"qtyDefect" binding:"gte=0"`
	UnitDefect   float64 `json:"unitDefect"`
	WorkOrderNo  string  `json:"workOrderNo"`
	Remark       string  `json:"remark"`
}

type CreateTeamAllocationReq struct {
	ReportID       uint64                  `json:"reportId" binding:"required"`
	AllocationRule int                     `json:"allocationRule" binding:"required,oneof=1 2 3"`
	Members        []AllocationMemberReq   `json:"members" binding:"required,min=1"`
}

type AllocationMemberReq struct {
	WorkerID    uint64   `json:"workerId" binding:"required"`
	WeightRatio float64  `json:"weightRatio" binding:"required,gt=0"`
	WorkHours   *float64 `json:"workHours"`
}

type PriceQueryReq struct {
	ProcessID  uint64 `form:"processId"`
	GradeLevel string `form:"gradeLevel"`
	Date       string `form:"date"`
}

type ReportQueryReq struct {
	WorkerID   uint64 `form:"workerId"`
	TeamID     uint64 `form:"teamId"`
	ProcessID  uint64 `form:"processId"`
	StartDate  string `form:"startDate"`
	EndDate    string `form:"endDate"`
	Status     int    `form:"status"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"pageSize,default=20"`
}

type WageSummaryQueryReq struct {
	WorkerID     uint64 `form:"workerId"`
	SummaryMonth string `form:"summaryMonth"`
	Page         int    `form:"page,default=1"`
	PageSize     int    `form:"pageSize,default=20"`
}

type WorkerWageDetailReq struct {
	WorkerID  uint64 `form:"workerId" binding:"required"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
	WageDate  string `form:"wageDate"`
}

type ProductCreateReq struct {
	ProductCode string `json:"productCode" binding:"required"`
	ProductName string `json:"productName" binding:"required"`
	Spec        string `json:"spec"`
}

type ProcessStepCreateReq struct {
	ProcessCode string `json:"processCode" binding:"required"`
	ProcessName string `json:"processName" binding:"required"`
	ProductID   uint64 `json:"productId" binding:"required"`
	Difficulty  int    `json:"difficulty" binding:"required,oneof=1 2 3 4"`
	Description string `json:"description"`
	IsShared    int    `json:"isShared" binding:"oneof=0 1"`
}

type TeamCreateReq struct {
	TeamName string `json:"teamName" binding:"required"`
	TeamCode string `json:"teamCode" binding:"required"`
	LeaderID *uint64 `json:"leaderId"`
}

type TeamMemberReq struct {
	TeamID   uint64 `json:"teamId" binding:"required"`
	UserID   uint64 `json:"userId" binding:"required"`
	JoinDate string `json:"joinDate" binding:"required"`
}
