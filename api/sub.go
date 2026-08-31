// api/subcription.go

package api

import (
"log"
"strconv"
"strings"
"sublink/models"
"sublink/utils"

"github.com/gin-gonic/gin"
)

func SubTotal(c *gin.Context) {
var Sub models.Subcription
subs, err := Sub.List()
if err != nil {
utils.ErrorWithStatus(c, 500, "50001", "取得订阅总数失败")
return
}
utils.Success(c, len(subs), "取得订阅总数")
}

func SubGet(c *gin.Context) {
var Sub models.Subcription
Subs, err := Sub.List()
if err != nil {
utils.ErrorWithStatus(c, 500, "50001", "node list error")
return
}
utils.Success(c, Subs, "node get")
}

func SubAdd(c *gin.Context) {
name := c.PostForm("name")
configs := c.PostForm("config")
nodes := c.PostForm("nodes")

if name == "" || nodes == "" {
utils.Error(c, "40001", "订阅名称或节点不能为空")
return
}

var NodesData []models.Node
for _, nodeName := range strings.Split(nodes, ",") {
if strings.TrimSpace(nodeName) == "" {
continue
}
firstNode := models.Node{Name: nodeName}
result := models.DB.Model(models.Node{}).Where("name = ?", firstNode.Name).First(&firstNode)
if result.Error != nil {
log.Println(result.Error)
utils.Error(c, "40002", result.Error.Error())
return
}
NodesData = append(NodesData, firstNode)
}

sub := models.Subcription{
Name:      name,
Config:    configs,
NodeOrder: nodes,
Nodes:     NodesData,
}
if err := sub.Add(); err != nil {
utils.Error(c, "40003", "添加订阅失败："+err.Error())
return
}
utils.Success(c, nil, "添加订阅成功")
}

func SubUpdate(c *gin.Context) {
newName := c.PostForm("name")
oldName := c.PostForm("oldname")
configs := c.PostForm("config")
nodes := c.PostForm("nodes")

if newName == "" || nodes == "" {
utils.Error(c, "40001", "订阅名称或节点不能为空")
return
}

var NodesData []models.Node
for _, nodeName := range strings.Split(nodes, ",") {
if strings.TrimSpace(nodeName) == "" {
continue
}
firstNode := models.Node{Name: nodeName}
result := models.DB.Model(models.Node{}).Where("name = ?", firstNode.Name).First(&firstNode)
if result.Error != nil {
log.Println(result.Error)
utils.Error(c, "40002", result.Error.Error())
return
}
NodesData = append(NodesData, firstNode)
}

oldSub := models.Subcription{Name: oldName}
newSub := models.Subcription{
Name:      newName,
Config:    configs,
NodeOrder: nodes,
Nodes:     NodesData,
}

if err := oldSub.Update(&newSub); err != nil {
utils.Error(c, "40003", "更新订阅失败："+err.Error())
return
}
utils.Success(c, nil, "更新订阅成功")
}

func SubDel(c *gin.Context) {
var sub models.Subcription
id := c.Query("id")
if id == "" {
utils.Error(c, "40001", "id 不能为空")
return
}
x, err := strconv.Atoi(id)
if err != nil {
utils.Error(c, "40002", "无效的 ID: "+err.Error())
return
}
sub.ID = x
if err := sub.Find(); err != nil {
utils.Error(c, "40003", "查找订阅失败："+err.Error())
return
}
if err := sub.Del(); err != nil {
utils.Error(c, "40004", "删除订阅失败："+err.Error())
return
}
utils.Success(c, nil, "删除订阅成功")
}
