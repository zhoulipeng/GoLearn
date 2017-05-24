package main
import (
    "io"
    "fmt"
    "strings"
    "strconv"
    "time"
    "net/url"
    "net/http"
    "io/ioutil"
    "crypto/md5" 
    "encoding/hex"
    "log"
)
import "go-simplejson"


func hash_sum(raw string) string{
    hasher := md5.New()
    hasher.Write([]byte(raw))
    return hex.EncodeToString(hasher.Sum(nil))
}
func auth_timestamp(time_str string) bool{
    cur_gmt := time.Now().Unix()
    fmt.Println("cur_gmt is:", cur_gmt);
    req_gmt, err := strconv.ParseInt(time_str, 10, 64)

    if err != nil {
       fmt.Println("request's gmt time can't convert to int")
       return false
    }
    fmt.Println("req_gmt is:", req_gmt);
    if req_gmt < cur_gmt {
        return false
    }
    return true
}
func auth_token(u *url.URL, m url.Values) bool{
    safe_code := "##svi&&lgslb##"
    hash_raw := u.Path +
        "?user_id=" + strings.Join(m["user_id"], "") +
        "&gmt="  + strings.Join(m["gmt"], "") +
        "#" + safe_code 
    fmt.Println(hash_raw)
    hash_s := hash_sum(hash_raw)
    hash_req := strings.Join(m["token"], "")
    fmt.Println("request hash: ", hash_req)
    fmt.Println("safe hash is: ", hash_s)
    if(hash_req == hash_s){
        fmt.Println("auth success.")
        return true
    }else{
        fmt.Println("auth failed.")
        return false
    }
    return true
}
func do_auth(result []byte) bool{
    js, err := simplejson.NewJson(result);
    if err != nil {
        fmt.Println("bad json error: ", err) //raw url
        return false
    }
    tc := js.Get("tcUrl").MustString()
    fmt.Println(tc)
        
    //解析tcUrl
    u, err := url.Parse(tc)
    if err != nil {
        fmt.Println(err)
        return false
    }
    fmt.Println(u) //raw url
    m, err := url.ParseQuery(u.RawQuery)
    if err != nil {
        fmt.Println(err)
        return false
    }
    if !auth_timestamp(strings.Join(m["gmt"], "")) {
        fmt.Println("auth timestamp failed")
        return false
    }
    return auth_token(u, m)
}

func HelloServer(w http.ResponseWriter, req *http.Request) {

    req.ParseForm() //解析参数，默认是不会解析的 
    fmt.Println("method:", req.Method, "time is: ", time.Now())
    if req.Method == "POST" {
        result, _:= ioutil.ReadAll(req.Body)
        req.Body.Close()
        fmt.Printf("%s\n", result)
        if do_auth(result) {
            io.WriteString(w, "0")
            fmt.Println("method:", req.Method, "time is: ", time.Now())
            return
        }
    } else {
        fmt.Println("bad request method")
    }
    io.WriteString(w, "1")
}
func main() {
    fmt.Println("start server at : "time.Now().Format("2006-01-02 15:04:05"))
    http.HandleFunc("/api/guard", HelloServer)
    fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
    err := http.ListenAndServe(":12345", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
} 
