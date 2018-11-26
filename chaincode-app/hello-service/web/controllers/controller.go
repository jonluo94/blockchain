package controllers

import (
	"github.com/astaxie/beego"
)

type WebController struct {
	beego.Controller
}

func (this *WebController) Get() {
	this.TplName = "index.html"
	this.Render()
}

