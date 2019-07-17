/*
 * Copyright © 2019 Hedzr Yeh.
 */

package gql_test

import (
	"encoding/json"
	"github.com/hedzr/awesome-tool/ags/gql"
	"io/ioutil"
	"testing"
)

func TestGhResultJson(t *testing.T) {
	b, err := ioutil.ReadFile("../../ags.result.json")
	if err != nil {
		t.Log(err)
		return
	}

	var respData = make(gql.GhRes)
	err = json.Unmarshal(b, &respData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(respData["Comcast_gaad"].CreatedAt)
}

func TestIsDigit(t *testing.T) {
	if true != gql.StartsWithDigit("360ff") {
		t.Fatal("err")
	}
	if false != gql.StartsWithDigit("fandu") {
		t.Fatal("err")
	}
}
