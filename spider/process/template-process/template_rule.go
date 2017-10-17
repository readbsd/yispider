package template_process

import (
	"github.com/PuerkitoBio/goquery"
	"bytes"
	"strings"
	"YiSpider/spider/logger"
	"YiSpider/spider/process/filter"
	"YiSpider/spider/model"
	url2 "net/url"
	"fmt"
	"net/http"
	"io/ioutil"
	"YiSpider/spider/common"
	"encoding/json"
)

func TemplateRuleProcess(process *model.Process,response *http.Response) (*model.Page,error){
	htmlBytes,err := ioutil.ReadAll(response.Body)
	if err != nil{
		return nil,err
	}
	defer response.Body.Close()

	rule := process.TemplateRule.Rule

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(htmlBytes))
	if err != nil {
		logger.Error("NewDocumentFromReader fail,",err)
		return nil,err
	}

	urls := []*model.Request{}

	if len(process.RegUrl) > 0{
		doc.Find("a").Each(func(i int, sel *goquery.Selection) {
			href, _ := sel.Attr("href")
			href = getComplateUrl(response.Request.URL, href)
			if filter.Filter(href, process) {
				urls = append(urls, &model.Request{Url:href,Method:"get",})
			}
		})
	}

	resultType := "map"
	rootSel := ""
	page := &model.Page{Urls:urls}

	v,ok := rule["node"]
	if ok{
		contentInfo := strings.Split(v,"|")
		resultType = contentInfo[0]
		rootSel = contentInfo[1]
	}

	if resultType == "array"{
		result := []map[string]interface{}{}

		doc.Find(rootSel).Each(func(i int, s *goquery.Selection) {
			data := getMapFromDom(rule,s)
			if data == nil{
				return
			}
			if len(process.AddQueue) > 0{
					page.Urls = append(page.Urls,common.PraseReq(process.AddQueue,data)...)
			}
			result = append(result,data)
		})
		page.Result = result
		page.ResultCount = len(result)
	}

	if resultType == "map"{
		data := getMapFromDom(rule,doc.Selection)
		if len(process.AddQueue) > 0{
			page.Urls = append(page.Urls,common.PraseReq(process.AddQueue,data)...)
		}
		page.Result = data
		page.ResultCount = 1
	}

	return page,nil
}

func getMapFromDom(rule map[string]string,node *goquery.Selection) map[string]interface{}{


	result := make(map[string]interface{})

	isNull := true

	for key,value := range rule{

		if key == "node"{
			continue
		}

		rules := strings.Split(value,"|")
		ValueType := strings.Split(rules[0],".")

		if len(rules) < 2{
			continue
		}

		s := node.Find(rules[1])
		switch ValueType[0] {
			case "text":
				result[key] = s.Text()
			case "html":
				result[key],_ = s.Html()
			case "attr":
				if len(ValueType) < 2{
					continue
				}
				result[key],_ = s.Attr(ValueType[1])
		    case "texts":
				arr := []string{}
				s.Each(func(i int,sel *goquery.Selection){
					text := sel.Text()
					arr = append(arr,text)
				})
				j,_ := json.Marshal(arr)
				result[key] = string(j)
			case "htmls":
				arr := []string{}
				s.Each(func(i int,sel *goquery.Selection){
					html,_ := s.Html()
					arr = append(arr,html)
				})
				j,_ := json.Marshal(arr)
				result[key] = string(j)
			case "attrs":
				arr := []string{}
				attr := ""
				s.Each(func(i int,sel *goquery.Selection){
					if len(ValueType) >= 2{
						attr,_ = sel.Attr(ValueType[1])
						arr = append(arr,attr)
					}
				})
				//j,_ := json.Marshal(arr)
				result[key] = arr
			default:
				result[key] = ""
		}

		if len(result[key].(string)) != 0{
			isNull = false
		}
	}

	if isNull == true{
		return nil
	}

	return result
}

func getComplateUrl(url *url2.URL,href string) string{

	if strings.HasPrefix(href,"/"){
		newHref := fmt.Sprintf("%s://%s%s",url.Scheme,url.Host,href)
		return newHref
	}

	newHref := fmt.Sprintf("%s://%s/%s",url.Scheme,url.Host,href)
	return newHref
}
