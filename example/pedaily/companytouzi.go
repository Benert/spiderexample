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
var dir = "./data/companytouzi/raw"
var dirresult = "./data/companytouzi/result"

func inittouzi() {
	var e error = nil
	client, e = spider.NewSpider(nil)
	if e != nil {
		panic(e.Error())
	}

	util.MakeDir(dir)
	util.MakeDir(dirresult)
}
func main() {
	inittouzi()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Please input url: ")
		keyword, err := inputReader.ReadString('\n')
		keyword = strings.Replace(keyword, "\n", "", -1)
		keyword = strings.Replace(keyword, "\r", "", -1)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		rootdirt := strings.Split(keyword, "company/")
		rootdir := strings.Split(rootdirt[len(rootdirt)-1], "/")[0]
		util.MakeDir(dir + "/" + rootdir)
		result, numt, e := featchcompanytouzi(keyword)
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
				temp, _, e := featchcompanytouzi(keyword + "/" + util.IS(i))
				if e != nil {
					fmt.Println(e.Error())
				} else {
					result = append(result, temp...)
				}
			}
			fmt.Printf("total comepanytouzi:%d\n", numt)
			txt := []string{}
			for _, k := range result {
				stemp := k["title"] + "," + k["abbr"] + "," + k["href"] + "," + k["hreft"] + "," + k["location"] + "," + k["desc"]
				txt = append(txt, stemp)
			}

			util.SaveToFile(dirresult+"/"+keyword+".csv", []byte(strings.Join(txt, "\n")))
		}

		fmt.Println("------------")
	}
}

func test2() {
	b, _ := util.ReadfromFile(dir + "/" + "search.html")
	a, num, _ := parsecompanytouzi(b)
	fmt.Printf("%v,%#v", num, a)
}
func featchcompanytouzi(keyword string) ([]map[string]string, int, error) {
	returnmap := []map[string]string{}
	url := keyword
	rootdirt := strings.Split(keyword, "company/")
	rootdir := strings.Split(rootdirt[len(rootdirt)-1], "/")[0]
	hashmd := util.Md5(url)
	keep := dir + "/" + rootdir + "/" + hashmd + ".html"
	if util.FileExist(keep) {
		dudu, _ := util.ReadfromFile(keep)
		return parsecompanytouzi(dudu)
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
	return parsecompanytouzi(body)
}

func parsecompanytouzi(body []byte) ([]map[string]string, int, error) {
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
		location := node.Find(".txt .location").Text()
		desc := strings.Replace(node.Find(".desc").Text(), ",", "", -1)
		desc = strings.Replace(desc, "\n", "", -1)
		desc = strings.TrimSpace(desc)
		desc = strings.Replace(desc, "\r", "", -1)
		temp["title"] = title
		temp["abbr"] = abbr
		temp["href"] = href
		hreft := strings.Split(href, "show")
		if len(hreft) == 0 {
			return
		}
		temp["hreft"] = "http://zdb.pedaily.cn/company/" + hreft[len(hreft)-1] + "vc/"
		temp["location"] = location
		temp["desc"] = desc
		fmt.Printf("%s,%s,%s,%s,%s\n", title, href, temp["hreft"], abbr, location)
		returnmap = append(returnmap, temp)
	})

	return returnmap, num, nil
}
