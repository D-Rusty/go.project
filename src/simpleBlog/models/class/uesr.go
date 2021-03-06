package class

import (
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"crypto/sha1"
	"crypto/md5"
	"encoding/base64"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"path"
	"context"
)

type User struct {
	UId        string                         //id 自动增长
	UserName   string    `orm:"pk;"`          //用户名
	LogoImgUrl string                         //个人头像
	Describe   string                         //个人描述
	About      string    `orm:"null"`         //兴趣爱好
	Email      string    `orm:"unique"`       //邮箱地址
	Password   string                         //登录密码
	RegTime    time.Time `orm:"auto_now_add"` //用户注册时间
}

/**
 * 插入一条用户记录
 */
func (u User) Insert() (err error) {
	o := orm.NewOrm()
	_, err = o.Insert(&u)
	return
}

/*
 * 添加默认数据
 */
func InitData() {
	u := &User{
		UId:        GenerateRandomUserId(),
		UserName:   "drusty",
		LogoImgUrl: "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1529722737856&di=d5c76addc8b84d0de0457327268e2c1c&imgtype=0&src=http%3A%2F%2Fattachments.gfan.com%2Fforum%2F201509%2F23%2F140610ebbookekbr1qbte1.jpg",
		Describe:   "个人博客",
		About:      "爬山，游泳",
		Email:      "onepice2014@sina.com",
		Password:   PwGen("123456"),
		RegTime:    time.Now(),
	}

	if u.Query() != nil {
		u.Insert()
	}

}

/**
 * 查询用户信息
 */
func (u *User) Query() (err error) {
	o := orm.NewOrm()
	err = o.Read(u)
	return
}

/**
 *更新指定用户信息
 */
func (u User) Update() (err error) {
	o := orm.NewOrm()
	_, err = o.Update(&u)
	return
}

/**
 * 删除用户记录
 */
func (u User) Delete() (err error) {
	err = u.Update()
	return
}

/**
 * 判断该用户是否已经存在
 */
func (u User) ExistId() bool {
	o := orm.NewOrm()
	if err := o.Read(&u); err == orm.ErrNoRows {
		return false
	}
	return true
}

/**
 * 判断该邮箱是否已经注册
 */
func (u User) ExistEmail() bool {
	o := orm.NewOrm()
	return o.QueryTable("user").Filter("Email", u.Email).Exist()
}

/*
 * 生成随机用户id
 */
func GenerateRandomUserId() string {
	return strconv.FormatInt(time.Now().UnixNano()%9000+1000, 10)
}

/**
 * 进行数据加密
 */
func PwGen(pass string) string {
	//依据时间产生一个4位长度的随机字符串
	salt := strconv.FormatInt(time.Now().UnixNano()%9000+1000, 10)
	//1.将密码明文先经过md5加密，2.将md5加密后的密文+随机字符串，3.对字符串进行sha1加密+随机字符串
	return Base64Encode(Sha1(Md5(pass)+salt) + salt)
}

/**
 * Sha1加密
 */
func Sha1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

/**
 * md5加密
 */
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

/**
 * base64 编码
 */
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

/**
 * base64 解码
 */
func Base64Decode(s string) string {
	res, _ := base64.StdEncoding.DecodeString(s)
	return string(res)
}

/**
 * 验证登录密码是否一致
 */
func PwCheck(pwd, saved string) bool {

	saved = Base64Decode(saved)

	if len(saved) < 4 {
		return false
	}

	salt := saved[len(saved)-4:]

	return Sha1(Md5(pwd)+salt)+salt == saved
}

/**
 * 简单的上传文件
 */
func SimeUploadFile(filename string) (err error) {

	path := path.Join("static/img", filename)
	accessKey := "HlE45UT8wRJBPWBb4HIup2dKn33cWcBaq6Wo-jye"
	secretKey := "IqPCJAY-0Q90VX9vF7BNSg2a_uzGlVH8TwvOi_j0"

	localFile := path

	key := filename

	bucket := "drustydatarepo"

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	mac := qbox.NewMac(accessKey, secretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuanan

	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExTtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExTtra)

	return err

}

/**
 * 覆盖已存在的文件
 * conversImgName 需要被覆盖的文件名称
 * uploadingImgName
 *
*/
func CoversimeUploadFile(conversImgName string, uploadingImgName string) (err error) {
	path := path.Join("static/img", uploadingImgName)
	accessKey := "HlE45UT8wRJBPWBb4HIup2dKn33cWcBaq6Wo-jye"
	secretKey := "IqPCJAY-0Q90VX9vF7BNSg2a_uzGlVH8TwvOi_j0"
	localFile := path
	bucket := "drustydatarepo"
	key := conversImgName
	keyToOverwrite := conversImgName

	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", bucket, keyToOverwrite),
	}

	mac := qbox.NewMac(accessKey, secretKey)

	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}

	cfg.Zone = &storage.ZoneHuanan

	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}

	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	err = formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)

	return err
}
