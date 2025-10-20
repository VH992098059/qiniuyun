package voice

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"golang/api/voice/v1"
)

func (c *ControllerV1) Asr(ctx context.Context, req *v1.AsrReq) (res *v1.AsrRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
