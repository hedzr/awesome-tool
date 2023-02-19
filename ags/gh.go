/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ags

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/hedzr/awesome-tool/ags/gql"
	"github.com/hedzr/cmdr"
	"github.com/hedzr/log/dir"
)

func writeJsonFile(filename string, obj interface{}) {
	b, err := json.Marshal(obj)
	if err == nil {
		_ = ioutil.WriteFile(filename, b, 0644)
	}
}

func readJsonFile(filename string, obj interface{}) {
	b, err := ioutil.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(b, obj)
	}
}

func sortMap(rrr map[string]*gql.Ghr) (keys []string) {
	keys = make([]string, 0, len(rrr))
	for key := range rrr {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if rrr[keys[i]].Ghf == nil {
			cmdr.Logger.Warnf("rrr[keys[i]].Ghf == nil, i=%v, keys[i]=%v, rrr[keys[i]]=%v", i, keys[i], rrr[keys[i]])
			return false
		}
		if rrr[keys[j]].Ghf == nil {
			cmdr.Logger.Warnf("rrr[keys[j]].Ghf == nil, j=%v, keys[j]=%v, rrr[keys[j]]=%v", j, keys[j], rrr[keys[j]])
			return true
		}
		return rrr[keys[i]].Ghf.Stargazers.TotalCount > rrr[keys[j]].Ghf.Stargazers.TotalCount
	})
	return
}

func httpGet(url, filepath string, cacheTime string) (err error) {
	var (
		dur time.Duration
		fi  os.FileInfo
	)
	dur, err = time.ParseDuration(cacheTime)
	if dir.FileExists(filepath) {
		if fi, err = os.Stat(filepath); err == nil {
			cmdr.Logger.Debugf("fi: %v, now: %v, fi+%v: %v", fi.ModTime(), time.Now(), cacheTime, fi.ModTime().Add(dur))
			if fi.ModTime().Add(dur).Before(time.Now()) {
				cmdr.Logger.Infof("cache expired, force http download: %v", url)
				goto goAhead
			}
		}
		cmdr.Logger.Infof("skip http download: %v", url)
		return
	}

goAhead:
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	var written int64
	written, err = io.Copy(out, resp.Body)
	if err == nil {
		cmdr.Logger.Infof("%v downloaded, %v bytes.", filepath, written)
	}
	return
}
