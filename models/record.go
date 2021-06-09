package models

import (
	"GolangApiTest/config/logging"
	"GolangApiTest/config/msgCode"
	"GolangApiTest/config/util"
	"github.com/jinzhu/gorm"
	"sync"
)

var mutex sync.Mutex
var syncApiMap = make(map[int]string)

type TestApi struct {
	Model
	Project string `gorm:"column:project"`
	Sid     int    `gorm:"column:sid"`
	Yid     int    `gorm:"column:yid"`
	Method  string `gorm:"column:method"`
	Title   string `gorm:"column:title"`
	Path    string `gorm:"column:path"`
}

func SyncSvcApi(id int) int {
	var token Token
	var bp string
	db.Where("id = ?", id).First(&token)
	// 获取服务配置的基础路径
	svcInfo, err := util.GetSvcBasicInfo(token.Token)
	if err != nil {
		return msgCode.TaskError
	}
	if svcInfo != "" {
		bp = svcInfo
	}
	// 获取服务所有api
	rs, err := util.GetSvcApis(token.Sid, token.Token)
	if err != nil {
		return msgCode.TaskError
	}
	token.ApiTotal = len(rs)
	err = db.Model(&token).Where("id = ?", id).Updates(&token).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return msgCode.TaskError
	}
	// 加锁
	mutex.Lock()
	defer mutex.Unlock()
	s := syncApiMap[token.Sid]
	if s != "" {
		return msgCode.TaskRunning
	}
	syncApiMap[token.Sid] = token.SvcName
	// 异步更新
	go func(token Token, basePath string, apis []util.ApiContent) {
		for _, api := range apis {
			var testApi = new(TestApi)
			db.Where("yid = ?", api.ApiId).First(&testApi)
			if testApi.ID != 0 {
				// todo: 是否需要更新api
				logging.Info("已存在项目接口", api.ApiId)
				continue
			}
			testApi.Yid = api.ApiId
			testApi.Sid = api.ProjectId
			testApi.Method = api.Method
			testApi.Title = api.Title
			testApi.Project = token.ProjectName
			// todo: 需要处理path parameter 不一致问题(:act_id, {act_id})
			if basePath != "" {
				testApi.Path = basePath + api.Path
			} else {
				testApi.Path = api.Path
			}
			if err := db.Create(&testApi).Error; err != nil {
				logging.Error("add test api error %v", err)
			}
		}
		delete(syncApiMap, token.Sid)
	}(token, bp, rs)
	return msgCode.TaskSuccess
}

func reportTestApi() {

}
