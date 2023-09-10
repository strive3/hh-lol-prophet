package hh_lol_prophet

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/strive3/hh-lol-prophet/conf"
	"github.com/strive3/hh-lol-prophet/global"
	ginApp "github.com/strive3/hh-lol-prophet/pkg/gin"
	"github.com/strive3/hh-lol-prophet/services/db/models"
	"github.com/strive3/hh-lol-prophet/services/lcu"
)

type (
	Api struct {
		p *Prophet
	}
	summonerNameReq struct {
		SummonerName string `json:"summonerName"`
	}
)

func (api Api) ProphetActiveMid(c *gin.Context) {
	app := ginApp.GetApp(c)
	if !api.p.lcuActive {
		app.ErrorMsg("请检查lol客户端是否已启动")
		return
	}
	c.Next()
}
func (api Api) QueryHorseBySummonerName(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &summonerNameReq{}
	if err := c.ShouldBind(d); err != nil {
		app.ValidError(err)
		return
	}
	summonerName := strings.TrimSpace(d.SummonerName)
	var summonerID int64 = 0
	if summonerName == "" {
		if api.p.currSummoner == nil {
			app.ErrorMsg("系统错误")
			return
		}
		summonerID = api.p.currSummoner.SummonerId
	} else {
		info, err := lcu.QuerySummonerByName(summonerName)
		if err != nil || info.SummonerId <= 0 {
			app.ErrorMsg("未查询到召唤师")
			return
		}
		summonerID = info.SummonerId
	}
	scoreInfo, err := GetUserScore(summonerID)
	if err != nil {
		app.CommonError(err)
		return
	}
	scoreCfg := global.GetScoreConf()
	clientCfg := global.GetClientConf()
	var horse string
	for i, v := range scoreCfg.Horse {
		if scoreInfo.Score >= v.Score {
			horse = clientCfg.HorseNameConf[i]
			break
		}
	}
	app.Data(gin.H{
		"score":   scoreInfo.Score,
		"currKDA": scoreInfo.CurrKDA,
		"horse":   horse,
	})
}

func (api Api) CopyHorseMsgToClipBoard(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Success()
}
func (api Api) GetAllConf(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(global.GetClientConf())
}
func (api Api) UpdateClientConf(c *gin.Context) {
	app := ginApp.GetApp(c)
	d := &conf.UpdateClientConfReq{}
	if err := c.ShouldBind(d); err != nil {
		app.ValidError(err)
		return
	}
	cfg := global.SetClientConf(*d)
	bts, _ := json.Marshal(cfg)
	m := models.Config{}
	err := m.Update(models.LocalClientConfKey, string(bts))
	if err != nil {
		app.CommonError(err)
		return
	}
	app.Success()
}
func (api Api) DevHand(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(gin.H{
		"buffge": 23456,
	})
}
func (api Api) GetAppInfo(c *gin.Context) {
	app := ginApp.GetApp(c)
	app.Data(global.AppBuildInfo)
}
func (api Api) GetLcuAuthInfo(c *gin.Context) {
	app := ginApp.GetApp(c)
	port, token, err := lcu.GetLolClientApiInfo()
	if err != nil {
		app.CommonError(err)
		return
	}
	app.Data(gin.H{
		"token": token,
		"port":  port,
	})
}
