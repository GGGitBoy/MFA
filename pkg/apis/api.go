package apis

import (
	"encoding/json"
	"fmt"
	"googleAuthenticator"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

func RegisterAPIs() *restful.Container {
	container := restful.NewContainer()

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"OPTIONS", "POST"},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      container}

	docxWs := new(restful.WebService).Path("/").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	container.Add(docxWs)

	container.Filter(cors.Filter)

	container.Filter(container.OPTIONSFilter)

	docxWs.Route(docxWs.GET("/index").To(index))
	docxWs.Route(docxWs.POST("/login").To(login))
	docxWs.Route(docxWs.POST("/validate").To(validate))

	docxWs.Route(docxWs.GET("/QRcode").To(createQRcode))
	docxWs.Route(docxWs.POST("/enableMFA").To(enableMFA))

	return container
}

func enableMFA(req *restful.Request, resp *restful.Response) {

	fmt.Println("test")
	bodyData, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		logrus.Errorf("Read req body err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	var mfa MFA
	if err := json.Unmarshal(bodyData, &mfa); err != nil {
		return
	}

	fmt.Println(mfa.Status)

	return
}

func index(req *restful.Request, resp *restful.Response) {

	loginHTML, _ := ioutil.ReadFile("pkg/html/index.html")
	resp.Write(loginHTML)

	return
}

func login(req *restful.Request, resp *restful.Response) {

	bodyData, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		logrus.Errorf("Read req body err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	var user User
	if err := json.Unmarshal(bodyData, &user); err != nil {
		return
	}

	if user.Username == "" || user.Password == "" {
		logrus.Errorf("用户名和密码不能为空")
		err = resp.WriteErrorString(400, "用户名和密码不能为空")
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}

	if user.Username == "admin" {
		var message string
		var errCode int
		if user.Password == "123456" {
			message = "账号密码匹配"
			errCode = 200
			loginHTML, _ := ioutil.ReadFile("pkg/html/validate.html")
			resp.Write(loginHTML)
		} else {
			message = "密码错误"
			errCode = 400
		}
		logrus.Infof("%s 登陆情况: %s", user.Username, message)
		err = resp.WriteErrorString(errCode, message)
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
	}

	return
}

func createQRcode(req *restful.Request, resp *restful.Response) {

	ga := googleAuthenticator.NewGAuth()

	secret := createSecret(ga)
	fmt.Println(secret)

	// go showCode(ga, secret)

	aa, _ := qrcode.Encode("otpauth://totp/kisexu@gmail.com?secret="+secret, qrcode.Medium, 256)
	resp.Write(aa)

	return
}

func showCode(ga *googleAuthenticator.GAuth, secret string) {
	for {
		fmt.Println("Code: " + getCode(ga, secret))
		time.Sleep(time.Duration(2) * time.Second)
	}

}

func validate(req *restful.Request, resp *restful.Response) {

	bodyData, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		logrus.Errorf("Read req body err:%v", err)
		err = resp.WriteErrorString(400, err.Error())
		if err != nil {
			logrus.Errorf("Failed to write error string err:%v", err)
		}
		return
	}
	var auth Auth
	if err := json.Unmarshal(bodyData, &auth); err != nil {
		return
	}

	ga := googleAuthenticator.NewGAuth()

	auth.Secret = createSecret(ga)

	if auth.Number == getCode(ga, auth.Secret) {
		logrus.Infof("%s 验证成功", auth.Number)
	}

	return
}

func GenValidateCode(width int) string {
	numeric := [9]byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}

	return sb.String()
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	Number string `json:"number"`
	Secret string `json:"secret"`
}

type MFA struct {
	Status string `json:"status"`
}

func createSecret(ga *googleAuthenticator.GAuth) string {
	secret, err := ga.CreateSecret(16)
	if err != nil {
		return ""
	}
	return secret
}

func getCode(ga *googleAuthenticator.GAuth, secret string) string {
	code, err := ga.GetCode(secret)
	if err != nil {
		return "*"
	}
	return code
}
