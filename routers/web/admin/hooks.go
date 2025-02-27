// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/context"
	"code.gitea.io/gitea/modules/setting"
)

const (
	// tplAdminHooks template path to render hook settings
	tplAdminHooks base.TplName = "admin/hooks"
)

// DefaultOrSystemWebhooks renders both admin default and system webhook list pages
func DefaultOrSystemWebhooks(ctx *context.Context) {
	var err error

	ctx.Data["PageIsAdminSystemHooks"] = true
	ctx.Data["PageIsAdminDefaultHooks"] = true

	def := make(map[string]interface{}, len(ctx.Data))
	sys := make(map[string]interface{}, len(ctx.Data))
	for k, v := range ctx.Data {
		def[k] = v
		sys[k] = v
	}

	sys["Title"] = ctx.Tr("admin.systemhooks")
	sys["Description"] = ctx.Tr("admin.systemhooks.desc")
	sys["Webhooks"], err = models.GetSystemWebhooks()
	sys["BaseLink"] = setting.AppSubURL + "/admin/hooks"
	sys["BaseLinkNew"] = setting.AppSubURL + "/admin/system-hooks"
	if err != nil {
		ctx.ServerError("GetWebhooksAdmin", err)
		return
	}

	def["Title"] = ctx.Tr("admin.defaulthooks")
	def["Description"] = ctx.Tr("admin.defaulthooks.desc")
	def["Webhooks"], err = models.GetDefaultWebhooks()
	def["BaseLink"] = setting.AppSubURL + "/admin/hooks"
	def["BaseLinkNew"] = setting.AppSubURL + "/admin/default-hooks"
	if err != nil {
		ctx.ServerError("GetWebhooksAdmin", err)
		return
	}

	ctx.Data["DefaultWebhooks"] = def
	ctx.Data["SystemWebhooks"] = sys

	ctx.HTML(http.StatusOK, tplAdminHooks)
}

// DeleteDefaultOrSystemWebhook handler to delete an admin-defined system or default webhook
func DeleteDefaultOrSystemWebhook(ctx *context.Context) {
	if err := models.DeleteDefaultSystemWebhook(ctx.FormInt64("id")); err != nil {
		ctx.Flash.Error("DeleteDefaultWebhook: " + err.Error())
	} else {
		ctx.Flash.Success(ctx.Tr("repo.settings.webhook_deletion_success"))
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"redirect": setting.AppSubURL + "/admin/hooks",
	})
}
