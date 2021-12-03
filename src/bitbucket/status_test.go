package bitbucket

import (
	"encoding/json"
	"testing"

	"github.com/evry-ace/link-unfurl-slack-bot/src/utils"
	"gotest.tools/assert"
)

func TestStatusFromJSON(t *testing.T) {
	var statusList StatusList

	rawJSON := utils.ReadTestdataFile("bitbucket-build-status-654382.json")

	if err := json.Unmarshal([]byte(rawJSON), &statusList); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1, len(statusList.Values))
	assert.Equal(t, "ad887ae7b1528d2058db02a5162c1ece", statusList.Values[0].Key)
}

func TestStatusListState(t *testing.T) {
	t.Run("should return failed if any staus is failed", func(t *testing.T) {
		s := StatusList{
			Values: []Status{
				{
					State: StatusSuccess,
				},
				{
					State: StatusInProgress,
				},
				{
					State: StatusFailed,
				},
			},
		}

		assert.Equal(t, StatusFailed, s.State())
	})

	t.Run("should return in progress if any status is in progress", func(t *testing.T) {
		s := StatusList{
			Values: []Status{
				{
					State: StatusSuccess,
				},
				{
					State: StatusInProgress,
				},
				{
					State: StatusSuccess,
				},
			},
		}

		assert.Equal(t, StatusInProgress, s.State())
	})

	t.Run("should return success if all statuses is success", func(t *testing.T) {
		s := StatusList{
			Values: []Status{
				{
					State: StatusSuccess,
				},
				{
					State: StatusSuccess,
				},
				{
					State: StatusSuccess,
				},
			},
		}

		assert.Equal(t, StatusSuccess, s.State())
	})
}
