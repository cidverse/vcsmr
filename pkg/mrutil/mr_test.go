package mrutil

import (
	"testing"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMRContext_RenovatePR(t *testing.T) {
	mr := api.MergeRequest{
		Title: "Renovate: Update spring-context to 6.2.8",
		Description: `
This PR contains the following updates:

| Package | Change | Age | Adoption | Passing | Confidence |
|---|---|---|---|---|---|
| [org.junit:junit-bom](https://junit.org/junit5/) ([source](https://redirect.github.com/junit-team/junit5)) | 5.13.0 -> 5.13.1 | [![age](https://developer.mend.io/api/mc/badges/age/maven/org.junit:junit-bom/5.13.1?slim=true)](https://docs.renovatebot.com/merge-confidence/) | [![adoption](https://developer.mend.io/api/mc/badges/adoption/maven/org.junit:junit-bom/5.13.1?slim=true)](https://docs.renovatebot.com/merge-confidence/) | [![passing](https://developer.mend.io/api/mc/badges/compatibility/maven/org.junit:junit-bom/5.13.0/5.13.1?slim=true)](https://docs.renovatebot.com/merge-confidence/) | [![confidence](https://developer.mend.io/api/mc/badges/confidence/maven/org.junit:junit-bom/5.13.0/5.13.1?slim=true)](https://docs.renovatebot.com/merge-confidence/) |
`,
		Labels: []string{"dependencies", "renovate"},
	}
	got := GenerateMRContext(mr, api.MergeRequestDiff{})

	assert.Equal(t, mr.Title, got["title"])
	assert.Equal(t, mr.Description, got["description"])
	assert.Equal(t, mr.Labels, got["labels"])
	assert.Equal(t, "maven", got["dependencyType"])
	assert.Equal(t, "org.junit:junit-bom", got["dependencyCoordinate"])
	assert.Equal(t, "5.13.0", got["dependencyFrom"])
	assert.Equal(t, "5.13.1", got["dependencyTo"])
}

func TestGenerateMRContext_NonRenovatePR(t *testing.T) {
	mr := api.MergeRequest{
		Title:       "Fix: Improve logging for auth service",
		Description: "This MR improves the debug logs for auth failures.",
		Labels:      []string{"enhancement"},
	}
	got := GenerateMRContext(mr, api.MergeRequestDiff{})

	assert.Equal(t, mr.Title, got["title"])
	assert.Equal(t, mr.Description, got["description"])
	assert.Equal(t, mr.Labels, got["labels"])
}
