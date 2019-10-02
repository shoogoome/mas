package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "GetFileInfo",
            Router: `/api/file/info`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "ChunkUpload",
            Router: `/api/file/upload/chuck`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "Download",
            Router: `/api/file/upload/download`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "Finish",
            Router: `/api/file/upload/finish`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "InitFileInfo",
            Router: `/api/file/upload/init`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "UploadSingle",
            Router: `/api/file/upload/single`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "GenerateDownloadToken",
            Router: `/api/token/download`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"] = append(beego.GlobalControllerRouter["mas/controllers/file:FileSystemController"],
        beego.ControllerComments{
            Method: "GenerateUploadToken",
            Router: `/api/token/upload`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
