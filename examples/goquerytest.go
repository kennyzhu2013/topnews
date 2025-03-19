package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Spider struct {
	DataType string
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
		fmt.Println("抓取" + spider.DataType + "失败1" + err.Error())
		return []map[string]interface{}{}
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	res, err := client.Do(request)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败2" + err.Error())
		return []map[string]interface{}{}
	}
	defer res.Body.Close()
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("抓取" + spider.DataType + "失败3" + err.Error())
		return []map[string]interface{}{}
	}
	var allData []map[string]interface{}
	document.Find(".item_title").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		println("GetV2EX find text:" + text)
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": "https://www.v2ex.com" + url})
		}
	})
	return allData
}

// GetZhiHu 知乎: 需要单独设置登录cookie.
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

	// 设置登录cookie即可， 每次都要单独更新.
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

	//println(document.Text())
	//document.Find(".HotList-list .HotItem-content").Each(func(i int, selection *goquery.Selection) {
	//	url, boolUrl := selection.Find("a").Attr("href")
	//	text := selection.Find("h2").Text()
	//	println("GetZhiHu find text=========================================:\n")
	//	println(text)
	//	if boolUrl {
	//		allData = append(allData, map[string]interface{}{"title": text, "url": url})
	//	}
	//})
	document.Contents().Each(func(i int, selection *goquery.Selection) {
		// url, boolUrl := selection.Filter("excerptArea").Attr("text")
		// text := selection.Find("h2").Text()
		// boolUrl := false
		println("GetZhiHu find text========================[" + url + "]=================:\n")
		if selection.Children() != nil {
			// json文本内容在这里打印
			selection.Children().Each(func(i int, selection *goquery.Selection) {
				println("GetZhiHu find grandson text=========================================:\n")
				findedStr := selection.Text()
				// 正则表达式查找 "hotList" 开头的子字符串
				re := regexp.MustCompile(`"hotList":\s*(\[[^]]*\])`)
				// re := regexp.MustCompile(`\{"hotList":\[(.*?)\]}`)
				matches := re.FindStringSubmatch(findedStr)
				if len(matches) > 1 {
					fmt.Println("Found hotList 1:", matches[1])

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
	return allData
}

// projectDesc contains a url, description for a project.
// 本代码单独定制的都是由于各种原因被墙掉了才单独用chromedp测试.
// 适合所用chromedp爬取的项目.
type projectDesc struct {
	URL, text string
}

// GetDouBan 豆瓣， 需要chromedp抓取，
func (spider Spider) GetDouBan() []map[string]interface{} {
	urlDouban := "https://www.douban.com/group/explore"

	// create context
	ctx, cancelChrome := chromedp.NewContext(context.Background())
	defer cancelChrome()

	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 等待topics div加载
	sel := "//*[@id=\"content\"]/div/div[1]/div[1]"
	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(urlDouban)); err != nil {
		log.Printf("could not navigate to weibo: %v\n", err)
		return []map[string]interface{}{}
	}

	// wait visible
	if err := chromedp.Run(ctx, chromedp.WaitVisible(sel), chromedp.Sleep(1*time.Second)); err != nil {
		log.Printf("could not get section: %v sel:%s\n", err, sel)
		return []map[string]interface{}{}
	}

	// 通过following-sibling函数我们可以提取指定元素目录tbody下的指定元素tr/td的所有同级元素，即获取目标元素的所有兄弟节点。
	// 同级下的ul/li节点.
	sib := sel + `//following-sibling::div/h3`

	// get h3 links and description text
	// node()匹配任意节点.
	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/a`, &linksAndDescriptions)); err != nil {
		log.Printf("could not get links and descriptions: %v\n", err)
		return []map[string]interface{}{}
	}

	// process data, add all into one struct.
	// output the values to all data
	var allData []map[string]interface{}
	for i := 0; i < len(linksAndDescriptions); i++ {
		url := linksAndDescriptions[i].AttributeValue("href")
		if len(url) == 0 {
			continue
		}
		// 直接从节点的子节点中获取文本内容
		var textContent string
		for _, child := range linksAndDescriptions[i].Children {
			textContent += child.NodeValue
		}
		fmt.Printf("Title: %s, Url: %s\n", textContent, url)
		// log.Printf("project %s (%s): '%s'", v, res[v].URL, res[v].text)
		allData = append(allData, map[string]interface{}{"title": textContent, "url": url})
	}
	return allData
}

func (spider Spider) GetBaiDu() []map[string]interface{} {
	urlBaidu := "https://top.baidu.com/board?tab=realtime"

	// create context
	ctx, cancelChrome := chromedp.NewContext(context.Background())
	defer cancelChrome()

	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 等待content_1YWBm加载,following-sibling函数和text()都只支持XPath
	sel := "//*[@id=\"sanRoot\"]/main/div[2]/div/div[2]/div/div[2]"
	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(urlBaidu)); err != nil {
		log.Printf("could not navigate to weibo: %v\n", err)
		return []map[string]interface{}{}
	}

	// wait visible
	if err := chromedp.Run(ctx, chromedp.WaitVisible(sel), chromedp.Sleep(1*time.Second)); err != nil {
		log.Printf("could not get section: %v sel:%s\n", err, sel)
		return []map[string]interface{}{}
	}

	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sel+`/a`, &linksAndDescriptions)); err != nil {
		log.Printf("could not get links and descriptions: %v\n", err)
		return []map[string]interface{}{}
	}

	// TODO: 为啥无法运行[1]/text()找到对应text？.
	// 结合/following-sibling::tr/td使用
	// /child::a/text()
	// ByQuery系列方法： all方法不支持XPath
	// //*[@id="sanRoot"]/main/div[2]/div/div[2]/div/div[2]/a/div[1]/text()
	// div.container.right-container_2EFJr .content_1YWBm a .c-single-text-ellipsis
	var divText []string
	//if err := chromedp.Run(ctx, chromedp.Nodes(`div.container.right-container_2EFJr .content_1YWBm a .c-single-text-ellipsis`, &divText, chromedp.ByQueryAll)); err != nil {
	//	log.Printf("could not get div text: %v\n", err)
	//
	//	return []map[string]interface{}{}
	//}
	if err := chromedp.Run(ctx, chromedp.Evaluate(`Array.from(document.querySelectorAll('div.container.right-container_2EFJr .content_1YWBm a .c-single-text-ellipsis')).map(el => el.textContent);`, &divText)); err != nil {
		log.Printf("could not get div text: %v\n", err)

		return []map[string]interface{}{}
	}

	// output the values to all data
	var allData []map[string]interface{}
	for i := 0; i < len(linksAndDescriptions); i++ {
		url := linksAndDescriptions[i].AttributeValue("href")
		if len(url) == 0 || i >= len(divText) {
			continue
		}
		fmt.Println(divText[i])
		fmt.Printf("Title: %s, Url: %s\n", divText[i], url)
		allData = append(allData, map[string]interface{}{"title": divText[i], "url": url})
	}
	return allData
}

// GetHuXiu TODO: 调试这里,使用json解析.
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

	// json格式，参考GetZhiHu
	fmt.Println(document.Text())
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

func main() {
	var spider Spider
	spider.GetHuXiu()
}

/*
GbkToUtf8  部分热榜标题需要转码，transform.NewReader读取的内容就是utf8.
*/
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
