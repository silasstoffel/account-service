package database

func buildPaginationParams(limit, page int) (int, int, int) {
	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return offset, limit, page
}
