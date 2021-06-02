package persist

import (
	"crawler/engine"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
)

func ItemSaver() chan interface{} {

	out := make(chan interface{})
	go func() {

		for {
			item := <-out
			_ = item
			_, err := engine.DB.Table(engine.Table).Replace(item)
			if err != nil {
				glog.Error(err)
				continue
			}
			glog.Infof("itemsaver saving rID: #%s",gconv.MapStrStr(item)["id"])
		}
	}()
	return out

}
