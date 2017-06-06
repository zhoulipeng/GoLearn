package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)
import "github.com/bitly/go-simplejson"
import "github.com/takama/daemon"

const (

	// name of the service
	name        = "myservice"
	description = "vlss_guard"

	// port which daemon should be listen
	port = ":9977"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

func startHttpServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world\n")
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic,
			// because this
			// probably is an
			// intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	// returning
	// reference
	// so
	// caller
	// can
	// call
	// Shutdown()
	return srv
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	/*fmt.Println("start server at : " + time.Now().Format("2006-01-02 15:04:05"))
	go func() {
		http.HandleFunc("/api/guard", HelloServer)
		errlog.Println("set handler ok!")
		err := http.ListenAndServe("127.0.0.1:12345", nil)
		if err != nil {
			//return "ListenAndServe: failed ", err
		}
	}()*/

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)
	
	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
		go handleClient(conn)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ")
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
		}
	}
}


// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}

func hash_sum(raw string) string {
	hasher := md5.New()
	hasher.Write([]byte(raw))
	return hex.EncodeToString(hasher.Sum(nil))
}

func auth_timestamp(time_str string) bool {
	cur_gmt := time.Now().UTC().Unix()
	fmt.Println("cur_gmt is:", cur_gmt)
	req_gmt, err := strconv.ParseInt(time_str, 10, 64)

	if err != nil {
		fmt.Println("request's gmt time can't convert to int")
		return false
	}
	fmt.Println("req_gmt is:", req_gmt)
	if req_gmt < cur_gmt {
		return false
	}
	return true
}

func auth_token(u *url.URL, m url.Values) bool {
	safe_code := "##svi&&lgslb##"
	hash_raw := u.Path +
		"?user_id=" + strings.Join(m["user_id"], "") +
		"&gmt=" + strings.Join(m["gmt"], "") +
		"#" + safe_code
	fmt.Println(hash_raw)
	hash_s := hash_sum(hash_raw)
	hash_req := strings.Join(m["token"], "")
	fmt.Println("request hash: ", hash_req)
	fmt.Println("safe hash is: ", hash_s)
	return hash_req == hash_s
}

func do_auth(result []byte) bool {
	js, err := simplejson.NewJson(result)
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
		result, _ := ioutil.ReadAll(req.Body)
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
func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

/*func main() {
	fmt.Println("start server at : " + time.Now().Format("2006-01-02 15:04:05"))
	http.HandleFunc("/api/guard", HelloServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}*/
func main() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error1: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError2: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
