package repo

// Deployment represents a persist-able deployment.
type Deployment struct {
	ID         int
	Name       string
	ServiceIDs string
	Template   string
}

// Persister is the interface representing necessary persistence actions.
type Persister interface {
	FindByID(int) (Deployment, error)
	All() ([]Deployment, error)
	Save(*Deployment) error
	Remove(int) error
}
