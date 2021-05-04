package docker

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestContainers(t *testing.T) {
	cs := NewContainers()
	assert.Len(t, cs.byId, 0)
	assert.Len(t, cs.byListen, 0)
	var c1 types.ContainerJSON
	err := json.Unmarshal([]byte(`
	{
		"Id": "beuha",
		"NetworkSettings": {
			"Networks": {
				"bridge": {
					"IPAddress": "172.17.0.2"
				}
			},
			"Ports": {
				"5000/tcp": null
			}
		}
	}
	`), &c1)
	assert.NoError(t, err)

	fmt.Println(c1)
	ok := cs.Add(&c1)
	assert.True(t, ok)
	assert.Len(t, cs.byId, 1)
	assert.Len(t, cs.byListen, 1)

	c := cs.ByListen("172.17.0.2:5000")
	assert.NotNil(t, c)
	fmt.Println(cs.byListen)
	assert.Equal(t, "beuha", c.ID)

	ok = cs.Remove("beuha")
	assert.True(t, ok)
	assert.Len(t, cs.byId, 0)
	assert.Len(t, cs.byListen, 0)
	ok = cs.Remove("aussi")
	assert.False(t, ok)

}
