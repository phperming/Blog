package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:8088"

	//1.请求 --模拟用户访问浏览器
	var(
		resp *http.Response
		err error
	)

	resp, err = http.Get(baseURL)

	//2.检测是否无错，且200
	assert.NoError(t,err,"有错误发生，err不为空")
	assert.Equal(t,200,resp.StatusCode,"应返回状态码 200")
}

func TestAboutPage(t *testing.T) {
	baseURL := "http://localhost:8088"

	//1.请求 模拟客户浏览器请求
	var (
		resp *http.Response
		err error
	)

	resp, err = http.Get(baseURL)

	//2.检测  是否无错误 且 200
	assert.NoError(t,err,"有错误发生，err 不为空")
	assert.Equal(t,200,resp.StatusCode,"应该返回状态码 200")
}

func TestAllPage(t *testing.T) {
	baseURL := "http://localhost:8088"

	//1.声明加初始化测试数据
	var tests = []struct{
		method string
		url string
		expected int
	}{
		{"GET","/",200},
		{"GET","/about",200},
		{"GET","/notfound",404},
		{"GET","/articles",200},
		{"GET","/articles/create",200},
		{"GET","/articles/3",200},
		{"GET","/articles/3/edit",200},
		{"POST","/articles/3",200},
		{"POST","/articles",200},
		{"POST","/articles/1/delete",404},
	}

	//2.遍历所有测试
	for _,test := range tests {
		t.Logf("当前请求URL:%v \n",test.url)

		var (
			resp *http.Response
			err error
		)

		switch {
		case test.method == "GET" :
			resp,err = http.Get(baseURL+ test.url)
		default:
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL+test.url, data)
		}

		assert.NoError(t,err,"请求 "+test.url+" 报错")
		assert.Equal(t,200,resp.StatusCode,test.url + " 应该返回状态码 " + strconv.Itoa(test.expected))
	}
}