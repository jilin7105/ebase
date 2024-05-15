package EBHttpRequest

import (
	"context"
	"github.com/jilin7105/ebase/util/LinkTracking"
	"github.com/levigross/grequests"
)

func Get(ctx context.Context, url string, options *grequests.RequestOptions) (*grequests.Response, error) {
	if options == nil {
		options = &grequests.RequestOptions{}
	}
	SetRequestId(ctx, options)
	return grequests.Get(url, options)
}

func Post(ctx context.Context, url string, options *grequests.RequestOptions) (*grequests.Response, error) {
	if options == nil {
		options = &grequests.RequestOptions{}
	}
	SetRequestId(ctx, options)
	return grequests.Post(url, options)
}

func SetRequestId(ctx context.Context, options *grequests.RequestOptions) {
	if LinkTracking.GetIsOpen() {
		//写入链路追踪id
		//获取 id
		request_id := LinkTracking.GetEbaseRequestID(ctx)
		if options.Headers == nil {
			options.Headers = map[string]string{
				"EbaseRequestID": request_id,
			}
		} else {
			options.Headers["EbaseRequestID"] = request_id
		}
	}
	return

}
