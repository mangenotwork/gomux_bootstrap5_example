package main

import (
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/mangenotwork/common/log"
	"path/filepath"
)

const maxSize = 8 * iris.MB

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	// 设置recover从panics恢复，设置log记录
	app.Use(recover.New())
	app.Use(logger.New())

	app.Handle("GET", "/hello", func(ctx iris.Context) {
		ctx.HTML("<h1>Hello Iris!</h1>")

	})
	app.Handle("GET", "/getjson", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "your msg"})
	})

	// 注册模板

	// Parse all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.
	tmpl := iris.HTML("../views/iris", ".html")
	// Set custom delimeters.
	tmpl.Delims("{{", "}}")
	// Enable re-build on local template files changes.
	tmpl.Reload(true)

	// Default template funcs are:
	//
	// - {{ urlpath "myNamedRoute" "pathParameter_ifNeeded" }}
	// - {{ render "header.html" . }}
	// and partial relative path to current page:
	// - {{ render_r "header.html" . }}
	// - {{ yield . }}
	// - {{ current }}
	// Register a custom template func:
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})

	// Register the view engine to the views,
	// this will load the templates.
	app.RegisterView(tmpl)

	app.Get("/", func(ctx iris.Context) {
		// 设置模板中"message"的参数值
		ctx.ViewData("message", "Hello world!")
		// 加载模板
		ctx.View("hello.html")
	}).Use(myMiddleware) // 中间件 myMiddleware

	// 参数获取
	app.Post("/login", func(ctx iris.Context) {
		username := ctx.FormValue("username")
		password := ctx.FormValue("password")
		ctx.JSON(iris.Map{
			"Username": username,
			"Password": password,
		})
	})

	// 分组
	booksAPI := app.Party("/books")
	{
		// GET: http://localhost:8080/books
		booksAPI.Get("/", list)
	}

	// 上传文件
	app.Post("/upload", func(ctx iris.Context) {
		// Set a lower memory limit for multipart forms (default is 32 MiB)
		ctx.SetMaxRequestBodySize(maxSize)
		// OR
		// app.Use(iris.LimitRequestBodySize(maxSize))
		// OR
		// OR iris.WithPostMaxMemory(maxSize)
		// single file
		_, fileHeader, err := ctx.FormFile("file")
		if err != nil {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}
		// Upload the file to specific destination.
		dest := filepath.Join("./uploads", fileHeader.Filename)
		ctx.SaveFormFile(fileHeader, dest)
		ctx.Writef("File: %s uploaded!", fileHeader.Filename)
	})

	// 上传多个文件
	app.Post("/uploads", func(ctx iris.Context) {
		files, n, err := ctx.UploadFormFiles("./uploads")
		if err != nil {
			ctx.StopWithStatus(iris.StatusInternalServerError)
			return
		}
		ctx.Writef("%d files of %d total size uploaded!", len(files), n)
	})

	// 绑定网址路径参数
	app.Get("/{name}/{age:int}/{tail:path}", func(ctx iris.Context) {
		var p myParams
		if err := ctx.ReadParams(&p); err != nil {
			ctx.StopWithError(iris.StatusInternalServerError, err)
			return
		}
		_, _ = ctx.Writef("myParams: %#v", p)
	})

	// 绑定标头
	app.Get("/head", func(ctx iris.Context) {
		var hs myHeaders
		if err := ctx.ReadHeaders(&hs); err != nil {
			ctx.StopWithError(iris.StatusInternalServerError, err)
			return
		}
		ctx.JSON(hs)
	})

	// 重 定向
	app.Get("/a", func(ctx iris.Context) {
		ctx.Redirect("/b", iris.StatusFound)
	})

	// 从处理程序发出本地路由器重定向
	app.Get("/test", func(ctx iris.Context) {
		r := ctx.Request()
		r.URL.Path = "/test2"
		ctx.Application().ServeHTTPC(ctx)
		// OR
		// ctx.Exec("GET", "/test2")
	})

	// cookie
	app.Get("/cookie", func(ctx iris.Context) {
		value := ctx.GetCookie("my_cookie")
		if value == "" {
			value = "NotSet"
			ctx.SetCookieKV("my_cookie", value)
			// Alternatively: ctx.SetCookie(&http.Cookie{...})
			ctx.SetCookieKV("test", "test")
			// Alternatively: ctx.SetCookie(&http.Cookie{...})
			//
			// If you want to set custom the path:
			// ctx.SetCookieKV(name, value, iris.CookiePath("/custom/path/cookie/will/be/stored"))
			//
			// If you want to be visible only to current request path:
			// (note that client should be responsible for that if server sent an empty cookie's path, all browsers are compatible)
			// ctx.SetCookieKV(name, value, iris.CookieCleanPath /* or iris.CookiePath("") */)
			// More:
			//                              iris.CookieExpires(time.Duration)
			//                              iris.CookieHTTPOnly(false)
		}
		ctx.Writef("Cookie value: %s \n", value)
	})

	app.Run(iris.Addr("localhost:18080"))
}

func myMiddleware(ctx iris.Context) {
	log.Info("myMiddleware...")
	ctx.Application().Logger().Infof("Runs before %s", ctx.Path())
	ctx.Next()
}

func list(ctx iris.Context) {
	books := []string{
		"Mastering Concurrency in Go",
		"Go Design Patterns",
		"Black Hat Go",
	}
	ctx.JSON(books)
	// TIP: negotiate the response between server's prioritizes
	// and client's requirements, instead of ctx.JSON:
	// ctx.Negotiation().JSON().MsgPack().Protobuf()
	// ctx.Negotiate(books)
}

type myParams struct {
	Name string   `param:"name"`
	Age  int      `param:"age"`
	Tail []string `param:"tail"`
}

type myHeaders struct {
	RequestID      string `header:"X-Request-Id,required"`
	Authentication string `header:"Authentication,required"`
}

// 强制将文件发送到客户端
func case1(ctx iris.Context) {
	src := "./files/first.zip"
	ctx.SendFile(src, "client.zip")
}

// 将下载速度限制为 ~50Kb/s，突发 100KB：
func case2(ctx iris.Context) {
	src := "./files/big.zip"
	// optionally, keep it empty to resolve the filename based on the "src".
	dest := ""
	limit := 50.0 * iris.KB
	burst := 100 * iris.KB
	ctx.SendFileWithRate(src, dest, limit, burst)
}
