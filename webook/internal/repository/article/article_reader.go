package article

import (
	"context"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
)

type ArticleReaderRepository interface {
	// Save 有就更新，没有就新建
	Save(ctx context.Context, art domain.Article) (int64, error)
	//Update(ctx context.Context, art domain.Article) error
}
