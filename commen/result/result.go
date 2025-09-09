package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Status struct {
	HTTPcode   int
	StatusCode int
	Message    string
}

func (s *Status) httpcode() int {
	return s.HTTPcode
}

func (s *Status) statuscode() int {
	return s.StatusCode
}

func (s *Status) message() string {
	return s.Message
}

func newstatus(httpcode int, statuscode int, message string) Status {
	return Status{
		HTTPcode:   httpcode,
		StatusCode: statuscode,
		Message:    message,
	}
}

func Sucess(c *gin.Context, data interface{}) {
	h := gin.H{
		"status":  0,
		"message": "success",
	}
	//r := make(R)
	//r.ToMap(data)
	//h["data"] = r
	h["data"] = data
	c.JSON(http.StatusOK, h)
}

func Error(c *gin.Context, s Status) {
	c.JSON(s.httpcode(), gin.H{
		"status":  s.StatusCode,
		"message": s.Message,
	})
}
