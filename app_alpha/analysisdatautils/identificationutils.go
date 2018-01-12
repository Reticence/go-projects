// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-16 09:42
// @Version : 0.1
// @Software: GoLand

package analysisdatautils

import (
	"fmt"
	"strconv"
)

const (
	identifiedY = "Y"
	identifiedN = "N"
	matrimonyY  = "Y"
	matrimonyN  = "N"
	genderMan   = "男"
	genderWoman = "女"
)

var (
	customerIdMap map[string]*IdentificationInfo
	idcardMap     map[string]*IdentificationInfo

	nextRelatedIdChan chan int
	nextPersonIdChan  chan int
)

type IdentificationInfo struct {
	RelatedId  int    // 体检报告关联id（meridian）
	Identified string // 身份识别（meridian）
	PersonId   int    // 人员编号（meridian）
	HospitalId int    // 医院编号
	CustomerId string // 客户编号
	Idcard     string // 身份证件号码
}

func (ii *IdentificationInfo) getHCKey() string {
	return strconv.Itoa(ii.HospitalId) + "_" + ii.CustomerId
}

func (ii *IdentificationInfo) GetInfo() {
	key := ii.getHCKey()
	ii1, ok1 := idcardMap[ii.Idcard]
	ii2, ok2 := customerIdMap[key]
	fmt.Println(ok1, ok2)
	if ok1 && ok2 {
		if ii1 != ii2 {
			if ii1.CustomerId == "" {
				ii1.HospitalId, ii1.CustomerId = ii2.HospitalId, ii2.CustomerId
				customerIdMap[key] = ii1
			}
		}
		*ii = *ii1
	} else if ok1 {
		if ii.CustomerId != "" {
			ii1.HospitalId, ii1.CustomerId = ii.HospitalId, ii.CustomerId
			customerIdMap[key] = ii1
		}
		*ii = *ii1
	} else if ok2 {
		if ii.Idcard != "" {
			ii2.Idcard = ii.Idcard
			idcardMap[ii.Idcard] = ii2
		}
		*ii = *ii2
	} else {
		ii.Identified = identifiedN
		ii.PersonId = <-nextPersonIdChan
		ii3 := &IdentificationInfo{}
		*ii3 = *ii
		if ii.CustomerId != "" {
			customerIdMap[key] = ii3
		}
		if ii.Idcard != "" {
			idcardMap[ii.Idcard] = ii3
		}
	}
	ii.RelatedId = <-nextRelatedIdChan
}

func Test() {
	var ii IdentificationInfo

	fmt.Println("-------------------------------")
	ii = IdentificationInfo{Idcard: "110"}
	ii.GetInfo()
	fmt.Println(ii)

	fmt.Println("-------------------------------")
	ii = IdentificationInfo{HospitalId: 1, CustomerId: "tj1111"}
	ii.GetInfo()
	fmt.Println(ii)

	fmt.Println("-------------------------------")
	ii = IdentificationInfo{Idcard: "110", HospitalId: 1, CustomerId: "tj1111"}
	ii.GetInfo()
	fmt.Println(ii)

	fmt.Println("-------------------------------")
	ii = IdentificationInfo{HospitalId: 1, CustomerId: "tj1111"}
	ii.GetInfo()
	fmt.Println(ii)
}
