package tool

// 获取灵签json
import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Lucky struct {
	Key     string   `json:"key"`
	Number  string   `json:"number"`
	Content []string `json:"content"`
}

type LinQian struct {
	Type      string `json:"type"`
	Url       string `json:"url"`
	Number    int    `json:"number"`
	FileName  string `json:"file_name"`
	Start     int    `json:"start"`      //爬虫的第一个元素位置
	StartPage int    `json:"start_page"` //爬虫开始的第一页
}

var LuckyListChan = make(chan Lucky)
var LuckyList = make([]Lucky, 0)
var lock sync.Mutex
var collyList = map[string]*colly.Collector{}
var linqian = LinQian{}

func GetJson(tt string) {
	setLinqianType(tt)

	if !FileExist(linqian.FileName) {
		os.Create(linqian.FileName)
	}

	var wg sync.WaitGroup
	i := 0
	for i = linqian.StartPage; i < linqian.Number; i++ {
		wg.Add(1)
		go func(ii int) {
			getGuanYin(ii)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func setLinqianType(tt string) {

	switch tt {

	case "观音":
		linqian.Url = "https://www.51chouqian.com/guanyinlingqian/"
		linqian.Number = 101
		linqian.FileName = "guanyinlingqian.json"
		linqian.Type = "观音"
		linqian.Start = 1
		linqian.StartPage = 1
		break

	case "佛祖":
		linqian.Url = "https://www.51chouqian.com/fozulingqian/"
		linqian.Number = 52
		linqian.FileName = "fozulingqian.json"
		linqian.Type = "佛祖"
		linqian.Start = 3
		linqian.StartPage = 0
		break

	case "月老":
		linqian.Url = "https://www.51chouqian.com/yuelaolingqian/"
		linqian.Number = 102
		linqian.FileName = "yuelaolingqian.json"
		linqian.Type = "月老"
		linqian.Start = 1
		linqian.StartPage = 0
		break
	}
}

func Do(ll Lucky, index string, i int) {
	lock.Lock()
	has := false
	for _, ccc := range LuckyList {
		if ccc.Number == ll.Number {
			has = true
			ccc = ll
			break
		}
	}
	if !has {
		LuckyList = append(LuckyList, ll)
	}
	println(len(LuckyList))
	if len(LuckyList) == linqian.Number-1 {
		//写入文件
		println("写入文件")
		f, err := os.OpenFile(linqian.FileName, os.O_WRONLY|os.O_TRUNC, 0600)

		defer f.Close()

		if err != nil {

			fmt.Println(err.Error())

		} else {

			bb, _ := json.Marshal(LuckyList)

			_, err = f.Write(bb)

			if err != nil {

				fmt.Println(err.Error())

			} else {
				println("数据写入完成")
			}
			os.Exit(0)
		}
	}
	lock.Unlock()
}
func getGuanYin(idx int) {

	index := strconv.Itoa(idx)

	collyList[index] = colly.NewCollector()

	collyList[index].OnHTML("#details_cnt", func(ee *colly.HTMLElement) {
		lucky := Lucky{}
		lucky.Key = index
		qLen := 0
		ee.ForEach("p", func(i int, element *colly.HTMLElement) {
			qLen++
		})
		ee.ForEach("p", func(i int, element *colly.HTMLElement) {
			element.Text = strings.ReplaceAll(element.Text, "\t", "")
			element.Text = strings.ReplaceAll(element.Text, "\n", "")
			if i == linqian.Start {
				lucky.Number = element.Text
			} else if i == qLen-1 {
				lucky.Content = append(lucky.Content, element.Text)
				Do(lucky, index, i)
			} else if i > linqian.Start {
				lucky.Content = append(lucky.Content, element.Text)
			}
		})
	})

	// Before making a request print "Visiting ..."

	collyList[index].OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org

	collyList[index].Visit(linqian.Url + index + ".html")
}
