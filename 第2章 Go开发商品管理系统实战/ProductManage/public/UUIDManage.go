package public

import (
	//"fmt"
	"strconv"
	"strings"

	//"time"
	"github.com/bwmarrin/snowflake"
	"github.com/zheng-ji/goSnowFlake"
)

func GetUUIDS() string {
	workerid := "1"
	//workerid := GetWorkerId()
	wid, errr := strconv.ParseInt(workerid, 10, 64)
	uuids := "##"
	if errr != nil {
		wid = 1
	}

	for i := 0; i < 1; i++ {
		node, err := snowflake.NewNode(wid)
		if err == nil {
			str := node.Generate().String()
			//fmt.Println(str)
			str = str[2:len(str)]
			uuids += "," + str
		}

	}
	uuids = strings.Replace(uuids, "##,", "", -1)
	uuids = strings.Replace(uuids, "##", "", -1)

	return uuids
}

func GetUUIDS22(num string) string {
	workerid := "1"
	//workerid := GetWorkerId()
	wid, errr := strconv.ParseInt(workerid, 10, 64)
	uuids := "##"
	if errr != nil {
		wid = 1
	}
	//t := time.Now()
	//fmt.Println(t.Format("20060102150405"))
	//fmt.Println(time.Now().Format(time.RFC850))
	//TimeStamp, errrt := strconv.ParseInt(t.Format("20060102150405"), 10, 64)

	//if errrt != nil {
	//	TimeStamp = -1
	//}

	iw, err := goSnowFlake.NewIdWorker(wid)
	if err != nil {
		return ""
	}

	inum, _ := strconv.Atoi(num)
	for i := 0; i < inum; i++ {
		if id, err := iw.NextId(); err != nil {
			//fmt.Println(err)
		} else {
			//fmt.Println(id)
			sid := strconv.FormatInt(id, 10)
			uuids += "," + sid
		}
	}

	uuids = strings.Replace(uuids, "##,", "", -1)
	uuids = strings.Replace(uuids, "##", "", -1)

	return uuids
}
