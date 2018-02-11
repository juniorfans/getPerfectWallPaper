package main

import (
    "fmt"
    "strings"
  //  "strconv"
    "net/http"
    "net/url"
    "io/ioutil"
    "os"
    //"log"
  //  "runtime"
  //  "flag"
  //  "github.com/PuerkitoBio/goquery"
    "runtime"
)

var (
    descript_file string
    img_dir string
    base_url string
    signal chan bool
)

func main(){
    runtime.GOMAXPROCS(8)

    img_dir="D:/pic/"
    fmt.Println("save to ", img_dir)
    base_url = "https://d3cbihxaqsuq0s.cloudfront.net/"
    descript_file="E:/Go/src/crawl/data/descriptor.xml"
    var descripor string = readFileToString(descript_file)
    var keys []string = strings.Split(descripor,"<Key>")

    urlSlice := make([]string, 0)

    for _, v := range keys {
        innerKeys := strings.Split(v,"</Key>")
        if(2 == len(innerKeys)){
            image := innerKeys[0]
            if(len(image) > len("images/")){
                //fmt.Println(image)
                urlSlice=append(urlSlice, image)
            }
        }
    }

    lastOne := false
    for i,v := range  urlSlice {

        lastOne = (i == len(urlSlice)-1);
        if(lastOne){
            fmt.Println("LAST ONE")
        }

        go saveImages(base_url + v, lastOne)
    }

    <-signal
}

func readFileToString(filename string) string{
    dat, err := ioutil.ReadFile(filename)
    if err != nil {
        panic(err)
    }
    return string(dat);
}

//下载图片
func saveImages(img_url string, lastone bool){
    fmt.Println(img_url)
    u, err := url.Parse(img_url)
    if err != nil {
        fmt.Println("parse url failed:", img_url, err)
        return 
    }

    fmt.Println("check")
    //去掉最左边的'/'
    tmp := strings.TrimLeft(u.Path, "/")
    filename := img_dir + strings.ToLower(strings.Replace(tmp, "/", "-", -1))

    fmt.Println("filename: ", filename)

    exists := isExists(filename)
    if exists {
        fmt.Println("Oops, has been exsits")
        //return
    }
	//set your own proxy if it must be so
	proxy := func(_ *http.Request) (*url.URL, error) {
            //注意下面的 url 一定要包含 http:// 否则运行报错
        return url.Parse("http://proxy.xxx.com:8080")	//your proxy addr
        }

    fmt.Println("set proxy")
    transport := &http.Transport{Proxy: proxy}

    client := &http.Client{Transport: transport}
    response, err := client.Get(img_url)

    fmt.Println("Get finished")
    //response, err := http.Get(img_url)
    if err != nil {
        fmt.Println("get img_url failed:", err)
        return 
    }

    defer response.Body.Close()

    fmt.Println("read")

    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        fmt.Println("read data failed:", img_url, err)
        return 
    }

    var a []byte  = data;
    fmt.Println("img data size: ", len(a))

    image, err := os.Create(filename)
    if err != nil {
        fmt.Println("create file failed:", filename, err)
        return 
    }

    defer image.Close()
    image.Write(data)
    if(lastone){
        signal <- true
    }
}

func isExists(filename string) bool {
    _, err := os.Stat(filename)
    return err != nil
}