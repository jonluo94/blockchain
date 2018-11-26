package web

import (
	"github.com/astaxie/beego"
	"gitee.com/chaincode-app/hello-service/web/controllers"
	"gitee.com/chaincode-app/hello-service/blockchain"
	"github.com/astaxie/beego/context"
	"fmt"
)


func Run(sdk *blockchain.FabricSetup) {

	beego.LoadAppConfig("ini", "web/conf/app.conf")

	beego.Router("/", &controllers.WebController{})

	beego.Get("/get",func(ctx *context.Context){
		val,err := sdk.Query("hello")
		if err != nil {
			fmt.Errorf(err.Error())
		}
		ctx.Output.Body([]byte(val))
	})

	beego.Get("/put",func(ctx *context.Context){
		var value string
		ctx.Input.Bind(&value, "value")
		val,err := sdk.Invoke("set","hello",value)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		ctx.Output.Body([]byte(val))
	})


	beego.Run()
}
