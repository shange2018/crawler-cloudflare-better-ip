package fetcher

import (
	"bufio"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/glog"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var rateLimiter = time.Tick(1 * time.Millisecond)

func Fetch(url string) (time.Duration, []byte, error) {

	var err error
	var resp *http.Response
	var elapsed time.Duration
	var retry = 2	// http.Get失败后重试的次数

	for i:=0; i<retry; i++ {
		<-rateLimiter
		startTime := time.Now()
		resp, err = http.Get(url)
		elapsed = time.Since(startTime)
		if err == nil {
			break
		}
		glog.Infof("http.Get %s, %dst retry",url, i + 1)
	}
	if err != nil {
		return elapsed, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			glog.Error(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err = gerror.NewCodef(resp.StatusCode,"error StatusCode")
		return elapsed, nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader,
		e.NewDecoder())
	result, err := ioutil.ReadAll(utf8Reader)
	return elapsed, result, err
}

// 读取网页源码前1024个字节，返回网页编码格式
func determineEncoding(r *bufio.Reader) encoding.Encoding {

	bytes, err := r.Peek(1024)
	if err != nil && err != io.EOF {
		log.Printf("Peek error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(
		bytes, "")
	return e
}