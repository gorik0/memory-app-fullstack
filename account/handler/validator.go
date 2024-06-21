package handler

import (
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"memory-app/account/models/apprerrors"
)

type invalidArg struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func BindData(ctx *gin.Context, req SignUpRequest) bool {

	err := ctx.ShouldBind(&req)
	if err != nil {
		log.Printf("Bad request, coudln't unmarshal to USER  ::: %s \n", err)
		e := apprerrors.NewBadRequest(err.Error())
		ctx.JSON(e.Status(), gin.H{"error": e})
		return false
	}

	v := validator.New()
	err = v.Struct(req)
	if err != nil {
		var errors validator.ValidationErrors
		if as := errors2.As(err, &errors); as {
			var invalidArguments []invalidArg
			for _, fieldError := range errors {
				arg := invalidArg{}
				arg.Field = fieldError.Field()
				arg.Value = fieldError.Value().(string)
				arg.Tag = fieldError.Tag()
				arg.Param = fieldError.Param()
				invalidArguments = append(invalidArguments, arg)
			}
			log.Printf("Bad request, invalid fields  ::: %s \n", invalidArguments)
			e := apprerrors.NewBadRequest(err.Error())
			ctx.JSON(e.Status(), gin.H{"error": e, "invalidArguments": invalidArguments})
			return false

		}

		return false
	}
	return true
}
