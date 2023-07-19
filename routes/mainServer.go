package routes

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"main/model"
	u "main/utils"
	"os"
	"strings"
	"time"
)

//region struct definitions

// User contains the login user information
type User struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
}

// TransferItem contains the information of the transfer item
type TransferItem struct {
	To     string `json:"to" xml:"to" form:"to"`
	Amount int    `json:"amount" xml:"amount" form:"amount"`
}

type MarkdownPage struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

//endregion

var sessionStore = session.New()
var csrfActivated = true

func init() {
	sessionStore.RegisterType(fiber.Map{})
	// this mean, csrf is activated
	csrfActivated = len(os.Args) > 1 && os.Args[1] == "withoutCsrf"
}

// Add CSRF protection middleware.
// Should be done AFTER session middleware.
var csrfProtection = csrf.New(csrf.Config{
	// only to control the switch whether csrf is activated or not
	Next: func(c *fiber.Ctx) bool {
		return csrfActivated
	},
	KeyLookup:      "form:_csrf",
	CookieName:     "csrf_",
	CookieSameSite: "Strict",
	Expiration:     1 * time.Hour,
	KeyGenerator:   utils.UUID,
	ContextKey:     "token",
})

// RegisterRoutes registers the routes and middlewares necessary for the server
func RegisterRoutes(app *fiber.App) {
	// Super simple login system.
	// This is not how real login systems should work.
	validLogins := []User{
		{Username: "admin", Password: "123456"},
	}
	// Simple accounts ledger.
	// This information would normally be stored in a database like MySQL, PostgreSQL, etc.

	app.Use(recover.New())

	app.Get("/dashboard", requireLogin, monitor.New())

	app.Get("/datago", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		return c.Render("views/datago", fiber.Map{})
	})

	app.Get("/", csrfProtection, func(c *fiber.Ctx) error {
		return c.Render("views/home", fiber.Map{})
		//currSession, err := sessionStore.Get(c)
		//if err != nil {
		//	return err
		//}
		//sessionUser := currSession.Get("User").(fiber.Map)
		//// release the currSession
		//err = currSession.Save()
		//if err != nil {
		//	return err
		//}
		//
		//if sessionUser["Name"] == "" {
		//	return c.Status(fiber.StatusBadRequest).SendString("User is empty")
		//}
		//username := sessionUser["Name"].(string)
		//
		//return c.Render("views/home", fiber.Map{
		//	"username":  username,
		//	"balance":   accounts[username],
		//	"csrfToken": c.Locals("token"),
		//})
	})

	app.Get("/lists", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		return c.Render("views/lists", fiber.Map{})
	})

	app.Get("/markdown-page", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		return c.Render("views/mdpage", fiber.Map{
			"csrfToken": c.Locals("token"),
		})
	})

	app.Get("/markdown-page/:id", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		mdLog := model.GetMdLog2(c)
		ref := c.Query("ref", "")
		if ref != "" {
			splits := strings.Split(ref, "-")
			exist := model.GetMdClickByRandom(splits[1])
			if !exist {
				model.NewMdClickForInner(mdLog.ID, splits[0], 1, splits[1])
			}
		}
		return c.Render("views/mdpage-edit", fiber.Map{
			"csrfToken": c.Locals("token"),
			"title":     mdLog.Title,
			"content":   mdLog.Content,
			"id":        mdLog.ID,
		})
	})

	app.Post("/markdown-add", requireLogin, csrfProtection, model.NewMdLog)
	app.Post("/markdown-update/:id", requireLogin, csrfProtection, model.UpdateMdLog)
	app.Get("/markdown-new5", requireLogin, csrfProtection, model.GetMdLogs5)
	app.Get("/markdown-search/:keyword", requireLogin, csrfProtection, model.GetMdLogByKeyword)
	app.Get("/markdown-latest-100", requireLogin, csrfProtection, model.GetMdLogLatest100)
	app.Delete("/markdown/:id", requireLogin, model.DeleteMdLog)

	app.Get("/markdown-list", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		return c.Render("views/mdlists", fiber.Map{})
	})

	app.Post("/upload-attachment-jiwekfsdsiqwokesdfsucvjoekewloxcvhwsdmlauwio", func(c *fiber.Ctx) error {
		// Get first fileUploader from form field "document":
		fileUploader, err := c.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return err
		}
		open, err := fileUploader.Open()
		if err != nil {
			fmt.Println(err)
			return err
		}
		// 创建OSSClient实例。
		// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
		// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
		client, err := oss.New("https://oss-us-west-1.aliyuncs.com", "LTAI5tLn6jhY3Nt9TXGpuekK", "")
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 填写存储空间名称，例如examplebucket。
		bucket, err := client.Bucket("notea-us")
		if err != nil {
			fmt.Println(err)
			return err
		}

		defer open.Close()

		// 将文件流上传至exampledir目录下的exampleobject.txt文件。
		prefix := "datago_pic_jkmelixx"
		date := time.Now().Format("2006-01-02")
		path := prefix + "/" + date + "/" + utils.UUID() + "/" + utils.UUID() + "." + fileUploader.Filename

		err = bucket.PutObject(path, open, oss.ContentType(fileUploader.Header.Get("Content-Type")))
		if err != nil {
			fmt.Println(err)
			return err
		}

		//meta, err := bucket.GetObjectDetailedMeta(path)
		//if err != nil {
		//	return err
		//}
		//fmt.Println(meta)
		// Save fileUploader to root directory:
		//https://github.com/xiaogao67/gin-cloud-storage/blob/master/controller/upload.go#L101
		attachment := model.NewAttachmentForInner("https://notea-us.oss-us-west-1.aliyuncs.com/" + path)
		return c.JSON(map[string]string{
			"fileUrl": attachment.FilePath,
		})
	})

	app.Get("/pop/api/alioss/policy", requireLogin, func(c *fiber.Ctx) error {
		return c.JSON(map[string]interface{}{
			"code":    0,
			"success": true,
			"msg":     "签名成功",
			"data": map[string]interface{}{
				"accessid":  "XXXXX",
				"host":      "https://datago2023.oss-cn-hangzhou.aliyuncs.com",
				"policy":    "XXXX==",
				"signature": "XXXX=",
				"expire":    1754851252,
			},
		})
	})

	app.Get("/markdown-lists-json", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		p := u.NewPagination(c)
		md, err := model.GetMdLogByPage(p)
		if err != nil {
			return c.JSON(err)
		}
		type PT struct {
			Code  int           `json:"code"`
			Count int           `json:"count"`
			Data  []model.MdLog `json:"data"`
		}
		ps := new(PT)
		ps.Code = 0
		ps.Count = p.Total
		ps.Data = md
		return c.JSON(ps)
	})

	app.Get("/lists-json", requireLogin, csrfProtection, func(c *fiber.Ctx) error {
		p := u.NewPagination(c)
		fish, err := model.GetFishAddressByPage(p)
		if err != nil {
			return c.JSON(err)
		}
		type PT struct {
			Code  int                 `json:"code"`
			Count int                 `json:"count"`
			Data  []model.FishAddress `json:"data"`
		}
		ps := new(PT)
		ps.Code = 0
		ps.Count = p.Total
		ps.Data = fish
		return c.JSON(ps)
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		currSession, err := sessionStore.Get(c)
		defer func(currSession *session.Session) {
			err := currSession.Save()
			if err != nil {
				panic(err)
			}
		}(currSession)
		if err != nil {
			return err
		}
		err = currSession.Destroy()
		if err != nil {
			return err
		}
		return c.Redirect("/")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("views/login", fiber.Map{})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := &User{}
		err := c.BodyParser(user)
		if err != nil {
			return err
		}

		if user.Username == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Username is required.")
		}

		if user.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Password is required.")
		}

		if !findUser(validLogins, user) {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid username or password.")
		}

		// Valid login.
		// Create a new currSession and save their user data in the currSession.
		currSession, err := sessionStore.Get(c)
		defer func(currSession *session.Session) {
			err := currSession.Save()
			if err != nil {
				panic(err)
			}
		}(currSession)
		if err != nil {
			return err
		}
		err = currSession.Regenerate()
		if err != nil {
			return err
		}
		currSession.Set("User", fiber.Map{"Name": user.Username})

		return c.Redirect("/")
	})
}

// Create a helper function to require login for some routes.
func requireLogin(c *fiber.Ctx) error {
	currSession, err := sessionStore.Get(c)
	if err != nil {
		return err
	}
	user := currSession.Get("User")
	defer func(currSession *session.Session) {
		err := currSession.Save()
		if err != nil {
			panic(err)
		}
	}(currSession)

	if user == nil {
		// This request is from a user that is not logged in.
		// Send them to the login page.
		return c.Redirect("/login")
	}

	// If we got this far, the request is from a logged-in user.
	// Continue on to other middleware or routes.
	return c.Next()
}

func findUser(list []User, compareUser *User) bool {
	for _, item := range list {
		if item.Username == compareUser.Username && item.Password == compareUser.Password {
			return true
		}
	}
	return false
}
