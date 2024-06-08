package api

import (
	"fiber-admin/internal/pkg/api/v1/admin"
	"fiber-admin/internal/pkg/api/v1/common"
)

type Api struct {
	AdminApi  *admin.Admin
	CommonApi *common.Common
}
