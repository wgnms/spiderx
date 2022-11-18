package sx
import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/color"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)
var text string
var num int
var flt float32
func Green(Text string){
	//color.Cyan("蓝绿色.")
	//color.Blue("蓝色.")
	//color.Red("红色.")
	//color.Magenta("品平.")
	//color.White("白色.")
	//color.Black("黑色.")
	//color.Yellow("黄色")
	//SWarn := color.New(color.Bold, color.FgYellow).PrintlnFunc()
	//SError := color.New(color.Bold, color.FgRed).PrintlnFunc()
	//SInfo := color.New(color.Bold, color.FgWhite).PrintlnFunc()
	//Faint := color.New(color.Faint, color.FgHiWhite).PrintlnFunc()
	//Italic := color.New(color.Italic, color.FgHiWhite).PrintlnFunc()
	//BlinkSlow := color.New(color.BlinkSlow, color.FgHiWhite).PrintlnFunc()
	//BlinkRapid := color.New(color.BlinkRapid, color.FgHiWhite).PrintlnFunc()
	//ReverseVideo := color.New(color.ReverseVideo, color.FgHiWhite).PrintlnFunc()
	//Concealed := color.New(color.Concealed, color.FgHiWhite).PrintlnFunc()
	//CrossedOut := color.New(color.CrossedOut, color.FgHiWhite).PrintlnFunc()
	color.New(color.Bold, color.FgGreen).PrintlnFunc()(Text)
}
func Red(Text string){
	color.New(color.Bold, color.FgRed).PrintlnFunc()(Text)
}
func Warn(Text string){
	color.New(color.Bold, color.FgYellow).PrintlnFunc()(Text)
}
func Info(Text string){
	color.New(color.Bold, color.FgWhite).PrintlnFunc()(Text)
}
func CatchError(){
	err := recover()
	if err != nil {
		_, file, lineNo, _ := runtime.Caller(0)
		color.Yellow(fmt.Sprintf("文件:%s",file))
		color.Red(fmt.Sprintf("行号:%d",lineNo))
		color.Red(fmt.Sprintf("错误:%s",err))
	}
}
func Gbk_to_utf8(s string) (rt string, err error) {
	// transform GBK bytes to UTF-8 bytes
	r := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return string(b),err
	}
	return string(b),nil
}
type Bytes struct {C []byte }
func (b *Bytes) Gbk() string{
	res,_:=Gbk_to_utf8(string(b.C))
	return res
}
func (b *Bytes) Utf8() string{
	return string(b.C)
}
func (b *Bytes) Content() []byte{
	return b.C
}
type Response struct {content []byte }
func (resp *Response)Text(encoding string) (string) {
	encoding=strings.ToLower(encoding)
	if (encoding=="utf-8" || encoding=="utf8"){
		return string(resp.content)
	}else{
		rt,_:=Gbk_to_utf8(string(resp.content))
		return rt
	}
}
func (resp *Response)Json() *simplejson.Json {
	res,_:=simplejson.NewJson(resp.content)
	return res
}
func (resp *Response)Content() []byte {
	return resp.content
}
func Get_requests(url string, headers string) (Response, error) {
	req,_:=http.NewRequest("GET",url,nil)
	arr:= strings.Split(headers,"\n")
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if headers!=""{
		for _, x := range arr {
			kv:=strings.SplitN(x,":",2) //分成两个
			fmt.Println(kv)
			req.Header.Set(kv[0],kv[1])
		}
	}

	client:=http.Client{}
	client.Timeout=15*time.Second  //超时
	resp,err:=client.Do(req)
	defer CatchError()
	defer resp.Body.Close()
	//获取HTML
	content,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return Response{content: content},err
	}
	return Response{content: content}, nil
}
func Post_requests(url string, dataStr string, headersStr string) (Response, error) {
	req,_:=http.NewRequest("POST", url,strings.NewReader(dataStr))
	arr:= strings.Split(headersStr,"\n")
	for _, x := range arr {
		kv:=strings.SplitN(x,":",2) //分成两个
		req.Header.Set(kv[0],kv[1])
	}
	client:=http.Client{}
	client.Timeout=15*time.Second  //超时
	resp,err:=client.Do(req)
	defer func(){
		err := recover()
		fmt.Println(err)
	}()
	defer resp.Body.Close()
	//获取HTML
	text,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return Response{content: text}, err
	}
	return Response{content: text},nil
}
func Print_progress_bar(str string,i int,count int,speed string,size int,char string,backchar string){
	var s []string
	i=i+1
	for k := 0; k < size; k++ {
		if k >=(i*size/count){
			s= append(s, backchar)
		}else{
			s= append(s, char)
		}
	}
	ss:=strings.Join(s,"")
	bfb:=(i*100)/count
	fmt.Fprintf(os.Stdout,"\r%s %3d%% %s %3d/%d %s",str, bfb,ss,i,count,speed)
}
type WriteCounter struct {
	Current float64
	Total float64
	Start_time time.Time
}
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Current += float64(n)
	c:=wc.Current/1024/1024
	t:=time.Now().Second()-wc.Start_time.Second()
	var speed string="0.0mb/s"
	if t!=0{
		speed=fmt.Sprintf("%.2fmb/s",c/float64(t))
	}
	fmt.Fprintf(os.Stdout,"\r%s %3d%% %s %.2fmb","下载进度",int(c*100/wc.Total),speed,c)
	return n, nil}
func Download_file_progress_size(filepath string, url string,headers string) (bool,error) {
	req,_:=http.NewRequest("GET",url,nil)
	arr:= strings.Split(headers,"\n")
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if headers!=""{
		for _, x := range arr {
			kv:=strings.SplitN(x,":",2) //分成两个
			fmt.Println(kv)
			req.Header.Set(kv[0],kv[1])
		}
	}
	client:=http.Client{}
	client.Timeout=15*time.Second  //超时
	resp,_:=client.Do(req)
	defer CatchError()
	defer resp.Body.Close()
	file,_:=os.Create(filepath)
	fmt.Println("保存路径:"+filepath)
	counter := &WriteCounter{}
	counter.Total= float64(resp.ContentLength/1024/1024)
	counter.Start_time=time.Now()
	io.Copy(file,io.TeeReader(resp.Body, counter))
	defer file.Close()
	return true,nil

}
func Download_file(save_path string, url string,headers string){
	resp,_:= Get_requests(url,headers)
	Save_file(save_path,resp.content)
}
func Download_file_progress_pool(save_path string, url string,headers string) {
	req,_:=http.NewRequest("GET",url,nil)
	arr:= strings.Split(headers,"\n")
	req.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36")
	if headers!=""{
		for _, x := range arr {
			kv:=strings.SplitN(x,":",2) //分成两个
			fmt.Println(kv)
			req.Header.Set(kv[0],kv[1])
		}
	}
	client:=http.Client{}
	client.Timeout=15*time.Second  //超时
	resp,_:=client.Do(req)
	size:=int(resp.ContentLength)
	fmt.Println(size)
	var block int =1024*1024
	n:=Get_page_count(block,size)
	fmt.Println(n)
	for i := 0; i <n; i++ {
		r, _ :=fmt.Printf("bytes=%d-%d",i*block,(i+1)*block)
		req.Header.Set("range", string(r))
		res,_:=client.Do(req)
		fmt.Println(string(r))
		content,_:=ioutil.ReadAll(res.Body)
		fmt.Println(len(content))
	}
	defer CatchError()
	defer resp.Body.Close()

	//获取HTML
	//content,_:=ioutil.ReadAll(resp.Body)

	//for i := 0; i <100; i++ {
	//	Print_progress_bar("下载文件",i,100,"10m/s",30,"#","_")
	//	time.Sleep(time.Second/20)
	//}
}
func Get_page_count(num int,count int) int{
	if count%num==0{
		return int(count/num)
	}else{
		return int(count/num)+1
	}

}
type IP_obj struct {
	City    string `json:"cname"`
	Ip   string `json:"cip"`
}//获取外网IP
func Get_wlan_ip() (IP_obj){
	resp,_:= Get_requests("http://pv.sohu.com/cityjson?ie=utf-8","")
	compile := regexp.MustCompile("({.*?})")
	str :=compile.FindString(resp.Text("utf-8"))
	var rt IP_obj
	json.Unmarshal([]byte(str),&rt)
	return rt

}

type Xpath struct {
	html string
	root *html.Node
}
func (sp *Xpath) Html(html string){
	sp.html =html
	root,_:=htmlquery.Parse(strings.NewReader(sp.html))
	sp.root =root
}
func (sp *Xpath)Get_title() string{
	//取title
	x:=htmlquery.FindOne(sp.root, "//title")
	return htmlquery.InnerText(x)
}
func (sp *Xpath)Get_attrs(path string, attr string) []string {
	//取属性
	var lst  []string
	var nodes = htmlquery.Find(sp.root, path)
	for _, x := range nodes {
		if attr==""{
		}else if attr=="text()"{
			lst = append(lst, htmlquery.InnerText(x))
		}else{
			lst = append(lst, string(htmlquery.SelectAttr(x, attr)))
		}
	}
	return lst
}
func (sp *Xpath)Get_attr(path string, attr string) (string) {
	//取第一个
	first:=htmlquery.FindOne(sp.root,path)

	if attr==""{
		return ""
	}else if attr == "text()"{
		return  htmlquery.InnerText(first)
	}else{
		return  htmlquery.SelectAttr(first, attr)
	}
}
func (sp *Xpath)Get_node_attrs(node *html.Node,path string, attr string) []string {
	//取属性
	var lst  []string
	var nodes = htmlquery.Find(node, path)
	for _, x := range nodes {
		if attr==""{
		}else if attr=="text()"{
			lst = append(lst, htmlquery.InnerText(x))
		}else{
			lst = append(lst, string(htmlquery.SelectAttr(x, attr)))
		}
	}
	return lst
}
func (sp *Xpath)Get_node_attr(node *html.Node, path string, attr string) (string) {
	//取第一个
	first:=htmlquery.FindOne(node,path)
	if attr==""{
		return ""
	}else if attr == "text()"{
		return  htmlquery.InnerText(first)
	}else{
		return  htmlquery.SelectAttr(first, attr)
	}
}
func (sp *Xpath)Get_html(node *html.Node) string {
	return htmlquery.OutputHTML(node,true)
}
func (sp *Xpath)Get_node(path string)(node *html.Node){
	return htmlquery.FindOne(sp.root,path)
}
func (sp *Xpath)Get_nodes(path string) (node []*html.Node) {
	return htmlquery.Find(sp.root,path)
}

func Save_file(path string,content []byte){
	file,_:=os.Create(path)
	defer CatchError()
	file.Write(content)
	file.Close()
}//保存文件
func Load_file(path string) ([]byte,error){
	file,_:=os.Open(path)
	defer CatchError()
	return ioutil.ReadAll(file)
}//加载文件
func Re_findall(rule string, s string) []string {
	demo, _ := regexp.Compile(rule)
	fmt.Println(demo.FindAllString(s,-1)) //[foo]
	return demo.FindAllString(s,-1)
}//正则查询
func Re_search(rule string,s string) string{
	demo, _ := regexp.Compile(rule)
	return demo.FindString(s)
}//正则查询
func Map_from_data_str(s string) map[string]string{
	data := make(map[string]string)
	for _, x := range strings.Split(s,"&") {
		kv:=strings.SplitN(x,"=",2)
		data[kv[0]]=kv[1]
	}
	return data
}
func Map_from_headers_str(s string) map[string]string{
	data := make(map[string]string)
	for _, x := range strings.Split(s,"&") {
		kv:=strings.SplitN(x,":",2)
		data[kv[0]]=kv[1]
	}
	return data
}
func Get_config(name string, data map[string]string) map[string]string {
	cfg, err :=ini.Load(name)
	if err !=nil{
		//文件不存在  data写入ini
		os.Create(name)
		cfg, _ :=ini.Load(name)
		for k, v := range data {
			cfg.Section("conf").Key(k).SetValue(v)
		}
		cfg.SaveTo(name)
		return data
	}else{
		//存在 读取ini
		fmt.Println(2)
		d:=make(map[string]string)
		ss:= cfg.SectionStrings()
		cfg.Section("conf").Keys()
		fmt.Println(ss[1])
		for _, x := range cfg.ChildSections("conf") {
			fmt.Println(x)
			xx:=x.Keys()
			fmt.Println(xx,1)
			fmt.Println(d,2)
		}
		//fmt.Println(cfg)


		return d
	}

}//加载配置文件
func Load_json_content(content []byte) *simplejson.Json {
	res,_:=simplejson.NewJson(content)
	return res
}//加载json
func Load_json_file(path string) *simplejson.Json {
	content,_:= Load_file(path)
	return Load_json_content(content)

}//加载json文件
func Save_json_file(path string,json_obj *simplejson.Json){
	content,_:=json_obj.MarshalJSON()
	Save_file(path,content)
}//保存json文件
func Exejs(str string,func_name string,args... interface{}) string{
	vm:=otto.New()
	vm.Run(str)
	value,_:=vm.Call(func_name,nil,args...)
	return value.String()
}
func Bytes_from_hex(str string) []byte {
	rt, _ := hex.DecodeString(str)
	return rt
}
func Run_command(cmd string) (Bytes, error) {
	s:=strings.SplitN(cmd," ",2)
	name:=s[0]
	var args []string
	if len(s)>1{
		args=strings.Split(s[1]," ")
	}
	fmt.Println(args)
	res := exec.Command(name,args...)
	// 执行命令，并返回结果
	output,err := res.Output()
	if err != nil {
		return Bytes{}, nil

	}
	return Bytes{output},nil
}

//AES
type Aes_128 struct{
	Length int
}
func (a *Aes_128)Encrypt(content []byte,key []byte,iv []byte) []byte {
	a.Length = 16 - (len(content) % 16)
	for i := 0; i < a.Length; i++ {
		content = append(content, byte(0))
	}
	c,_:=aes.NewCipher(key)
	if iv==nil {
		iv=[]byte("0000000000000000")
	}
	aes_obj:=cipher.NewCBCEncrypter(c,iv)
	rtbuf := make([]byte, len(content))
	aes_obj.CryptBlocks(rtbuf,content)
	return rtbuf
}
func (a *Aes_128)Decrypt(content []byte,key []byte,iv []byte) (string,error) {
	if (len(content) < 16) || (len(content)%16 != 0) {
		return "", errors.New("数据长度非16的倍数")
	}
	c,_:=aes.NewCipher(key)
	aes_obj:=cipher.NewCBCDecrypter(c,iv)
	rtbuf := make([]byte, len(content))
	aes_obj.CryptBlocks(rtbuf,content)
	return string(rtbuf[:(len(content)-a.Length)]), nil
}

//线程池
type Job struct { //任务类型
	Id int
	Args map[string]interface{}
}
type Result struct {  //返回值类型
	Id int
	Res interface{}
}
type Pool struct {
	ThreadNum int //线程数
	Func func(a map[string]interface{}) interface{} //函数接口
	Args []map[string]interface{} //参数类
	Show bool //显示线程
}
func (p *Pool)Work(workName string, jobs chan Job, results chan Result) {
	for job := range jobs {
		if p.Show{fmt.Println(fmt.Sprintf("开始:%s 第%d个任务",workName, job.Id))}
		rt:=p.Func(job.Args)
		results <- Result{job.Id,rt}
		if p.Show{fmt.Println(fmt.Sprintf("结束:%s 第%d个任务",workName, job.Id))}
	}
}
func (p *Pool)Results() []Result {
	jobs_length:= len(p.Args)
	jobs := make(chan Job, jobs_length)
	results := make(chan Result, jobs_length)
	for i := 0; i < p.ThreadNum; i++ {
		go p.Work(fmt.Sprintf("第%d个线程", i+1), jobs, results)
	}
	//写入所有jobs
	for i := 0; i < jobs_length; i++ {
		jobs <- Job{i+1,p.Args[i]}

	}
	//取出所有results
	var res []Result
	for i := 0; i < jobs_length; i++ {
		res=append(res,<-results )
	}
	close(results)
	close(jobs)
	return res
}

//go get "xxx"
//go mod tidy 安装缺少的模块
//go mod verify 校验依赖
//go mod why
//go mod init
//go mod download
//go mod edit
//go mod verdor 依赖复制到vendor下
//go list