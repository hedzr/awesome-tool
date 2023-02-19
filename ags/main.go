/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ags

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hedzr/awesome-tool/ags/gh"
	"github.com/hedzr/awesome-tool/ags/gql"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
	"gopkg.in/hedzr/errors.v3"

	// "github.com/russross/blackfriday"
	"github.com/russross/blackfriday/v2"

	// graphql "github.com/graph-gophers/graphql-go"
	// "github.com/graph-gophers/graphql-go/relay"

	"github.com/machinebox/graphql"
)

// type query struct{}
//
// func (_ *query) Hello() string { return "Hello, world!" }

func Main() (err error) {
	// logrus.SetLevel(logrus.DebugLevel)
	// logex.Enable()

	topUrl := cmdr.GetStringR("build.one.source")
	name := cmdr.GetStringR("build.one.name")
	if len(name) == 0 {
		if i := strings.LastIndex(topUrl, "/"); i >= 0 {
			name = topUrl[i+1:]
		} else {
			return errors.New("'name' must NOT be empty string.")
		}
	}

	workDir := path.Join(cmdr.GetStringR("build.one.work-dir"), name)
	README := fmt.Sprintf("%v/raw/master/README.md", topUrl)
	readmeLocal := fmt.Sprintf("%v/%v", workDir, "index.md")
	outLocal := fmt.Sprintf("%v/%v", workDir, "output.md")
	ignoreUsernamesFilename := path.Join(workDir, "ignore.users.json")
	ignoreReposFilename := path.Join(workDir, "ignore.repos.json")

	var (
		ignoreUsernames []string
		ignoreRepos     []string
		input           []byte
		out             *os.File
		loop            = 0
	)

	if err = dir.EnsureDir(workDir); err != nil {
		return
	}

	fmt.Printf(`
    name: %v
     url: %v
work dir: %v
`, name, topUrl, workDir)

	infoMd := path.Join(workDir, "info.md")
	_ = ioutil.WriteFile(infoMd, []byte(fmt.Sprintf("# %v\n\n%v", name, topUrl)), 0644)

	if err = httpGet(README, readmeLocal, "72h"); err != nil {
		return
	}

	if input, err = ioutil.ReadFile(readmeLocal); err != nil {
		return
	}

	readJsonFile(ignoreUsernamesFilename, &ignoreUsernames)
	readJsonFile(ignoreReposFilename, &ignoreRepos)

	r := newMarkdownRenderer()
	output := blackfriday.Run(input, blackfriday.WithRenderer(r))

	cmdr.Logger.Infof("output: %v", output)
	// cmdr.Logger.Infof("r: %v", r)

	if out, err = os.Create(outLocal); err != nil {
		return
	}
	defer out.Close()

	for _, sec := range r.(*mdRenderer).sections {
		// sp := strings.Repeat("  ", sec.Level-1)
		hd := strings.Repeat("#", sec.Level)
		_, _ = fmt.Fprintf(out, "\n%v %v\n\n", hd, sec.Header)

		rrr, cnt, ql, respData := make(map[string]*gql.Ghr), 0, "", make(gql.GhRes)
		jsonFile := fmt.Sprintf("%v/pre/result.%v.json", workDir, sec.Header)
		_ = dir.EnsureDir(path.Dir(jsonFile))

		if dir.FileExists(jsonFile) {
			cmdr.Logger.Debugf("  ..loading existed json result from %v", jsonFile)
			var b []byte
			b, err = ioutil.ReadFile(jsonFile)
			if err != nil {
				// cmdr.Logger.Fatal(err)
				return
			}

			err = json.Unmarshal(b, &respData)
			if err != nil {
				// cmdr.Logger.Fatal(err)
				return
			}

		} else {
		RetryQl:
			ql, rrr, cnt = gql.BuildQl(sec, ignoreUsernames, ignoreRepos)
			if cnt == 0 {
				continue
			}

			respData = make(gql.GhRes)

			cmdr.Logger.Debugf("  ..querying with github api. %v", sec.Header)

			// schema := graphql.MustParseSchema(ql, &query{})
			// http.Handle("/query", &relay.Handler{Schema: schema})
			// cmdr.Logger.Fatal(http.ListenAndServe(":8080", nil))

			gqlFile := fmt.Sprintf("%v/gql/%v", workDir, "query.graphql")
			_ = dir.EnsureDir(path.Dir(gqlFile))
			_ = ioutil.WriteFile(gqlFile, []byte(ql), 0644)

			client := graphql.NewClient(gh.ApiEntryV4)
			client.Log = func(s string) { cmdr.Logger.Debugf("%v", s) }
			req := gh.ApplyToken(graphql.NewRequest(ql))
			// req.Var("key", "value")
			cmdr.Logger.Debugf("  ..querying now")
			if err = client.Run(context.Background(), req, &respData); err != nil {
				// graphql: Could not resolve to a User with the username 'Equilibrium-Games'
				if strings.Index(err.Error(), "Could not resolve to a User with the username") >= 0 {
					if i := strings.Index(err.Error(), "'"); i > 0 {
						n := err.Error()[i+1:]
						if i := strings.Index(n, "'"); i > 0 {
							n = n[0:i]
						}
						ignoreUsernames = append(ignoreUsernames, n)
						writeJsonFile(ignoreUsernamesFilename, ignoreUsernames)
						goto RetryQl
					}
				} else if strings.Index(err.Error(), "Could not resolve to a Repository with the name") >= 0 {
					if i := strings.Index(err.Error(), "'"); i > 0 {
						n := err.Error()[i+1:]
						if i := strings.Index(n, "'"); i > 0 {
							n = n[0:i]
						}
						ignoreRepos = append(ignoreRepos, n)
						writeJsonFile(ignoreReposFilename, ignoreRepos)
						goto RetryQl
					}
				}
				cmdr.Logger.Errorf("error: %v", err)
				return
			}

			cmdr.Logger.Debugf("    resp: %v", respData)
			b, _ := json.MarshalIndent(respData, "", "  ")
			if err = ioutil.WriteFile(jsonFile, b, 0644); err != nil {
				// cmdr.Logger.Fatal(err)
				return
			}
		}

		// connect rrr and respData
		for k, v := range respData {
			aa := strings.Split(k, gh.SepLine)
			index, owner, repo := 0, aa[0], aa[1]

			for i, s := range sec.List {
				if strings.HasPrefix(s.Url, "https://github.com/") {
					sa := strings.Split(s.Url[len("https://github.com/"):], "/")
					if len(sa) > 1 {
						user, rep := strings.ReplaceAll(sa[0], "-", "_"), strings.ReplaceAll(sa[1], "-", "_")
						if gql.StartsWithDigit(user) {
							user = gh.SepDigit + user
						}
						if user == owner && rep == repo {
							index = i
							break
						}
					}
				}
			}

			xv := new(gql.GhFrag)
			*xv = v
			rrr[k] = &gql.Ghr{Owner: owner, Repo: repo, RepoUrl: "", Index: index, Ghf: xv}
		}

		// sort rrr by stars
		keys := sortMap(rrr)

		_, _ = fmt.Fprint(out, `<!--
| No.  | Name                                          | Category                 | Star | Fork | Commits | Contributors | License |
| ---- | --------------------------------------------- | ------------------------ | ---- | ---- | ------- | ------------ | ------- |
-->
<table><tr><th>No.</th><th>Name</th><th>Category</th><th>Star</th><th>Fork</th><th>Commits</th><th>Contributors</th><th>License</th><th>Updated At</th></tr>

`)

		for i, k := range keys {
			v := rrr[k]
			s := sec.List[v.Index]
			if v.Repo == "cmdr" {
				cmdr.Logger.Printf("")
			}
			v.RepoUrl = fmt.Sprintf("%s/%s", v.Owner, v.Repo)

			_, _ = fmt.Fprintf(out, `
<tr><td>%d</td><td><a href='%v'>%v</a></td><td><font size='-2'>%v</font></td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td><font size='-2'>%v</font></td></tr>
<tr><td colspan='9'><i>%v</i></td></tr>
`,
				i+1, s.Url, v.RepoUrl,
				sec.Header,
				v.Ghf.Stargazers.TotalCount,
				v.Ghf.ForkCount,
				v.Ghf.DefaultBranchRef.Target.History.TotalCount,
				v.Ghf.Collaborators.TotalCount,
				v.Ghf.LicenseInfo.NormalName(),
				v.Ghf.UpdatedAt,
				s.Desc)
		}
		_, _ = fmt.Fprint(out, "</table>\n\n")
		_ = out.Sync()

		loop++
		if stopLoops := cmdr.GetIntR("build.one.loops"); stopLoops <= 0 || loop < stopLoops {
			continue
		}
		os.Exit(0)

		//if cmdr.GetBoolR("build.one.first-loop") {
		//	if loop >= 1 {
		//		os.Exit(0)
		//	}
		//} else if cmdr.GetBoolR("build.one.2nd-loop") {
		//	if loop >= 2 {
		//		os.Exit(0)
		//	}
		//} else if cmdr.GetBoolR("build.one.5th-loop") {
		//	if loop >= 5 {
		//		os.Exit(0)
		//	}
		//}

	}

	return
}
