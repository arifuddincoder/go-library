package httpresponse

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Paginated struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}

func NewPaginated(data any, page, limit int, total int64) Paginated {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return Paginated{
		Data: data,
		Meta: Meta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
