package models

import (
	"github.com/astaxie/goredis"
)

var (
	client goredis.Client
)

const(
	URL_QUEUE = "url_queue"
	URL_VISIT_SET = "url_visit_set"
)

/**
连接redis
 */
func ConnectRedis(addr string){
	client.Addr = addr
}

/**
加入队列
 */
func PutInQueue(url string){
	client.Lpush(URL_QUEUE,[]byte(url))
}

/**
（带阻塞）从队列中返回url
 */
func PopFromQueue() string{
	res, err := client.Rpop(URL_QUEUE)
	if err != nil{
		panic(err)
	}
	return string(res)
}
/**
获取长度
 */
func GetQueueLength() int{
	length , err := client.Llen(URL_QUEUE)
	if err != nil{
		panic(err)
	}
	return length
}

/**
添加已访问url队列
 */
func AddToSet(url string) {
	client.Sadd(URL_VISIT_SET, []byte(url))
}

/**
判断是否访问过
 */
func IsVisit(url string) bool{
	IsVisit, err := client.Sismember(URL_VISIT_SET, []byte(url))
	if err != nil{
		return false
	}
	return IsVisit
}



