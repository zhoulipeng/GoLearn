package main
import (
    "io"
    "fmt"
    "time"
    "net/http"
    "log"
)
import "go-simplejson"

type GuardKey struct {
    key string
    timeout int
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
    
    //guard := &GuardKey{"PrivateKey123abc", 5600 + time.Now().Unix()}
    js, err := simplejson.NewJson(
        []byte(`{"action": "on_connect","client_id": 1985,
    "ip": "192.168.1.10", "vhost": "video.test.com", "app": "live",
    "tcUrl": "rtmp://video.test.com/live?key=1444435200-0-0-80cd3862d699b7118eed99103f2a3a4f",
    "pageUrl": "http://www.test.com/live.html"}`))
    ms := js.Get("tcUrl").MustString()
    fmt.Println(ms)
    fmt.Println(err)
    io.WriteString(w, "PrivateKey123abc")
}
func main() {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    fmt.Println("time is ", time.Now().Unix())
    http.HandleFunc("/api/guard", HelloServer)
    err := http.ListenAndServe("127.0.0.1:12345", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
} 
