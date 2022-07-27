package user

import (
	. "apiserver/handler"
	"apiserver/logger"
	"apiserver/pkg/errno"
	"apiserver/util"
	"fmt"

	"github.com/gin-gonic/gin"
)


func Create(c *gin.Context) {
	logger.Info("User Create function called.", logger.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	if err := r.checkParam(); err != nil {
		SendResponse(c, err, nil)
		return
	}

	admin2 := c.Param("username")
	logger.Infof("URL username: %s", admin2)

	desc := c.Query("desc")
	logger.Infof("URL key param desc: %s", desc)

	contentType := c.GetHeader("Content-Type")
	logger.Infof("Header Content-Type: %s", contentType)

	logger.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
	if r.Username == "" {
		SendResponse(c, errno.New(errno.ErrUserNotFound, fmt.Errorf("username can not found in db: xx.xx.xx.xx")), nil)
		return
	}

	if r.Password == "" {
		SendResponse(c, fmt.Errorf("password is empty"), nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	// Show the user information.
	SendResponse(c, nil, rsp)
}

func (r *CreateRequest) checkParam() error {
	if r.Username == "" {
		return errno.New(errno.ErrValidation, nil).Add("username is empty.")
	}

	if r.Password == "" {
		return errno.New(errno.ErrValidation, nil).Add("password is empty.")
	}

	return nil
}

