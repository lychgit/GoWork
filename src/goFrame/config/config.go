/*
@Time : 2019/3/28 10:48 
@Author : shilinqing
@File : config
*/
package config

import "github.com/astaxie/beego"

var bConfig = beego.BConfig

type serverConfig struct {
	Debug bool
}

//初始化配置
func init()  {
	//是否路由忽略大小写匹配，默认是 true，区分大小写
	bConfig.RouterCaseSensitive = true
	//beego 服务器默认在请求的时候输出 server 为 LYCH
	bConfig.ServerName = "LYCH"
	//是否异常恢复，默认值为 true，即当应用出现异常的情况，通过 recover 恢复回来，而不会导致应用异常退出
	bConfig.RecoverPanic = true
	//是否允许在 HTTP 请求时，返回原始请求体数据字节，默认为 false （GET or HEAD or 上传文件请求除外）
	bConfig.CopyRequestBody = true
	//文件上传默认内存缓存大小，默认值是 1 << 26(64M)
	bConfig.MaxMemory = 1 << 26
	//是否显示系统错误信息，默认为 true
	bConfig.EnableErrorsShow = true
	//是否将错误信息进行渲染，默认值为 true，即出错会提示友好的出错页面，对于 API 类型的应用可能需要将该选项设置为 false 以阻止在 dev 模式下不必要的模板渲染信息返回
	bConfig.EnableErrorsShow = true
	//是否开启进程内监控模块，默认 false 关闭。
	bConfig.Listen.EnableAdmin = false
	//监控程序监听的地址，默认值是 8088
	//bConfig.Listen.AdminPort = 8088
	//是否输出日志到 Log，默认在 prod 模式下不会输出日志，默认为 false 不输出日志。此参数不支持配置文件配置
	bConfig.Log.AccessLogs = false
	//是否在日志里面显示文件名和输出日志行号，默认 true。此参数不支持配置文件配置
	bConfig.Log.FileLineNum = true
	//日志输出配置，参考 logs 模块，console file 等配置，此参数不支持配置文件配置。
	bConfig.Log.Outputs = map[string]string{"console": ""} // or   bConfig.Log.Outputs["console"] = ""

	/******************************		web 配置	******************************/

	//是否模板自动渲染，默认值为 true，对于 API 类型的应用，应用需要把该选项设置为 false，不需要渲染模板。
	bConfig.WebConfig.AutoRender = true

	//当你设置了自动渲染，然后在你的 Controller 中没有设置任何的 TplName，那么 beego 会自动设置你的模板文件如下：
	//c.TplName = strings.ToLower(c.controllerName) + "/" + strings.ToLower(c.actionName) + "." + c.TplExt
	//也就是你对应的 Controller 名字+请求方法名.模板后缀，也就是如果你的 Controller 名是 AddController，
	// 请求方法是 POST，默认的文件后缀是 tpl，那么就会默认请求 /viewpath/AddController/post.tpl 文件。

	//静态文件目录设置，默认是static  可配置单个或多个目录:
	beego.SetStaticPath("/static", "static") //设置静态文件处理目录
	beego.SetStaticPath("/images", "static/images")
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/js", "static/js")

	//beego.SetViewsPath("templatePath") //设置模板目录
	beego.SetViewsPath("views") //设置模板目录

	//是否开启 XSRF，默认为 false，不开启
	//bConfig.WebConfig.EnableXSRF = false
	//XSRF 的 key 信息，默认值是 beegoxsrf。 EnableXSRF＝true 才有效
	//bConfig.WebConfig.XSRFKEY = "beegoxsrf"

	/******************************		监听配置	******************************/

	//是否启用 HTTP 监听，默认是 true
	bConfig.Listen.EnableHTTP = true
	//应用监听地址，默认为空，监听所有的网卡 IP
	bConfig.Listen.HTTPAddr = "127.0.0.1"
	//应用监听端口，默认为 8080
	bConfig.Listen.HTTPPort = 8080
	//是否启用 HTTPS，默认是 false 关闭。当需要启用时，先设置 EnableHTTPS = true，并设置 HTTPSCertFile 和 HTTPSKeyFile
	bConfig.Listen.EnableHTTPS = false
	//应用监听地址，默认为空，监听所有的网卡 IP
	bConfig.Listen.HTTPSAddr = ""
	//应用监听端口，默认为 10443
	bConfig.Listen.HTTPSPort = 10443
	//开启 HTTPS 后，ssl 证书路径，默认为空
	bConfig.Listen.HTTPSCertFile = "conf/ssl.crt"
	//开启 HTTPS 之后，SSL 证书 keyfile 的路径
	bConfig.Listen.HTTPSKeyFile = "conf/ssl.key"

	//是否启用 fastcgi ， 默认是 false
	bConfig.Listen.EnableFcgi = false
	//通过fastcgi 标准I/O，启用 fastcgi 后才生效，默认 false
	bConfig.Listen.EnableStdIo = false

	/******************************	Session配置	******************************/

	bConfig.WebConfig.Session.SessionOn = true

	//读取不同模式下配置参数的方法是“模式::配置参数名”，比如：beego.AppConfig.String(“dev::mysqluser”)。
	//对于自定义的参数，需使用 beego.GetConfig(typ, key string, defaultVal interface{}) 来获取指定 runmode 下的配置（需 1.4.0 以上版本），typ 为参数类型，key 为参数名, defaultVal 为默认值

	//我们的程序往往期望把信息输出到 log 中，现在设置输出到文件很方便，如下所示：
	beego.SetLogger("file", `{"filename": "runtime/logs/debug.log"}`);

}