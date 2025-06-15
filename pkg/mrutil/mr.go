package mrutil

import (
	"regexp"
	"strings"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
)

type DependencyUpdate struct {
	PackageType string
	Coordinate  string
	From        string
	To          string
}

// GenerateMRContext builds the merge request context map for user-provided rule evaluation
func GenerateMRContext(mr api.MergeRequest, diff api.MergeRequestDiff) map[string]interface{} {
	mrContext := map[string]interface{}{
		"title":               mr.Title,
		"description":         mr.Description,
		"labels":              mr.Labels,
		"sourceBranch":        mr.SourceBranch,
		"targetBranch":        mr.TargetBranch,
		"state":               string(mr.State),
		"isMerged":            mr.IsMerged,
		"isLocked":            mr.IsLocked,
		"isDraft":             mr.IsDraft,
		"authorId":            mr.Author.ID,
		"authorName":          mr.Author.Username,
		"repositoryNamespace": mr.Repository.Namespace,
		"repositoryName":      mr.Repository.Name,
		"repositoryPath":      mr.Repository.Path,
		"repositoryUrl":       mr.Repository.URL,
	}

	if strings.Contains(mr.Description, "| Package |") {
		updates := extractDependencyUpdates(mr.Description)
		/*
			if len(updates) > 0 {
				mrContext["dependencies"] = updates
			}
		*/
		if len(updates) == 1 {
			mrContext["dependencyType"] = updates[0].PackageType
			mrContext["dependencyCoordinate"] = updates[0].Coordinate
			mrContext["dependencyFrom"] = updates[0].From
			mrContext["dependencyTo"] = updates[0].To
		} else {
			mrContext["dependencyType"] = ""
			mrContext["dependencyCoordinate"] = ""
			mrContext["dependencyFrom"] = ""
			mrContext["dependencyTo"] = ""
		}
	}

	return mrContext
}

func extractDependencyUpdates(description string) []DependencyUpdate {
	var updates []DependencyUpdate

	re := regexp.MustCompile(`https://developer\.mend\.io/api/mc/badges/compatibility/([^/]+)/([^/]+)/([^/]+)/([^/?]+)`)

	matches := re.FindAllStringSubmatch(description, -1)
	for _, match := range matches {
		if len(match) >= 5 {
			updates = append(updates, DependencyUpdate{
				PackageType: match[1],
				Coordinate:  match[2],
				From:        match[3],
				To:          match[4],
			})
		}
	}

	return updates
}
