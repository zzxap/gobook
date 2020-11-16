package HttpBusiness

import (
	"net/http"

	"os"
	"os/exec"
	"runtime"
	"strings"

	"ProductManage/dal"
	"ProductManage/grace"
	"ProductManage/public"
	"github.com/gorilla/mux"
	//"ProductManage/dal"
)

var localip = ""
var ConfigServerName = ""
var ConfigLocalServerName = ""
var oldHost = ""
var ch = make(chan int)

//启动一个http server
func StartHttpServer() {

	CopyFileToRunningPath()

	PthSep := string(os.PathSeparator)
	if PthSep == "\\" {
		go open("http://127.0.0.1:8090/html/index.html")
	}

	dal.InitDB()
	myHandler := mux.NewRouter()

	myHandler.HandleFunc("/api/addproduct", AddProductTask)
	myHandler.HandleFunc("/api/getlist", GetProductListTask)
	myHandler.HandleFunc("/api/delete", DeleteTask)
	s := http.StripPrefix("/html/", http.FileServer(http.Dir("./html/")))
	myHandler.PathPrefix("/html/").Handler(s)

	public.Log("start http server")
	errr := grace.ListenAndServe(":8090", Middleware(myHandler))
	if errr != nil {
		public.Log("ListenAndServe  error: %v"+public.GetCurDateTime(), errr)
		//panic("http server stop exit at" + public.GetCurDateTime())
	} else {
		public.Log("ListenAndServe success")
	}

}

//http请求的中间件 跨域处理
func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin,Authorization,Origin, X-Requested-With, Content-Type, Accept,common")

		h.ServeHTTP(w, r)

		if r.Method == "OPTIONS" {
			return
		}
	})
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

//添加产品
func AddProductTask(w http.ResponseWriter, r *http.Request) {
	public.Log("AddProductTask")

	Id := r.FormValue("id")
	UserId := r.FormValue("userid")
	CategoryId := r.FormValue("categoryid")
	ProductName := r.FormValue("productname")
	Price := r.FormValue("price")
	Remark := r.FormValue("remark")
	AddTime := public.GetCurDateTime()
	mode := r.FormValue("mode")
	m := &dal.ProductModel{
		Id,
		UserId,
		CategoryId,
		ProductName,
		Price,
		Remark,
		AddTime,
	}
	public.Log("mode=" + mode)
	if mode == "add" {
		if dal.ExistsProductName(ProductName) {
			writeJsonValue(w, "1", "exists", "")
		} else if dal.AddProduct(m) > 0 {
			writeJsonValue(w, "0", "success", "")
		} else {
			writeJsonValue(w, "1", "fail", "")
		}
	} else {
		if dal.UpdateProduct(m) > 0 {
			writeJsonValue(w, "0", "success", "")
		} else {
			writeJsonValue(w, "1", "fail", "")
		}
	}

}

func writeJsonValue(w http.ResponseWriter, code, message, data string) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.Contains(data, "{") {
		w.Write([]byte(data))
	} else {
		if data == "" {
			data = "[]"
		}
		ret := "{\"code\":" + code + ",\"message\":\"" + message + "\",\"data\":" + data + "}"

		w.Write([]byte(ret))
	}

}

//删除
func DeleteTask(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	ret := dal.DelProduct(id)
	if ret > 0 {
		writeJsonValue(w, "0", "success", "")
	} else {
		writeJsonValue(w, "1", "fail", "")
	}

}

//获取产品列表
func GetProductListTask(w http.ResponseWriter, r *http.Request) {

	pageid := r.FormValue("pageid")
	pagesize := r.FormValue("pagesize")
	userid := r.FormValue("userid")
	productname := r.FormValue("productname")
	str := dal.GetProductList(pageid, pagesize, userid, productname)
	writeJsonValue(w, "0", "success", str)

}

//把文件拷贝到运行环境方便测试调试
func CopyFileToRunningPath() {
	curpath := public.GetCurDir()
	public.Log("curpath=" + curpath)
	PthSep := string(os.PathSeparator)
	curdir := public.GetCurRunPath()
	public.Log("curdir=" + curdir)
	if PthSep == "\\" {
		//windows
		public.Log("copy cer ")

		datapath := curpath + PthSep + "db"
		htmlpath := curpath + PthSep + "html"

		public.CreatePath(datapath)
		public.CreatePath(htmlpath)

		public.CopyFiles("D:\\mybook\\bookSource\\ProductManage\\html", curpath+"\\html")
		public.CopyFile("D:\\mybook\\bookSource\\ProductManage\\db\\db.sqlite", datapath+"\\db.sqlite")

	} else {
		//其它操作系统
	}

}
