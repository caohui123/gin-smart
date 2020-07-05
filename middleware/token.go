package middleware

import (
	"net/http"

	"github.com/jangozw/gin-smart/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/jangozw/gin-smart/errcode"

	"github.com/jangozw/gin-smart/pkg/app"
	"github.com/jangozw/go-api-facility/auth"
)

// header 中 token key
const TokenHeaderKey = "Authorization"

// 无论什么业务场景，只需要传递参数verifyInstance实现自验证token的接口 auth.VerifyIF 即可，这个方法的代码不用动
func CheckToken(verifyInstance auth.VerifyIF) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(TokenHeaderKey)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, app.ResponseFailByCode(errcode.ErrToken))
			return
		}
		jwtPayload, callback, err := auth.VerifyToken(verifyInstance, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, app.ResponseFailByCode(errcode.ErrToken))
			return
		}
		appLoginUser := appLoginUser(jwtPayload.AccountInfo, callback)
		c.Set(app.CtxKeyNeedLogin, true)
		c.Set(app.CtxKeyLoginUser, appLoginUser)
		// 继续下一步
		c.Next()
	}
}

func appLoginUser(accountInfo auth.AccountInfo, detail interface{}) app.LoginUser {
	// extra
	// 从 accountInfo.AccountExtra， 或 detail 里得到 extra信息，都是调用时候根据业务场景自定义的
	// accountInfo.AccountExtra 是在login时候自定义， detail 是在验证token 成功后回调的自定义数据
	accountID := util.StringNumber(accountInfo.AccountID)
	return app.LoginUser{
		ID:    accountID.Uint(),
		Extra: detail,
	}
}
