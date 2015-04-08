package agent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergedImagesKeepsOgENV(t *testing.T) {
	depB := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "bar"},
					},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name: "wp",
				},
			},
		},
	}

	mImgs := depB.MergedImages()

	e := []Environment{{Variable: "FOO", Value: "bar"}}

	assert.Equal(t, e, mImgs[0].Environment)
}

func TestMergedImagesOverridesENV(t *testing.T) {
	depB := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "bar"},
					},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "overridden"},
					},
				},
			},
		},
	}

	mImgs := depB.MergedImages()

	e := []Environment{{Variable: "FOO", Value: "overridden"}}

	assert.Equal(t, e, mImgs[0].Environment)
}

func TestMergedImagesAddsExtraENVs(t *testing.T) {
	depB := DeploymentBlueprint{
		Template: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "bar"},
					},
				},
			},
		},
		Override: Template{
			Images: []Image{
				{
					Name: "wp",
					Environment: []Environment{
						{Variable: "FOO", Value: "overridden"},
						{Variable: "MORE", Value: "stuff"},
					},
				},
			},
		},
	}

	mImgs := depB.MergedImages()

	e := []Environment{
		{Variable: "FOO", Value: "overridden"},
		{Variable: "MORE", Value: "stuff"},
	}

	assert.Equal(t, e, mImgs[0].Environment)
}

func TestEmptyImageProducesEmptyJSON(t *testing.T) {
	img := Image{}
	j, err := json.Marshal(img)

	assert.NoError(t, err)
	assert.Equal(t, "{}", string(j))
}

func TestFullImageProducesProperJSON(t *testing.T) {
	img := Image{
		Name:        "foo",
		Source:      "bar/foo",
		Command:     "./run.sh",
		Deployment:  DeploymentSettings{Count: FromIntOrString{2}},
		Links:       []Link{{Service: "bla", Alias: "b"}},
		Environment: []Environment{{Variable: "FOO", Value: "bar"}},
		Ports:       []Port{{HostPort: FromIntOrString{22}, ContainerPort: FromIntOrString{23}}},
		Expose:      []FromIntOrString{{33}, {44}},
		Volumes:     []Volume{{ContainerPath: "/var", HostPath: "/usr"}},
		VolumesFrom: []string{"/viz"},
	}
	j, err := json.Marshal(img)

	buf := []byte(`{
		"command":"./run.sh",
		"deployment":{"count":2},
		"environment":[{"variable":"FOO","value":"bar"}],
		"expose":[33,44],
		"links":[{"alias":"b","name":"bla"}],
		"name":"foo",
		"ports":[{"containerPort":23,"hostPort":22}],
		"source":"bar/foo",
		"volumes":[{"containerPath":"/var","hostPath":"/usr"}],
		"volumesFrom":["/viz"]
	}`)

	bb := bytes.Buffer{}
	json.Compact(&bb, buf)
	e, _ := ioutil.ReadAll(&bb)

	assert.NoError(t, err)
	assert.Equal(t, string(e), string(j))
}
