package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
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