package deployed

import (
	"github.com/mojocn/base64Captcha"
	"github.com/thinkgos/go-core-package/lib/password"
)

var Verify = password.BCrypt{}

// var DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 4,
// 	"234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"}).ConvertFonts()
var Captcha = base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, base64Captcha.DefaultMemStore)
