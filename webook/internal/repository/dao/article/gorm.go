package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Article, error) {
	var res []Article
	err := dao.db.WithContext(ctx).
		Where("utime<?", start.UnixMilli()).
		Order("utime DESC").Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (dao *GORMArticleDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	var arts []Article
	// SELECT * FROM XXX WHERE XX order by aaa
	// 在设计 order by 语句的时候，要注意让 order by 中的数据命中索引
	// SQL 优化的案例：早期的时候，
	// 我们的 order by 没有命中索引的，内存排序非常慢
	// 你的工作就是优化了这个查询，加进去了索引
	// author_id => author_id, utime 的联合索引
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("author_id = ?", author).
		Offset(offset).
		Limit(limit).
		// 升序排序。 utime ASC
		// 混合排序
		// ctime ASC, utime desc
		Order("utime DESC").
		//Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		//	{Column: clause.Column{Name: "utime"}, Desc: true},
		//	{Column: clause.Column{Name: "ctime"}, Desc: false},
		//}}).
		Find(&arts).Error
	return arts, err
}

func (dao *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var pub PublishedArticle
	err := dao.db.WithContext(ctx).
		Where("id = ?", id).
		First(&pub).Error
	return pub, err
}

func (dao *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ?", id).
		First(&art).Error
	return art, err
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{
		db: db,
	}
}

func (dao *GORMArticleDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? AND author_id = ?", id, author).
			Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			// 要么 ID 是错的， 要么作者不对
			// 后者情况下，要小心，可能有人在搞你
			// 没必要再用 ID 搜索数据库来区分这两种情况
			// 用 prometheus 打点，只要频繁出现，你就告警，然后手工介入排查
			return ErrPossibleIncorrectAuthor
		}

		res = tx.Model(&PublishedArticle{}).
			Where("id=? AND author_id = ?", id, author).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrPossibleIncorrectAuthor
		}
		return nil
	})
}

func (dao *GORMArticleDAO) Sync(ctx context.Context,
	art Article) (int64, error) {
	// tx => Transaction, trx, txn
	// 在事务内部，这里采用了闭包形态
	// GORM 帮助我们管理了事务的生命周期
	// Begin，Rollback 和 Commit 都不需要我们操心
	tx := dao.db.WithContext(ctx).Begin()
	now := time.Now().UnixMilli()
	defer tx.Rollback()
	txDAO := NewGORMArticleDAO(tx)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = txDAO.Insert(ctx, art)
	} else {
		err = txDAO.UpdateById(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	publishArt := PublishedArticle(art)
	publishArt.Utime = now
	publishArt.Ctime = now
	err = tx.Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		}),
	}).Create(&publishArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, tx.Error
}

func (dao *GORMArticleDAO) SyncClosure(ctx context.Context,
	art Article) (int64, error) {
	var (
		id = art.Id
	)
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		now := time.Now().UnixMilli()
		txDAO := NewGORMArticleDAO(tx)
		if id == 0 {
			id, err = txDAO.Insert(ctx, art)
		} else {
			err = txDAO.UpdateById(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		publishArt := art
		publishArt.Utime = now
		publishArt.Ctime = now
		return tx.Clauses(clause.OnConflict{
			// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":   art.Title,
				"content": art.Content,
				"utime":   now,
			}),
		}).Create(&publishArt).Error
	})
	return id, err
}

// 事务传播机制是指如果当前有事务，就在事物内部执行 Insert
// 如果没有事务
// 1. 开启事务，执行 Insert
// 2. 直接执行
// 3. 报错

func (dao *GORMArticleDAO) Insert(ctx context.Context,
	art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	// 返回自增主键
	return art.Id, err
}

// UpdateById 只更新标题、内容和状态
func (dao *GORMArticleDAO) UpdateById(ctx context.Context,
	art Article) error {
	now := time.Now().UnixMilli()
	// 依赖 gorm 忽略零值的特性，会用主键进行更新
	// 可读性很差
	res := dao.db.Model(&Article{}).WithContext(ctx).
		Where("id=? AND author_id = ? ", art.Id, art.AuthorId).
		// 当你用这种每次都指定被更新列的写法
		// 可读性强，但是每一次更新更多列的时候，你都要修改
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		})
	err := res.Error
	// 你要不要检查真的更新了没
	if err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

// Upsert : INSERT or UPDATE
//func (dao *GORMArticleDAO) Upsert(ctx context.Context, art PublishedArticle) error {
//	now := time.Now().UnixMilli()
//	art.Ctime = now
//	art.Utime = now
//	// 这个是插入
//	// OnConflict 的意思是数据冲突了
//	err := dao.db.Clauses(clause.OnConflict{
//		// SQL 2003 标准
//		// INSERT XXX ON CONFLICT(BBB) DO NOTHING
//		// INSERT XXX ON CONFLICT(BBB) DO UPDATES CCC WHERE DDD
//
//		// 哪些列冲突
//		//Columns: []clause.Column{clause.Column{Name: "id"}},
//		// 意思是数据冲突，啥也不干
//		//DoNothing:
//		// 数据冲突了，并且符合 WHERE 条件的就会执行 DO UPDATES
//		//Where:
//
//		// MySQL 只需要关心这里
//		DoUpdates: clause.Assignments(map[string]interface{}{
//			"titile":  art.Title,
//			"content": art.Content,
//			"status":  art.Status,
//			"utime":   art.Utime,
//		}),
//	}).Create(&art).Error
//	// MySQL 最终的语句 INSERT xxx ON DUPLICATE KEY UPDATE xxx
//	// 正常来说，一条 SQL 语句，都不需要开启事务
//	// auto commit: 意思是自动提交
//
//	return err
//}
