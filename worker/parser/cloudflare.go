package parser

import (
	"crawler/worker"
	"github.com/gogf/gf/util/gconv"
	"strings"
	"time"
)

func ParseCloudFlareIPTrace(id int, elapsed time.Duration, contents []byte) worker.ParseResult {

	ipTrace := make(map[string]string)
	body := strings.TrimSpace(string(contents))
	s := strings.Split(body, "\n")
	for _, v := range s {
		s1 := strings.Split(v, "=")
		if s1[0] != "ip" && s1[0] != "uag" {
			ipTrace[s1[0]] = s1[1]
		}
	}
	ipTrace["id"] = gconv.String(id)
	ipTrace["elapsed"] = gconv.String(elapsed)
	result := worker.ParseResult{Items: []interface{}{ipTrace}}
	return result

}
