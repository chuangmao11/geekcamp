package web

import (
	"geekcamp/webook/internal/domain"
	"geekcamp/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPatten    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPatten = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPatten, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPatten, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return

	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须超过8位")
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "账号密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	sess := sessions.Default(ctx)
	sess.Save()
	sess.Set("userId", user.Id)
	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type InfoReq struct {
		Email    string `json:"email"`
		NickName string `json:"nickname"`
		Birthday string `json:"birthday"`
		Info     string `json:"info"`
	}
	var req InfoReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	_, err := u.svc.Edit(ctx, domain.User{
		Id:       id,
		NickName: req.NickName,
		Birthday: req.Birthday,
		Info:     req.Info,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "补充修改信息成功")
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	userinfo, err := u.svc.Profile(ctx, domain.User{
		Id: id,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.JSON(http.StatusOK, userinfo)
}
