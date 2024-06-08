package admin

import (
	"fiber-admin/internal/pkg/service/admin/mods"
)

type Admin struct {
	DocumentationService mods.DocumentationService
	LogsService          mods.LogsService
	NoticeService        mods.NoticeService
	UserService          mods.UserService
}
