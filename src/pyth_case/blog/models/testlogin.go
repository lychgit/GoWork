// testlogin
package models

import (
	"github.com/astaxie/beego"
	"fmt"
)

type TestLoginController struct {
	beego.Controller
}

func (c *TestLoginController) SelfTest() {
	c.Ctx.WriteString("this is myself  controller!")
}

func (c *TestLoginController) Login() {
	fmt.Println("login")
	name := c.Ctx.GetCookie("name")
	password := c.Ctx.GetCookie("password")

	if name != "" {
		c.Ctx.WriteString("read cookie ok ! username:"+ name + "password:" + password)
	} else {
		formData := `<html><form action="/test_login" method="post">
         <input type="text" name="Username">
         <input type="password" name="Password">
        <input type="submit" value="post">
            </html>
        `
		c.Ctx.WriteString(formData)
	}
}

func (c *TestLoginController) PostData() {
	fmt.Println("PostData")
	u := User{}
	if err := c.ParseForm(&u); err != nil {

	}
	//c.Ctx.SetCookie("name", u.Username, 100, "/")  // 设置cookie
	//c.Ctx.SetCookie("password", u.Password, 100, "/")  // 设置cookie
	//c.Ctx.WriteString("username:" + u.Username + "   password:" + u.Password)
}