package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["mas/controllers/system:SystemController"] = append(beego.GlobalControllerRouter["mas/controllers/system:SystemController"],
        beego.ControllerComments{
            Method: "Signal",
            Router: `/api/server/signal`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
