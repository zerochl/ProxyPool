package storage

import (
	"ProxyPool/util"
)

// Config 全局配置文件
var ConfigFile = util.NewConfig()

// Storage struct is used for storeing persistent data of alerts
type FileStorage struct {
	fileName string
}

// NewStorage creates and returns new Storage instance
func NewFileStorage() *FileStorage {
	return &FileStorage{fileName:ConfigFile.LocalFile.FileName}
}

// Create insert new item
//func (s *FileStorage) Create(item interface{}) error {
//	//ses := s.GetDBSession()
//	//defer ses.Close()
//	//err := ses.DB(s.database).C(s.table).Insert(item)
//	//if err != nil {
//	//	return err
//	//}
//	//return nil
//
//}