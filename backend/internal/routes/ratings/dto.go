package ratings

type PostRatingBody struct {
	Rating     uint   `json:"rating" binding:"min=0,max=1"`
	Feedback   string `json:"feedback" binding:"required"`
	RevieweeID uint   `json:"reviewee_id" binding:"required"`
}

type PutRatingBody struct {
	Rating   uint   `json:"rating" binding:"min=0,max=1"`
	Feedback string `json:"feedback" binding:"required"`
}
