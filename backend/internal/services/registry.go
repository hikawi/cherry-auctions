package services

type ServiceRegistry struct {
	CaptchaService    *CaptchaService
	JWTService        *JWTService
	RandomService     *RandomService
	PasswordService   *PasswordService
	MiddlewareService *MiddlewareService
	S3Service         *S3Service
	MailerService     *MailerService
	OTPService        *OTPService
}
