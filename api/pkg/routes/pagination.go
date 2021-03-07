package routes

type PaginatedResponseDTO struct {
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}

func ToPaginatedDTO(totalPages int, items interface{}) *PaginatedResponseDTO {
	return &PaginatedResponseDTO{
		TotalPages: totalPages,
		Items:      items,
	}
}
