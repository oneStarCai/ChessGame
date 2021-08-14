package ChessGame

import (
	"context"
	"github.com/google/uuid"
	"github.com/orestonce/ChessGame/ent"
	"github.com/orestonce/ChessGame/ent/dsession"
	"github.com/orestonce/ChessGame/ent/duser"
	"github.com/orestonce/ChessGame/ymd/ymdTime"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const (
	ErrUsernameExists          = `用户名已注册`
	ErrUsernameNotExists       = `用户不存在`
	ErrPassword                = `密码错误`
	ErrRoomIdInvalid           = `房间号码不合法`
	ErrGameIsRunning           = `游戏正在运行`
	ErrUserNotEnough           = `玩家不足`
	ErrGameIsNotRunning        = `游戏未开始`
	ErrCannotMove              = `不能这样走`
	ErrNoPosition              = `没有位置了`
	ErrLoginRepeat             = `重复登陆`
	ErrNextLogin               = `你的账号在别处登录`
	ErrEmptyPasswordOrUsername = `账号或密码为空`
	ErrUnknown                 = `未知错误`
)

type PingRequest struct{}
type PingResponse struct{}

func (this *GameUser) RpcPing(req PingRequest) (resp PingResponse) {
	return
}

type RegisterRequest struct {
	Username string
	Password string
}

type RegisterResponse struct {
	ErrMsg string
}

func (this *GameUser) RpcRegister(req RegisterRequest) (resp RegisterResponse) {
	if req.Username == `` || req.Password == `` {
		resp.ErrMsg = ErrEmptyPasswordOrUsername
		return
	}
	exists, err := gDbClient.DUser.Query().Where(duser.Name(req.Username)).Exist(context.Background())
	if err != nil {
		resp.ErrMsg = ErrUnknown
		log.Println("GameUser.RpcRegister", err)
		return resp
	}
	if exists {
		resp.ErrMsg = ErrUsernameExists
		return resp
	}
	hashByte, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("GameUser.RpcRegister", err)
		resp.ErrMsg = ErrUnknown
		return resp
	}
	_, err = gDbClient.DUser.Create().SetID(uuid.NewString()).SetName(req.Username).SetPasswordHash(string(hashByte)).Save(context.Background())
	if err != nil {
		log.Println("GameUser.RpcRegister", err)
		resp.ErrMsg = ErrUnknown
		return resp
	}
	return
}

type LoginRequest struct {
	Username string
	Password string
	RoomId   string
}

type LoginResponse struct {
	UserId string
	ErrMsg string
}

func (this *GameUser) RpcLogin(req LoginRequest) (resp LoginResponse) {
	user, err := gDbClient.DUser.Query().Where(duser.Name(req.Username)).Only(context.Background())
	if err != nil {
		if _, ok := err.(*ent.NotFoundError); ok {
			resp.ErrMsg = ErrUsernameNotExists
			return resp
		}
		log.Println("GameUser.RpcLogin", err)
		resp.ErrMsg = ErrPassword
		return resp
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		resp.ErrMsg = ErrPassword
		return resp
	}
	if req.RoomId == `` {
		resp.ErrMsg = ErrRoomIdInvalid
		return
	}
	if this.Session.UserID != `` || this.Session.RoomID != `` {
		resp.ErrMsg = ErrLoginRepeat
		return
	}
	_, err = gDbClient.DSession.Update().Where(dsession.ID(this.Session.ID)).SetUserID(user.ID).SetRoomID(req.RoomId).Save(context.Background())
	if err != nil {
		log.Println("DSession.Update", err)
		resp.ErrMsg = ErrUnknown
		return
	}
	this.Session.UserID = user.ID
	this.Session.RoomID = req.RoomId
	// 踢掉上一个用户
	{
		for _, oldSession := range getSessionListBy(dsession.UserID(this.Session.UserID), dsession.IDNEQ(this.Session.ID)) {
			sendNoticeToSession(oldSession.ID, ServerKickYou{
				ErrMsg: ErrNextLogin,
			})
			kickSession(oldSession.ID)
		}
	}
	room := getRoomById(this.Session.RoomID)
	if room == nil {
		room = &GameRoom{
			RoomId: this.Session.RoomID,
			Data: &ent.DRoom{
				ID:         this.Session.RoomID,
				CreateTime: time.Now(),
			},
		}
		room.LoadPanelFromData()
		err = gDbClient.DRoom.Create().SetID(this.Session.RoomID).SetCreateTime(room.Data.CreateTime).Exec(context.Background())
		if err != nil {
			log.Println("GameRoomManager.OnUserLogin", err)
		}
	}
	room.sync2Client(nil)
	resp.UserId = user.ID
	return
}

type ChatRequest struct {
	Text string
}

type ChatResponse struct {
	ErrMsg string
}

type ChatMessage struct {
	TimeStr  string
	Username string
	UserId   string
	Text     string
}

func (this *GameRoom) RpcChat(session *ent.DSession, req ChatRequest) (resp ChatResponse) {
	now := time.Now()
	err := gDbClient.DChat.Create().SetID(uuid.NewString()).SetSessionID(session.ID).SetUserID(session.UserID).SetRoomID(session.RoomID).SetText(req.Text).SetCreateTime(now).Exec(context.Background())
	if err != nil {
		log.Println("GameRoom.RpcChat", err)
		resp.ErrMsg = ErrUnknown
		return resp
	}
	this.BroadcastToAll(ChatMessage{
		TimeStr:  ymdTime.DefaultFormat(now),
		UserId:   session.UserID,
		Username: getUserNameByIdIgnoreEmpty(session.UserID),
		Text:     req.Text,
	})
	return resp
}

func (this *GameRoom) getSessionList() []*ent.DSession {
	return getSessionListBy(dsession.RoomID(this.RoomId))
}

type SyncPanelMessage struct {
	PanelFull      string
	WUserName      string
	WUserId        string
	BUserName      string
	BUserId        string
	IsGameRunning  bool
	NextTurnUserId string
	ShowReGame     bool
	ShowSiteDown   bool
}

type TakeSiteRequest struct{}

type TakeSiteResponse struct {
	ErrMsg string
}

func (this *GameRoom) RpcTakeSite(session *ent.DSession, req TakeSiteRequest) (resp TakeSiteResponse) {
	if this.Data.IsGameRunning {
		resp.ErrMsg = ErrGameIsRunning
		return
	}
	userId := session.UserID

	if this.Data.WUserID == `` && this.Data.BUserID != userId {
		this.Data.WUserID = userId
	} else if this.Data.BUserID == `` && this.Data.WUserID != userId {
		this.Data.BUserID = userId
	} else {
		resp.ErrMsg = ErrNoPosition
		return
	}
	if this.Data.BUserID != `` && this.Data.WUserID != `` {
		this.Data.IsGameRunning = true
		this.Data.NextTurnUserID = this.Data.WUserID
		this.Data.Panel = ""
		this.LoadPanelFromData()
	}
	this.sync2Client(nil)

	return
}

type PiecePoint struct {
	Line   int32
	Column int32
}

type MovePieceRequest struct {
	FromPoint PiecePoint
	ToPoint   PiecePoint
}

type MovePieceResponse struct {
	ErrMsg string
}

func (this *GameRoom) RpcMovePiece(session *ent.DSession, req MovePieceRequest) (resp MovePieceResponse) {
	if !this.Data.IsGameRunning {
		resp.ErrMsg = ErrGameIsNotRunning
		return
	}
	if !this.CanMove(session, req.FromPoint, req.ToPoint) {
		resp.ErrMsg = ErrCannotMove
		return
	}
	this.DoMove(req.FromPoint, req.ToPoint)
	return
}

type GetSuggestionRequest struct {
	FromPoint PiecePoint
}

type GetSuggestionResponse struct {
	FromPoint     PiecePoint
	CanMoveToList []PiecePoint
}

func (this *GameRoom) RpcGetSuggestion(session *ent.DSession, req GetSuggestionRequest) (resp GetSuggestionResponse) {
	for line := int32(0); line < LINE_COUNT; line++ {
		for column := int32(0); column < COLUMN_COUNT; column++ {
			to := PiecePoint{Line: line, Column: column}
			if this.CanMove(session, req.FromPoint, to) {
				resp.CanMoveToList = append(resp.CanMoveToList, to)
			}
		}
	}
	resp.FromPoint = req.FromPoint
	return
}

type ReGameRequest struct{}
type ReGameResponse struct {
	ErrMsg string
}

func (this *GameRoom) RpcReGame(session *ent.DSession, req ReGameRequest) (resp ReGameResponse) {
	if this.Data.IsGameRunning {
		resp.ErrMsg = ErrGameIsRunning
		return
	}
	if this.Data.BUserID == `` || this.Data.WUserID == `` {
		resp.ErrMsg = ErrUserNotEnough
		return
	}
	this.Data.IsGameRunning = true
	this.Data.Panel = ""
	this.LoadPanelFromData()
	this.Data.NextTurnUserID = this.Data.WUserID
	this.sync2Client(nil)
	return
}

func (this *GameRoom) sync2Client(user *GameUser) {
	var s SyncPanelMessage
	this.formatPanelFull(&s)
	var sessionList []*ent.DSession
	if user != nil {
		sessionList = append(sessionList, user.Session)
	} else {
		sessionList = getSessionListBy(dsession.RoomID(this.RoomId))
	}
	for _, session := range sessionList {
		this.formatShowStatus(&s, session.UserID)
		sendNoticeToSession(session.ID, s)
	}
	this.SaveRoomDataToDb()
}

func (this *PiecePoint) IsValid() bool {
	return 0 <= this.Line && this.Line < 10 && 0 <= this.Column && this.Column < 9
}
