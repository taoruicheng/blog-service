package service

import (
	"context"

	"github.com/taoruicheng/blog-service/global"
	"github.com/taoruicheng/blog-service/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	//TODO
	svc.dao = dao.New(global.DBEngine)
	return svc
}
