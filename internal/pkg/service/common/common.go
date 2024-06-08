package common

import (
	"fiber-admin/internal/pkg/service/common/mods"
)

type Common struct {
	AuthService          mods.AuthService
	DocumentationService mods.DocumentationService
	NoticeService        mods.NoticeService
	ProfileService       mods.ProfileService
}
