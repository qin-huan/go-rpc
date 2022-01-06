package main

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
)

const (
	UserAccountsTableName = "user_accounts"
	GrantType             = "client_credential"
	AppId                 = "wxe6470960bc87cc21"
	Secret                = "1c5be9373cceb13998f90a384899c71d"
)

type UserAccounts struct {
	UserId          int64  `xorm:"int(11) pk  autoincr  notnull 'user_id'"`
	SubscribeWechat int64  `xorm:"int(11) null 'subscribe_wechat'"`
	WechatOpenId    string `xorm:"varchar(191) null 'wechat_open_id'"`
	WechatUnionId   string `xorm:"varchar(255) null 'wechat_union_id'"`
	Mobile          string `xorm:"varchar(191) null 'mobile'"`
	PlatformId      int64  `xorm:"int(11) null 'platform_id'"`
	DeviceNo        string `xorm:"varchar(255) null 'device_no'"`
	JpushNo         string `xorm:"varchar(255) null 'jpush_no'"`
	PlatformName    string `xorm:"varchar(255) null 'platform_name'"`
	BrandName       string `xorm:"varchar(255) null 'brand_name'"`
	BrandDevice     string `xorm:"varchar(255) null 'brand_device'"`
	BrandDeviceNo   string `xorm:"varchar(255) null 'brand_device_no'"`
	BackupOpenId    string `xorm:"longtext null 'backup_open_id'"`
}

func main() {
	engine, err := xorm.NewEngine("mysql", "user_dev:MysqlDev2020^ddd1dad@tcp(rm-uf66387h69xbq0heq.mysql.rds.aliyuncs.com:3306)/hanfu")
	if err != nil {
		panic(err)
	}

	sum, err := engine.Table(UserAccountsTableName).Count()
	if err != nil {
		panic(err)
	}

	var (
		limit  = 100
		offset = 0
	)

	token, err := getAccessToken(GrantType, AppId, Secret)
	if err != nil {
		panic(err)
	}
	for int64(offset) < sum {
		var userAccountList = make([]*UserAccounts, 0)
		if err := engine.Table(UserAccountsTableName).Limit(limit, offset).Asc("user_id").Find(&userAccountList); err != nil {
			panic(err)
		}

		var (
			count     = 0
			num       = 0
			minFirst  = int(math.Min(float64(len(userAccountList)), 25))
			minSecond = int(math.Min(float64(len(userAccountList)), 50))
			minThird  = int(math.Min(float64(len(userAccountList)), 75))
			minFourth = int(math.Min(float64(len(userAccountList)), 100))
		)

		var wait sync.WaitGroup
		wait.Add(1)

		go func() {
			n, m := updateUnionId(userAccountList[0:minFirst], token, engine)
			num += n
			count += m
			wait.Done()
		}()

		if len(userAccountList) >= 25 {
			wait.Add(1)
			go func() {
				n, m := updateUnionId(userAccountList[25:minSecond], token, engine)
				num += n
				count += m
				wait.Done()
			}()
		}

		if len(userAccountList) >= 50 {
			wait.Add(1)
			go func() {
				n, m := updateUnionId(userAccountList[50:minThird], token, engine)
				num += n
				count += m
				wait.Done()
			}()
		}

		if len(userAccountList) >= 75 {
			wait.Add(1)
			go func() {
				n, m := updateUnionId(userAccountList[75:minFourth], token, engine)
				num += n
				count += m
				wait.Done()
			}()
		}

		wait.Wait()

		//for _, userAccount := range userAccountList {
		//	if userAccount.WechatOpenId != "" && userAccount.WechatUnionId == "" {
		//		unionId, err := getUnionId(token, userAccount.WechatOpenId)
		//		if err != nil {
		//			panic(err)
		//		}
		//		if unionId != "" {
		//			newUserAccount := &UserAccounts{
		//				WechatUnionId: unionId,
		//			}
		//			lines, err := engine.Table(UserAccountsTableName).Where("user_id = ?", userAccount.UserId).Cols("wechat_union_id").Update(newUserAccount)
		//			if err != nil {
		//				panic(err)
		//			}
		//			count += int(lines)
		//		}
		//		num++
		//	}
		//}
		fmt.Printf("update userAccounts successfully, sum: %v, num: %v, count: %v\n", len(userAccountList), num, count)
		offset += len(userAccountList)

		if offset % 10000 == 0 {
			token, err = getAccessToken(GrantType, AppId, Secret)
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("done!")
}

func updateUnionId(accounts []*UserAccounts, token string, engine *xorm.Engine) (int, int) {
	var num, count int
	for _, userAccount := range accounts {
		if userAccount.WechatOpenId != "" && userAccount.WechatUnionId == "" {
			unionId, err := getUnionId(token, userAccount.WechatOpenId)
			if err != nil {
				panic(err)
			}
			if unionId != "" {
				newUserAccount := &UserAccounts{
					WechatUnionId: unionId,
				}
				lines, err := engine.Table(UserAccountsTableName).Where("user_id = ?", userAccount.UserId).Cols("wechat_union_id").Update(newUserAccount)
				if err != nil {
					panic(err)
				}
				count += int(lines)
			}
			num++
		}
	}
	return num, count
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

func getAccessToken(grantType, appId, secret string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=%s&appid=%s&secret=%s", grantType, appId, secret)
	get, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer get.Body.Close()

	var resp AccessTokenResponse
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(all, &resp); err != nil {
		return "", err
	}

	if resp.ErrCode != 0 {
		return "", errors.New(resp.ErrMsg)
	}

	return resp.AccessToken, nil
}

type UnionInfo struct {
	Subscribe      int           `json:"subscribe"`
	Openid         string        `json:"openid"`
	Nickname       string        `json:"nickname"`
	Sex            int           `json:"sex"`
	Language       string        `json:"language"`
	City           string        `json:"city"`
	Province       string        `json:"province"`
	Country        string        `json:"country"`
	Headimgurl     string        `json:"headimgurl"`
	SubscribeTime  int           `json:"subscribe_time"`
	Unionid        string        `json:"unionid"`
	Remark         string        `json:"remark"`
	Groupid        int           `json:"groupid"`
	TagidList      []interface{} `json:"tagid_list"`
	SubscribeScene string        `json:"subscribe_scene"`
	QrScene        int           `json:"qr_scene"`
	QrSceneStr     string        `json:"qr_scene_str"`
}

func getUnionId(accessToken, openId string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN参数说明", accessToken, openId)
	get, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer get.Body.Close()

	var resp UnionInfo
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(all, &resp); err != nil {
		return "", err
	}

	return resp.Unionid, nil
}
