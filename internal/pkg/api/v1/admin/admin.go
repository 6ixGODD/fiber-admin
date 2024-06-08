package admin

import (
	"fiber-admin/internal/pkg/api/v1/admin/mods"
)

type Admin struct {
	UserApi          *mods.UserApi
	NoticeApi        *mods.NoticeApi
	DocumentationApi *mods.DocumentationApi
	LogsApi          *mods.LogsApi
}
