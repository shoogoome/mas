package models

type FileToken struct {
	// 令牌类型 0-上传 1-下载
	TokenType int `json:"token_type" yaml:"token_type"`
	// 创建时间
	CreateTime int64 `json:"create_time" yaml:"create_time"`
	// 文件hash
	Hash string `json:"hash" yaml:"hash"`
}

