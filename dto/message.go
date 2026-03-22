package dto

type GetMessageListReqDto struct {
	UserOneId string `form:"user1"`
	UserTwoId string `form:"user2"`
}

type GetGroupMessageListReqDto struct {
	GroupId string `form:"group_id"`
}
