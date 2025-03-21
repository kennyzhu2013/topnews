package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// projectDesc contains a url, description for a project.
// 本代码单独定制的都是由于各种原因被墙掉了才单独用chromedp测试.
type projectDesc struct {
	URL, text string
}
type ByName []string

func (a ByName) Len() int { return len(a) }
func (a ByName) Less(i, j int) bool {
	indexa, erra := strconv.Atoi(a[i])
	if erra != nil {
		return true
	}

	indexb, errb := strconv.Atoi(a[j])
	if errb != nil {
		return false
	}

	return indexa < indexb
}
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// GetWeiBo 微博: 单独定制.
func (spider Spider) GetWeiBo() []map[string]interface{} {
	// create context
	ctx, cancelChrome := chromedp.NewContext(context.Background())
	defer cancelChrome()

	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 等待热点第十个加载
	sel := "//*[@id=\"pl_top_realtimehot\"]/table/tbody/tr[10]/td[2]/a"

	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(`https://s.weibo.com/top/summary`)); err != nil {
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
	sib := "//*[@id=\"pl_top_realtimehot\"]/table/tbody/tr" + `/following-sibling::tr/td`

	// get project link text
	var texts []*cdp.Node
	// child子节点下的a节点.
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/child::a/text()`, &texts)); err != nil {
		log.Printf("could not get projects: %v\n", err)
		return []map[string]interface{}{}
	}

	// child子节点下的内容. ranktop.//*[@id="pl_top_realtimehot"]/table/tbody/tr[2]/td[1]
	var textLocks []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`[1]/text()`, &textLocks)); err != nil {
		log.Printf("could not get projects: %v\n", err)
		return []map[string]interface{}{}
	}

	// get links and description text
	// node()匹配任意节点.
	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sib+`/a`, &linksAndDescriptions)); err != nil {
		log.Printf("could not get links and descriptions: %v\n", err)
		return []map[string]interface{}{}
	}

	// check length
	if len(texts) != len(linksAndDescriptions) || len(texts) != len(textLocks) {
		log.Printf("projects and links and descriptions lengths do not match (%d != %d)\n", len(texts), len(linksAndDescriptions))
		return []map[string]interface{}{}
	}

	// process data, add all into one struct.
	res := make(map[string]projectDesc)
	for i := 0; i < len(textLocks); i++ {
		url := linksAndDescriptions[i].AttributeValue("href")
		urlGuanxuan := linksAndDescriptions[i].AttributeValue("href_to")
		if linksAndDescriptions[i].AttributeValue("href_to") != "" {
			url = urlGuanxuan
		}
		res[textLocks[i].NodeValue] = projectDesc{
			URL:  url,
			text: texts[i].NodeValue,
		}
	}

	// 提取键
	keys := make(ByName, 0, len(res))
	for k, _ := range res {
		keys = append(keys, k)
	}

	// 排序键
	sort.Sort(keys)

	// output the values to all data
	var allData []map[string]interface{}
	for _, v := range keys {
		// log.Printf("project %s (%s): '%s'", v, res[v].URL, res[v].text)
		if len(res[v].URL) > 0 {
			allData = append(allData, map[string]interface{}{"title": res[v].text, "url": "https://s.weibo.com" + res[v].URL})
		}
	}

	return allData
}

// GetDouBan 豆瓣， 需要chromedp抓取，单独定制，
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

// GetHuPu 虎扑: 使用chromedp
func (spider Spider) GetHuPu() []map[string]interface{} {
	baseUrl := "https://bbs.hupu.com/"
	urlHupu := baseUrl + "all-gambia"

	// create context
	ctx, cancelChrome := chromedp.NewContext(context.Background())
	defer cancelChrome()

	// force max timeout of 15 seconds for retrieving and processing the data
	var cancel func()
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 等待topics div加载
	sel := "div.text-list-model  .t-info"
	// navigate
	if err := chromedp.Run(ctx, chromedp.Navigate(urlHupu)); err != nil {
		log.Printf("could not navigate to weibo: %v\n", err)
		return []map[string]interface{}{}
	}

	// wait visible
	if err := chromedp.Run(ctx, chromedp.WaitVisible(sel), chromedp.Sleep(1*time.Second)); err != nil {
		log.Printf("could not get section: %v sel:%s\n", err, sel)
		return []map[string]interface{}{}
	}

	var linksAndDescriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sel+` a`, &linksAndDescriptions)); err != nil {
		log.Printf("could not get links and descriptions: %v\n", err)
		return []map[string]interface{}{}
	}

	var spanText []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(sel+` a>span`, &spanText)); err != nil {
		log.Printf("could not get span text: %v\n", err)
		return []map[string]interface{}{}
	}

	// output the values to all data
	var allData []map[string]interface{}
	for i := 0; i < len(linksAndDescriptions); i++ {
		url := linksAndDescriptions[i].AttributeValue("href")
		if len(url) == 0 || i >= len(spanText) {
			continue
		}

		// 直接从节点的子节点中获取文本内容
		// 这里不知道为啥span没有当做a的节点.无法直接取bug?
		//var textContent string
		//for _, child := range linksAndDescriptions[i].Children {
		//	fmt.Println("================")
		//	textContent += child.NodeValue
		//}
		// textContent := spanText[i].NodeValue
		var textContent string
		for _, child := range spanText[i].Children {
			textContent += child.NodeValue
		}

		allData = append(allData, map[string]interface{}{"title": textContent, "url": baseUrl + url})
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
	var divText []string
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
		//fmt.Println(divText[i])
		//fmt.Printf("Title: %s, Url: %s\n", divText[i], url)
		allData = append(allData, map[string]interface{}{"title": divText[i], "url": url})
	}
	return allData
}
