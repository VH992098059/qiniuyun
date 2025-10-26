package voice

import (
	"context"
	"golang/internal/logic"

	"golang/api/voice/v1"
)

func (c *ControllerV1) Text(ctx context.Context, req *v1.TextReq) (res *v1.TextRes, err error) {
	logic.RobotAutoApplication(ctx, req.Text)
	return &v1.TextRes{ResText: "已完成"}, nil
}
