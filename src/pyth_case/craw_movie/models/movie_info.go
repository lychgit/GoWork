package models

import (
	"github.com/astaxie/beego/orm"
	_"github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

var (
	db orm.Ormer
)

type MovieInfo struct {
	Id int64  `orm:"auto"`
	Movie_id int64
	Movie_name string
	Movie_pic string
	Movie_director string
	Movie_writer string
	Movie_country string
	Movie_language string
	Movie_main_character string
	Movie_type string
	Movie_on_time string
	Movie_span string
	Movie_grade	string
	Remark string
	_Create_time string
}

func init(){
	orm.Debug = true
	orm.RegisterDataBase("default", "mysql","root:123456@tcp(127.0.0.1:3306)/tprbac?charset=utf8")
	orm.RegisterModel(new(MovieInfo))
	db = orm.NewOrm()
}

func AddMovie(movie_info *MovieInfo) (int64 ,error){
	movie_info.Id = 0
	id ,err := db.Insert(movie_info)
	return id ,err
}

/*正则匹配*/
//电影导演
func GetMovieDirector(html string) string{
	if html == ""{
		return ""
	}
	//'<a.*rel="v:directedBy">(.*)</a>'
	//.* : 匹配除换行以外的任意字符 重复0次或更多次
	reg := regexp.MustCompile(`<a.*rel="v:directedBy">(.*)</a>`)
	result := reg.FindAllStringSubmatch(html, -1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])

}
//电影名称
func GetMovieName(html string) string{
	if html == ""{
		return ""
	}
	//'<span property="v:itemreviewed">秒速5厘米 秒速5センチメートル</span>'
	//空格 : \s*
	reg := regexp.MustCompile(`<span.*?property="v:itemreviewed">(.*)</span>`)
	result := reg.FindAllStringSubmatch(html, -1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}

//主演
func GetMovieMainCharacters(html string) string{
	if html == ""{
		return ""
	}
	//<a href="/celebrity/1015108/" rel="v:starring">水桥研二</a>
	// ? : 重复0次或1次
	reg := regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(html, -1)
	mainCharacters := ""
	for _,v := range result{
		mainCharacters += v[1] + "/"
	}
	//截取最后一位
	return strings.Trim(mainCharacters,"/")
}

//类型
func GetMovieGenre(html string) string{
	if html == ""{
		return ""
	}
	//<span property="v:genre">剧情</span>
	reg := regexp.MustCompile(`<span.*?property="v:genre">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(html, -1)
	movieTypes := ""
	for _,v := range result{
		movieTypes += v[1] + "/"
	}
	return strings.Trim(movieTypes,"/")

}

//国家
func GetMovieCountry(html string) string{
	//<span class="pl">制片国家/地区:</span> 日本<br/>
	reg := regexp.MustCompile(`<span class="pl">制片国家/地区:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(html,-1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}

//封面
func GetMoviePic(html string) string{
	//<img src="https://img3.doubanio.com/view/photo/s_ratio_poster/public/p982896012.webp" title="点击看更多海报" alt="秒速5センチメートル" rel="v:image" />
	reg := regexp.MustCompile(`<img src="(.*?)".*?rel="v:image" />`)
	result := reg.FindAllStringSubmatch(html,-1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}

//语言
func GetMovieLanguage(html string) string{
	//<span class="pl">语言:</span> 日语<br/>
	reg := regexp.MustCompile(`<span class="pl">语言:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(html,-1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}

//上映时间
func GetMovieOnTime(html string) string{
	//<span property="v:initialReleaseDate" content="2007-03-03(日本)">2007-03-03(日本)</span>
	reg := regexp.MustCompile(`<span.*?property="v:initialReleaseDate".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(html, -1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}

//时长
func GetMovieRunningTime(html string) string{
	//<span property="v:runtime" content="63">63分钟</span>
	reg := regexp.MustCompile(`<span.*?property="v:runtime".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(html, -1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}


//评分
func GetMovieGrade(html string) string{
	//<strong class="ll rating_num" property="v:average">8.3</strong>
	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*?)</strong>`)
	result := reg.FindAllStringSubmatch(html, -1)
	if len(result) == 0{
		//判断长度是否等于0，内容是否为空
		return ""
	}
	return string(result[0][1])
}

//其他电影链接
func GetMovieUrls(html string) []string{
	//<a href="https://movie.douban.com/subject/1937946/?from=subject-page" class="" >穿越时空的少女</a>
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/.*?)"`)
	result := reg.FindAllStringSubmatch(html, -1)
	var movieSets []string
	for _, v := range result{
		movieSets = append(movieSets, v[1])
	}
	return movieSets

}