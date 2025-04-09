package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type AticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

func NewGORMArticleDAO(db *gorm.DB) AticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// Article 这是制作库的
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 长度 1024
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	// 如何设计索引
	// 在帖子这里，什么样的查询场景？
	// 对于创作者来说，是不是看草稿箱，看到所有自己的文章？
	// SELECT * FROM articles WHERE author_id = 123 ORDER BY `ctime` DESC;
	// 产品经理告诉你，要按照创建时间的倒叙排序
	// 单独查询某一篇 SELECT * FROM articles WHERE id = 1
	// 在查询接口，我们深入讨论这个问题
	// 最佳选择，就是要在 author_id 和 ctime 上创建联合索引
	// - 在 author_id 上创建索引
	AuthorId int64 `gorm:"index"`
	//AuthorId int64 `gorm:"index=aid_ctime"`
	//Ctime    int64 `gorm:"index=aid_ctime"`
	Ctime int64
	Utime int64
}
