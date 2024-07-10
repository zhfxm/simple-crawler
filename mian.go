package main

import (
	"fmt"
	"time"

	"github.com/zhfxm/simple-crawler/collect"
	"github.com/zhfxm/simple-crawler/log"
	"github.com/zhfxm/simple-crawler/parse/doubangroup"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("logger init end")

	cookie := "viewed=\"1007305\"; bid=QGFx5rUryE0; _pk_id.100001.8cb4=c8407de68c359e04.1720513239.; __utmc=30149280; dbcl2=\"270053320:gi+WD2F5+Ic\"; ck=LlMn; push_noty_num=0; push_doumail_num=0; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1720576818%2C%22https%3A%2F%2Faccounts.douban.com%2F%22%5D; _pk_ses.100001.8cb4=1; __utma=30149280.471027257.1718694576.1720513239.1720576818.3; __utmz=30149280.1720576818.3.2.utmcsr=accounts.douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; __utmt=1; __utmv=30149280.27005; __utmb=30149280.24.5.1720576899469"

	var worklist []*collect.Request
	for i := 0; i <= 0; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d&type=new", i)
		worklist = append(worklist, &collect.Request{
			Url: str,
			ParseFunc: doubangroup.ParseURL,
			Log: logger,
			Cookie: cookie,
		})
	}
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 30 * time.Second,
	}

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				logger.Error("read context failed", zap.Error(err))
				continue
			}
			res := item.ParseFunc(body, item)
			for _, it := range res.Items {
				logger.Info("result", zap.String("get url", it.(string)))
			}
			if len(res.Requests) > 0 {
				worklist = append(worklist, res.Requests...)	
			}
		}
	}
}
