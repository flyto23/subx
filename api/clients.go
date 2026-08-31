package api

import (
"crypto/md5"
"encoding/hex"
"encoding/json"
"fmt"
"io"
"log"
"net/http"
"net/url"
"strings"
"sublink/models"
"sublink/node"

"github.com/gin-gonic/gin"
)

var sunName string // 使用小写命名

// md5Hash MD5 加密
func md5Hash(src string) string {
m := md5.New()
m.Write([]byte(src))
return hex.EncodeToString(m.Sum(nil))
}

// GetClient 获取客户端配置
func GetClient(c *gin.Context) {
token := c.Query("token")
clientIndex := c.Query("client")

if token == "" {
log.Println("token 为空")
c.Writer.WriteString("token 为空")
return
}

sub := new(models.Subcription)
list, err := sub.List()
if err != nil {
log.Printf("获取订阅列表失败：%v", err)
c.Writer.WriteString("获取订阅列表失败")
return
}

// 查找匹配的订阅
for _, sub := range list {
sunName = sub.Name
if md5Hash(sunName) == strings.ToLower(token) {
handleClientRequest(c, clientIndex)
return
}
}

log.Println("未找到匹配的订阅")
c.Writer.WriteString("未找到匹配的订阅")
}

// handleClientRequest 处理客户端请求
func handleClientRequest(c *gin.Context, clientIndex string) {
// 如果指定了客户端类型，直接返回对应配置
switch clientIndex {
case "clash":
GetClash(c)
return
case "surge":
GetSurge(c)
return
case "v2ray":
GetV2ray(c)
return
}

// 自动识别客户端
userAgent := c.GetHeader("User-Agent")
if userAgent == "" {
log.Println("User-Agent 为空")
GetV2ray(c)
return
}

clientList := []string{"clash", "surge"}
for _, client := range clientList {
if strings.Contains(strings.ToLower(userAgent), strings.ToLower(client)) {
switch client {
case "clash":
GetClash(c)
return
case "surge":
GetSurge(c)
return
}
}
}

// 默认返回 V2Ray 配置
GetV2ray(c)
}

// GetV2ray 获取 V2Ray 配置
func GetV2ray(c *gin.Context) {
var sub models.Subcription
if sunName == "" {
c.Writer.WriteString("订阅名为空")
return
}

sub.Name = sunName
if err := sub.Find(); err != nil {
c.Writer.WriteString("找不到这个订阅:" + sunName)
return
}

baselist := buildNodeLinks(sub.Nodes)

c.Set("subname", sunName)
filename := fmt.Sprintf("%s.txt", sunName)
encodedFilename := url.QueryEscape(filename)
c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
c.Writer.WriteString(node.Base64Encode(baselist))
}

// GetClash 获取 Clash 配置
func GetClash(c *gin.Context) {
var sub models.Subcription
sub.Name = sunName

if err := sub.Find(); err != nil {
c.Writer.WriteString("找不到这个订阅:" + sunName)
return
}

models.DB.Model(sub).Preload("Nodes").Find(&sub)
log.Println("订阅名:", sub.Nodes)

urls := buildNodeURLs(sub.Nodes)
log.Println("urls", urls)

var configs node.SqlConfig
if err := json.Unmarshal([]byte(sub.Config), &configs); err != nil {
c.Writer.WriteString("配置读取错误")
return
}

decodeClash, err := node.EncodeClash(urls, configs)
if err != nil {
c.Writer.WriteString(err.Error())
return
}

c.Set("subname", sunName)
filename := fmt.Sprintf("%s.yaml", sunName)
encodedFilename := url.QueryEscape(filename)
c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
c.Writer.WriteString(string(decodeClash))
}

// GetSurge 获取 Surge 配置
func GetSurge(c *gin.Context) {
var sub models.Subcription
sub.Name = sunName

if err := sub.Find(); err != nil {
c.Writer.WriteString("找不到这个订阅:" + sunName)
return
}

urls := buildNodeURLs(sub.Nodes)

var configs node.SqlConfig
if err := json.Unmarshal([]byte(sub.Config), &configs); err != nil {
c.Writer.WriteString("配置读取错误")
return
}

decodeClash, err := node.EncodeSurge(urls, configs)
if err != nil {
c.Writer.WriteString(err.Error())
return
}

c.Set("subname", sunName)
filename := fmt.Sprintf("%s.conf", sunName)
encodedFilename := url.QueryEscape(filename)
c.Writer.Header().Set("Content-Disposition", "inline; filename*=utf-8''"+encodedFilename)
c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

host := c.Request.Host
reqURL := c.Request.URL.String()

// 如果包含头部更新信息，直接返回
if strings.Contains(decodeClash, "#!MANAGED-CONFIG") {
c.Writer.WriteString(decodeClash)
return
}

// 否则插入头部更新信息
interval := fmt.Sprintf("#!MANAGED-CONFIG %s interval=86400 strict=false", host+reqURL)
c.Writer.WriteString(interval + "\n" + decodeClash)
}

// buildNodeLinks 构建节点链接列表 (V2Ray)
func buildNodeLinks(nodes []models.Node) string {
var baselist strings.Builder
for _, v := range nodes {
switch {
case strings.Contains(v.Link, ","):
links := strings.Split(v.Link, ",")
for _, link := range links {
baselist.WriteString(link + "\n")
}
case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
resp, err := http.Get(v.Link)
if err != nil {
log.Println(err)
continue
}
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()
baselist.WriteString(node.Base64Decode(string(body)) + "\n")
default:
baselist.WriteString(v.Link + "\n")
}
}
return baselist.String()
}

// buildNodeURLs 构建节点 URL 列表 (Clash/Surge)
func buildNodeURLs(nodes []models.Node) []string {
var urls []string
for _, v := range nodes {
switch {
case strings.Contains(v.Link, ","):
links := strings.Split(v.Link, ",")
urls = append(urls, links...)
case strings.Contains(v.Link, "http://") || strings.Contains(v.Link, "https://"):
resp, err := http.Get(v.Link)
if err != nil {
log.Println(err)
continue
}
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()
decoded := node.Base64Decode(string(body))
links := strings.Split(decoded, "\n")
urls = append(urls, links...)
default:
urls = append(urls, v.Link)
}
}
return urls
}
