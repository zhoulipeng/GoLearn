package main
import (
    "io"
    "fmt"
    "time"
    "net/http"
    "log"
)
type GuardKey struct {
    key string
    timeout int
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
    
    //guard := &GuardKey{"PrivateKey123abc", 5600 + time.Now().Unix()}
 
    io.WriteString(w, "PrivateKey123abc")
}
func main() {
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    fmt.Println("time is ", time.Now().Unix())
    http.HandleFunc("/guard", HelloServer)
    err := http.ListenAndServe("127.0.0.1:12345", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
} 
