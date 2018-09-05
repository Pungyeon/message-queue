package queue

/* NO TIME FOR THIS!!! */

// ExchangeList is a struct for managing Queue Exchanges
type ExchangeList struct {
	list map[string]string
}

// NewExchangeList will return an initialised Exchangelist
func NewExchangeList() *ExchangeList {
	return &ExchangeList{
		list: make(map[string]string),
	}
}
