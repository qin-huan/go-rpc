package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"github.com/yanyiwu/gojieba"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

type JsonData struct {
	User  User    `json:"user"`
	Weibo []Weibo `json:"weibo"`
}
type User struct {
	ID             string `json:"id"`
	Nickname       string `json:"nickname"`
	Gender         string `json:"gender"`
	Location       string `json:"location"`
	Birthday       string `json:"birthday"`
	Description    string `json:"description"`
	VerifiedReason string `json:"verified_reason"`
	Talent         string `json:"talent"`
	Education      string `json:"education"`
	Work           string `json:"work"`
	WeiboNum       int    `json:"weibo_num"`
	Following      int    `json:"following"`
	Followers      int    `json:"followers"`
}
type Weibo struct {
	ID               string      `json:"id"`
	UserID           string      `json:"user_id"`
	Content          string      `json:"content"`
	ArticleURL       string      `json:"article_url"`
	OriginalPictures string      `json:"original_pictures"`
	RetweetPictures  interface{} `json:"retweet_pictures"`
	Original         bool        `json:"original"`
	VideoURL         string      `json:"video_url"`
	PublishPlace     string      `json:"publish_place"`
	PublishTime      string      `json:"publish_time"`
	PublishTool      string      `json:"publish_tool"`
	UpNum            int         `json:"up_num"`
	RetweetNum       int         `json:"retweet_num"`
	CommentNum       int         `json:"comment_num"`
}

var keys = []string{
	"冠状病毒",
	"肺炎",
	"疫情",
	"防疫",
	"新冠",
}

func main() {
	//if err := readJson(); err != nil {
	//	panic(err)
	//}

	if err := readCsv(); err != nil {
		panic(err)
	}
}

func readJson() error {
	dirPath := "./json/"
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	xlsx := excelize.NewFile()
	xlsx_key := excelize.NewFile()
	for _, d := range dir {
		now := time.Now()
		if d.IsDir() {
			continue
		}

		var (
			filepath  = dirPath + d.Name()
			data      = &JsonData{}
			yqNum     int
			dateSum   = make(map[string]int)
			yqDateNum = make(map[string]int)
		)
		file, err := os.Open(filepath)
		if err != nil {
			continue
		}
		fmt.Println("filepath: ", filepath)

		encoder := json.NewDecoder(file)
		if err = encoder.Decode(data); err != nil {
			continue
		}

		// 自然语言处理
		fmt.Println("nlp start...", d.Name())
		nlp(data, d.Name(), xlsx_key)
		fmt.Println("nlp done...", d.Name())

		for _, weibo := range data.Weibo {
			str := strings.Split(weibo.PublishTime, " ")
			date := str[0]

			_, ok := dateSum[date]
			if !ok {
				dateSum[date] = 1
			} else {
				dateSum[date]++
			}

			for _, key := range keys {
				if strings.Contains(weibo.Content, key) {
					yqNum++
					_, ok := yqDateNum[date]
					if !ok {
						yqDateNum[date] = 1
					} else {
						yqDateNum[date]++
					}
					break
				}
			}
		}

		_ = file.Close()

		str := strings.Split(d.Name(), ".")
		var sheetName = str[0]
		sheet := xlsx.NewSheet(sheetName)
		if err = xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{
			"总微博数", "与疫情相关微博数量",
		}); err != nil {
			return err
		}
		if err = xlsx.SetSheetRow(sheetName, "A2", &[]interface{}{
			len(data.Weibo), yqNum,
		}); err != nil {
			return err
		}
		if err = xlsx.SetSheetRow(sheetName, "A4", &[]interface{}{
			"日期", "微博数", "与疫情相关微博数",
		}); err != nil {
			return err
		}

		var index = 5
		for key, val := range dateSum {
			var yq int
			if v, ok := yqDateNum[key]; ok {
				yq = v
			}
			if err := xlsx.SetSheetRow(sheetName, fmt.Sprintf("%s%d", "A", index), &[]interface{}{
				key, val, yq,
			}); err != nil {
				continue
			}
			index++
		}
		xlsx.SetActiveSheet(sheet)
		fmt.Println("has run: ", time.Since(now))
	}

	fmt.Println("save file...")
	if err = xlsx.SaveAs("./xlsx/jiangsu.xlsx"); err != nil {
		return err
	}
	if err = xlsx_key.SaveAs("./xlsx/jiangsu_key.xlsx"); err != nil {
		return err
	}
	return nil
}

var stopWords = make(map[string]bool)
func init() {
	file, err := os.Open("./txt/stop_words.txt")
	if err != nil {
		panic(err)
	}

	all, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	strs := strings.Split(string(all), "\n")
	for _, str := range strs {
		stopWords[strings.TrimSpace(str)] = true
	}
	stopWords["原图"] = true
}
var signal = []string{
	"【","】","{","}","，","。","？","；","：","\"","”","[0-9a-zA-Z]","...","!","！","、","（","）","》","《",
}

func nlp(data *JsonData, name string, xlsx *excelize.File) {
	var wordsMap = make(map[string]int)
	var strChan = make(chan[]string, 0)
	wg := &errgroup.Group{}

	wg.Go(func() error {
		var index int
		foreach:
		for {
			select {
			case strs, ok := <- strChan:
				if !ok {
					break foreach
				}
				index++
				fmt.Printf("%d:%d\t", len(strs), index)
				for _, s := range strs {
					if len([]rune(s)) <= 1 {
						continue
					}
					if _, ok := stopWords[s]; ok {
						continue
					}
					if isMatch, _ := regexp.MatchString("[0-9a-zA-Z]+", s); isMatch {
						continue
					}
					if _, ok := wordsMap[s]; ok {
						wordsMap[s]++
					} else {
						wordsMap[s] = 1
					}
				}
			}
		}
		return nil
	})

	for _, weibo := range data.Weibo {
		var (
			content = weibo.Content
			isContained = false
		)
		for _, s := range signal {
			content = strings.ReplaceAll(content, s, "")
		}

		for _, key := range keys {
			if strings.Contains(content, key) {
				isContained = true
				break
			}
		}

		if isContained {
			//wg.Go(func() error {
				strs := jieba(content)
				strChan <- strs
			//for _, s := range strs {
			//	if len([]rune(s)) <= 1 {
			//		continue
			//	}
			//	if _, ok := stopWords[s]; ok {
			//		continue
			//	}
			//	if isMatch, _ := regexp.MatchString("[0-9a-zA-Z]+", s); isMatch {
			//		continue
			//	}
			//	if _, ok := wordsMap[s]; ok {
			//		wordsMap[s]++
			//	} else {
			//		wordsMap[s] = 1
			//	}
			//}
			//	return nil
			//})
		}
	}

	close(strChan)
	if err := wg.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("len(wordsMap): ", len(wordsMap))

	split := strings.Split(name, ".")
	sheetName := split[0]
	sheet := xlsx.NewSheet(sheetName)

	if err := xlsx.SetSheetRow(sheetName, "A1", &[]interface{}{
		"key", "num",
	}); err != nil {
		fmt.Println("SetSheetRow error: ", err)
		panic(err)
	}

	var index = 2
	for key, val := range wordsMap {
		if err := xlsx.SetSheetRow(sheetName, fmt.Sprintf("A%d", index), &[]interface{}{
			key, val,
		}); err != nil {
			fmt.Println(err)
			continue
		}
		index++
	}
	xlsx.SetActiveSheet(sheet)
	fmt.Println(sheetName+" done")
}

func jieba(str string) []string {
	x := gojieba.NewJieba()
	defer x.Free()
	return x.Cut(str, true)
}

func readCsv() error {
	dirPath := "./csv/"
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, d := range dir {
		if d.IsDir() {
			continue
		}

		var filepath  = dirPath + d.Name()
		file, err := os.Open(filepath)
		if err != nil {
			return err
		}

		all, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Printf("err: %v, filepath: %v", err, filepath)
			continue
		}

		newCsv, err := os.Create("./newCsv/" + d.Name())
		if err != nil {
			fmt.Printf("err: %v, filepath: %s\n", err, filepath)
			continue
		}
		writer := csv.NewWriter(newCsv)
		lines := strings.Split(string(all), "\n")
		for index, l := range lines {
			str := strings.Split(strings.TrimSpace(l), ",")
			if index == 0 {
				if err := writer.Write([]string{str[1]}); err != nil {
					return err
				}
				continue
			}
			if len(str) <= 1 {
				continue
			}
			content := str[1]
			if content == "微博正文" {
				continue
			}

			for _, key := range keys {
				if strings.Contains(content, key) {
					if err := writer.Write([]string{content}); err != nil {
						fmt.Printf("err: %v, filepath: %s\n", err, filepath)
					}
					break
				}
			}
		}

		writer.Flush()
		newCsv.Close()
	}
	return nil
}

func readTxt() {
	file, err := os.Open("./sj/txt/2539961154.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	data := string(all)
	lines := strings.Split(data, "\n")
	var (
		context  = make([]string, 0)
		date     = make([]string, 0)
		isNext   = true
		dateLine int
	)

	for index, line := range lines {
		if index < 8 {
			continue
		}
		if strings.TrimSpace(line) == "" {
			isNext = true
			continue
		}
		if isNext {
			dateLine = index + 2
			context = append(context, strings.TrimSpace(line))
			isNext = false
		}
		if index == dateLine {
			date = append(date, strings.TrimSpace(line))
		}
	}

	var iData = make([]string, 0)
	for _, line := range context {
		for _, key := range keys {
			if strings.Contains(line, key) {
				iData = append(iData, key)
				break
			}
		}
	}
	fmt.Printf("len(context): %d, len(iData): %d", len(context), len(iData))
}
