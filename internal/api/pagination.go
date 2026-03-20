package api

// PaginationInfo holds pagination metadata for list responses.
type PaginationInfo struct {
	StartAt    int  `json:"startAt"`
	MaxResults int  `json:"maxResults"`
	Total      int  `json:"total"`
	IsLast     bool `json:"isLast"`
}

// ToMap converts pagination info to a map for JSON embedding.
func (p *PaginationInfo) ToMap() map[string]any {
	return map[string]any{
		"startAt":    p.StartAt,
		"maxResults": p.MaxResults,
		"total":      p.Total,
		"isLast":     p.IsLast,
	}
}

// PaginationFromResponse extracts pagination info from a typical Atlassian response.
func PaginationFromResponse(data map[string]any) *PaginationInfo {
	p := &PaginationInfo{}

	if v, ok := data["startAt"].(float64); ok {
		p.StartAt = int(v)
	}
	if v, ok := data["maxResults"].(float64); ok {
		p.MaxResults = int(v)
	}
	if v, ok := data["total"].(float64); ok {
		p.Total = int(v)
	}
	if v, ok := data["isLast"].(bool); ok {
		p.IsLast = v
	} else {
		// Compute isLast from offset pagination
		p.IsLast = p.StartAt+p.MaxResults >= p.Total
	}

	return p
}
