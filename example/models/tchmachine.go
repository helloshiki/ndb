package models

import "encoding/json"

type TchMachine struct {
	ID  int
	GID int
	UID int
}

var (
	indexs = map[string][]string{
		"guid": {"GID", "UID"},
		"gid":  {"GID"},
		"uid":  {"UID"},
	}
)

//索引列的值暂不支持修改，如需修改，可以先删除记录，再插入新记录
func (t *TchMachine) Index() map[string][]string {
	return indexs
}

func (t *TchMachine) GetUID() int {
	return t.ID
}

func (u *TchMachine) Encode() []byte {
	bs, _ := json.Marshal(u)
	return bs
}

func (u *TchMachine) Category() string {
	return "TchMachine"
}
