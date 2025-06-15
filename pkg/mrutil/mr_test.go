package mrutil

import (
	"reflect"
	"testing"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
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
	diff := api.MergeRequestDiff{} // Not used in current logic

	got := GenerateMRContext(mr, diff)

	want := map[string]interface{}{
		"title":  mr.Title,
		"body":   mr.Description,
		"labels": mr.Labels,
		"dependencies": []DependencyUpdate{
			{
				PackageType: "maven",
				Coordinate:  "org.junit:junit-bom",
				From:        "5.13.0",
				To:          "5.13.1",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GenerateMRContext() = %v\nwant = %v", got, want)
	}
}

func TestGenerateMRContext_NonRenovatePR(t *testing.T) {
	mr := api.MergeRequest{
		Title:       "Fix: Improve logging for auth service",
		Description: "This MR improves the debug logs for auth failures.",
		Labels:      []string{"enhancement"},
	}
	diff := api.MergeRequestDiff{}

	got := GenerateMRContext(mr, diff)

	want := map[string]interface{}{
		"title":  mr.Title,
		"body":   mr.Description,
		"labels": mr.Labels,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GenerateMRContext() = %v\nwant = %v", got, want)
	}
}
