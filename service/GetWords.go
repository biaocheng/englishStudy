package service

import (
	"goquery"
	"github.com/astaxie/beego"
	"github.com/biaocheng/englishStudy/models"
	"strings"
	"fmt"
	"path"
	"os"
	"net/http"
	"io"
	"github.com/astaxie/beego/logs"
	"regexp"
)

func GetWordS()  {
	doc,err:= goquery.NewDocument("http://word.iciba.com/?action=index&reselect=y")
	if err!=nil{
		beego.Error(err)
	}
	doc.Find(".c_nav").Each(func(i int, selection *goquery.Selection) {
		topWordType := new(models.WordType)
		topWordType.Name = strings.TrimSpace(selection.Find("p").Text())

		logs.Debug("============one==============")
		models.InsertWordType(topWordType)

		childClassId,_:= selection.Attr("cid")
		doc.Find(".main_l").Each(func(i int, selection2 *goquery.Selection) {

			cFilter,_ := selection2.Attr("c_filter")


			if cFilter==childClassId{

				selection2.Find(".cl>li").Each(func(i int, selection3 *goquery.Selection) {
					twoWordType := new(models.WordType)
					twoWordType.Parent = topWordType
					twoWordType.Name = strings.TrimSpace(selection3.Find("h3").Text())
					models.InsertWordType(twoWordType)

					if selection3.Find(".main_l_box").Size()==0{
						classId,_ := selection3.Attr("class_id")
						getUnitPage(twoWordType,classId,"http://word.iciba.com/?action=courses&classid="+classId)
						return
					}

					selection3.Find(".main_l_box .nobt li").Each(func(i int, selection4 *goquery.Selection) {
						threeWordType := new(models.WordType)
						threeWordType.Name = strings.TrimSpace(selection4.Find("a h4").Text())
						threeWordType.Parent = twoWordType

						logs.Debug("============three==============")
						models.InsertWordType(threeWordType)

						classId,_ := selection4.Attr("class_id")
						classId = strings.TrimSpace(classId)
						a := selection4.Find("a")
						url,_ := a.Attr("href")
						url = "http://word.iciba.com"+url

						getUnitPage(threeWordType,classId,url)

					})

				})
			}
		})
	})
}

func getUnitPage(typeWord *models.WordType,classId string,url string){
	//logs.Debug(url)
	doc2 ,err := goquery.NewDocument(url)
	if err!=nil{
		beego.Error(err)
	}
	doc2.Find(".mid li").Each(func(i int, selection5 *goquery.Selection) {
		fourWordType := new(models.WordType)
		fourWordType.Name = strings.TrimSpace(selection5.Find("h4").Text())
		fourWordType.Parent = typeWord

		logs.Debug("============four==============")
		models.InsertWordType(fourWordType)

		wordsUrlId,_ := selection5.Attr("course_id")
		wordsUrlId = strings.TrimSpace(wordsUrlId)
		getWords(fourWordType,classId,wordsUrlId)

	})
}

func getWords(worldType *models.WordType,classId string,wordsUrlId string){

	wordsUrl := fmt.Sprintf("http://word.iciba.com/?action=words&class=%s&course=%s",classId,wordsUrlId)
	doc3,err := goquery.NewDocument(wordsUrl)
	if err!=nil{
		beego.Error(err)
	}
	doc3.Find(".word_main_list").Each(func(i int, selection6 *goquery.Selection) {
		selection6.Find("li").Each(func(i int, selection7 *goquery.Selection) {
			word := new(models.Word)
			defer func(){
				models.InsertWord(word)

				if err:=recover();err!=nil{
					logs.Error(err) // 这里的err其实就是panic传入的内容，55
				}

			}()
			//获得单词
			wordName := strings.TrimSpace(selection7.Find(".word_main_list_w span").Text())
			//获得音标
			yinbiao := strings.TrimSpace(selection7.Find(".word_main_list_y strong").Text())
			//获得释义
			shiyi,exists := selection7.Find(".word_main_list_s span").Attr("title")
			if exists{
				shiyi = strings.TrimSpace(shiyi)
			}else {
				shiyi = strings.TrimSpace(selection7.Find(".word_main_list_s span").Text())
			}
			wordMusicUrl,exists:= selection7.Find(".word_main_list_y a").Attr("id")

			word.Word = wordName
			word.Phonetic = yinbiao
			word.Interpretation = shiyi
			word.WordType = worldType

			reg,err := regexp.Compile(`https?://.*\.mp3`)

			if err!=nil{
				logs.Error(err)
			}
			//判断是否存在读音
			if exists && reg.MatchString(wordMusicUrl){
				wordMusicUrl = strings.TrimSpace(wordMusicUrl)

				p := path.Join("d:/music",getPath(worldType),wordName+".mp3")
				os.MkdirAll(path.Dir(p),777)
				file,err := os.Create(p)
				if err!=nil{
					beego.Error(err)
				}

				resp,_ := http.Get(wordMusicUrl)

				data := make([]byte,1024)
				for{
					len,err := resp.Body.Read(data)
					if len<=0{
						break
					}
					if err!=nil && err!=io.EOF{
						beego.Error(err)
					}
					file.Write(data[0:len])
				}
				resp.Body.Close()
				file.Close()
				word.MusicUrl = path.Join(getPath(worldType),wordName+".mp3")
			}

			logs.Debug("============word==============")

		})
	})
}

func getPath(wordType *models.WordType) string{
	if wordType.Parent!=nil{
		return path.Join(getPath(wordType.Parent),wordType.Name)
	}else{
		return wordType.Name
	}
}
