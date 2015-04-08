package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type servicesClient struct {
	endpoint  string
	netClient *http.Client
}

// MakeClient returns a Client connected with the adapter.
func MakeClient(ep string) Client {
	hc := &http.Client{}

	c := servicesClient{
		netClient: hc,
		endpoint:  ep,
	}

	return c
}

func (sc servicesClient) CreateServices(buf *bytes.Buffer) ([]Service, error) {
	resp, err := sc.netClient.Post(sc.servicesPath(""), "application/json", buf)
	defer resp.Body.Close()
	if err != nil {
		return []Service{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return []Service{}, fmt.Errorf("Failed to create services, resp code: %d, body: %s", resp.StatusCode, string(body))
	}

	ars := &[]Service{}
	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(ars)
	if err != nil {
		return []Service{}, err
	}

	return *ars, nil
}

func (sc servicesClient) GetService(sid string) (Service, error) {
	resp, err := sc.netClient.Get(sc.servicesPath(sid))
	defer resp.Body.Close()
	if err != nil {
		return Service{}, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return Service{ID: sid, ActualState: "not found"}, nil
	} else if resp.StatusCode != http.StatusOK {
		return Service{ID: sid, ActualState: "error"}, nil
	}

	srvc := &Service{}
	jd := json.NewDecoder(resp.Body)
	if err := jd.Decode(srvc); err != nil {
		return Service{}, err
	}

	return *srvc, nil
}

func (sc servicesClient) DeleteService(sid string) error {
	req, err := http.NewRequest("DELETE", sc.servicesPath(sid), nil)

	if err != nil {
		return err
	}

	_, err = sc.netClient.Do(req)

	return err
}

func (sc servicesClient) FetchMetadata() (interface{}, error) {
	res, err := sc.netClient.Get(sc.endpoint + "/v1/metadata")

	if err != nil {
		return nil, err
	}

	var r interface{}
	jd := json.NewDecoder(res.Body)
	if err := jd.Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}

func (sc servicesClient) servicesPath(id string) string {
	if id != "" {
		return sc.endpoint + "/v1/services/" + id
	}
	return sc.endpoint + "/v1/services"
}
