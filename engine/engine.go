package engine

import (
	"crawler/scheduler"
	"crawler/worker"
	"crawler/worker/parser"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
)

type Engine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemChan    chan interface{} // 数据持久化通道
}

type Scheduler interface {
	scheduler.ReadyNotifier             // worker空闲通知
	Submit(worker.Request)              // worker提交新请求
	GetWorkerChan() chan worker.Request // 引擎向调度器获取当前可用worker通道
	Run()
}

var DB gdb.DB
var Table string
var Field string

func  init(){

	glog.SetFlags(glog.F_TIME_TIME)
	glog.Stack(false)
	_ = g.Cfg().SetPath("config")
	g.Cfg().SetFileName("config.toml")

	DB = g.DB("production")
	dblog := glog.New()
	dblog.SetFlags(glog.F_TIME_TIME)
	DB.SetLogger(dblog)

	Table = "ip_trace"
	Field = "id,h"

}

func (e *Engine) Run(seeds ...worker.Request) {

	out := make(chan worker.ParseResult) // 工作完成后返回值通道，共用
	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		createWorker(i, e.Scheduler.GetWorkerChan(), out, e.Scheduler)
	}

	e.getSeeds()

	for {
		parseResult := <-out // 工作完成后的返回内容
		for _, item := range parseResult.Items {
			go func() { e.ItemChan <- item }() // 返回内容送到Item持久化通道
		}

		for _, request := range parseResult.Requests {
			e.Scheduler.Submit(request) // worker新产生的request提交到调度器
		}
	}
}

func createWorker(workerID int, in chan worker.Request, out chan worker.ParseResult, ready scheduler.ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in
			result, err := worker.Worker(workerID, request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

// 循环读取数据库，获得种子请求
func (e *Engine) getSeeds() {

	go func() {

		limit := 2000  // 每次取多少条记录
		offset := 0 // 跳过多少条记录开始取数
		for {
			res, _ := DB.Table(Table).
				Fields(Field).
				Limit(limit).
				Offset(offset).
				All()

			if res.IsEmpty() {
				return
			}
			offset = offset + len(res.List())

			for _, v := range res.List() {
				e.Scheduler.Submit(worker.Request{
					ID:         gconv.Int(v["id"]),
					Url:        "http://" + gconv.String(v["h"]) + "/cdn-cgi/trace",
					ParserFunc: parser.ParseCloudFlareIPTrace,
				})
			}

		}
	}()

}
