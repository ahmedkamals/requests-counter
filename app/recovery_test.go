package app

import (
	"testing"
	testingUtil "requests-counter/testing_util"
	"encoding/json"
	"reflect"
)

func TestShouldUnmarshallCorrectly(t *testing.T) {

	data := `{"index":33,"time_stamp":1505810337,"data":[1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}`
	expectedData := [60]uint64{}
	expectedData[0] = 1

	testCases := []testingUtil.TestCase{
		{
			Id: "Should unmarshall recovered data correctly.",
			Input: []byte(data),
			Expected: NewRecovery(33, 1505810337, expectedData),
		},
	}

	for _, testCase := range testCases {
		input := testCase.Input.([]byte)
		expected := testCase.Expected.(*Recovery)
		actual := NewEmptyRecovery()
		err := json.Unmarshal(input, actual)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Error(testingUtil.Format(testCase.Id, expected, actual))
		}
	}
}
