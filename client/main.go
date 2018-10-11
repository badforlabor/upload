/**
 * Auth :   liubo
 * Date :   2018/10/11 15:26
 * Comment: post一个文件，同时带上一些额外参数，比如账号、密码
 */

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// 传输文件，亦或是传输内存中的内容
var gTransFile = false

func postFile() error {

	testfile := "test.bin"
	url := "http://127.0.0.1:9090/upload"

	var pendingReader io.Reader = nil

	//打开文件句柄操作
	file, err := os.Open(testfile)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer file.Close()

	buff := bytes.NewBufferString("memory content.")
	memReader := bufio.NewReader(buff)

	// 传输文件时，即可以是文件，也可以是内存中的内容。
	if gTransFile {
		pendingReader = file
	} else {
		pendingReader = memReader
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", testfile)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}


	//iocopy
	_, err = io.Copy(fileWriter, pendingReader)
	if err != nil {
		return err
	}


	////设置其他参数
	//params := map[string]string{
	//	"user": "test",
	//	"password": "123456",
	//}
	//
	////这种设置值得仿佛 和下面再从新创建一个的一样
	//for key, val := range params {
	//	_ = bodyWriter.WriteField(key, val)
	//}

	//和上面那种效果一样
	//建立第二个fields
	if fileWriter, err = bodyWriter.CreateFormField("user"); err != nil  {
		fmt.Println(err, "----------4--------------")
	}
	if _, err = fileWriter.Write([]byte("test")); err != nil {
		fmt.Println(err, "----------5--------------")
	}
	//建立第三个fieds
	if fileWriter, err = bodyWriter.CreateFormField("password"); err != nil  {
		fmt.Println(err, "----------4--------------")
	}
	if _, err = fileWriter.Write([]byte("123456")); err != nil {
		fmt.Println(err, "----------5--------------")
	}


	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()


	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil
}

// sample usage
func main() {
	postFile()
}

