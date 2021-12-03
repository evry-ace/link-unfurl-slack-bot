package bitbucket

import (
	"testing"

	"gotest.tools/assert"
)

func TestCommitOpionsToQueryString(t *testing.T) {
	t.Run("should return empty string if no options are defined", func(t *testing.T) {
		options := CommitOptions{}
		assert.Equal(t, "", options.ToQueryString())
	})

	t.Run("should return query string if a single option is defined", func(t *testing.T) {
		options := CommitOptions{
			Since: "foo",
		}
		assert.Equal(t, "since=foo", options.ToQueryString())
	})

	t.Run("should return query string if multiple options are defined", func(t *testing.T) {
		options := CommitOptions{
			Since: "foo",
			Path:  "bar",
		}
		assert.Equal(t, "path=bar&since=foo", options.ToQueryString())
	})
}
