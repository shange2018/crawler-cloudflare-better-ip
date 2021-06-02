package worker

import (
	"crawler/worker/fetcher"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"time"
)

type Request struct { // 请求类型
	ID			int
	Url			string                                  // 请求地址
	ParserFunc	func(int, time.Duration, []byte) ParseResult // 处理请求的解析函数
}

type ParseResult struct { // 解析网页后的返回值类型
	Requests []Request     // 新产生的请求
	Items    []interface{} // 解析出来的要持久化的数据
}

func Worker(workerID int, r Request) (ParseResult, error) {

	completionRate := gconv.Float32(r.ID / 1524706)
	elapsed, body, err := fetcher.Fetch(r.Url)
	if err != nil {
		glog.Errorf("wID: #%d, rID: #%d, Fetcher error: %v", workerID, r.ID, err)
		return ParseResult{}, err
	} else {
		glog.Infof("wID: #%-5d, rID: #%-7d, cRate: %10.7f, Fetcher ok %s", workerID, r.ID, completionRate, r.Url)
		return r.ParserFunc(r.ID, elapsed, body), nil
	}

}
