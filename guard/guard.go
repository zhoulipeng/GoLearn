package main
import (
    "io"
    "fmt"
    "strings"
    "time"
    "net/http"
    "io/ioutil"
    "log"
)
import "go-simplejson"

type GuardKey struct {
    key string
    timeout int
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
    
    //guard := &GuardKey{"PrivateKey123abc", 5600 + time.Now().Unix()}
    req.ParseForm() //解析参数，默认是不会解析的 
    //fmt.Fprintf(w, "Hi, I love you %s", html.EscapeString(req.URL.Path[1:]))
    fmt.Println("method:", req.Method) //获取请求的方法 
    if req.Method == "GET" {

        fmt.Println("username", req.Form["username"]) 
        fmt.Println("password", req.Form["password"]) 

        for k, v := range req.Form {
            fmt.Print("key:", k, "; ")
            fmt.Println("val:", strings.Join(v, ""))
        }
        io.WriteString(w, "1")
    } else if req.Method == "POST" {
        result, _:= ioutil.ReadAll(req.Body)
        req.Body.Close()
        fmt.Printf("%s\n", result)
        js, err := simplejson.NewJson(result);
        ms := js.Get("tcUrl").MustString()
        fmt.Println(ms)
        fmt.Println(err)
        io.WriteString(w, "0")
        return
    }
    /*
    js, err := simplejson.NewJson(
        []byte(`{"action": "on_connect","client_id": 1985,
    "ip": "192.168.1.10", "vhost": "video.test.com", "app": "live",
    "tcUrl": "rtmp://video.test.com/live?key=1444435200-0-0-80cd3862d699b7118eed99103f2a3a4f",
    "pageUrl": "http://www.test.com/live.html"}`))
    ms := js.Get("tcUrl").MustString()
    fmt.Println(ms)
    fmt.Println(err)
    */
    io.WriteString(w, "1")
}
func main() {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    fmt.Println("time is ", time.Now().Unix())
    http.HandleFunc("/api/guard", HelloServer)
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
} 
