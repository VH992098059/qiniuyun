package voice

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"golang/api/voice/v1"
)

func (c *ControllerV1) Voice(ctx context.Context, req *v1.VoiceReq) (res *v1.VoiceRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
