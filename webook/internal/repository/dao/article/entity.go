package article

// Article 这是制作库的
// 准备在 articles 表中准备十万/一百万条数据，author_id 各不相同（或者部分相同）
// 准备 author_id = 123 的， 插入两百条数据
// 执行 SELECT * FROM articles WHERE author_id = 123 ORDER BY ctime DESC
// 比较两种索引的性能
type Article struct {
	//model
	Id int64 `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	// 标题的长度
	// 正常都不会超过这个长度
	Title   string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	Content string `gorm:"type=BLOB" bson:"content,omitempty"`
	// 如何设计索引
	// 在帖子这里，什么样的查询场景？
	// 对于创作者来说，是不是看草稿箱，看到所有自己的文章？
	// SELECT * FROM articles WHERE author_id = 123 ORDER BY `ctime` DESC;
	// 产品经理告诉你，要按照创建时间的倒叙排序
	// 单独查询某一篇 SELECT * FROM articles WHERE id = 1
	// 在查询接口，我们深入讨论这个问题
	// 最佳选择，就是要在 author_id 和 ctime 上创建联合索引
	// - 在 author_id 上创建索引
	// 作者
	AuthorId int64 `gorm:"index" bson:"author_id,omitempty"`
	//AuthorId int64 `gorm:"index=aid_ctime"`
	//Ctime    int64 `gorm:"index=aid_ctime"`

	// 有些人考虑到，经常用状态来查询
	// WHERE status = xxx AND
	// 在 status 上和别的列混在一起，创建一个联合索引
	// 要看别的列究竟是什么列。
	Status uint8 `bson:"status,omitempty"`
	Ctime  int64 `bson:"ctime,omitempty"`
	Utime  int64 `bson:"utime,omitempty"`
}

// PublishedArticle 这个代表的是线上表
// PublishedArticle 衍生类型，偷个懒
type PublishedArticle Article

// PublishedArticleV1 s3 演示专属

type PublishedArticleV1 struct {
	Id       int64  `gorm:"primaryKey,autoIncrement" bson:"id,omitempty"`
	Title    string `gorm:"type=varchar(4096)" bson:"title,omitempty"`
	AuthorId int64  `gorm:"index" bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}

//func (u *Article) BeforeCreate(tx *gorm.DB) (err error) {
//	startTime := time.Now()
//	tx.Set("start_time", startTime)
//	slog.Default().Info("这是 BeforeCreate 钩子函数")
//	return nil
//}

//func (u *Article) AfterCreate(tx *gorm.DB) (err error) {
//	// 我要计算执行时间，我怎么拿到 before 里面的 startTime?
//	val, _ := tx.Get("start_time")
//	startTime, ok := val.(time.Time)
//	if !ok {
//		return nil
//	}
//	// 执行时间就出来了
//	duration := time.Since(startTime)
//	slog.Default().Info("这是 AfterCreate 钩子函数")
//	return nil
//}

//type model struct {
//}
//
//func (u model) BeforeSave(tx *gorm.DB) (err error) {
//	startTime := time.Now()
//	tx.Set("start_time", startTime)
//	slog.Default().Info("这是 BeforeCreate 钩子函数")
//	return nil
//}

//func (u model) AfterSave(tx *gorm.DB) (err error) {
//	// 我要计算执行时间，我怎么拿到 before 里面的 startTime?
//	val, _ := tx.Get("start_time")
//	startTime, ok := val.(time.Time)
//	if !ok {
//		return nil
//	}
//	// 执行时间就出来了
//	duration := time.Since(startTime)
//	slog.Default().Info("这是 AfterCreate 钩子函数")
//	return nil
//}
