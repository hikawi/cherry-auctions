package repositories

type RepositoryRegistry struct {
	CategoryRepository     *CategoryRepository
	UserRepository         *UserRepository
	RoleRepository         *RoleRepository
	RefreshTokenRepository *RefreshTokenRepository
	ProductRepository      *ProductRepository
	QuestionRepository     *QuestionRepository
	ChatSessionRepository  *ChatSessionRepository
	TransactionRepository  *TransactionRepository
	RatingRepostory        *RatingRepostory
}
