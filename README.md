# Golang爬虫示例包

文件结构

```
自己用Golang原生包封装了一个爬虫库，源码见go get -u -v github.com/hunterhug/go_tool/spider
---- data   存放数据
---- example 爬虫例子
	--- pedaily 投资界爬虫 
```

使用说明:

```
go get -u -v github.com/hunterhug/spiderexample
```

## 一.投资界爬虫pedaily(pedaily.cn)

1. companysearch.exe可通过关键字查找一家机构的简单信息
2. companytouzi.exe可通过公司代号查找一家机构的投资情况


### 关键字查找一家机构的简单信息

<img src='https://raw.githubusercontent.com/hunterhug/spiderexample/master/img/pedaily1.png' />

<img src='https://raw.githubusercontent.com/hunterhug/spiderexample/master/img/pedaily2.png' />

<img src='https://raw.githubusercontent.com/hunterhug/spiderexample/master/img/pedaily3.png' />


关键字查找一家机构的简单信息,代码如下：

```
package main
import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/go_tool/spider"
	"github.com/hunterhug/go_tool/spider/query"
	"github.com/hunterhug/go_tool/util"
	"math"
	"os"
	"strings"
)

var client *spider.Spider
var dir = "./data/company/raw"
var dirdetail = "./data/company/detailraw"
var dirresult = "./data/company/result"

func main() {
	welcome()
	initx()
	mainx()
	//fmt.Printf("%#v", detail("http://zdb.pedaily.cn/company/show587/"))
}
func welcome(){
	fmt.Println(`
************************************************************

		投资界关键字查找公司信息

		1.输入关键字
		首先翻页后抓取详情页信息保存在data文件夹中

		2.查看结果
		查看data/company/result中csv文件

		作者:一只尼玛
		联系:569929309

		Golang大法

		/*
		go get -u -v github.com/hunterhug/spiderexample
		go build *.go，然后点击exe运行或go run *.go
		*/
************************************************************
`)
}
func initx() {
	var e error = nil
	client, e = spider.NewSpider(nil)
	if e != nil {
		panic(e.Error())
	}

	util.MakeDir(dir)
	util.MakeDir(dirresult)
	util.MakeDir(dirdetail)
}
func mainx() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Please input kyword: ")
		//fmt.Scanln(&keyword)
		keyword, err := inputReader.ReadString('\n')
		keyword = strings.Replace(keyword, "\n", "", -1)
		keyword = trip(strings.Replace(keyword, "\r", "", -1))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		e := util.MakeDir(dir + "/" + keyword)
		if e != nil {
			fmt.Println(e.Error())
		}
		//keyword := "投资管理股份有限公司"
		result, numt, e := featchcompany(keyword)
		if e != nil {
			fmt.Println(e.Error())
			continue
		}
		if len(result) == 0 {
			fmt.Println("empty")
			continue
		} else {
			num := int(math.Ceil(float64(numt) / 20.0))
			for i := 2; i <= num; i++ {
				temp, _, e := featchcompany(keyword + "/" + util.IS(i))
				if e != nil {
					fmt.Println(e.Error())
				} else {
					result = append(result, temp...)
				}
			}
			fmt.Printf("total comepany:%d\n", numt)
			txt := []string{}
			txt = append(txt, "公司名,英语名称,简称,详情URL,投资URL,资本类型,机构性质,注册地点,机构总部,投资阶段,成立时间,官网,简介,联系方式")
			for _, k := range result {
				detailr := detail(k["href"])

				stemp := k["title"] + "," + detailr["english"] + "," + k["abbr"] + "," + k["href"] + "," + k["hreft"] + "," + detailr["zitype"] + "," + detailr["jx"]
				stemp = stemp + "," + detailr["rloca"] + "," + detailr["tloca"] + "," + detailr["tstage"] + "," + detailr["date"] + "," + detailr["website"] + "," + detailr["desc"] + "," + detailr["contact"]
				fmt.Println(stemp)
				txt = append(txt, stemp)
			}

			e := util.SaveToFile(dirresult+"/"+keyword+".csv", []byte(strings.Join(txt, "\n")))
			if e != nil {
				fmt.Println(e.Error())
			}
		}

		fmt.Println("------------")
	}
}

func detail(url string) map[string]string {
	returnmap := map[string]string{
		"english": "", //英文名称
		"zitype":  "", //资本类型
		"jx":      "", //机构性质
		"rloca":   "", //注册地点
		"tloca":   "", //机构总部
		"tstage":  "", //投资阶段
		"date":    "", //成立时间
		"website": "", //官方网站
		"desc":    "", //简介
		"contact": "", //联系方式
	}
	hashmd := util.Md5(url)
	keep := dirdetail + "/" + hashmd + ".html"
	body := []byte("")
	var e error = nil
	if util.FileExist(keep) {
		body, e = util.ReadfromFile(keep)
	} else {
		client.Url = url
		body, e = client.Get()
	}
	if e != nil {
		return returnmap
	}
	util.SaveToFile(keep, body)
	doc, e := query.QueryBytes(body)
	if e != nil {
		return returnmap
	}
	returnmap["english"] = trip(doc.Find(".info h2").Text())
	returnmap["contact"] = strings.Replace(trip(doc.Find("#contact").Text()), "\n", "<br/>", -1)
	returnmap["desc"] = strings.Replace(trip(doc.Find("#desc").Text()), "\n", "<br/>", -1)
	returnmap["website"] = trip(doc.Find("li.link a").Text())
	info := ""
	doc.Find(".info ul li").Each(func(i int, node *goquery.Selection) {
	temp:=node.Text()
		temp=trip(strings.Replace(temp,"\n","",-1))
		temp=strings.Replace(temp,"&nbsp;","",-1)
		temp=strings.Replace(temp,"　","",-1)
		info=info+"\n"+temp
	})

	dudu := tripemptyl(strings.Split(info, "\n"))
	for _, r := range dudu {
		rr := strings.Split(r, "：")
		dd := ""
		if len(rr) == 2 {
			dd = rr[1]
		} else {
			continue
		}
		if strings.Contains(r, "资本类型") {
			returnmap["zitype"] = dd
		} else if strings.Contains(r, "机构性质") {
			returnmap["jx"] = dd
		} else if strings.Contains(r, "注册地点") {
			returnmap["rloca"] = dd
		} else if strings.Contains(r, "成立时间") {
			returnmap["date"] = dd
		} else if strings.Contains(r, "机构总部") {
			returnmap["tloca"] = dd
		} else if strings.Contains(r, "投资阶段") {
			returnmap["tstage"] = dd
		} else {
		}
	}
	return returnmap
}

func trip(s string) string {
	return strings.TrimSpace(strings.Replace(s, ",", "", -1))
}

func tripemptyl(dudu []string) []string {
	returnlist := []string{}
	for _, r := range dudu {
		if trip(r) != "" {
			returnlist = append(returnlist, trip(r))
		}
	}
	return returnlist
}
func featchcompany(keyword string) ([]map[string]string, int, error) {
	returnmap := []map[string]string{}
	url := "http://zdb.pedaily.cn/company/w" + keyword
	rootdir := strings.Split(keyword, "/")[0]
	hashmd := util.Md5(url)
	keep := dir + "/" + rootdir + "/" + hashmd + ".html"
	if util.FileExist(keep) {
		dudu, _ := util.ReadfromFile(keep)
		return parsecompany(dudu)
	}
	fmt.Printf("featch:%s\n", url)
	client.Url = url
	body, err := client.Get()
	if err != nil {
		return returnmap, 0, err
	}
	e := util.SaveToFile(keep, body)
	if e != nil {
		fmt.Println(url + ":" + e.Error())
	}
	return parsecompany(body)
}

func parsecompany(body []byte) ([]map[string]string, int, error) {
	returnmap := []map[string]string{}
	d, e := query.QueryBytes(body)
	if e != nil {
		return returnmap, 0, e
	}
	total := d.Find(".total").Text()
	num, e := util.SI(total)
	if e != nil {
		return returnmap, 0, nil
	}
	d.Find(".company-list li").Each(func(i int, node *goquery.Selection) {
		temp := map[string]string{}
		content := node.Find(".txt a.f16")
		abbr := content.Next().Text()
		title := content.Text()
		href, ok := content.Attr("href")
		if !ok {
			return
		} else {
			href = "http://zdb.pedaily.cn" + href
		}
		//location := node.Find(".txt .location").Text()
		//desc := strings.Replace(node.Find(".desc").Text(), ",", "", -1)
		//desc = strings.Replace(desc, "\n", "", -1)
		//desc = strings.TrimSpace(desc)
		//desc = strings.Replace(desc, "\r", "", -1)
		temp["title"] = title
		temp["abbr"] = abbr
		temp["href"] = href
		hreft := strings.Split(href, "show")
		if len(hreft) == 0 {
			return
		}
		temp["hreft"] = "http://zdb.pedaily.cn/company/" + hreft[len(hreft)-1] + "vc/"
		//temp["location"] = location
		//temp["desc"] = desc
		//fmt.Printf("%s,%s,%s,%s,%s\n", title, href, temp["hreft"], abbr, location)
		returnmap = append(returnmap, temp)
	})

	return returnmap, num, nil
}

```

### 公司代号查找一家机构的投资情况

<img src='https://raw.githubusercontent.com/hunterhug/spiderexample/master/img/pedaily4.png' />

<img src='https://raw.githubusercontent.com/hunterhug/spiderexample/master/img/pedaily5.png' />

<img src='https://raw.githubusercontent.com/hunterhug/spiderexample/master/img/pedaily6.png' />

公司代号查找一家机构的投资情况代码如下：

```
package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunterhug/go_tool/spider"
	"github.com/hunterhug/go_tool/spider/query"
	"github.com/hunterhug/go_tool/util"
	"math"
	"os"
	"regexp"
	"strings"
)

var tclient *spider.Spider
var tresult = "./data/companyt/result"
var traw = "./data/companyt/raw"

func main() {
	dudu()
	inittouzi()
	tmain()
	//b, _ := util.ReadfromFile(tresult + "/3392.html")
	//l := parsetouzi(b)
	//fmt.Printf("%#v", l)
}

func parset(body []byte) ([]string, string) {
	returnlist := []string{}
	doc, _ := query.QueryBytes(body)
	total := doc.Find(".total").Text()
	doc.Find("#inv-list li").Each(func(i int, node *goquery.Selection) {
		href, ok := node.Find("dt.view a").Attr("href")
		if !ok {
			return
		}
		href = "http://zdb.pedaily.cn" + href
		fmt.Printf("Inv: %s:%s\n", node.Find(".company a").Text(), href)
		returnlist = append(returnlist, href)
	})
	return returnlist, total
}
func inittouzi() {
	var e error = nil
	tclient, e = spider.NewSpider(nil)
	if e != nil {
		panic(e.Error())
	}

	util.MakeDir(tresult)
	util.MakeDir(traw)
}
func tmain() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Please input url: ")
		//fmt.Scanln(&keyword)
		keyword, err := inputReader.ReadString('\n')
		keyword = strings.Replace(keyword, "\n", "", -1)
		keyword = trip(strings.Replace(keyword, "\r", "", -1))
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		r, _ := regexp.Compile(`[\d]+`)
		mark := r.FindString(keyword)
		if mark == "" {
			fmt.Println("找不到")
			continue
		}
		urls := []string{}
		loop := []string{}
		loop = append(loop, "y-2004")
		tyear, _ := util.SI(util.TodayString(1))
		for i := 2014; i <= tyear; i++ {
			loop = append(loop, "y"+util.IS(i))
		}

		result := []map[string]string{}
		for _, pp := range loop {
			url := "http://zdb.pedaily.cn/company/" + mark + "/vc/" + pp
			body, e := fetchpage(url)
			if e != nil {
				fmt.Println(e.Error())
				continue
			} else {
				fmt.Println("fetch " + url)
			}
			//e = util.SaveToFile(tresult+"/"+mark+".html", body)
			//if e != nil {
			//	fmt.Println(e.Error())
			//}
			l, t := parset(body)
			fmt.Printf("%s:total:%s\n", pp, t)
			total, e := util.SI(t)
			if e != nil {
				fmt.Println(e.Error())
				continue
			}
			if total == 0 {
				fmt.Println("empty")
				continue
			}
			urls = append(urls, l...)
			page := int(math.Ceil(float64(total) / 20.0))
			for i := 2; i <= page; i++ {
				url := "http://zdb.pedaily.cn/company/" + mark + "/vc/" + pp + "/" + util.IS(i)
				body, e = fetchpage(url)
				if e != nil {
					fmt.Println(e.Error())
					continue
				} else {
					fmt.Println("fetch " + url)
				}
				temp, _ := parset(body)
				urls = append(urls, temp...)
			}
		}
		if len(urls) == 0 {
			fmt.Println("empty")
			continue
		}
		//fmt.Printf("%#v\n", urls)
		for _, url := range urls {
			body := []byte("")
			var e error = nil
			keep := traw + "/" + util.Md5(url) + ".html"
			if util.FileExist(keep) {
				body, e = util.ReadfromFile(keep)
			} else {
				body, e = fetchpage(url)
			}
			if e != nil {
				fmt.Println(e.Error())
				continue
			} else {
				fmt.Println("fetch " + url)
			}
			util.SaveToFile(keep, []byte(body))
			dududu := parsetouzi(body)
			dududu["url"] = url
			result = append(result, dududu)
		}

		if len(result) == 0 {
			fmt.Println("empty")
			continue
		}
		s := []string{"页面,事件名称,融资方,投资方,金额,融资时间,轮次,所属行业,简介"}
		for _, jinhan := range result {
			s = append(s, jinhan["url"]+","+jinhan["name"]+","+jinhan["rf"]+","+jinhan["tf"]+","+jinhan["money"]+","+jinhan["date"]+","+jinhan["times"]+","+jinhan["han"]+","+jinhan["desc"])
		}

		util.SaveToFile(tresult+"/"+mark+".csv", []byte(strings.Join(s, "\n")))
	}
}

func parsetouzi(body []byte) map[string]string {
	returnmap := map[string]string{
		"name":  "",
		"rf":    "",
		"tf":    "",
		"money": "",
		"date":  "",
		"times": "",
		"han":   "",
		"desc":  "",
	}
	doc, _ := query.QueryBytes(body)
	returnmap["name"] = doc.Find(".info h1").Text()
	returnmap["desc"] = strings.Replace(trip(doc.Find("#desc").Text()), "\n", "<br/>", -1)

	info := ""
	doc.Find(".info ul li").Each(func(i int, node *goquery.Selection) {
		temp := node.Text()
		temp = trip(strings.Replace(temp, "\n", "", -1))
		temp = strings.Replace(temp, "&nbsp;", "", -1)
		temp = strings.Replace(temp, "　", "", -1)
		info = info + "\n" + temp
	})

	//fmt.Println(info)
	dudu := tripemptyl(strings.Split(info, "\n"))
	//fmt.Printf("%#v\n", dudu)
	for _, r := range dudu {
		rr := strings.Split(r, "：")
		dd := ""
		if len(rr) == 2 {
			dd = strings.Replace(rr[1], " ", "", -1)
		} else {
			continue
		}
		if strings.Contains(r, "融") && strings.Contains(r, "方") {
			returnmap["rf"] = dd
		} else if strings.Contains(r, "投") && strings.Contains(r, "方") {
			returnmap["tf"] = dd
		} else if strings.Contains(r, "金") && strings.Contains(r, "额") {
			returnmap["money"] = dd
		} else if strings.Contains(r, "融资时间") {
			returnmap["date"] = dd
		} else if strings.Contains(r, "轮") && strings.Contains(r, "次") {
			returnmap["times"] = dd
		} else if strings.Contains(r, "所属行业") {
			returnmap["han"] = dd
		} else {
		}
	}
	return returnmap
}
func fetchpage(url string) ([]byte, error) {
	tclient.Url = url
	return tclient.Get()
}
func trip(s string) string {
	return strings.TrimSpace(strings.Replace(s, ",", "", -1))
}

func tripemptyl(dudu []string) []string {
	returnlist := []string{}
	for _, r := range dudu {
		if trip(r) != "" {
			returnlist = append(returnlist, trip(r))
		}
	}
	return returnlist
}
func dudu() {
	fmt.Println(`
************************************************************

		投资界根据公司查找投资案例

		1.输入URL或者输入数字587等
		如：http://zdb.pedaily.cn/company/587/vc/

		2.查看结果
		查看data/company/tresult中csv文件

		作者:一只尼玛
		联系:569929309

		Golang大法

		/*
		go get -u -v github.com/hunterhug/spiderexample
		go build *.go，然后点击exe运行或go run *.go
		*/
************************************************************
`)
}

```

