package converter_test

import (
	"testing"
	"workshop_demo/client"
	"workshop_demo/converter"
	"workshop_demo/model"
	"github.com/stretchr/testify/assert"
)
// all test functions must start with Test; go test ./... will run all tests that can be found in the project
func TestConverter(t *testing.T) {
	// assert.Equal(t,1,2)
	// t.Run()

	expected := model.ServerResponse{}
	actual := converter.ConvertModels(client.DNSResponse{}, client.DBaaSResponse{})

	assert.Equal(t, expected, actual)

}

func TestOne(t *testing.T) {
	assert.Equal(t,1,1)
}
