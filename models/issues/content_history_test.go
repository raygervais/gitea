// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/modules/timeutil"

	"github.com/stretchr/testify/assert"
)

func TestContentHistory(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())

	dbCtx := db.DefaultContext
	dbEngine := db.GetEngine(dbCtx)
	timeStampNow := timeutil.TimeStampNow()

	_ = SaveIssueContentHistory(dbEngine, 1, 10, 0, timeStampNow, "i-a", true)
	_ = SaveIssueContentHistory(dbEngine, 1, 10, 0, timeStampNow.Add(2), "i-b", false)
	_ = SaveIssueContentHistory(dbEngine, 1, 10, 0, timeStampNow.Add(7), "i-c", false)

	_ = SaveIssueContentHistory(dbEngine, 1, 10, 100, timeStampNow, "c-a", true)
	_ = SaveIssueContentHistory(dbEngine, 1, 10, 100, timeStampNow.Add(5), "c-b", false)
	_ = SaveIssueContentHistory(dbEngine, 1, 10, 100, timeStampNow.Add(20), "c-c", false)
	_ = SaveIssueContentHistory(dbEngine, 1, 10, 100, timeStampNow.Add(50), "c-d", false)
	_ = SaveIssueContentHistory(dbEngine, 1, 10, 100, timeStampNow.Add(51), "c-e", false)

	h1, _ := GetIssueContentHistoryByID(dbCtx, 1)
	assert.EqualValues(t, 1, h1.ID)

	m, _ := QueryIssueContentHistoryEditedCountMap(dbCtx, 10)
	assert.Equal(t, 3, m[0])
	assert.Equal(t, 5, m[100])

	/*
		we can not have this test with real `User` now, because we can not depend on `User` model (circle-import), so there is no `user` table
		when the refactor of models are done, this test will be possible to be run then with a real `User` model.
	*/
	type User struct {
		ID   int64
		Name string
	}
	_ = dbEngine.Sync2(&User{})

	list1, _ := FetchIssueContentHistoryList(dbCtx, 10, 0)
	assert.Len(t, list1, 3)
	list2, _ := FetchIssueContentHistoryList(dbCtx, 10, 100)
	assert.Len(t, list2, 5)

	h6, h6Prev, _ := GetIssueContentHistoryAndPrev(dbCtx, 6)
	assert.EqualValues(t, 6, h6.ID)
	assert.EqualValues(t, 5, h6Prev.ID)

	// soft-delete
	_ = SoftDeleteIssueContentHistory(dbCtx, 5)
	h6, h6Prev, _ = GetIssueContentHistoryAndPrev(dbCtx, 6)
	assert.EqualValues(t, 6, h6.ID)
	assert.EqualValues(t, 4, h6Prev.ID)

	// only keep 3 history revisions for comment_id=100
	keepLimitedContentHistory(dbEngine, 10, 100, 3)
	list1, _ = FetchIssueContentHistoryList(dbCtx, 10, 0)
	assert.Len(t, list1, 3)
	list2, _ = FetchIssueContentHistoryList(dbCtx, 10, 100)
	assert.Len(t, list2, 3)
	assert.EqualValues(t, 7, list2[0].HistoryID)
	assert.EqualValues(t, 6, list2[1].HistoryID)
	assert.EqualValues(t, 4, list2[2].HistoryID)
}
