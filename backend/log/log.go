package log

import (
	"context"
	"encoding/json"

	"github.com/vvakame/sdlog/aelog"
)

func Info(ctx context.Context, v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		aelog.Errorf(ctx, "failed json.Marshal(). err=%+v", err)
		return
	}
	aelog.Infof(ctx, "%+v", string(j))
}
