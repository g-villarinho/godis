package domain

type Item struct {
	Value     string
	ExpiresAt *int64
}

func NewItem(value string, expiresAt *int64) *Item {
	return &Item{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func (i *Item) IsExpired(now int64) bool {
	if i.ExpiresAt == nil {
		return false
	}

	return now > *i.ExpiresAt
}
