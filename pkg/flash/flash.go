package flash

import (
	"Blog/pkg/session"
	"encoding/gob"
)

// Flashes  Flash 数组类型用于在会话中存储Map
type Flashes map[string]interface{}

//存入会话数据里的key
var flashKey = "_flashes"

func init()  {
	//在gorilla/session中存储map和Struct需要提前注册gob方便后续gob的序列化编码、解码
	gob.Register(Flashes{})
}

func Info(message string) {
	addFlash("info",message)
}

func Warning(message string)  {
	addFlash("warning",message)
}

func Success(message string)  {
	addFlash("success",message)
}

//添加danger消息提示
func Danger(message string)  {
	addFlash("danger",message)
}

//获取所有消息
func All() Flashes {
	val := session.Get(flashKey)

	//读取是必须做类型检测的
	flashMessages,ok := val.(Flashes)
	if !ok {
		return nil
	}

	//读取及销毁，直接删除
	session.Forget(flashKey)
	return flashMessages
}


//私有的方法  新增一条提示
func addFlash(key string , message string) {
	Flashes := Flashes{}
	Flashes[key] = message
	session.Put(flashKey,Flashes)
	session.Save()
}


