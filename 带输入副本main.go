package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func main() {
	os.Remove("list.txt")

	var articleNum, pageNum int
	fmt.Println("请输入小说页面id:_________________________________________________________")
	fmt.Scanln(&articleNum)
	fmt.Println("请输入小说网页目录页面数量 ：________________________________________________")
	fmt.Scanln(&pageNum)

	fmt.Printf("你输入的小说id为： %d ,小说网页目录页数数量为： %d\n", articleNum, pageNum)

	var a []string
	for i := 1; i <= 9; i++ {
		url := `https://www.biqugeso.org/indexlist/` + strconv.Itoa(articleNum) + `/`
		url = url + strconv.Itoa(i)
		read(url, a)
	}

	URLAndArtileTitle(articleNum)
}
func URLAndArtileTitle(articleNum int) {
	const URLinfo = `https://www.biqugeso.org`
	n, err := ioutil.ReadFile("./list.txt")
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`<a href=".*" rel="chapter">`)
	urlResult := re.FindAllStringSubmatch(string(n), -1)

	for i := 0; i < len(urlResult); i++ {
		urlResult[i][0] = URLinfo + urlResult[i][0][9:len(urlResult[i][0])-21]
	}
	/*fmt.Println(urlResult)*/
	re = regexp.MustCompile(`<dd>.*</dd>`)
	titleResult := re.FindAllStringSubmatch(string(n), -1)

	for i := 0; i < len(titleResult); i++ {
		titleResult[i][0] = titleResult[i][0][4 : len(titleResult[i][0])-5]
	}
	fmt.Println("正在爬取文章....")

	//提取标题,,设置目录
	resp, err := http.Get(`https://www.biqugeso.org/indexlist/` + strconv.Itoa(articleNum) + `/`)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	titleURL := `my章节/`
	doc.Find(".box_con").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h1").Text()
		titleURL = titleURL + title
	})
	fmt.Println(titleURL)
	os.MkdirAll(titleURL, os.ModePerm)

	getArtileContent(urlResult, titleResult, titleURL)
}

func writeAritleList(a []string) { //写入文件
	/*fmt.Println(a)*/
	list, err := os.OpenFile("./list.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	defer list.Close()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(a); i++ {
		list.WriteString(a[i] + "\n")
	}

}
func read(url string, a []string) { //读取所有章节信息
	//获取网址
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`<a href=".*\.html" rel="chapter"><dd>第.*章.*</dd></a>`)
	result := re.FindAllStringSubmatch(string(body), -1)
	for i := 0; i < len(result); i++ {
		a = append(a, result[i][0])
	}
	writeAritleList(a)

}

func getArtileContent(urlresult, titleResult [][]string, titleURL string) {

	/*os.Mkdir("my章节", os.ModePerm)*/

	for i := 0; i < len(urlresult); i++ {
		for j := 1; j <= 2; j++ {
			resp, err := http.Get(urlresult[i][0] + `_` + strconv.Itoa(j) + `.html`)
			if err != nil {
				panic(err)
			}
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error code: %d", resp.StatusCode)
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				panic(err)
			}

			txtfile, err := os.OpenFile("./"+titleURL+`/`+titleResult[i][0]+`.txt`, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
			if err != nil {
				panic(err)
			}
			doc.Find(".box_con #content").Each(func(i int, s *goquery.Selection) {
				Content := s.Find("p").Text()
				txtfile.WriteString(Content + "\n")
			})
		}
	}
	fmt.Println("爬取文章结束...")

}
