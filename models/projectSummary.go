package models

import (
	"GoGin/config/logging"
	"fmt"
	"github.com/jinzhu/gorm"
	"sort"
	"strconv"
)

type Member struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	Product string `json:"product"`
}

type Summary struct {
	Total    int     `json:"total"`
	Coverage float64 `json:"coverage"`
}

type Product struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Summary
}

type Person struct {
	Owner string `json:"owner"`
	Summary
}

func OwnerApiNumber() ([]*Person, error) {
	var trueResult []*Person
	rows, err := db.Table("yapi_token").
		Select("owner, sum(total) as total,  CONVERT(SUM(coverage)/(count(coverage) * 100), DECIMAL ( 15, 2 )) as coverage").
		Where("prodect IN (?)", []string{"live", "enterprise", "cloud", "class"}).Group("owner").Rows()
	if err != nil {
		return nil, err
	} else {
		for rows.Next() {
			p := new(Person)
			err := rows.Scan(&p.Owner, &p.Total, &p.Coverage)
			if err != nil {
				logging.Error("获取业务线个人数据接口，操作数据库出错了", err)
				continue
			}
			p.Coverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", p.Coverage), 64)
			trueResult = append(trueResult, p)
		}
	}
	sort.Slice(trueResult, func(i, j int) bool { return trueResult[i].Total > trueResult[j].Total })

	return trueResult, nil

}

func ExistGroup(prodect string) bool {
	var groupResult Token
	err := db.Select("id").Where("prodect = ?", prodect).Find(&groupResult).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	if groupResult.ID > 0 {
		return true
	} else {
		return false
	}
}

func GetProductSummary() (rs []*Product, err error) {
	rows, err := db.Table("yapi_token").
		Select("id, prodect, sum(total) as total,  CONVERT(SUM(coverage)/(count(coverage) * 100), DECIMAL ( 15, 2 )) as coverage").
		Where("prodect IN (?)", []string{"live", "enterprise", "cloud", "class"}).Group("prodect").Rows()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := new(Product)
		err := rows.Scan(&p.Id, &p.Name, &p.Total, &p.Coverage)
		if err != nil {
			logging.Error("获取业务线数据接口，操作数据库出错了", err)
			continue
		}
		p.Coverage, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", p.Coverage), 64)
		rs = append(rs, p)
	}
	return rs, nil
}

func GetGroupMembers(prodect string) (rs []*Member) {
	err := db.Table("yapi_test_member").Where("product = ? and status = 0", prodect).Find(&rs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return rs
	}
	return rs
}
