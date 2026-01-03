package questions

type PostQuestionBody struct {
	ProductID uint   `json:"product_id" binding:"required,gt=0,number"`
	Content   string `json:"content" binding:"required,min=2"`
}

type PutQuestionBody struct {
	Answer string `json:"answer" binding:"required,min=2"`
}
