package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/stretchr/testify/assert"
)

type fakeStore struct {
	callbackForAll    func() ([]Deployment, error)
	callbackForFind   func(int) (Deployment, error)
	callbackForSave   func(*Deployment) error
	callbackForRemove func(int) error
}

func (fs fakeStore) FindByID(qID int) (Deployment, error) {
	if fs.callbackForFind != nil {
		return fs.callbackForFind(qID)
	}
	return Deployment{}, nil
}
func (fs fakeStore) All() ([]Deployment, error) {
	if fs.callbackForAll != nil {
		return fs.callbackForAll()
	}
	return []Deployment{}, nil
}
func (fs fakeStore) Save(d *Deployment) error {
	if fs.callbackForSave != nil {
		return fs.callbackForSave(d)
	}
	return nil
}
func (fs fakeStore) Remove(qID int) error {
	if fs.callbackForRemove != nil {
		return fs.callbackForRemove(qID)
	}
	return nil
}

type fakeAdapterClient struct {
	callbackForCreateServices func(*bytes.Buffer) []adapter.Service
	callbackForGetService     func(string) adapter.Service
	callbackForDeleteService  func(string) error
	callbackForFetchMetadata  func() (interface{}, error)
}

func (fc fakeAdapterClient) CreateServices(b *bytes.Buffer) ([]adapter.Service, error) {
	if fc.callbackForCreateServices != nil {
		return fc.callbackForCreateServices(b), nil
	}
	return []adapter.Service{}, nil
}
func (fc fakeAdapterClient) GetService(sID string) (adapter.Service, error) {
	if fc.callbackForGetService != nil {
		return fc.callbackForGetService(sID), nil
	}
	return adapter.Service{}, nil
}
func (fc fakeAdapterClient) DeleteService(sID string) error {
	if fc.callbackForDeleteService != nil {
		return fc.callbackForDeleteService(sID)
	}
	return nil
}
func (fc fakeAdapterClient) FetchMetadata() (interface{}, error) {
	if fc.callbackForFetchMetadata != nil {
		return fc.callbackForFetchMetadata()
	}
	return nil, nil
}

func TestSuccessfullListDeployments(t *testing.T) {
	fs := fakeStore{
		callbackForAll: func() ([]Deployment, error) {
			drs := []Deployment{
				{
					ID:         1,
					Name:       "booyah",
					ServiceIDs: `["wp-pod", "db-pod"]`,
					Template:   `{"name": "boom"}`,
				},
			}
			return drs, nil
		},
	}
	fc := fakeAdapterClient{}
	dm := MakeDeploymentManager(fs, fc, "v1")

	drs, err := dm.ListDeployments()

	ex := []DeploymentResponseLite{
		{
			ID:           1,
			Name:         "booyah",
			Redeployable: true,
			ServiceIDs:   []string{"wp-pod", "db-pod"},
		},
	}

	assert.Equal(t, ex, drs)
	assert.NoError(t, err)
}

func TestErroredListDeployments(t *testing.T) {
	fs := fakeStore{
		callbackForAll: func() ([]Deployment, error) {
			return []Deployment{}, errors.New("something failed")
		},
	}
	fc := fakeAdapterClient{}
	dm := MakeDeploymentManager(fs, fc, "v1")

	drs, err := dm.ListDeployments()

	ex := []DeploymentResponseLite{}

	assert.Equal(t, ex, drs)
	assert.EqualError(t, err, "something failed")
}

func TestSuccessfulGetDeployment(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			dr := Deployment{
				ID:         7,
				Name:       "booyah",
				ServiceIDs: `["wp-pod", "db-pod"]`,
				Template:   `{"name": "boom"}`,
			}

			return dr, nil
		},
	}
	dm := MakeDeploymentManager(fs, fakeAdapterClient{}, "v1")

	dr, err := dm.GetDeployment(7)

	ex := DeploymentResponseLite{
		ID:           7,
		Name:         "booyah",
		Redeployable: true,
		ServiceIDs:   []string{"wp-pod", "db-pod"},
	}

	assert.Equal(t, ex, dr)
	assert.NoError(t, err)
}

func TestSuccessfulGetDeploymentWithNoTemplate(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			dr := Deployment{
				ID:         8,
				Name:       "NO TEMPLATE",
				ServiceIDs: "",
			}

			return dr, nil
		},
	}
	fc := fakeAdapterClient{
		callbackForGetService: func(qID string) adapter.Service {
			return adapter.Service{
				ID:          "wp-pod",
				ActualState: "running",
			}
		},
	}

	dm := MakeDeploymentManager(fs, fc, "v1")

	dr, err := dm.GetDeployment(7)

	ex := DeploymentResponseLite{
		ID:           8,
		Name:         "NO TEMPLATE",
		Redeployable: false,
	}

	assert.Equal(t, ex, dr)
	assert.NoError(t, err)
}

func TestErroredGetDeployment(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			return Deployment{}, errors.New("something failed")
		},
	}
	dm := MakeDeploymentManager(fs, fakeAdapterClient{}, "v1")

	dr, err := dm.GetDeployment(7)

	assert.Equal(t, DeploymentResponseLite{}, dr)
	assert.EqualError(t, err, "something failed")
}

func TestSucessfulGetFullDeployment(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			dr := Deployment{
				ID:         7,
				Name:       "Full booyah",
				ServiceIDs: `["wp-pod"]`,
			}

			return dr, nil
		},
	}
	fc := fakeAdapterClient{
		callbackForGetService: func(qID string) adapter.Service {
			return adapter.Service{
				ID:          "wp-pod",
				ActualState: "running",
			}
		},
	}
	dm := MakeDeploymentManager(fs, fc, "v1")

	dr, err := dm.GetFullDeployment(7)

	ex := DeploymentResponseFull{
		ID:   7,
		Name: "Full booyah",
		Status: Status{
			Services: []Service{
				{
					ID:          "wp-pod",
					ActualState: "running",
				},
			},
		},
	}

	assert.Equal(t, ex, dr)
	assert.NoError(t, err)
}

func TestErroredGetFullDeployment(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			return Deployment{}, errors.New("something failed")
		},
	}
	dm := MakeDeploymentManager(fs, fakeAdapterClient{}, "v1")

	dr, err := dm.GetFullDeployment(7)

	assert.Equal(t, DeploymentResponseFull{}, dr)
	assert.EqualError(t, err, "something failed")
}

func TestSuccessfulDeleteDeployment(t *testing.T) {
	var rmIDs []int
	var delIDs []string
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			dr := Deployment{
				ID:         7,
				Name:       "To be deleted",
				ServiceIDs: `["wp-pod"]`,
			}

			return dr, nil
		},
		callbackForRemove: func(qID int) error {
			rmIDs = append(rmIDs, qID)
			return nil
		},
	}

	fc := fakeAdapterClient{
		callbackForDeleteService: func(sID string) error {
			delIDs = append(delIDs, sID)
			return nil
		},
	}

	dm := MakeDeploymentManager(fs, fc, "v1")

	err := dm.DeleteDeployment(7)

	assert.Equal(t, []string{"wp-pod"}, delIDs)
	assert.Equal(t, []int{7}, rmIDs)
	assert.NoError(t, err)
}

func TestDeleteDeploymentWhenItDoesNotExit(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {

			return Deployment{}, errors.New("deployment does not exist")
		},
	}

	dm := MakeDeploymentManager(fs, fakeAdapterClient{}, "v1")

	err := dm.DeleteDeployment(7)

	assert.EqualError(t, err, "deployment does not exist")
}

func TestDeleteDeploymentWhenServiceDeletionFails(t *testing.T) {
	fc := fakeAdapterClient{
		callbackForDeleteService: func(sID string) error {
			return errors.New("failed to delete service")
		},
	}
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			dr := Deployment{
				ID:         7,
				Name:       "To be deleted",
				ServiceIDs: `["wp-pod"]`,
			}

			return dr, nil
		},
	}

	dm := MakeDeploymentManager(fs, fc, "v1")

	err := dm.DeleteDeployment(7)

	assert.EqualError(t, err, "failed to delete service")
}

func TestDeleteDeploymentRepoDeletionFails(t *testing.T) {
	fs := fakeStore{
		callbackForFind: func(qID int) (Deployment, error) {
			dr := Deployment{
				ID:         7,
				Name:       "To be deleted",
				ServiceIDs: `["wp-pod"]`,
			}

			return dr, nil
		},
		callbackForRemove: func(qID int) error {
			return errors.New("failed")
		},
	}

	dm := MakeDeploymentManager(fs, fakeAdapterClient{}, "v1")

	err := dm.DeleteDeployment(7)

	assert.EqualError(t, err, "failed")
}

func TestCreateDeploymentPersistedTheMergedTemplate(t *testing.T) {
	var passedDep Deployment
	fs := fakeStore{
		callbackForSave: func(dep *Deployment) error {
			passedDep = *dep
			return nil
		},
	}

	dm := MakeDeploymentManager(fs, fakeAdapterClient{}, "v1")

	d := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{
					Name:       "wp",
					Deployment: DeploymentSettings{Count: FromIntOrString{1}},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name:       "wp",
					Deployment: DeploymentSettings{Count: FromIntOrString{2}},
				},
			},
		},
	}

	_, err := dm.CreateDeployment(d)

	var tpl Template
	json.Unmarshal([]byte(passedDep.Template), &tpl)

	assert.NoError(t, err)
	assert.Equal(t, FromIntOrString{2}, tpl.Images[0].Deployment.Count)
}
