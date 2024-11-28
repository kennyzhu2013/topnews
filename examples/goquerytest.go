package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
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

	// 设置登录cookie即可.
	request.Header.Add("Cookie", `_9755xjdesxxd_=32; _zap=dda53e17-1735-4484-92ef-615d6a2f6d4f; _xsrf=VxlKKV2n3wnGZN61TZg1VonI9pjA7kxr; d_c0=ABDaXipSihiPTsemKgU2KKePE-FbLQOSJwo=|1714379422; q_c1=0e5b8e45b9fb4510a3e545c52f8e8e51|1716973986000|1648615377000; z_c0=2|1:0|10:1729681488|4:z_c0|80:MS4xN05FQkJnQUFBQUFtQUFBQVlBSlZUY2J3OFdjRnREWFprV3dBdDAwSEgxbkw0SDEzLTdzaml3PT0=|9ed8db772357c4dc429119916b05cbcb6c6d4afd809d848ce27ed3b8c0a4e09c; Hm_lvt_98beee57fd2ef70ccdd5ca52b9740c49=1728357049,1729483043,1730248319; HMACCOUNT=2E7A3AC10D710648; __zse_ck=003_bZJen2g9eZH0BF2t7KCNLY6J0tbLG/kE7cWFHU4=Zp4Nv6PbSTgehJ2mo8i+qj/Hi1LQT3Rd8+=VyDZWt1gPZKKULUgpdgI=E6MLzO55j17Y; _ga=GA1.2.513319920.1730270985; _gid=GA1.2.1337799469.1730270985; tst=h; SESSIONID=tURUQZMbjzVyfbGZmFGcoPSo7lEfirsDqx1gY5oiJWu; JOID=VVkUAUqDYutgZcbaO4aCMOzB5t0u8Te6NjKal1rvHr0aAKKbYZy2ygJlwNk7EtE_V5YKjvsxm-JiV5tHSNdYPqE=; osd=UlgWAU-EY-lgYMHbOYaHN-3D5tgp8DW6MzWblVrqGbwYAKecYJ62zwVkwtk-FdA9V5MNj_kxnuVjVZtCT9ZaPqQ=; _ga_MKN1TB7ZML=GS1.2.1730270986.1.1.1730271130.0.0.0; BEC=244e292b1eefcef20c9b81b1d9777823; Hm_lpvt_98beee57fd2ef70ccdd5ca52b9740c49=1730282270`)
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

	println(document.Text())
	document.Find(".HotList-list .HotItem-content").Each(func(i int, selection *goquery.Selection) {
		url, boolUrl := selection.Find("a").Attr("href")
		text := selection.Find("h2").Text()
		println("GetZhiHu find text=========================================:\n")
		println(text)
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": text, "url": url})
		}
	})
	return allData
}

func main() {
	var spider Spider
	spider.GetZhiHu()
}
