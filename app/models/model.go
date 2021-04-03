package models

import "Blog/pkg/types"

type BaseModel struct {
	ID uint64
}

//获取ID的字符串
func (a BaseModel) GetStringID() string {
	return types.Uint64ToString(a.ID)
}
