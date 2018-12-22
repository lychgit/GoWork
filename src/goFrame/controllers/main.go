package controllers

import "fmt"

type MainController struct {
	BaseController
}

func (this *MainController) Index() {
	fmt.Println("123456")
	this.TplName = "index.html"
}