package models

import (
	"github.com/astaxie/beego"
	"fmt"
)

type LoginaController struct {
	beego.Controller
}


func (con *LoginaController) Logina(){
	fmt.Println("login")
	//get cookie
	name := con.Ctx.GetCookie("name")
	password := con.Ctx.GetCookie("password")

	//do verify password
	if name != "" {
		con.Ctx.WriteString("read cookie ok ! username:"+ name + "password:" + password)

	}else{
		con.Ctx.WriteString(
			`<html><form action="http://127.0.0.1:8080/login" method="post">
			<input type="text" name="Username" />
			<input type="password" name="Password" />
			<input type="submit" value="submit" />
		</form></html>`)
	}
}

func (con *LoginaController) Posta(){
	fmt.Println("Post")
	u:= User{}
	if err := con.ParseForm(&u); err != nil{
		//panic(err)
	}
	////set cookie
	//con.Ctx.SetCookie("name",u.Username,1000,"/")
	//con.Ctx.SetCookie("password",u.Password,1000,"/")
	//name := con.Ctx.GetCookie("name")
	//password := con.Ctx.GetCookie("password")
	//con.Ctx.WriteString("cookie username:"+ name + "password:" + password)

}