package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
	"github.com/CenturyLinkLabs/panamax-remote-agent-go/repo"
	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	callbackForAll    func() ([]repo.Deployment, error)
	callbackForFind   func(int) (repo.Deployment, error)
	callbackForSave   func(*repo.Deployment) error
	callbackForRemove func(int) error
}

func (fr fakeRepo) FindByID(qID int) (repo.Deployment, error) {
	if fr.callbackForFind != nil {
		return fr.callbackForFind(qID)
	}
	return repo.Deployment{}, nil
}
func (fr fakeRepo) All() ([]repo.Deployment, error) {
	if fr.callbackForAll != nil {
		return fr.callbackForAll()
	}
	return []repo.Deployment{}, nil
}
func (fr fakeRepo) Save(d *repo.Deployment) error {
	if fr.callbackForSave != nil {
		return fr.callbackForSave(d)
	}
	return nil
}
func (fr fakeRepo) Remove(qID int) error {
	if fr.callbackForRemove != nil {
		return fr.callbackForRemove(qID)
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
	fr := fakeRepo{
		callbackForAll: func() ([]repo.Deployment, error) {
			drs := []repo.Deployment{
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
	dm := MakeDeploymentManager(fr, fc)

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
	fr := fakeRepo{
		callbackForAll: func() ([]repo.Deployment, error) {
			return []repo.Deployment{}, errors.New("something failed")
		},
	}
	fc := fakeAdapterClient{}
	dm := MakeDeploymentManager(fr, fc)

	drs, err := dm.ListDeployments()

	ex := []DeploymentResponseLite{}

	assert.Equal(t, ex, drs)
	assert.EqualError(t, err, "something failed")
}

func TestSuccessfulGetDeployment(t *testing.T) {
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			dr := repo.Deployment{
				ID:         7,
				Name:       "booyah",
				ServiceIDs: `["wp-pod", "db-pod"]`,
				Template:   `{"name": "boom"}`,
			}

			return dr, nil
		},
	}
	dm := MakeDeploymentManager(fr, fakeAdapterClient{})

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
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			dr := repo.Deployment{
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

	dm := MakeDeploymentManager(fr, fc)

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
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			return repo.Deployment{}, errors.New("something failed")
		},
	}
	dm := MakeDeploymentManager(fr, fakeAdapterClient{})

	dr, err := dm.GetDeployment(7)

	assert.Equal(t, DeploymentResponseLite{}, dr)
	assert.EqualError(t, err, "something failed")
}

func TestSucessfulGetFullDeployment(t *testing.T) {
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			dr := repo.Deployment{
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
	dm := MakeDeploymentManager(fr, fc)

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
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			return repo.Deployment{}, errors.New("something failed")
		},
	}
	dm := MakeDeploymentManager(fr, fakeAdapterClient{})

	dr, err := dm.GetFullDeployment(7)

	assert.Equal(t, DeploymentResponseFull{}, dr)
	assert.EqualError(t, err, "something failed")
}

func TestSuccessfulDeleteDeployment(t *testing.T) {
	var rmIDs []int
	var delIDs []string
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			dr := repo.Deployment{
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

	dm := MakeDeploymentManager(fr, fc)

	err := dm.DeleteDeployment(7)

	assert.Equal(t, []string{"wp-pod"}, delIDs)
	assert.Equal(t, []int{7}, rmIDs)
	assert.NoError(t, err)
}

func TestDeleteDeploymentWhenItDoesNotExit(t *testing.T) {
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {

			return repo.Deployment{}, errors.New("deployment does not exist")
		},
	}

	dm := MakeDeploymentManager(fr, fakeAdapterClient{})

	err := dm.DeleteDeployment(7)

	assert.EqualError(t, err, "deployment does not exist")
}

func TestDeleteDeploymentWhenServiceDeletionFails(t *testing.T) {
	fc := fakeAdapterClient{
		callbackForDeleteService: func(sID string) error {
			return errors.New("failed to delete service")
		},
	}
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			dr := repo.Deployment{
				ID:         7,
				Name:       "To be deleted",
				ServiceIDs: `["wp-pod"]`,
			}

			return dr, nil
		},
	}

	dm := MakeDeploymentManager(fr, fc)

	err := dm.DeleteDeployment(7)

	assert.EqualError(t, err, "failed to delete service")
}

func TestDeleteDeploymentRepoDeletionFails(t *testing.T) {
	fr := fakeRepo{
		callbackForFind: func(qID int) (repo.Deployment, error) {
			dr := repo.Deployment{
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

	dm := MakeDeploymentManager(fr, fakeAdapterClient{})

	err := dm.DeleteDeployment(7)

	assert.EqualError(t, err, "failed")
}

func TestCreateDeploymentPersistedTheMergedTemplate(t *testing.T) {
	var passedDep repo.Deployment
	fr := fakeRepo{
		callbackForSave: func(dep *repo.Deployment) error {
			passedDep = *dep
			return nil
		},
	}

	dm := MakeDeploymentManager(fr, fakeAdapterClient{})

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
