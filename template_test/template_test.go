package template_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html/template"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	type User struct {
		Name string
	}

	/*
		模板基本语法
		- 使用{{}}来包裹渲染数据
		- 使用.访问数据，.代表当前作用域的当前对象，类似this、self关键词
	*/
	tpl := template.New("hello-world")         //模板名
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, User{Name: "breeze"}) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `Hello, breeze`, buffer.String())
}

func TestMapData(t *testing.T) {
	tpl := template.New("hello-world")         //模板名
	tpl, err := tpl.Parse(`Hello, {{ .Name}}`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, map[string]string{"Name": "breeze"}) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `Hello, breeze`, buffer.String())
}

func TestSlice(t *testing.T) {
	tpl := template.New("hello-world")            //模板名
	tpl, err := tpl.Parse(`Hello, {{index . 0}}`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, []string{"breeze"}) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `Hello, breeze`, buffer.String())
}

func TestBasic(t *testing.T) {
	tpl := template.New("hello-world")    //模板名
	tpl, err := tpl.Parse(`Hello, {{.}}`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, 123) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `Hello, 123`, buffer.String())
}

/*
声明变量
- 去除空格和换行： -，注意要与其它元素空格隔开
- 声明变量：使用$，$xxx := value
- 执行方法调用：“调用者.方法 参数1 参数2”
*/

func TestFunc(t *testing.T) {
	tpl := template.New("hello-world") //模板名
	tpl, err := tpl.Parse(`
切片长度：{{len .Slice}}
{{printf "%.2f" 1.2345}}
Hello, {{.Hello "breeze"}}`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b", "c"},
	}) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `
切片长度：3
1.23
Hello, breeze///`, buffer.String())
}

type FuncCall struct {
	Slice []string
}

func (f FuncCall) Hello(name string) string {
	return name + "///"
}

/*
循环
- 使用range关键字：range $idx,$elem:=【slice】
不支持for i，也不支持for true
*/
func TestForLoop(t *testing.T) {
	tpl := template.New("hello-world") //模板名
	tpl, err := tpl.Parse(`
{{- range $idx, $elem := .Slice -}}
{{.}}
{{$idx}}-{{$elem}}
{{end}}`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, FuncCall{
		Slice: []string{"a", "b", "c"},
	}) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `a
0-a
b
1-b
c
2-c
`, buffer.String())
}

/*
eq: ==
ne: !=
lt: <
gt: >
le: <=
ge: >=
*/
func TestIfElse(t *testing.T) {
	tpl := template.New("hello-world") //模板名
	tpl, err := tpl.Parse(`
{{- if and (gt .Age 0) (le .Age 6)}}
儿童
{{else if and (gt .Age 6) (le .Age 18)}}
少年
{{else}}
成年
{{end -}}
`) //模板
	require.NoError(t, err)

	buffer := &bytes.Buffer{} //定义输出

	err = tpl.Execute(buffer, map[string]any{"Age": 12}) //渲染数据
	require.NoError(t, err)

	assert.Equal(t, `
少年
`, buffer.String())
}
