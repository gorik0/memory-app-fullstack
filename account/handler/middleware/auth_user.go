package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"strings"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func AuthUser(tokenServ models.TokenServiceI) gin.HandlerFunc {
	return func(c *gin.Context) {

		h := authHeader{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {

				var invalidArgs []*invalidArgument
				for _, er := range errs {
					invalidArgs = append(invalidArgs, &invalidArgument{
						Field: er.Field(),
						Value: er.Value().(string),
						Tag:   er.Tag(),
						Param: er.Param(),
					})

				}
				err := apprerrors.NewBadRequest("See invalid arguments")
				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})

				c.Abort()
				return
			}

			err := apprerrors.NewInternal()
			c.JSON(err.Status(), gin.H{
				"error": err,
			})

			c.Abort()
			return
		}
		//	::: EXTRACTing token from HEADER
		headerSplitted := strings.Split(h.IDToken, "Bearer ")
		if len(headerSplitted) < 2 {
			err := apprerrors.NewBadRequest("Token is not with correct format :::")
			c.JSON(err.Status(), gin.H{
				"error": err,
			})

			c.Abort()
			return
		}
		user, err := tokenServ.ValidateIDToken(headerSplitted[1])
		if err != nil {
			log.Println("token is invalid!!! :::", err)

			c.Abort()
			return

		}
		c.Set("user", user)
		c.Next()

	}

}
