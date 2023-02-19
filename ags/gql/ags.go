/*
 * Copyright © 2019 Hedzr Yeh.
 */

package gql

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/hedzr/awesome-tool/ags/gh"
)

type (
	Ghr struct {
		Owner   string
		Repo    string
		RepoUrl string
		Index   int
		Ghf     *GhFrag
	}

	GhRes map[string]GhFrag

	GhFrag struct {
		DefaultBranchRef DefaultBranchRef
		Collaborators    Collaborators
		ForkCount        int
		Issues           Issues
		Name             string
		Stargazers       Stargazers
		Watchers         Watchers
		LicenseInfo      LicenseInfo
		CreatedAt        time.Time
		UpdatedAt        time.Time
		PushedAt         time.Time
	}

	LicenseInfo struct {
		Featured               bool
		Key, Name, SpdxId, Url string
	}

	DefaultBranchRef struct {
		Name   string
		Target Target
	}
	Target struct {
		History History
		Id      string
	}
	History struct {
		TotalCount int
	}
	Issues struct {
		TotalCount int
	}
	Stargazers struct {
		TotalCount int
	}
	Watchers struct {
		TotalCount int
	}
	Collaborators struct {
		TotalCount int
	}

	Section struct {
		Header string
		Level  int
		List   []*ListItem
	}
	ListItem struct {
		Title string
		Url   string
		Desc  string
	}
)

func (s LicenseInfo) NormalName() string {
	if len(s.SpdxId) > 0 {
		return s.SpdxId
	}
	if len(s.Key) > 0 {
		return s.Key
	}
	return s.Name
}

func BuildQl(sec *Section, ignoreUsernames, ignoreRepos []string) (ql string, rrr map[string]*Ghr, count int) {
	rrr = make(map[string]*Ghr)
	ql = "{\n"
	for ix, s := range sec.List {
		if strings.HasPrefix(s.Url, "https://github.com/") {
			sa := strings.Split(s.Url[len("https://github.com/"):], "/")
			if len(sa) < 2 {
				continue
			}
			if strings.Index(sa[1], "#") > 0 {
				// https://github.com/Boris-Em/BEMCheckBox#sample-app
				continue
			}
			user, repo := sa[0], sa[1]

			igr := false
			for _, un := range ignoreUsernames {
				if un == user {
					igr = true
				}
			}
			if !igr {
				for _, un := range ignoreRepos {
					if un == repo {
						igr = true
					}
				}
			}
			if igr {
				continue
			}

			safeName := SafeGqlFieldName(user, repo)
			rrr[safeName] = &Ghr{Owner: user, Repo: repo, RepoUrl: s.Url, Index: ix}
			ql += fmt.Sprintf(`  %v: repository(owner: "%v", name: "%v") {
    ...RepoFragment
  }
`, safeName, user, repo)
			count++
		}
	}
	if count == 0 {
		ql = ""
		return
	}
	ql += "}\n\n"

	ql += `
fragment RepoFragment on Repository {
  name
  licenseInfo{
    featured, key, 
    name, spdxId, url
  }
  # collaborators{
  #  totalCount
  # }
  forkCount
  stargazers{
    totalCount
  }
  watchers{
    totalCount
  }
  issues{
    totalCount
  }
  defaultBranchRef {
    name
    target {
      ... on Commit {
        id
        history(first: 0) {
          totalCount
        }
      }
    }
  }
  createdAt
  updatedAt
  pushedAt
}
`
	return
}

func SafeGqlFieldName(name, repo string) string {
	if StartsWithDigit(name) {
		name = gh.SepDigit + name
	}

	s := fmt.Sprintf("%v%s%v", name, gh.SepLine, repo)
	return strings.ReplaceAll(strings.ReplaceAll(s, "-", "_"), ".", "_")
}

func StartsWithDigit(s string) bool {
	for _, c := range s {
		return unicode.IsDigit(c)
	}
	return false
}

func StartsWithDigitHeavy(s string) bool {
	m, err := regexp.MatchString("^\\d", s)
	if err != nil {
		return false
	}
	return m
}

// var (
// 	ignoreUsernames = []string{
// 		"aurelien-rainone",
// 		"themester",
// 		"Obaied",
// 		"go-rtc",
// 		"pions",
// 		"PromonLogicalis",
// 		"tuvistavie",
// 		"trending",
// 		"dietsche",
// 		"tockins",
// 	}
//
// 	ignoreRepos = []string{
// 		"zerver",
// 	}
// )
