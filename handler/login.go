package handler

import (
	"time"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/sso/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"fmt"
)

func Login(c echo.Context) error {
	// upload param
	mobile := c.FormValue("mobile")
	password := c.FormValue("password")
	// auth
	sql := `SELECT id, password, mobile, nickname,salt FROM account WHERE mobile = ? AND status = 1`
	rows, _ := global.DB.Query(sql, mobile)
	if len(rows) != 1 {
		return utils.Error(c, "帐号或密码不存在", nil)
	}
	userInfo := rows[0]
	pwd, _ := convert.ToString(userInfo["password"])
	salt, _ := convert.ToString(userInfo["salt"])
	if utils.Sha1Encode(password+salt) != pwd {
		return utils.Error(c, "密码错误", nil)
	}
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return utils.Error(c, "服务器异常", nil)
	}

	return utils.Success(c, "操作成功", map[string]string{
		"token": t,
	})
}

func Register(c echo.Context) error {
	mobile := c.FormValue("mobile")
	smsCode := c.FormValue("sms_code")
	password := c.FormValue("password")

	if mobile == "" {
		return utils.Error(c, "待发送短信的手机号码不能为空", nil)
	}
	if smsCode == "" {
		return utils.Error(c, "短信验证码不能为空", nil)
	}
	if password == "" {
		return utils.Error(c, "登陆密码不能为空", nil)
	}
	if len(password) > 20 {
		return utils.Error(c, "登陆密码最长不能超过20位", nil)
	}
	if utils.CheckRegexp(password, "/^[0-9A-Z_-]*$/i") {
		return utils.Error(c, "登陆密码仅包含字母数字字符，包括破折号、下划线", nil)
	}
	if !utils.CheckMobile(mobile) {
		return utils.Error(c, "手机号码格式错误", nil)
	}
	if CheckMobile(mobile) {
		return utils.Error(c, "手机号码已注册", nil)
	}
	rdSmsCode, setStrErr := global.RD.GetString(mobile + "_sms_code")
	if setStrErr != nil {
		global.Log.Error("注册帐号验证码Redis存储错误：" + setStrErr.Error())
	}
	if rdSmsCode != smsCode {
		return utils.Error(c, "短信验证码错误", nil)
	}

	sql := "INSERT INTO account (id,mobile,password,salt,status) VALUES (?,?,?,?,1)"
	iw, _ := utils.NewIdWorker(1)
	id, idErr := iw.NextId()
	if idErr != nil {
		return utils.Error(c, "ID生成器发生错误", nil)
	}
	_, err := global.DB.Insert(sql, id, mobile, utils.Sha1Encode(password+smsCode), smsCode)
	if err != nil {
		return utils.Error(c, "注册失败，"+err.Error(), nil)
	}

	global.RD.DelKey(mobile + "_sms_code")
	return utils.Success(c, "注册成功", nil)
}

var (
	gatewayUrl      = "http://dysmsapi.aliyuncs.com/"
	accessKeyId     = ""
	accessKeySecret = ""
	signName        = ""
	templateCode    = ""
	templateParam   = "{\"code\":\"%s\"}"
)

func RegSendSms(c echo.Context) error {
	mobile := c.FormValue("mobile")
	code := c.FormValue("code")
	if mobile == "" {
		return utils.Error(c, "待发送短信的手机号码不能为空", nil)
	}
	if code == "" {
		return utils.Error(c, "待发送短信需要的图形验证码不能为空", nil)
	}

	if !utils.CheckMobile(mobile) {
		return utils.Error(c, "手机号码格式错误", nil)
	}

	//短信接口数量限制

	//注册帐号限制
	if CheckMobile(mobile) {
		return utils.Error(c, "手机号码已注册", nil)
	}

	smsCode := utils.NewLen(4)
	templateParam = fmt.Sprintf(templateParam, smsCode)
	smsClient := utils.NewSmsClient(gatewayUrl)
	if result, err := smsClient.Execute(accessKeyId, accessKeySecret, mobile, signName, templateCode, templateParam); err != nil {
		fmt.Println("error:", err.Error())
		return utils.Error(c, "发送失败"+err.Error(), nil)
	} else {
		resultCode := fmt.Sprintf("%s", result["Code"])
		if resultCode == "OK" {
			_, setStrErr := global.RD.SetAndExpire(mobile+"_sms_code", smsCode, global.SmsCodeExpire)
			if setStrErr != nil {
				global.Log.Error("注册帐号验证码Redis存储错误：" + setStrErr.Error())
			}
			return utils.Success(c, "短信发送成功", nil)
		} else {
			return utils.Error(c, "短信发送失败", nil)
		}
	}

}

func CheckRegMobile(c echo.Context) error {
	mobile := c.FormValue("mobile")
	if mobile == "" || !utils.CheckMobile(mobile) {
		return utils.Success(c, "", nil)
	}
	if CheckMobile(mobile) {
		return utils.Error(c, "手机号码已注册", nil)
	}
	return utils.Success(c, "", nil)
}

func CheckMobile(mobile string) bool {
	if mobile == "" {
		return false
	}
	sql := `SELECT mobile FROM account WHERE mobile = ? `
	rows, err := global.DB.Query(sql, mobile)
	if err != nil {
		return false
	}
	if len(rows) >= 1 {
		return true
	}
	return false
}
