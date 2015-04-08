package adapter

import (
	"bytes"
)

// Client is the interface representing all the use cases required
// to interact with an adapter.
type Client interface {
	CreateServices(*bytes.Buffer) ([]Service, error)
	GetService(string) (Service, error)
	DeleteService(string) error
	FetchMetadata() (interface{}, error)
}

// Service is the representation of the entity coming back from the adapter.
type Service struct {
	ID          string `json:"id"`
	ActualState string `json:"actualState"`
}
