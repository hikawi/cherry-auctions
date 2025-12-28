package repositories

type RepositoryRegistry struct {
	CategoryRepository     *CategoryRepository
	UserRepository         *UserRepository
	RoleRepository         *RoleRepository
	RefreshTokenRepository *RefreshTokenRepository
}
