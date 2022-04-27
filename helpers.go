package ginCore

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func FileExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func IsMatch(text string, filter string) bool {
	reg := regexp.MustCompile(filter)
	result := reg.FindAllString(text, -1)
	if len(result) > 0 {
		return true
	} else {
		return false
	}
}

func In(haystack interface{}, needle interface{}) bool {
	sVal := reflect.ValueOf(haystack)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == needle {
				return true
			}
		}

		return false
	}

	return false
}

func StrVal(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func IntVal(value interface{}) int {
	var key = 0
	if value == nil {
		return 0
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = int(ft)
	case float32:
		ft := value.(float32)
		key = int(ft)
	case int:
		it := value.(int)
		key = it
	case uint:
		it := value.(uint)
		key = int(it)
	case int8:
		it := value.(int8)
		key = int(it)
	case uint8:
		it := value.(uint8)
		key = int(it)
	case int16:
		it := value.(int16)
		key = int(it)
	case uint16:
		it := value.(uint16)
		key = int(it)
	case int32:
		it := value.(int32)
		key = int(it)
	case uint32:
		it := value.(uint32)
		key = int(it)
	case int64:
		it := value.(int64)
		key = int(it)
	case uint64:
		it := value.(uint64)
		key = int(it)
	case string:
		it := value.(string)
		key, _ = strconv.Atoi(it)
	case []byte:
		it := value.([]byte)
		key, _ = strconv.Atoi(string(it))
	default:
		newValue, _ := json.Marshal(value)
		it := string(newValue)
		key, _ = strconv.Atoi(it)
	}
	return key
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

func SampleRequestGet(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func GetMapFromContext(context *gin.Context) map[string]interface{} {
	result := make(map[string]interface{})
	json := make(map[string]interface{})
	for i, _ := range context.Request.URL.Query() {
		temp := context.DefaultQuery(i, "")
		if temp != "" {
			result[i] = temp
		}
	}
	for j, _ := range context.Request.Form {
		temp := context.DefaultPostForm(j, "")
		if temp != "" {
			result[j] = temp
		}
	}
	for k, _ := range context.Request.PostForm {
		temp := context.DefaultPostForm(k, "")
		if temp != "" {
			result[k] = temp
		}
	}
	_ = context.ShouldBindBodyWith(&json, binding.JSON)
	for index, value := range json {
		result[index] = value
	}

	return result
}

func GetIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" || ip == "::1" {
		ip = "127.0.0.1"
	}

	return ip
}

func GetToken(ctx *gin.Context) string {
	var token = ""
	bear := ctx.Request.Header.Get("Authorization")
	if bear != "" {
		token = strings.Replace(bear, "Bearer ", "", 1)
	} else {
		token = ctx.DefaultQuery("token", "")
		if token == "" {
			token = ctx.DefaultPostForm("token", "")
		}
		if token == "" {
			token, _ = ctx.Cookie("token")
		}
		if token == "" {
			token = ctx.GetHeader("X-TOKEN")
		}
	}
	return token
}
