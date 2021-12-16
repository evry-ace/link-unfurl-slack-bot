package unfurl

import (
	"context"
	"fmt"
	"testing"

	"github.com/bndr/gojenkins"
	"github.com/jarcoal/httpmock"
	"gotest.tools/assert"
)

func TestJenkinsBuildLink(t *testing.T) {
	t.Parallel()
	t.Run("should build a link", func(t *testing.T) {
		t.Parallel()

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		// This is called by the Inti() function
		httpmock.RegisterResponder("GET", "https://jenkins.corp.org/api/json",
			httpmock.NewStringResponder(200, ""))

		// Mock the Jenkins job
		jobUrl := "https://jenkins.corp.org/job/my-proj/job/my-repo/job/master/api/json"
		jobJson := fmt.Sprintf("%s/%s", testdataDir, "jenkins-build.json")

		httpmock.RegisterResponder("GET", jobUrl,
			httpmock.NewStringResponder(200, httpmock.File(jobJson).String()))

		// Mock the Jenkins build
		// There is a bug in the jenkins client that adds an extra slash between the
		// job and the build number hence the `master//789`.
		buildUrl := "https://jenkins.corp.org/job/my-proj/job/my-repo/job/master//789/api/json"
		buildJson := fmt.Sprintf("%s/%s", testdataDir, "jenkins-build-789.json")

		httpmock.RegisterResponder("GET", buildUrl,
			httpmock.NewStringResponder(200, httpmock.File(buildJson).String()))

		// Initialize the Jenkins client
		ctx := context.Background()
		jenkins, jenkinsErr := gojenkins.CreateJenkins(nil, "https://jenkins.corp.org/").Init(ctx)

		if jenkinsErr != nil {
			t.Errorf("Error creating Jenkins client: %s", jenkinsErr)
		}

		// Unfurl the jenkins build link
		u := Unfurl{
			Jenkins: jenkins,
		}
		a, err := u.jenkinsBuildLink("my-proj", "my-repo", "master", 789)

		if err != nil {
			t.Errorf("Error building link: %v", err)
		}

		assert.Equal(t, a.Title, "My Proj » my-repo » master #789")
		assert.Equal(t, a.TitleLink, "https://jenkins.corp.org/job/my-proj/job/my-repo/job/master/789/")
	})
}
