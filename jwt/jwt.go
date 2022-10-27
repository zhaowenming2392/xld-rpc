package jwt

import (
	"errors"
	"fmt"

	jwtGo "github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"helpers.zhaowenming.cn/strs"
	"helpers.zhaowenming.cn/time"
)

var defaultJwtKey = "xiaolandou.com"

type Jwt struct {
	Key string
}

func NewJwt(key string) *Jwt {
	if key == "" {
		key = defaultJwtKey
	}
	return &Jwt{Key: key}
}

//创建jwt字符串
func (j *Jwt) Create(c *Claim, f func(string)) (string, error) {
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, c)
	ss, err := token.SignedString([]byte(j.Key))
	if err != nil {
		fmt.Printf("%v %v", ss, err)
	}

	//TODO 看需要可以通过此函数保存数据库等
	if f != nil {
		f(ss)
	}

	return ss, err
}

//解析jwt字符串
func (j *Jwt) Parse(tokenStr string) (*Claim, error) {
	c := new(Claim)
	token, err := jwtGo.ParseWithClaims(tokenStr, c, func(token *jwtGo.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(j.Key), nil
	})

	if token.Valid {
		fmt.Printf("%+v\n", token.Claims)
		if claims, ok := token.Claims.(*Claim); ok {
			return claims, nil
		} else {
			return nil, errors.New("Couldn't handle this token1")
		}
	} else if ve, ok := err.(*jwtGo.ValidationError); ok {
		if ve.Errors&jwtGo.ValidationErrorMalformed != 0 {
			return nil, errors.New("That's not even a token")
		} else if ve.Errors&jwtGo.ValidationErrorExpired != 0 {
			// Token is either expired or not active yet
			return nil, errors.New("Token is either expired")
		} else if ve.Errors&jwtGo.ValidationErrorNotValidYet != 0 {
			// Token is either expired or not active yet
			return nil, errors.New("Token not active yet")
		} else {
			return nil, err
			//return nil, errors.New("Couldn't handle this token2")
		}
	} else {
		return nil, errors.New("Couldn't handle this token3")
	}
}

//刷新旧的令牌声明
//注意必须要使用 旧jwt来请求--》解析出旧Claim-->修改到新的过期时间节点--》给出新的有效jwt
//如果旧jwt已经失效，则需要重新登录啦
func (j *Jwt) Refresh(oldTokenStr string, f func(string)) (string, error) {
	oldClaim, err := j.Parse(oldTokenStr)
	if err != nil {
		return "", err
	}
	oldClaim.ExpiresAt = time.Now() + 30*24*60*60
	return j.Create(oldClaim, f)
}

type Claim struct {
	UserID   int    `json:"userid"`
	Username string `json:"username"`
	jwtGo.StandardClaims
}

//新颁发一个声明
func NewClaim(userID int, username string) *Claim {
	uid := uuid.NewV4()
	return &Claim{
		userID,
		username,
		jwtGo.StandardClaims{
			Audience:  "xld",
			ExpiresAt: time.Now() + 30*24*60*60,
			Id:        strs.HexToStr(uid.Bytes()),
			IssuedAt:  time.Now(),
			Issuer:    "xiaolandou.com",
			NotBefore: time.Now(),
			Subject:   "Authorization",
		},
	}
}

func (c *Claim) Valid() error {
	//继承StandardClaims里面的有效期的验证
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	//颁发时相关声明验证
	//也可以采用 c.StandardClaims.VerifyAudience("",false) 参考更加规范的认证
	if c.Audience != "xld" {
		return errors.New("This is not a valid token1")
	}
	if c.Issuer != "xiaolandou.com" {
		return errors.New("This is not a valid token2")
	}
	if c.Subject != "Authorization" {
		return errors.New("This is not a valid token3")
	}

	//复杂验证
	//TODO 查找用户是否存在

	//TODO 查找jwt是否在黑名单中

	return nil
}
