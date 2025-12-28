package services

type ServiceRegistry struct {
	CaptchaService    *CaptchaService
	JWTService        *JWTService
	RandomService     *RandomService
	PasswordService   *PasswordService
	MiddlewareService *MiddlewareService
}
