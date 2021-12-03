package bitbucket

import (
	"testing"
	"time"

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

func TestCommitString(t *testing.T) {
	t.Run("should return string representation of commit", func(t *testing.T) {
		commit := Commit{
			ID:              "1234567890",
			Message:         "foo",
			AuthorTimestamp: time.Now().UnixMilli(),
			Author: User{
				DisplayName: "bar",
			},
		}
		assert.Equal(t, "foo (1234567890) by bar about a second ago", commit.String())
	})
}

func TestCommitTimeAgo(t *testing.T) {
	t.Run("should return empty string if no timeago is defined", func(t *testing.T) {
		commit := Commit{}
		assert.Equal(t, "", commit.TimeAgo())
	})

	t.Run("should return human readable one month ago", func(t *testing.T) {
		commit := Commit{
			//AuthorTimestamp: 1636629378000,
			AuthorTimestamp: time.Now().UnixMilli() - 30*24*60*60*1000,
		}
		assert.Equal(t, "one month ago", commit.TimeAgo())
	})

	t.Run("should return human readable one year ago", func(t *testing.T) {
		commit := Commit{
			//AuthorTimestamp: 1636629378000,
			AuthorTimestamp: time.Now().UnixMilli() - 365*24*60*60*1000,
		}
		assert.Equal(t, "one year ago", commit.TimeAgo())
	})
}
