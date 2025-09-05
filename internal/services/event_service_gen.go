package services

func (c CreateEvent) SQL() ([]string, []any, error) {
	columns := []string{
		"OwnerUserId",
		"Label",
	}
	params := []any{
		c.OwnerUserId,
		c.Label,
	}

	if c.CoverPhotoPath != nil {
		columns = append(columns, "CoverPhotoPath")
		params = append(params, &c.CoverPhotoPath)
	}

	return columns, params, nil
}
