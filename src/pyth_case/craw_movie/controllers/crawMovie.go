package controllers

import (
	"github.com/astaxie/beego"
	"goFrame/admin/models"
	"github.com/astaxie/beego/httplib"
	"time"
	"fmt"
	)

type CrawMovieController struct{
	beego.Controller
}

func (c *CrawMovieController) CrawMovie(){
	//连接redis
	models.ConnectRedis("127.0.0.1:6379")
	//爬虫入口url
	url := "https://movie.douban.com/subject/2043546/?tag=%E6%97%A5%E6%9C%AC&from=gaia_video"
	//加入队列
	models.PutInQueue(url)
	//实例化结构体
	var movieInfo models.MovieInfo

	for {
		//获取队列长度
		length :=  models.GetQueueLength()
		if length == 0{
			//如果队列为空，则退出循环
			break
		}
		//从队列中获取url
		url = models.PopFromQueue()
		//判断url是否已经被访问过
		if models.IsVisit(url){
			continue
		}

		rsp := httplib.Get(url)
		html ,err := rsp.String()
		if err != nil{
			panic(err)
		}
		movieInfo.Movie_name = models.GetMovieName(html)
		if movieInfo.Movie_name != ""{
			//存储数据
			movieInfo.Movie_director = models.GetMovieDirector(html)
			movieInfo.Movie_main_character = models.GetMovieMainCharacters(html)
			movieInfo.Movie_type = models.GetMovieGenre(html)
			movieInfo.Movie_on_time = models.GetMovieOnTime(html)
			movieInfo.Movie_grade = models.GetMovieGrade(html)
			movieInfo.Movie_span = models.GetMovieRunningTime(html)
			movieInfo.Movie_pic = models.GetMoviePic(html)
			movieInfo.Movie_country = models.GetMovieCountry(html)
			movieInfo.Movie_language = models.GetMovieLanguage(html)
			//传变量参数指针
			id , _:= models.AddMovie(&movieInfo)
			fmt.Println(id)
			c.Ctx.WriteString(string(id))

		}
		//提取该页面的所有链接
		urls := models.GetMovieUrls(html)
		for _, sUrl := range urls{
			models.PutInQueue(sUrl) //放入redis队列

		}
		//提取成功后，记录到已访问的队列中
		models.AddToSet(url)
		time.Sleep(time.Second)

	}
	c.Ctx.WriteString("end of crawl!")

}
