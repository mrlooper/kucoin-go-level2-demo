package web

import (
	"net/http"
	"os/exec"
	"runtime"
	"time"

	order_book "github.com/mrlooper/kucoin-go-level2-demo"
)

type Router struct {
	port        string
	l2OrderBook *order_book.Level2OrderBook
}

func NewRouter(port string, l2OrderBook *order_book.Level2OrderBook) *Router {
	return &Router{
		port:        port,
		l2OrderBook: l2OrderBook,
	}
}

func (router *Router) index(w http.ResponseWriter, r *http.Request) {
	data, err := router.l2OrderBook.SnapshotBytes()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(data)
	return
}

func (router *Router) Handle() {
	http.HandleFunc("/", router.index)

	if runtime.GOOS == "darwin" {
		go func() {
			time.Sleep(time.Second)
			exec.Command("open", "http://localhost:"+router.port).Run()
		}()
	}

	if err := http.ListenAndServe(":"+router.port, nil); err != nil {
		panic(err)
	}
}
