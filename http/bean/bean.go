package bean

type PageRo struct {
	Page int `query:"page" validate:"required,gte=1" default:"1"`
	Size int `query:"size" validate:"required,lte=100,gte=1" default:"10"`
}

type TimeRo struct {
	// 开始时间（一般指CreatedAt） 时间戳 秒
	Begin int64 `query:"begin" json:"begin"`
	// 结束时间 同上
	End int64 `query:"end" json:"end"`
}
