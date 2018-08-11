package App

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//初始化模块
func Run() {
	//用户列队
	go RunUserList()
	//消息列队
	RunMessage()
}

var(
	//在线用户
	UserList []UserListType = []UserListType{}
)

//用户在线列队类型
type UserListType struct{
	//token
	Token string `json:"token"`
	//用户名称
	NiceName string `json:"nice_name"`
	//上次访问时间
	UpdateTime int64 `json:"update_time"`
}

//初始化
func RunUserList(){
	//循环遍历，如果超过24小时无响应，自动清理
	for{
		lastTime := time.Now().Unix() - 24 * 60 * 60
		newUserList := []UserListType{}
		for _,v := range UserList{
			if v.UpdateTime >= lastTime{
				newUserList[len(newUserList)] = v
			}
		}
		UserList = newUserList
		time.Sleep(time.Hour * 3)
	}
}

//获取用户列表
func GetUserList() ([]UserListType){
	return UserList
}

//登记一个新用户
func Login(c *gin.Context,niceName string) (string,error) {
	//检查是否存在该用户？
	for k,v := range UserList{
		if v.NiceName == niceName{
			UserList[k].UpdateTime = time.Now().Unix()
			return v.Token,nil
		}
	}
	//找不到，则创建
	token := c.MustGet("cookie").(string)
	UserList[len(UserList)] = UserListType{
		token,
		niceName,
		time.Now().Unix(),
	}
	return token,nil
}

//检查用户并更新在线信息
func CheckLogin(token string) error{
	//遍历检查用户是否存在？
	for k,v := range UserList{
		if v.Token == token{
			UserList[k].UpdateTime = time.Now().Unix()
			return nil
		}
	}
	return errors.New("cannot find user.")
}

var(
	//消息列队
	Message []MessageType = []MessageType{}
)

//消息列队类型
type MessageType struct{
	//发送人
	SendUserToken string `json:"send_user_token"`
	//收件人
	PostUserToken string `json:"post_user_token"`
	//创建时间
	CreateTime int64 `json:"create_time"`
	//消息内容
	Content string `json:"content"`
}

//清理超出30天数据
func RunMessage(){
	for{
		lastTime := time.Now().Unix() - 30 * 24 * 60 * 60
		newMessage := []MessageType{}
		for _,v := range Message{
			if v.CreateTime >= lastTime{
				newMessage[len(newMessage)] = v
			}
		}
		Message = newMessage
		time.Sleep(time.Hour * 24)
	}
}

//根据发件人获取信息列表
func GetMessageListBySend(token string) []MessageType{
	resMessage := []MessageType{}
	for _,v := range Message{
		if v.SendUserToken == token{
			resMessage[len(resMessage)] = v
		}
	}
	return resMessage
}

//根据收件人获取信息列表
func GetMessageListByPost(token string) []MessageType{
	resMessage := []MessageType{}
	for _,v := range Message{
		if v.PostUserToken == token{
			resMessage[len(resMessage)] = v
		}
	}
	return resMessage
}

//发送一个信息
func SendMessage(sendToken string,postToken string,message string) error{
	Message[len(Message)] = MessageType{
		sendToken,
		postToken,
		time.Now().Unix(),
		message,
	}
	return nil
}