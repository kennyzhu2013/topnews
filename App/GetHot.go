package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/bitly/go-simplejson"
	"github.com/tophubs/TopList/Common"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HotData struct {
	Code    int
	Message string
	Data    interface{}
}

type Spider struct {
	DataType string
}

func SaveDataToJson(data interface{}) string {
	Message := HotData{}
	Message.Code = 0
	Message.Message = "获取成功"
	Message.Data = data
	jsonStr, err := json.Marshal(Message)
	if err != nil {
		log.Fatal("序列化json错误")
	}
	println(string(jsonStr))
	return string(jsonStr)

}

// GetV2EX V2EX
func (spider Spider) GetV2EX() []map[string]interface{} {
	url := "https://www.v2ex.com/?tab=hot"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败1")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败2")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败3")
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	document.Find(".item_title").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		println(url + " " + text)
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://www.v2ex.com" + url})
		}
	})
	return allData
}

func (spider Spider) GetITHome() []map[string]interface{} {
	url := "https://www.ithome.com/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	document.Find(".hot-list .bx ul li").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": url})
		}
	})
	return allData
}

// GetZhiHu 知乎:知乎是需要登录的, 每次需要单独cookie.
// 无法直接用cookie获取爬取，goquery不支持.
func (spider Spider) GetZhiHu() []map[string]interface{} {
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	url := "https://www.zhihu.com/hot"
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}

	// 设置登录cookie即可.
	request.Header.Add("Cookie", `_9755xjdesxxd_=32; _zap=dda53e17-1735-4484-92ef-615d6a2f6d4f; _xsrf=VxlKKV2n3wnGZN61TZg1VonI9pjA7kxr; d_c0=ABDaXipSihiPTsemKgU2KKePE-FbLQOSJwo=|1714379422; q_c1=0e5b8e45b9fb4510a3e545c52f8e8e51|1716973986000|1648615377000; _ga=GA1.2.513319920.1730270985; _ga_MKN1TB7ZML=GS1.2.1730270986.1.1.1730271130.0.0.0; __zse_ck=004_fycGi7SRFe2SN80i/T3AsTOeK8SeSpEGzwf1pH/aAZ7CuSe14Sl/ukCGNkF/4bXNeJGW=3KxxdVmrjECUvQulp88q0vQpMcrh687rYD3v77YcVBRpbeQRh4x8C52qdty-b0Im4GmxVdNaf6Ar6/09o7MMxwI1KbOzXqarkRDcObywzwGw4GbSz/76M+QN6q5AtxngcUJdnRBPPtoziiKwC1CAFMo2dkiAltjW339IchmD3TjQfAj8k4dysFkNr/jK; z_c0=2|1:0|10:1734934887|4:z_c0|80:MS4xN05FQkJnQUFBQUFtQUFBQVlBSlZUV1pQVm1pUjAtVHhOa2ota1kxUjNJanVrY2dOdkh5eHRBPT0=|9cbe3cc7c9aacbda5f54f94c0fe56b760d00b9e70757d461f6b69f7f5a93898b; Hm_lvt_98beee57fd2ef70ccdd5ca52b9740c49=1732515576,1733185966,1734934887; HMACCOUNT=2E7A3AC10D710648; tst=h; SESSIONID=wkd6kIiXkW9KmaJwmR1iMboP84ioLxa4yILC9cXMfmI; JOID=UlAdAEjXiXP83tT6P9Ripnhx8vYkqOIcyraVk2y_1kq8s7a4SHSPEZjV3PE034G1rf1aGe3q-RUSKU9QvqR7Bhw=; osd=VF8QAE3Rhn7829L1MtRnoHd88vMip-8cz7Canmy60EWxs7O-R3mPFJ7a0fEx2Y64rfhcFuDq_BMdJE9VuKt2Bhk=; BEC=ec64a27f4feb1b29e8161db426d61998; Hm_lpvt_98beee57fd2ef70ccdd5ca52b9740c49=1734938115`)
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)

	res, err := client.Do(request)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}

	// println(document.Text())
	// 知乎没法直接爬取，返回的的是一堆json数据没法解析，直接正则查找字符串.
	document.Contents().Each(func(i int, selection *goquery.Selection) {
		if selection.Children() != nil {
			// json文本内容在这里打印
			selection.Children().Each(func(i int, selection *goquery.Selection) {
				findedStr := selection.Text()
				// 正则表达式查找 "hotList" 开头的子字符串
				re := regexp.MustCompile(`"hotList":\s*(\[[^]]*\])`)
				// re := regexp.MustCompile(`\{"hotList":\[(.*?)\]}`)
				matches := re.FindStringSubmatch(findedStr)
				if len(matches) > 1 {
					// fmt.Println("Found hotList 1:", matches[1])

					// 使用正则表达式查找 excerptArea 中的 text 和 link 中的 url
					titlePattern := regexp.MustCompile(`"titleArea":\{"text":"(.*?)"\}`)
					linkPattern := regexp.MustCompile(`"link":\{"url":"(.*?)"\}`)

					// 查找所有的 excerpt 和 link
					titles := titlePattern.FindAllStringSubmatch(matches[1], -1)
					links := linkPattern.FindAllStringSubmatch(matches[1], -1)

					// 输出结果
					for i := 0; i < len(titles); i++ {
						// 如果截断了就跳过
						if i < len(links) {
							cleanedString := strings.ReplaceAll(links[i][1], "\\u002F", "/")
							fmt.Printf("Link: %s\n", cleanedString)
							fmt.Printf("Title: %s, Url: %s\n", titles[i][1], cleanedString)
							allData = append(allData, map[string]interface{}{"title": titles[i][1], "url": cleanedString})
						}
					}
				}
			})
		}
		println("GetZhiHu find text end============================:\n")
	})
	//document.Contents().Each(func(i int, selection *goquery.Selection) {
	//	url, boolUrl := selection.Find("a").Attr("href")
	//	text := selection.Find("h2").Text()
	//	println("GetZhiHu find text========================[" + url + "]=================:\n")
	//	if selection.Children() != nil {
	//		println("GetZhiHu find Children text=========================================:\n")
	//		selection.Children().Each(func(i int, selection *goquery.Selection) { println(selection.Text()) })
	//	}
	//	println("GetZhiHu find text end============================:\n")
	//
	//	if boolUrl {
	//		allData = append(allData, map[string]interface{}{"title": text, "url": url})
	//	}
	//})
	return allData
}

// GetTieBa 贴吧
func (spider Spider) GetTieBa() []map[string]interface{} {
	url := "http://tieba.baidu.com/hottopic/browse/topicList"
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	str, _ := ioutil.ReadAll(res.Body)
	js, err2 := simplejson.NewJson(str)
	if err2 != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	i := 1
	for i < 30 {
		test := js.Get("data").Get("bang_topic").Get("topic_list").GetIndex(i).MustMap()
		allData = append(allData, map[string]interface{}{"title": test["topic_name"], "url": test["topic_url"]})
		i++
	}
	return allData

}

// Github
func (spider Spider) GetGitHub() []map[string]interface{} {
	url := "https://github.com/trending"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}

	document.Find(".Box article").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find(".lh-condensed a")
		//desc := selection.Find(".col-9 .text-gray .my-1 .pr-4")
		//descText := desc.Text()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		descText := selection.Find("p").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "desc": descText, "url": "https://github.com" + url})
		}
	})
	return allData
}

func (spider Spider) GetBaiDu() []map[string]interface{} {
	url := "https://top.baidu.com/board?tab=realtime"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `top.baidu.com`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("table tr").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		MyText, _ := GbkToUtf8([]byte(text))
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": string(MyText), "url": url})
		}
	})
	return allData

}

func (spider Spider) Get36Kr() []map[string]interface{} {
	url := "https://36kr.com/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `36kr.com`)
	request.Header.Add("Referer", `https://36kr.com/`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".hotlist-item-toptwo").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := selection.Find("a p").Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://36kr.com" + url})
		}
	})
	document.Find(".hotlist-item-other-info").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://36kr.com" + url})
		}
	})
	return allData

}

func (spider Spider) GetQDaily() []map[string]interface{} {
	url := "https://www.qdaily.com/tags/29.html"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `www.qdaily.com`)
	request.Header.Add("Referer", `https://www.qdaily.com/tags/30.html`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".packery-item").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := selection.Find(".grid-article-bd h3").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.qdaily.com/" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetGuoKr() []map[string]interface{} {
	url := "https://www.guokr.com/scientific/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `www.guokr.com`)
	request.Header.Add("Referer", `https://www.guokr.com/scientific/`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("div .article").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h3 a")
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
			}
		}
	})
	return allData
}

func (spider Spider) GetHuXiu() []map[string]interface{} {
	url := "https://www.huxiu.com/article"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Host", `www.guokr.com`)
	request.Header.Add("Referer", `https://www.huxiu.com/channel/107.html`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str,_ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".article-item--large__content").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Find("h5").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.huxiu.com" + url})
			}
		}
	})
	document.Find(".article-item__content").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Find("h5").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.huxiu.com" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetDBMovie() []map[string]interface{} {
	url := "https://movie.douban.com/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".slide-container").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a")
		url, boolUrl := s.Attr("href")
		text := s.Find("p").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.huxiu.com" + url})
			}
		}
	})
	return allData
}

// 每日知乎
func (spider Spider) GetZHDaily() []map[string]interface{} {
	url := "http://daily.zhihu.com/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".row .box").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Find("span").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://daily.zhihu.com" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetSegmentfault() []map[string]interface{} {
	url := "https://segmentfault.com/hottest"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".news-list .news__item-info").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a:nth-child(2)").First()
		url, boolUrl := s.Attr("href")
		text := s.Find("h4").Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://segmentfault.com" + url})
			}
		}
	})
	return allData
}

func (spider Spider) GetHacPai() []map[string]interface{} {
	url := "https://hacpai.com/domain/play"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".hotkey li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h2 a")
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
			}
		}
	})
	return allData
}

func (spider Spider) GetWYNews() []map[string]interface{} {
	url := "http://news.163.com/special/0001386F/rank_whole.html"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("table tr").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("td a").First()
		url, boolUrl := s.Attr("href")
		text, _ := GbkToUtf8([]byte(s.Text()))
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData
}

func (spider Spider) GetWaterAndWood() []map[string]interface{} {
	url := "https://www.newsmth.net/nForum/mainpage?ajax"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//sss,_ := GbkToUtf8([]byte(string(str)))
	//fmt.Println(string(sss))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	// topics
	document.Find("#top10 li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a:nth-child(2)").First()
		url, boolUrl := s.Attr("href")
		text, _ := GbkToUtf8([]byte(s.Text()))
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.newsmth.net" + url})
				}
			}
		}
	})
	document.Find(".topics").Find("li").Each(func(i int, selection *goquery.Selection) {
		if i > 10 {
			s := selection.Find("a:nth-child(2)").First()
			url, boolUrl := s.Attr("href")
			text, _ := GbkToUtf8([]byte(s.Text()))
			if len(text) != 0 {
				if boolUrl {
					if len(allData) <= 100 {
						allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.newsmth.net" + url})
					}
				}
			}
		}
	})
	return allData
}

// http://nga.cn/

func (spider Spider) GetNGA() []map[string]interface{} {
	url := "http://nga.cn/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("h2").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData
}

func (spider Spider) GetCSDN() []map[string]interface{} {
	url := "https://www.csdn.net/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("#feedlist_id li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h2 a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData
}

// https://weixin.sogou.com/?pid=sogou-wsse-721e049e9903c3a7&kw=
func (spider Spider) GetWeiXin() []map[string]interface{} {
	url := "https://weixin.sogou.com/?pid=sogou-wsse-721e049e9903c3a7&kw="
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".news-list li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("h3 a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData
}

//

func (spider Spider) GetKD() []map[string]interface{} {
	url := "http://www.kdnet.net/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".indexside-box-hot li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text, _ := GbkToUtf8([]byte(s.Text()))
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData
}

// http://www.mop.com/

func (spider Spider) GetMop() []map[string]interface{} {
	url := "http://www.mop.com/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find(".swiper-slide").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := selection.Find("h2").Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	document.Find(".tabel-right").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := selection.Find("h3").Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData[:15]
}

// https://www.chiphell.com/

func (spider Spider) GetChiphell() []map[string]interface{} {
	url := "https://www.chiphell.com/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("#frameZ3L5I7 li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	// portal_block_530_content
	document.Find("#portal_block_530_content dt").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	// frame-tab move-span cl
	document.Find("#portal_block_560_content dt").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	// portal_block_564_content
	document.Find("#portal_block_564_content dt").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	// portal_block_568_content
	document.Find("#portal_block_568_content dt").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	// portal_block_569_content
	document.Find("#portal_block_569_content dt").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	// portal_block_570_content
	document.Find("#portal_block_570_content dt").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": "https://www.chiphell.com/" + url})
				}
			}
		}
	})
	return allData
}

// http://jandan.net/

func (spider Spider) GetJianDan() []map[string]interface{} {
	url := "http://jandan.net/"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	res, err := client.Do(request)

	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	//str, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(str))
	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	document.Find("h2").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("a").First()
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if len(text) != 0 {
			if boolUrl {
				if len(allData) <= 100 {
					allData = append(allData, map[string]interface{}{"title": string(text), "url": url})
				}
			}
		}
	})
	return allData
}

// https://dig.chouti.com/

func (spider Spider) GetChouTi() []map[string]interface{} {
	url := "https://dig.chouti.com/top/24hr?_=" + strconv.FormatInt(time.Now().Unix(), 10) + "163"
	url2 := "https://dig.chouti.com/link/hot?afterTime=" + strconv.FormatInt(time.Now().Unix(), 10) + "026000" + "&_=" + strconv.FormatInt(time.Now().Unix(), 10) + "667"
	res, err := http.Get(url)
	res2, _ := http.Get(url2)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	str, _ := ioutil.ReadAll(res.Body)
	str2, _ := ioutil.ReadAll(res2.Body)
	js, err2 := simplejson.NewJson(str)
	js2, _ := simplejson.NewJson(str2)
	if err2 != nil {
		fmt.Println("抓取" + spider.DataType + "失败")
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	i := 1
	for i < 30 {
		test := js.Get("data").GetIndex(i).MustMap()
		if test["title"] != nil && test["url"] != nil {
			allData = append(allData, map[string]interface{}{"title": test["title"], "url": test["url"]})
		}
		i++
	}
	j := 1
	for j < 60 {
		test := js2.Get("data").GetIndex(j).MustMap()
		if test["title"] != nil && test["url"] != nil {
			allData = append(allData, map[string]interface{}{"title": test["title"], "url": test["url"]})
		}
		j++
	}
	return allData

}

/*
*
部分热榜标题需要转码
*/
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// ExecGetData 执行每个分类数据
func ExecGetData(spider Spider) {
	reflectValue := reflect.ValueOf(spider)
	dataType := reflectValue.MethodByName("Get" + spider.DataType)
	//println("ExecGetData is:" + spider.DataType)
	//println(dataType.Type())
	data := dataType.Call(nil)
	originData := data[0].Interface().([]map[string]interface{})
	start := time.Now()

	// update实际使用判断获取的数据是否更新了.
	if len(originData) > 0 {
		Common.MySql{}.GetConn().Where(map[string]string{"dataType": spider.DataType}).Update("hotData2", map[string]string{"str": SaveDataToJson(originData)})
	}
	println("ExecGetData begins======")
	// Common.MySql{}.GetConn().Where(map[string]string{"dataType": spider.DataType}).Update("hotData2", map[string]string{"str": SaveDataToJson(originData)})
	group.Done()
	seconds := time.Since(start).Seconds()
	fmt.Printf("耗费 %.2fs 秒完成抓取%s 记录：%d", seconds, spider.DataType, len(originData))
	fmt.Println()

}

var group sync.WaitGroup

func main() {
	allData := []string{
		"ZhiHu",
		"WeiBo", // chromedp实现
		"TieBa",
		"DouBan", // chromedp实现.
		"HuPu",   // chromedp实现.

		// 没法直接用goquery爬取，用chromedp方式显得太复杂.
		// 检查到这.
		// sql中执行delete from hotData2 where dataType='TianYa'
		// "TianYa",

		"BaiDu",
		"36Kr",
		"QDaily",
		"GuoKr",
		"HuXiu",
		"ZHDaily",
		"Segmentfault",
		"WYNews",
		"WaterAndWood",
		"HacPai",
		"KD",
		"NGA",
		"WeiXin",
		"Mop",
		"Chiphell",
		"JianDan",
		"ChouTi",
		"ITHome",

		// 下面是需要翻墙的
		"V2EX",
		"GitHub",

		// 以下网站直接从https://www.pxlet.com/获取。
	}
	fmt.Println("开始抓取" + strconv.Itoa(len(allData)) + "种数据类型")
	group.Add(len(allData))
	var spider Spider
	for _, value := range allData {
		fmt.Println("开始抓取" + value)
		spider = Spider{DataType: value}
		// go ExecGetData(spider)
		ExecGetData(spider)
		// break
	}
	group.Wait()
	fmt.Print("完成抓取")
}
