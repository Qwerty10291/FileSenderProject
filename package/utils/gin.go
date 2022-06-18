package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)



func GetJsonFromRequest(ctx *gin.Context, storage interface{}) error {
	data, err := ctx.GetRawData()
	if err != nil{
		return err
	}
	err = json.Unmarshal(data, storage)
	return err
}