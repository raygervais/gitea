// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"testing"

	"code.gitea.io/gitea/models/db"
	"github.com/stretchr/testify/assert"
)

func TestNotice_TrStr(t *testing.T) {
	notice := &Notice{
		Type:        NoticeRepository,
		Description: "test description",
	}
	assert.Equal(t, "admin.notices.type_1", notice.TrStr())
}

func TestCreateNotice(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())

	noticeBean := &Notice{
		Type:        NoticeRepository,
		Description: "test description",
	}
	db.AssertNotExistsBean(t, noticeBean)
	assert.NoError(t, CreateNotice(noticeBean.Type, noticeBean.Description))
	db.AssertExistsAndLoadBean(t, noticeBean)
}

func TestCreateRepositoryNotice(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())

	noticeBean := &Notice{
		Type:        NoticeRepository,
		Description: "test description",
	}
	db.AssertNotExistsBean(t, noticeBean)
	assert.NoError(t, CreateRepositoryNotice(noticeBean.Description))
	db.AssertExistsAndLoadBean(t, noticeBean)
}

// TODO TestRemoveAllWithNotice

func TestCountNotices(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())
	assert.Equal(t, int64(3), CountNotices())
}

func TestNotices(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())

	notices, err := Notices(1, 2)
	assert.NoError(t, err)
	if assert.Len(t, notices, 2) {
		assert.Equal(t, int64(3), notices[0].ID)
		assert.Equal(t, int64(2), notices[1].ID)
	}

	notices, err = Notices(2, 2)
	assert.NoError(t, err)
	if assert.Len(t, notices, 1) {
		assert.Equal(t, int64(1), notices[0].ID)
	}
}

func TestDeleteNotice(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())

	db.AssertExistsAndLoadBean(t, &Notice{ID: 3})
	assert.NoError(t, DeleteNotice(3))
	db.AssertNotExistsBean(t, &Notice{ID: 3})
}

func TestDeleteNotices(t *testing.T) {
	// delete a non-empty range
	assert.NoError(t, db.PrepareTestDatabase())

	db.AssertExistsAndLoadBean(t, &Notice{ID: 1})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 2})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 3})
	assert.NoError(t, DeleteNotices(1, 2))
	db.AssertNotExistsBean(t, &Notice{ID: 1})
	db.AssertNotExistsBean(t, &Notice{ID: 2})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 3})
}

func TestDeleteNotices2(t *testing.T) {
	// delete an empty range
	assert.NoError(t, db.PrepareTestDatabase())

	db.AssertExistsAndLoadBean(t, &Notice{ID: 1})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 2})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 3})
	assert.NoError(t, DeleteNotices(3, 2))
	db.AssertExistsAndLoadBean(t, &Notice{ID: 1})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 2})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 3})
}

func TestDeleteNoticesByIDs(t *testing.T) {
	assert.NoError(t, db.PrepareTestDatabase())

	db.AssertExistsAndLoadBean(t, &Notice{ID: 1})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 2})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 3})
	assert.NoError(t, DeleteNoticesByIDs([]int64{1, 3}))
	db.AssertNotExistsBean(t, &Notice{ID: 1})
	db.AssertExistsAndLoadBean(t, &Notice{ID: 2})
	db.AssertNotExistsBean(t, &Notice{ID: 3})
}
