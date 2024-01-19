// Code generated by goctl. DO NOT EDIT.
package types

type FeedRequest struct {
	LatestTime uint   `form:"latest_time"`
	Token      string `form:"token,optional"`
	UserId     uint   `form:"user_id,optional"`
}

type VideoListRequest struct {
	UserId uint `json:"userId,optional"`
}

type FeedResponse struct {
	NextTime uint64 `json:"next_time"`
	BaseResponse
	VideoList []VideoListRes `json:"video_list"`
}

type VideoListResponse struct {
	BaseResponse
	VideoList []VideoListRes `json:"video_list"`
}

type VideoListRes struct {
	ID            int       `json:"id"`
	Author        AuthorRes `json:"author"`
	PlayURL       string    `json:"play_url"`
	CoverURL      string    `json:"cover_url"`
	IsFavorite    bool      `json:"is_favorite"`
	CommentCount  int       `json:"comment_count"`
	FavoriteCount int64     `json:"favorite_count"`
}

type AuthorRes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserInfoRequest struct {
	Token string `form:"token"`
}

type UserInfoResponse struct {
	BaseResponse
	User UserInfo `json:"user"`
}

type UserInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserLoginRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type UserLoginResponse struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
	BaseResponse
}

type UserRegisterRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type UserRegisterResponse struct {
	BaseResponse
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type BaseResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type GetFriendRequest struct {
	Token  string `form:"token"`
	UserId uint   `form:"user_id"`
}

type GetFriendResponse struct {
	StatusCode uint   `json:"status_code"`
	UserList   []User `json:"user_list"`
}

type GetMessageRequest struct {
	Token    string `form:"token"`
	ToUserId uint   `form:"to_user_id"`
	UserId   uint   `form:"user_id,optional"`
}

type GetMessageResponse struct {
	StatusCode  int       `json:"status_code"`
	MessageList []Message `json:"message_list"`
}

type Message struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"create_time"`
}

type PostMessageRequest struct {
	Token      string `form:"token"`
	ToUserId   uint   `form:"to_user_id"`
	ActionType uint   `form:"action_type"`
	Content    string `form:"content"`
	UserId     uint   `form:"user_id,optional"`
}

type PostMessageResponse struct {
	BaseResponse
}

type GetFollowingRequest struct {
	Token  string `form:"token"`
	UserId uint   `form:"user_id"`
}

type GetFollowingResponse struct {
	StatusCode uint   `json:"status_code"`
	UserList   []User `json:"user_list"`
}

type GetFansRequest struct {
	Token  string `form:"token"`
	UserId uint   `form:"user_id"`
}

type GetFansResponse struct {
	StatusCode uint   `json:"status_code"`
	UserList   []User `json:"user_list"`
}

type PostFollowRequest struct {
	Token      string `form:"token"`
	ToUserId   uint   `form:"to_user_id"`
	ActionType uint   `form:"action_type"`
	UserId     uint   `form:"user_id,optional"`
}

type PostCommentRequest struct {
	Token       string `form:"token"`
	ActionType  uint   `form:"action_type"`
	VideoId     uint   `form:"video_id"`
	CommentText string `form:"comment_text"`
	CommentId   uint   `form:"comment_id,optional"`
	UserName    string `form:"username,optional"`
	UserId      uint   `form:"user_id,optional"`
}

type PostCommentResponse struct {
	StatusCode uint    `json:"status_code"`
	Comment    Comment `json:"comment"`
}

type GetCommentRequest struct {
	Token   string `form:"token"`
	VideoId uint   `form:"video_id"`
}

type GetCommentResponse struct {
	StatusCode  uint      `json:"status_code"`
	CommentList []Comment `json:"comment_list"`
}

type Comment struct {
	Id        uint   `json:"id"`
	User      User   `json:"user"`
	Content   string `json:"content"`
	CreatedAt string `json:"create_date"`
}

type User struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type GetLikeListResponse struct {
	BaseResponse
	VideoList []VideoListRes `json:"video_list"`
}

type GetLikeListRequest struct {
	UserId uint   `form:"user_id"`
	Token  string `form:"token"`
}

type LikeRequest struct {
	Token      string `form:"token"`
	VideoId    uint   `form:"video_id"`
	ActionType uint   `form:"action_type"`
	UserId     uint   `form:"user_id,optional"`
}

type PublishRequest struct {
	Token   string `form:"token"`
	Title   string `form:"title"`
	PlayURL string `form:"playUrl,optional"`
}

type PublishResponse struct {
	BaseResponse
}
