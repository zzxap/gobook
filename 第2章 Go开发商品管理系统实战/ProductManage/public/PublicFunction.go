package public

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"

	"github.com/otiai10/copy"

	//"strconv"
	//"net/http/cookiejar"
	"io/ioutil"
	//"log"
	//"path/filepath"
	//"path"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	//"github.com/kardianos/osext"
	"archive/zip"
	"bytes"
	"encoding/binary"
	//"github.com/tomasen/realip"
	//"github.com/satori/go.uuid"
	//"github.com/op/go-logging"
)

var IsShowLog = false

func GetRandom() string {

	return GetUUIDS()

}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
func GetMd5(str string) string {

	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	//Log("sign=" + md5str1)
	return strings.ToUpper(md5str1)
}
func Unzip(src_zip string) string {
	// 解析解压包名
	dest := strings.Split(src_zip, ".")[0]
	// 打开压缩包
	unzip_file, err := zip.OpenReader(src_zip)
	if err != nil {
		return "压缩包损坏"
	}
	// 创建解压目录
	os.MkdirAll(dest, 0755)
	// 循环解压zip文件
	for _, f := range unzip_file.File {
		rc, err := f.Open()
		if err != nil {
			return "压缩包中文件损坏"
		}
		path := filepath.Join(dest, f.Name)
		// 判断解压出的是文件还是目录
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			// 创建解压文件
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "创建本地文件失败"
			}
			// 写入本地
			_, err = io.Copy(f, rc)
			if err != nil {
				if err != io.EOF {
					return "写入本地失败"
				}
			}
			f.Close()
		}
	}
	unzip_file.Close()
	return "OK"
}

func UnzipToest(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			Log(err)

		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				Log(err)

			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					Log(err)

				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCurDir() string {
	dir, _ := GetCurrentPath()
	return dir
}
func GetCurrentPath() (dir string, err error) {
	//path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	path, err := exec.LookPath(os.Args[0])

	if err != nil {
		Log("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		Log("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}
	dir = filepath.Dir(absPath)
	return dir, nil
}
func GetCurRunPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}
func ExistsPath(fullpath string) bool {
	//dir, _ := GetCurrentPath() //os.Getwd() //当前的目录
	//fullpath := dir + "/" + path
	_, err := os.Stat(fullpath)

	//Log("fullpath==" + fullpath)
	return err == nil || os.IsExist(err)
}

func CreatePath(fullpath string) {
	//dir, _ := GetCurrentPath() //os.Getwd() //当前的目录
	//fullpath := dir + "/" + newPath
	//fullpath = strings.Replace(fullpath, "/", "\\", -1)
	//fullpath = strings.Replace(fullpath, " ", "", -1)

	//newPath = strings.Replace(newPath, " ", "", -1)
	_, errr := os.Stat(fullpath)
	if errr != nil && os.IsNotExist(errr) {
		//Log(ff, fullpath+"文件不存在 创建") //为什么打印nil 是这样的如果file不存在 返回f文件的指针是nil的 所以我们不能使用defer f.Close()会报错的
		/*
			var path string
			if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
				path = "\\"
			} else {
				path = "/"
			}
		*/

		if err := os.MkdirAll(fullpath, 0777); err != nil {
			if os.IsPermission(err) {
				Log("你不够权限创建文件")
			}

		} else {
			//Log("创建目录" + fullpath + "成功")
		}

		//err := os.Mkdir(fullpath, os.ModePerm) //在当前目录下生成md目录
		//if err != nil {
		//	Log(err)
		//}

	} else {
		//Log(ff, fullpath+"文件存在 ")
	}

}

func SetCookie(r *http.Request, name string, value string) {
	COOKIE_MAX_MAX_AGE := time.Hour * 24 / time.Second // 单位：秒。
	maxAge := int(COOKIE_MAX_MAX_AGE)

	uid_cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: false,
		MaxAge:   maxAge}
	r.AddCookie(uid_cookie)

}
func GetTotal(price string, num string) string {
	fPrice, err1 := strconv.ParseFloat(price, 64)
	fnum, err2 := strconv.ParseFloat(num, 64)
	if err1 == nil && err2 == nil {
		return fmt.Sprintf("%1.2f", fPrice*fnum)
	}
	return ""

}
func RemovePath(path string) bool {
	//Log("upload picture Task is running...")
	//curdir := GetCurDir()
	//fullPath := curdir + "/" + path + "/"
	if ExistsPath(path) {
		err := os.RemoveAll(path)
		if err != nil {

			Log("remove fail " + path)
			return false
		} else {
			//如果删除成功则输出 file remove OK!
			return true
		}
	} else {
		return false
	}

}
func RemoveFile(path string) bool {
	//Log("upload picture Task is running...")
	//curdir := GetCurDir()
	//fullPath := curdir + "/" + path + "/"
	if ExistsPath(path) {
		err := os.Remove(path) //删除文件test.txt
		if err != nil {

			Log("remove fail " + path)
			return false
		} else {
			//如果删除成功则输出 file remove OK!
			return true
		}
	} else {
		return false
	}

}
func SavePictureTask(res http.ResponseWriter, req *http.Request, path string, userid string, typeid string) string {
	//Log("upload picture Task is running...")
	curdir := GetCurDir()
	var fileNames string = "#"
	if req.Method == "GET" {

	} else {

		ff, errr := os.Open(curdir + "/" + path + "/")
		if errr != nil && os.IsNotExist(errr) {
			Log(ff, ""+path+"文件不存在,创建") //为什么打印nil 是这样的如果file不存在 返回f文件的指针是nil的 所以我们不能使用defer f.Close()会报错的
			CreatePath(curdir + "/" + path + "/")
		}

		var (
			status int
			err    error
		)
		defer func() {
			if nil != err {
				http.Error(res, err.Error(), status)
			}
		}()
		// parse request
		const _24K = (1 << 20) * 24
		if err = req.ParseMultipartForm(_24K); nil != err {
			status = http.StatusInternalServerError
			return ""
		}
		for _, fheaders := range req.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				var infile multipart.File
				if infile, err = hdr.Open(); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				filename := hdr.Filename

				if strings.Contains(strings.ToLower(filename), ".mp3") || strings.Contains(strings.ToLower(filename), ".mov") {
					//如果是音频文件，直接存到picture文件夹，不存temp文件夹
					path = "Picture/" + userid + "/" + typeid
					CreatePath(curdir + "/" + path + "/")
				}

				// open destination
				var outfile *os.File
				savePath := curdir + "/" + path + "/" + filename
				if outfile, err = os.Create(savePath); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				// 32K buffer copy
				//var written int64
				if _, err = io.Copy(outfile, infile); nil != err {
					status = http.StatusInternalServerError
					return ""
				}

				infile.Close()
				outfile.Close()
				//CreatePath(curdir + "/" + path + "/thumbnial")
				//ImageFile_resize(infile, curdir+"/"+path+"/thumbnial/"+hdr.Filename, 200, 200)
				fileNames += "," + hdr.Filename
				//outfile.Close()
				//res.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
			}
		}
	}
	fileNames = strings.Replace(fileNames, "#,", "", -1)
	fileNames = strings.Replace(fileNames, "#", "", -1)
	return fileNames
}
func SaveConfigTask(res http.ResponseWriter, req *http.Request, path string, filename string) string {
	//Log("upload picture Task is running...")
	curdir := GetCurDir()
	var fileNames string = "#"
	if req.Method == "GET" {

	} else {

		ff, errr := os.Open(curdir + "/" + path + "/")
		if errr != nil && os.IsNotExist(errr) {
			Log(ff, ""+path+"文件不存在,创建") //为什么打印nil 是这样的如果file不存在 返回f文件的指针是nil的 所以我们不能使用defer f.Close()会报错的
			CreatePath(curdir + "/" + path + "/")

		}

		var (
			status int
			err    error
		)
		defer func() {
			if nil != err {
				http.Error(res, err.Error(), status)
			}
		}()
		// parse request
		const _24K = (1 << 20) * 24
		if err = req.ParseMultipartForm(_24K); nil != err {
			status = http.StatusInternalServerError
			return ""
		}
		for _, fheaders := range req.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				var infile multipart.File
				if infile, err = hdr.Open(); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				//filename := hdr.Filename

				// open destination
				var outfile *os.File
				savePath := curdir + "/" + path + "/" + filename
				if outfile, err = os.Create(savePath); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				// 32K buffer copy
				//var written int64
				if _, err = io.Copy(outfile, infile); nil != err {
					status = http.StatusInternalServerError
					return ""
				}

				infile.Close()
				outfile.Close()
				//CreatePath(curdir + "/" + path + "/thumbnial")
				//ImageFile_resize(infile, curdir+"/"+path+"/thumbnial/"+hdr.Filename, 200, 200)
				fileNames += "," + hdr.Filename
				//outfile.Close()
				//res.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
			}
		}
	}
	fileNames = strings.Replace(fileNames, "#,", "", -1)
	fileNames = strings.Replace(fileNames, "#", "", -1)
	return fileNames
}
func SaveUploadPictureTask(res http.ResponseWriter, req *http.Request, path string) string {
	//Log("upload picture Task is running...")
	curdir := GetCurDir()
	var fileNames string = "#"
	if req.Method == "GET" {

	} else {

		defer func() {
			if err := recover(); err != nil {
				Log("SaveUploadPictureTask")
				Log(err)
			}
		}()

		ff, errr := os.Open(curdir + "/" + path + "/")
		if errr != nil && os.IsNotExist(errr) {
			Log(ff, ""+path+"文件不存在,创建") //为什么打印nil 是这样的如果file不存在 返回f文件的指针是nil的 所以我们不能使用defer f.Close()会报错的
			CreatePath(curdir + "/" + path + "/")

		}

		var (
			status int
			err    error
		)
		defer func() {
			if nil != err {
				http.Error(res, err.Error(), status)
			}
		}()
		// parse request
		const _24K = (1 << 20) * 24
		if err = req.ParseMultipartForm(_24K); nil != err {
			status = http.StatusInternalServerError
			return ""
		}
		for _, fheaders := range req.MultipartForm.File {
			for _, hdr := range fheaders {
				// open uploaded
				var infile multipart.File
				if infile, err = hdr.Open(); nil != err {
					status = http.StatusInternalServerError
					return ""
				}

				filename := hdr.Filename

				// open destination
				var outfile *os.File
				savePath := curdir + "/" + path + "/" + filename
				//如果文件存在就给一个随机文件名
				if ExistsPath(savePath) {
					filename = GetRandomFileName(hdr.Filename)
					savePath = curdir + "/" + path + "/" + filename
				}

				if outfile, err = os.Create(savePath); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				// 32K buffer copy
				//var written int64
				if _, err = io.Copy(outfile, infile); nil != err {
					status = http.StatusInternalServerError
					return ""
				}
				infile.Close()
				outfile.Close()
				//CreatePath(curdir + "/" + path + "/thumbnial")
				//ImageFile_resize(infile, curdir+"/"+path+"/thumbnial/"+hdr.Filename, 200, 200)
				fileNames += "," + filename
				//outfile.Close()
				//res.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
			}
		}
	}
	fileNames = strings.Replace(fileNames, "#,", "", -1)
	fileNames = strings.Replace(fileNames, "#", "", -1)
	return fileNames
}

func GetRandomFileName(name string) string {
	//name := hdr.Filename
	arr := strings.Split(name, ".")
	extent := arr[len(arr)-1]

	return GetRandom() + "." + extent
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		Log(err.Error())
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		Log(err.Error())
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		Log(err.Error())
		return err
	}
	if ExistsPath(dst) {
		//Log("copy success" + dst)
	}

	return out.Close()
}

//拷贝文件  要拷贝的文件路径 拷贝到哪里 "github.com/otiai10/copy"
func CopyFiles(source, dest string) bool {
	if source == "" || dest == "" {
		Log("source or dest is null")
		return false
	}
	err := copy.Copy(source, dest)
	if err == nil {
		return true
	} else {
		return false
	}
}
func NetWorkStatus() bool {
	cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
	fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	err := cmd.Run()
	fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println("Net Status , OK")
	}
	return true
}

func GetMapByJsonStr(jsonstr string) map[string]interface{} {
	if !strings.Contains(jsonstr, "{") {
		Log("bad json=" + jsonstr)
		return nil
	}
	jsonstr = strings.Replace(jsonstr, "\x00", "", -1)
	if len(jsonstr) > 4 {
		var d map[string]interface{}
		err := json.Unmarshal([]byte(jsonstr), &d)
		if err != nil {
			log.Printf("error decoding sakura response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			//log.Printf("sakura response: %q", resBody)
			Log("bad json" + jsonstr)
			Log(err)
			//panic("bad json")
			return nil
		}
		return d
	}
	return nil
}
func GetMessageMapByJsonKey(jsonstr string, keystr string) map[string]interface{} {
	//var jsonstr:='{\"data\": { \"mes\": [ {\"fromuserid\": \"25\", \"touserid\": \"56\",\"message\": \"hhhhhaaaaaa\",\"time\": \"2017-12-12 12:11:11\"}]}}';
	index := strings.IndexRune(jsonstr, '{')
	jsonstr = jsonstr[index : len(jsonstr)-index]
	if len(jsonstr) > 4 && strings.Index(jsonstr, "{") > -1 && strings.Index(jsonstr, "}") > -1 {
		mapp := GetMapByJsonStr(jsonstr)
		//Log(mapp)
		mappp := mapp[keystr]
		//Log(mappp)
		//kll := mapp.(map[string]interface{})[keystr]
		//Log(kll)
		mymap := mappp.(map[string]interface{})
		//Log(mymap["Fromuserid"])

		return mymap
	}
	return nil
}
func GetMessageMapByJson(jsonstr string) map[string]interface{} {
	//var jsonstr:='{\"data\": { \"mes\": [ {\"fromuserid\": \"25\", \"touserid\": \"56\",\"message\": \"hhhhhaaaaaa\",\"time\": \"2017-12-12 12:11:11\"}]}}';
	index := strings.IndexRune(jsonstr, '{')
	jsonstr = jsonstr[index : len(jsonstr)-index]
	if len(jsonstr) > 4 && strings.Index(jsonstr, "{") > -1 && strings.Index(jsonstr, "}") > -1 {
		mapp := GetMapByJsonStr(jsonstr)
		//Log(mapp)
		mappp := mapp["data"]
		//Log(mappp)
		kll := mappp.(map[string]interface{})["mes"]
		//Log(kll)
		mymap := kll.(map[string]interface{})
		//Log(mymap["fromuserid"])

		return mymap
	}
	return nil
}

func GetJsonStrByMap(MapList map[int]map[string]string) string {
	var str string = "##"

	sorted_keys := make([]int, 0)
	for k, _ := range MapList {
		sorted_keys = append(sorted_keys, k)
	}

	// sort 'string' key in increasing order
	sort.Ints(sorted_keys)

	for _, k := range sorted_keys {
		//fmt.Printf("k=%v, v=%v\n", k, MapList[k])

		jsonStr, err := json.Marshal(MapList[k])
		if err != nil {
			Log(err)
		}
		//Log("map to json", string(str))
		str += "," + string(jsonStr)

	}

	str = strings.Replace(str, "##,", "", -1)
	str = strings.Replace(str, "##", "", -1)
	return str
}
func ConverToStr(v interface{}) string {
	if v == nil {
		return ""
	}
	var str string = ""
	if reflect.TypeOf(v).Kind() == reflect.String {
		str = v.(string)
	} else if reflect.TypeOf(v).Kind() == reflect.Int {
		str = string(v.(int))
	} else if reflect.TypeOf(v).Kind() == reflect.Int8 {
		str = string(v.(int8))
	} else if reflect.TypeOf(v).Kind() == reflect.Int16 {
		str = string(v.(int16))
	} else if reflect.TypeOf(v).Kind() == reflect.Int32 {
		str = string(v.(int32))
	} else if reflect.TypeOf(v).Kind() == reflect.Int64 {
		str = string(v.(int64))
	} else if reflect.TypeOf(v).Kind() == reflect.Float32 {
		str = fmt.Sprintf("%f", v)
	} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
		str = fmt.Sprintf("%f", v)
	} else {
		str = v.(string)
	}

	return strings.Replace(str, ".000000", "", -1)
}
func GetCurDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func GetCurDay() string {

	return time.Now().Format("2006-01-02")
}
func GetNameSinceNow(after int) string {
	day := time.Now().AddDate(0, 0, after).Format("2006-01-02")
	day = strings.Replace(day, "-", "", -1)
	return day
}
func GetDaySinceNow(after int) string {
	return time.Now().AddDate(0, 0, after).Format("2006-01-02")
}
func ReplaceStr(str string) string {
	//过滤一下，防止sql注入
	str = strings.Replace(str, "'", "", -1)
	//str = strings.Replace(str, "-", "\\-", -1)
	str = strings.Replace(str, "exec", "exe.c", -1)
	return str //.Replace(str, ",", "", -1).Replace(str, "-", "\-", -1) //-1表示替换所有
}

var logfile *os.File
var oldFileName string

func Log(a ...interface{}) (n int, err error) {
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println(a...)

	return 1, nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP22() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipstr := ipnet.IP.String()
				index := strings.Index(ipstr, "127.0")
				if index > -1 {
					continue
				}

				index = strings.Index(ipstr, "192.168.")
				if index > -1 {
					return ipstr
					break
				}

				index = strings.Index(ipstr, "169.254.")
				if index > -1 {
					continue
				}

				return ipstr
			}
		}
	}
	return ""
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)

		return localAddr.IP.String()
	} else {
		return GetLocalIPP()
	}

}
func GetLocalIPP() string {
	//GetIpList()
	var ipstr string = ""
	//windows 获取IP
	host, _ := os.Hostname()
	addrss, err := net.LookupIP(host)
	if err != nil {
		Log("error", err.Error())
		//return ""
	}
	var ipArray []string

	for _, addr := range addrss {
		if ipv4 := addr.To4(); ipv4 != nil {
			Log("ippppp=: ", ipv4)
			ipstr = ipv4.String()

			if !strings.HasPrefix(ipstr, "127.0") && !strings.HasPrefix(ipstr, "169.254") && !strings.HasPrefix(ipstr, "172.16") {
				ipArray = append(ipArray, ipstr)
			}
		}
	}

	//提取公网IP
	//var pubIpArray []string
	for i := 0; i < len(ipArray); i++ {
		//Log("pubip===" + ipArray[i])
		if !strings.HasPrefix(ipArray[i], "10.") && !strings.HasPrefix(ipArray[i], "192.168") && !strings.HasPrefix(ipArray[i], "172.") {
			return ipArray[i]
			//pubIpArray = append(pubIpArray, ipstr)
		}
	}

	//如果没有公网IP 就返回一个本地IP

	if len(ipArray) > 0 {

		return ipArray[0]
	}

	//linux 获取IP
	if ipstr == "" {

		ifaces, errr := net.Interfaces()
		// handle err
		if errr != nil {
			Log("error", errr.Error())
			return ""
		}

		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			// handle err
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				// process IP address
				//Log("ip=", ip)
				ipstr = fmt.Sprintf("%s", ip)
				Log("ipstr=", ipstr)

				index := strings.Index(ipstr, "127.0")
				if index > -1 {
					continue
				}
				index = strings.Index(ipstr, "192.168.")
				if index > -1 {
					return ipstr
					break
				}

				index = strings.Index(ipstr, "169.254.")
				if index > -1 {
					continue
				}
				if len(ipstr) > 6 {
					array := strings.Split(ipstr, ".")
					if len(array) == 4 {
						return ipstr
					}

				}

			}
		}

	}

	return ""
}

func HttpPost(url string, paras string) string {
	//Log("url=" + url + " paras=" + paras)
	client := &http.Client{}

	req, err := http.NewRequest("POST",
		url,
		strings.NewReader(paras))

	if err != nil {
		// handle error
		return ""
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return ""
	}

	//Log(string(body))
	return string(body)
}

func HttpGet(url string) string {
	//Log("get =" + url)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		Log(err.Error())
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		Log(err.Error())
		return ""
	}

	//Log("response =" + string(body))
	return string(body)
}
func HttpDownloadFile(url string, toPath string) {
	//Log("get =" + url)
	res, err := http.Get(url)
	if err != nil {
		Log(err)
		return
	}
	f, err := os.Create(toPath)
	defer f.Close()
	if err != nil {
		Log(err)
		return
	}
	io.Copy(f, res.Body)

	//Log("size =" + size)
}

//整形转换成字节
func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

func RealIPHand(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if rip := RealIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xForwardedFor2 = http.CanonicalHeaderKey("x-forwarded-for")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")
var xRealIP2 = http.CanonicalHeaderKey("x-real-ip")
var xRealIP3 = http.CanonicalHeaderKey("x-real-client-ip")

var ProxyClientIP = http.CanonicalHeaderKey("Proxy-Client-IP")
var WLProxyClientIP = http.CanonicalHeaderKey("WL-Proxy-Client-IP")
var HTTPXFORWARDEDFOR = http.CanonicalHeaderKey("HTTP_X_FORWARDED_FOR")

func RealIP(r *http.Request) string {

	PrintHead(r)

	var ip string

	//clientIP := realip.FromRequest(r)
	//log.Println("GET / from", clientIP)

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		//Log(xff)
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xff := r.Header.Get(xForwardedFor2); xff != "" {
		//Log(xff)
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(xRealIP2); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(xRealIP3); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(ProxyClientIP); xrip != "" {
		ip = xrip
	} else if xrip := r.Header.Get(WLProxyClientIP); xrip != "" {
		ip = xrip
	} else {
		ip = r.RemoteAddr
	}

	return ip

	//return realip.FromRequest(r)
}

func PrintHead(r *http.Request) {

	realip := r.Header.Get(xForwardedFor)

	if len(realip) == 0 {
		realip = r.Header.Get("http_client_ip")
	}
	if len(realip) == 0 {
		//Log(xRealIP)
		realip = r.Header.Get(xRealIP)
	}

	if len(realip) == 0 {
		//Log(ProxyClientIP)
		realip = r.Header.Get(ProxyClientIP)
	}
	if len(realip) == 0 {
		//Log(WLProxyClientIP)
		realip = r.Header.Get(WLProxyClientIP)
	}
	if len(realip) == 0 {
		//Log(HTTPXFORWARDEDFOR)
		realip = r.Header.Get(HTTPXFORWARDEDFOR)
	}
	if len(realip) == 0 {

		realip = r.RemoteAddr
	}
	//Log("ip=" + r.RemoteAddr)
	//Log("realip=" + realip)

}
