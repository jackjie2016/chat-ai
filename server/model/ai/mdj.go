// 自动生成模板Mdj
package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// Mdj 结构体
type Mdj struct {
	global.GVA_MODEL
	User_id     *int   `json:"user_id" form:"user_id" gorm:"column:user_id;comment:;"`
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	ProxyURL    string `json:"proxy_url"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	ContentType string `json:"content_type"`
}

// TableName Mdj 表名
func (Mdj) TableName() string {
	return "mdj"
}

type MdjMsg struct {
	ID         string `json:"id" bson:"_id" `
	PID        string `json:"pid" bson:"pid" `
	RootHash   string `json:"rootHash" bson:"rootHash" `
	Status     bool   `json:"status" bson:"status" `
	Qustion    string `json:"qustion" bson:"qustion" `
	ChangeType string `json:"changeType" bosn:"changeType"`

	MdjID       string `json:"mdjID" bson:"mdjID" `
	ParentID    string `json:"parentID" bson:"parentID" `
	Index       int    `json:"index" bson:"index" `
	UserId      int    `json:"userId" bson:"userId" `
	AttachID    string `json:"attachID" bson:"attachID"`
	Type        int    `json:"type" bosn:"type"`
	Content     string `json:"content" bosn:"content"`
	Filename    string `json:"filename" bson:"filename"`
	Size        int    `json:"size" bson:"size"`
	OrgURL      string `json:"orgUrl" bson:"orgUrl"`
	URL         string `json:"url" bson:"URL"`
	ProxyURL    string `json:"proxy_url" bson:"proxyURL"`
	OssURL      string `json:"ossURL" bson:"ossURL"`
	Width       int    `json:"width" bson:"width"`
	Height      int    `json:"height" bson:"height"`
	ContentType string `json:"content_type" bson:"contentType"`
	CreateTime  string `json:"create_time" bson:"create_time"`
	UpdateTime  string `json:"update_time" bson:"update_time"`
	Change      Change `json:"change" bson:"change"`
}

type MdjMsg2 struct {
	ID         string `json:"id" bson:"_id" `
	HashID     string `json:"hashid" bson:"hashid" ` //程序生成
	ChannelId  string `json:"channel_id" bson:"channel_id" `
	GuildId    string `json:"guild_id" bson:"guild_id" `
	MsgId      string `json:"msg_id" bson:"msg_id" `
	AuthorId   string `json:"author_id" bson:"author_id" `
	ChangeType string `json:"changeType" bosn:"changeType"`

	ParentID string `json:"parentID" bson:"parentID" `
	Index    int    `json:"index" bson:"index" `
	UserId   int    `json:"userId" bson:"userId" `

	AttachID    string      `json:"attachID" bson:"attachID"`
	Type        int         `json:"type" bosn:"type"`
	Content     string      `json:"content" bosn:"content"`
	Filename    string      `json:"filename" bson:"filename"`
	Embeds      interface{} `json:"embeds" bson:"embeds"`
	Size        int         `json:"size" bson:"size"`
	OrgURL      string      `json:"orgUrl" bson:"orgUrl"`
	URL         string      `json:"url" bson:"URL"`
	ImgHash     string      `json:"img_hash" bson:"img_hash"`
	ProxyURL    string      `json:"proxy_url" bson:"proxyURL"`
	OssURL      string      `json:"ossURL" bson:"ossURL"`
	Width       int         `json:"width" bson:"width"`
	Height      int         `json:"height" bson:"height"`
	ContentType string      `json:"content_type" bson:"contentType"`
	CreateTime  string      `json:"create_time" bson:"create_time"`
	UpdateTime  string      `json:"update_time" bson:"update_time"`
	IsRead      bool        `json:"isRead" bson:"isRead"`
}

type ChatMsg struct {
	HashId     string `json:"hashId" bson:"hashId" `
	UserId     int    `json:"userId" bson:"userId" `
	IsSender   bool   `json:"isSender" bson:"isSender" `
	Content    string `json:"content" bson:"content" `
	Type       string `json:"type" bson:"type" `
	CreateTime string `json:"create_time" bson:"create_time"`
	UpdateTime string `json:"update_time" bson:"update_time"`
}

type Change struct {
	Upscales   Upscales   `json:"upscales" bson:"upscales" `
	Variations Variations `json:"variations" bson:"variations"`
}
type Variations struct {
	V1 bool `json:"v1" bosn:"v1"`
	V2 bool `json:"v2" bosn:"v2"`
	V3 bool `json:"v3" bosn:"v3"`
	V4 bool `json:"v4" bosn:"v4"`
}

type Upscales struct {
	U1 bool `json:"u1" bosn:"u1"`
	U2 bool `json:"u2" bosn:"u2"`
	U3 bool `json:"u3" bosn:"u3"`
	U4 bool `json:"u4" bosn:"u4"`
}
type NewMsgs struct {
	MdjMsg
	Upscales []*Upscale `json:"upscales"`
}
type MidMessages struct {
	ID               string           `json:"id" bson:"_id" `
	Type             int              `json:"type" bosn:"type"`
	Content          string           `json:"content" bosn:"content"`
	ChannelID        string           `json:"channel_id" bosn:"channelID"`
	Author           Author           `json:"author" bosn:"author"`
	Attachments      []Attachments    `json:"attachments" bosn:"attachments"`
	Embeds           []interface{}    `json:"embeds" bosn:"embeds"`
	Mentions         []interface{}    `json:"mentions" bosn:"mentions"`
	MentionRoles     []interface{}    `json:"mention_roles" bosn:"mentionRoles"`
	MessageReference MessageReference `json:"messageReference" bosn:"messageReference"`
	Pinned           bool             `json:"pinned" bosn:"pinned"`
	MentionEveryone  bool             `json:"mention_everyone" bosn:"mentionEveryone"`
	Tts              bool             `json:"tts" bosn:"tts"`
	Timestamp        time.Time        `json:"timestamp" bosn:"timestamp"`
	EditedTimestamp  interface{}      `json:"edited_timestamp" bosn:"editedTimestamp"`
	Flags            int              `json:"flags" bosn:"flags"`
	Components       []interface{}    `json:"components" bosn:"components"`
}
type MessageReference struct {
	ChannelId string `json:"channelId,omitempty"`
	GuildId   string `json:"guild_id,omitempty"`
	MessageId string `json:"message_id,omitempty"`
}
type Author struct {
}
type Attachments struct {
	ID          string `json:"attachID" bson:"attachID"`
	Filename    string `json:"filename" bson:"filename"`
	Size        int    `json:"size" bson:"size"`
	URL         string `json:"url" bson:"URL"`
	ProxyURL    string `json:"proxy_url" bson:"proxyURL"`
	Width       int    `json:"width" bson:"width"`
	Height      int    `json:"height" bson:"height"`
	ContentType string `json:"content_type" bson:"contentType"`
}

type UpscaleRequest struct {
	Type  string        `json:"type"`
	Param *UpscaleParam `json:"param"`
}
type UpscaleParam struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Index    int    `json:"index"`
	URL      string `json:"url"`
}

type Upscale struct {
	ID       string `json:"id" bson:"_id"`
	MsgId    string `json:"msgId" bson:"msgId"`
	URI      string `json:"uri" bson:"URI"`
	Hash     string `json:"hash" bson:"hash"`
	Content  string `json:"content" bson:"content"`
	Progress string `json:"progress" bson:"progress"`
}
