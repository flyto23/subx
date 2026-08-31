package api

import (
"errors"
"fmt"
"log"
"net/url"
"strconv"
"strings"
"sublink/models"
"sublink/node"
"sublink/utils"

"github.com/gin-gonic/gin"
"gorm.io/gorm"
)

func DocodeNodeName(nd *models.Node) (models.Node, error) {
if nd.Name == "" {
u, err := url.Parse(nd.Link)
if err != nil {
return *nd, err
}
switch u.Scheme {
case "ss":
ss, err := node.DecodeSSURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = ss.Name
case "ssr":
ssr, err := node.DecodeSSRURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = ssr.Qurey.Remarks
case "trojan":
trojan, err := node.DecodeTrojanURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = trojan.Name
case "vmess":
vmess, err := node.DecodeVMESSURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = vmess.Ps
case "vless":
vless, err := node.DecodeVLESSURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = vless.Name
case "hy", "hysteria":
hy, err := node.DecodeHYURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = hy.Name
case "hy2", "hysteria2":
hy2, err := node.DecodeHY2URL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = hy2.Name
case "tuic":
tuic, err := node.DecodeTuicURL(nd.Link)
if err != nil {
return *nd, err
}
nd.Name = tuic.Name
}
}
return *nd, nil
}

func NodeUpdadte(c *gin.Context) {
newName := c.PostForm("name")
newLink := c.PostForm("link")
id := c.PostForm("id")
group := c.PostForm("group")
groups := strings.Split(group, ",")

index, err := strconv.Atoi(id)
if err != nil {
utils.Error(c, "40001", "id 不能为空或者格式不正确")
return
}
if newName == "" || newLink == "" {
utils.Error(c, "40002", "节点名称或链接不能为空")
return
}

oldNode := &models.Node{ID: index}
newNode := &models.Node{Name: newName, Link: newLink}

var gns []models.GroupNode
if len(groups) > 0 {
for _, g := range groups {
if strings.TrimSpace(g) == "" {
continue
}
gns = append(gns, models.GroupNode{Name: strings.TrimSpace(g)})
}
}

if err := oldNode.UpdateGroup(gns); err != nil {
utils.Error(c, "40003", fmt.Sprintf("更新失败：%s", err.Error()))
return
}
if err := oldNode.UpdateNode(newNode); err != nil {
utils.Error(c, "40003", fmt.Sprintf("更新失败：%s", err.Error()))
return
}
utils.Success(c, nil, "更新成功")
}

func NodeGet(c *gin.Context) {
ns, err := models.GetNodeList()
if err != nil {
utils.ErrorWithStatus(c, 500, "50001", "node list error")
return
}
utils.Success(c, ns, "node get")
}

func GroupNodeGet(c *gin.Context) {
gns, err := models.GetGroupNodeList()
if err != nil {
utils.Error(c, "40001", err.Error())
return
}
var data []string
for _, g := range gns {
data = append(data, g.Name)
}
utils.Success(c, data, "GroupNode get")
}

func GroupNodeSet(c *gin.Context) {
var gns []models.GroupNode
var firstGroup models.GroupNode
name := c.PostForm("name")
group := c.PostForm("group")
groups := strings.Split(group, ",")

if len(groups) == 0 {
utils.Error(c, "40001", "分组不能为空")
return
}
log.Println("分组列表:", groups, "数组长度", len(groups))

for _, g := range groups {
if strings.TrimSpace(g) == "" {
log.Println("分组名为空，跳过")
continue
}
log.Println("分组名:", g)
firstGroup.Name = g
if err := firstGroup.Add(); err != nil {
log.Println("添加分组失败:", err)
utils.Error(c, "40002", err.Error())
return
}
result := models.DB.Model(models.GroupNode{}).Where("name = ?", g).First(&firstGroup)
log.Println("FirstGroup", firstGroup)
if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
log.Println(result.Error)
utils.Error(c, "40003", result.Error.Error())
return
}
gns = append(gns, firstGroup)
}

n := models.Node{Name: name}
if err := n.UpdateGroup(gns); err != nil {
utils.Error(c, "40004", err.Error())
return
}
utils.Success(c, nil, "更新关联分组成功")
}

func NodeAdd(c *gin.Context) {
link := c.PostForm("link")
name := c.PostForm("name")
group := c.PostForm("group")

if link == "" || !strings.Contains(link, "://") {
utils.Error(c, "40001", "link 不能为空或者格式不正确，请检查链接是否包含协议头")
return
}

n := models.Node{Name: name, Link: link}
n, err := DocodeNodeName(&n)
if err != nil {
log.Println("解码节点名称错误:", err)
utils.Error(c, "40002", "解码节点名称错误")
return
}

if err := n.Add(); err != nil {
log.Println("添加节点失败:", err)
utils.Error(c, "40003", err.Error())
return
}

if strings.TrimSpace(group) != "" {
groups := strings.Split(group, ",")
for _, g := range groups {
gn := &models.GroupNode{Name: g}
if err := gn.Add(); err != nil {
log.Println(err)
utils.Error(c, "40004", err.Error())
return
}
if err := gn.Ass(&n); err != nil {
log.Println("关联失败:", err)
utils.Error(c, "40005", err.Error())
return
}
}
}
utils.Success(c, nil, "添加成功")
}

func NodeDel(c *gin.Context) {
var n models.Node
id := c.Query("id")
if id == "" {
utils.Error(c, "40001", "id 不能为空")
return
}
x, err := strconv.Atoi(id)
if err != nil {
utils.Error(c, "40002", "无效的 ID")
return
}
n.ID = x
if err := n.Del(); err != nil {
utils.Error(c, "40003", "删除失败")
return
}
utils.Success(c, nil, "删除成功")
}

func NodesGroup(c *gin.Context) {
var gn models.GroupNode
id := c.Query("id")
if id == "" {
utils.Error(c, "40001", "id 不能为空")
return
}
x, err := strconv.Atoi(id)
if err != nil {
utils.Error(c, "40002", "无效的 ID")
return
}
gn.ID = x
if err := gn.Del(); err != nil {
utils.Error(c, "40003", "删除失败")
return
}
utils.Success(c, nil, "删除成功")
}

func NodesTotal(c *gin.Context) {
nodes, err := models.GetNodeList()
if err != nil {
utils.ErrorWithStatus(c, 500, "50001", "获取不到节点统计")
return
}
utils.Success(c, len(nodes), "取得节点统计")
}
