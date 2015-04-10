package agent

import (
	"bytes"
	"encoding/json"

	"github.com/CenturyLinkLabs/panamax-remote-agent-go/adapter"
)

// DeploymentStore is an interface that can persist information for an agent
// about deployments.
type DeploymentStore interface {
	FindByID(int) (Deployment, error)
	All() ([]Deployment, error)
	Save(*Deployment) error
	Remove(int) error
}

type deploymentManager struct {
	Store   DeploymentStore
	Client  adapter.Client
	version string
}

// MakeDeploymentManager returns a deploymentManager hydrated with a persister and adapter client.
func MakeDeploymentManager(s DeploymentStore, c adapter.Client, v string) Manager {
	return deploymentManager{
		Store:   s,
		Client:  c,
		version: v,
	}
}

// ListDeployments lists all available deployments in a repo.
func (dm deploymentManager) ListDeployments() ([]DeploymentResponseLite, error) {
	deps, err := dm.Store.All()
	if err != nil {
		return []DeploymentResponseLite{}, err
	}

	drs := make([]DeploymentResponseLite, len(deps), len(deps))

	for i, dep := range deps {
		dr := deploymentResponseLiteFromRawValues(
			dep.ID,
			dep.Name,
			dep.Template,
			dep.ServiceIDs,
		)

		drs[i] = dr
	}

	return drs, nil
}

// GetFullDeployment returns an extended representation of the deployment with the given ID.
func (dm deploymentManager) GetFullDeployment(qid int) (DeploymentResponseFull, error) {
	dep, err := dm.GetDeployment(qid)

	if err != nil {
		return DeploymentResponseFull{}, err
	}

	as := make([]Service, len(dep.ServiceIDs), len(dep.ServiceIDs))
	for i, sID := range dep.ServiceIDs {
		srvc, err := dm.Client.GetService(sID)
		if err != nil {
			return DeploymentResponseFull{}, err
		}
		as[i] = Service{
			ID:          srvc.ID,
			ActualState: srvc.ActualState,
		}
	}

	dr := DeploymentResponseFull{
		Name:         dep.Name,
		ID:           dep.ID,
		Redeployable: dep.Redeployable,
		Status:       Status{Services: as},
	}

	return dr, nil
}

// GetDeployment returns a representation of the deployment with the given ID.
func (dm deploymentManager) GetDeployment(qid int) (DeploymentResponseLite, error) {
	dep, err := dm.Store.FindByID(qid)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	drl := deploymentResponseLiteFromRawValues(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	return drl, nil
}

// DeleteDeployment deletes the deployment, with the given ID,
// from both the store and adapter.
func (dm deploymentManager) DeleteDeployment(qID int) error {
	dep, err := dm.Store.FindByID(qID)

	if err != nil {
		return err
	}

	var sIDs []string
	if err := json.Unmarshal([]byte(dep.ServiceIDs), &sIDs); err != nil {
		return err
	}

	for _, sID := range sIDs {
		if err := dm.Client.DeleteService(sID); err != nil {
			return err
		}
	}

	if err := dm.Store.Remove(qID); err != nil {
		return err
	}

	return err
}

// CreateDeployment creates a new deployment from a DeploymentBlueprint.
func (dm deploymentManager) CreateDeployment(depB DeploymentBlueprint) (DeploymentResponseLite, error) {

	mImgs := depB.MergedImages()

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(mImgs); err != nil {
		return DeploymentResponseLite{}, err
	}

	as, err := dm.Client.CreateServices(buf)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	tn := depB.Template.Name
	dep, err := makeDeployment(tn, mImgs, as)
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	if err := dm.Store.Save(&dep); err != nil {
		return DeploymentResponseLite{}, err
	}

	drl := deploymentResponseLiteFromRawValues(
		dep.ID,
		dep.Name,
		dep.Template,
		dep.ServiceIDs,
	)

	return drl, nil
}

// ReDeploy recreates a given deployment, by deleteing, then creating with the
// same template. The returned record will have a new ID.
func (dm deploymentManager) ReDeploy(ID int) (DeploymentResponseLite, error) {

	dep, err := dm.Store.FindByID(ID)

	var tpl Template
	if err := json.Unmarshal([]byte(dep.Template), &tpl); err != nil {
		return DeploymentResponseLite{}, err
	}

	if err := dm.DeleteDeployment(ID); err != nil {
		return DeploymentResponseLite{}, err
	}

	drl, err := dm.CreateDeployment(DeploymentBlueprint{Template: tpl})
	if err != nil {
		return DeploymentResponseLite{}, err
	}

	return drl, nil
}

// FetchMetadata returns metadata for the agent and adapter.
func (dm deploymentManager) FetchMetadata() (Metadata, error) {
	adapterMeta, _ := dm.Client.FetchMetadata()

	md := Metadata{
		Agent:   AgentMetadata{Version: dm.version},
		Adapter: adapterMeta,
	}

	return md, nil
}

func makeDeployment(tn string, mImgs []Image, as []adapter.Service) (Deployment, error) {
	ts, err := stringifyTemplate(tn, mImgs)
	if err != nil {
		return Deployment{}, err
	}

	ss, err := stringifyServiceIDs(as)

	if err != nil {
		return Deployment{}, err
	}

	return Deployment{
		Name:       tn,
		Template:   ts,
		ServiceIDs: ss,
	}, nil
}

func stringifyTemplate(tn string, imgs []Image) (string, error) {
	mt := Template{
		Name:   tn,
		Images: imgs,
	}
	b, err := json.Marshal(mt)

	return string(b), err
}

func stringifyServiceIDs(as []adapter.Service) (string, error) {
	sIDs := make([]string, len(as), len(as))

	for i, ar := range as {
		sIDs[i] = ar.ID
	}

	sb, err := json.Marshal(sIDs)

	return string(sb), err
}

func deploymentResponseLiteFromRawValues(id int, nm string, tpl string, sids string) DeploymentResponseLite {
	drl := &DeploymentResponseLite{
		ID:           id,
		Name:         nm,
		Redeployable: tpl != "",
	}
	// if this is an empty string it will fail to unmarshal, and we are fine with that.
	json.Unmarshal([]byte(sids), &drl.ServiceIDs)

	return *drl
}
