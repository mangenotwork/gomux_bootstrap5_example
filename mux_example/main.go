package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mangenotwork/common/log"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

/*
2. 页面渲染
3. 模板参数
*/

func main() {
	r := mux.NewRouter()

	// 自定义 404 页面
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	//普通路由
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/products", ProductsHandler)

	// 参数路由
	r.HandleFunc("/articles/{title}", TitleHandler).Methods("GET")

	// 正则路由参数，下面例子中限制为英文字母
	r.HandleFunc("/articles2/{title:[a-z]+}", TitleHandler2).Methods("GET")

	// 子路由, 分组路由
	user := r.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{name}", UsersHandler)

	// 结构体
	r.Handle("/products/{id}", &ProductsIdHandler{}).Methods("GET")

	// 自定义匹配 MatcherFunc()
	r.HandleFunc("/products2/matcher", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "FormValue: %s ", r.FormValue("func"))
	}).MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
		b := false
		if req.FormValue("func") == "matcherfunc" {
			b = true
		}
		return b
	})

	// test 例子
	r.HandleFunc("/health", HealthCheckHandler)

	// 静态资源
	s := "./static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s))))

	// 路由中间件
	r.Use(loggingMiddleware)
	user.Use(userMiddleware)

	// 跨域
	user.Use(mux.CORSMethodMiddleware(user))

	// Walking Routes 遍历注册的所有路由
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	log.Info("start http server :18080")

	// 优雅地关闭
	var wait time.Duration
	srv := &http.Server{
		Addr: "0.0.0.0:18080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Error("shutting down")
	os.Exit(0)

}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello world")
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello, Products")
}

func TitleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取参数
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "title: %v\n", vars["title"])
}

func TitleHandler2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // 获取参数
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "title: %v\n", vars["title"])
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r) //获取值
	name := vars["name"]
	fmt.Fprintf(w, "name: %s \r\n", name)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Do stuff here
		log.Info(r.RequestURI)
		//fmt.Fprintf(w, "%s\r\n", r.URL)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("请求的User")
		next.ServeHTTP(w, r)
	})
}

type ProductsIdHandler struct{}

func (handler *ProductsIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "products id: %s", vars["id"])
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}
