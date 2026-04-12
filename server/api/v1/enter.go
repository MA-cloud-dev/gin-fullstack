package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/campus"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	CampusApiGroup  campus.ApiGroup
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
}
