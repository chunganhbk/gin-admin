package schema

// StatusText
type StatusText string

func (t StatusText) String() string {
	return string(t)
}


const (
	OKStatus    StatusText = "OK"
	ErrorStatus StatusText = "ERROR"
	FailStatus  StatusText = "FAIL"
)

// StatusResult
type StatusResult struct {
	Status StatusText `json:"status"`
}



// ErrorItem
type ErrorResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListResult
type ListResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination,omitempty"`
}

// Pagination Result
type PaginationResult struct {
	Total    int  `json:"total"`
	Current  uint `json:"current"`
	PageSize uint `json:"pageSize"`
}

// Pagination Param
type PaginationParam struct {
	Pagination bool `form:"-"`
	OnlyCount  bool `form:"-"`
	Current    uint `form:"current,default=1"`
	PageSize   uint `form:"pageSize,default=10" binding:"max=100"`
}

// GetCurrent
func (a PaginationParam) GetCurrent() uint {
	return a.Current
}

// GetPageSize
func (a PaginationParam) GetPageSize() uint {
	pageSize := a.PageSize
	if a.PageSize == 0 {
		pageSize = 100
	}
	return pageSize
}

// OrderDirection
type OrderDirection int

const (
	// Order By ASC
	OrderByASC OrderDirection = 1
	// Order By DESC
	OrderByDESC OrderDirection = 2
)


func NewOrderFieldWithKeys(keys []string, directions ...map[int]OrderDirection) []*OrderField {
	m := make(map[int]OrderDirection)
	if len(directions) > 0 {
		m = directions[0]
	}

	fields := make([]*OrderField, len(keys))
	for i, key := range keys {
		d := OrderByASC
		if v, ok := m[i]; ok {
			d = v
		}

		fields[i] = NewOrderField(key, d)
	}

	return fields
}

// New Order Fields
func NewOrderFields(orderFields ...*OrderField) []*OrderField {
	return orderFields
}

// NewOrderField 创建排序字段
func NewOrderField(key string, d OrderDirection) *OrderField {
	return &OrderField{
		Key:       key,
		Direction: d,
	}
}

// OrderField 排序字段
type OrderField struct {
	Key       string
	Direction OrderDirection
}

// NewID Result
func NewIDResult(id string) *IDResult {
	return &IDResult{
		ID: id,
	}
}

// IDResult 响应唯一标识
type IDResult struct {
	ID string `json:"id"`
}
