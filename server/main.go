/**
 * Auth :   liubo
 * Date :   2018/10/11 15:35
 * Comment: 
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 将接收到的文件保存成文件，亦或是打印到控制台
var gSaveFile = false

// 处理/upload 逻辑
func upload(w http.ResponseWriter, r *http.Request)  {


	fmt.Println("method:", r.Method) //POST

	//因为上传文件的类型是multipart/form-data 所以不能使用 r.ParseForm(), 这个只能获得普通post
	r.ParseMultipartForm(32 << 20) //上传最大文件限制32M

	user := r.Form.Get("user")
	password := r.Form.Get("password")


	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err, "--------1------------")//上传错误
	}
	defer file.Close()

	// 保存成文件或者打印出来
	if gSaveFile {
		f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err == nil {
			defer f.Close()
			io.Copy(f, file)
		}
	} else {
		buff := bytes.NewBufferString("")
		writer := bufio.NewWriter(buff)
		io.Copy(writer, file)
		fmt.Println("接受到的文件内容是：", buff.String())
	}

	fmt.Println("接收到文件", user, password, handler.Filename) //test 123456 json.zip


}


func main() {

	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe("127.0.0.1:9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
