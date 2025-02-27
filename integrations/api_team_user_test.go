// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrations

import (
	"net/http"
	"testing"
	"time"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/modules/convert"
	api "code.gitea.io/gitea/modules/structs"
	"github.com/stretchr/testify/assert"
)

func TestAPITeamUser(t *testing.T) {
	defer prepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session)
	req := NewRequest(t, "GET", "/api/v1/teams/1/members/user1?token="+token)
	session.MakeRequest(t, req, http.StatusNotFound)

	req = NewRequest(t, "GET", "/api/v1/teams/1/members/user2?token="+token)
	resp := session.MakeRequest(t, req, http.StatusOK)
	var user2 *api.User
	DecodeJSON(t, resp, &user2)
	user2.Created = user2.Created.In(time.Local)
	user := db.AssertExistsAndLoadBean(t, &models.User{Name: "user2"}).(*models.User)

	expectedUser := convert.ToUser(user, user)

	// test time via unix timestamp
	assert.EqualValues(t, expectedUser.LastLogin.Unix(), user2.LastLogin.Unix())
	assert.EqualValues(t, expectedUser.Created.Unix(), user2.Created.Unix())
	expectedUser.LastLogin = user2.LastLogin
	expectedUser.Created = user2.Created

	assert.Equal(t, expectedUser, user2)
}
