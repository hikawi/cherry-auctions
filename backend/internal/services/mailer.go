package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/internal/config"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
)

type MailerService struct {
	cfg          *config.Config
	mailer       *gomail.Dialer
	productRepo  *repositories.ProductRepository
	questionRepo *repositories.QuestionRepository
}

const (
	fromHeader            = "CherryAuctions <noreply@luny.dev>"
	questionEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>New Question</title>
</head>
<body style="font-family: sans-serif;">
  <h2>New question on your auction</h2>

  <p>
    For product "%s", there's a question related.
  </p>

	<a href="%s">Link to product</a>

  <hr />

  <p style="white-space: pre-wrap;">
    %s
  </p>

  <hr />

  <p>
	Who asked: %s
  </p>

  <p style="color: #666; font-size: 12px;">
	This mail is automated, do not reply.
  </p>
</body>
</html>
`
	answerEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
</head>
<body style="font-family: sans-serif;">
  <h2>Question was answered</h2>

  <p>
    On a question to the product "<strong>%s</strong>", the seller has posted an answer.
  </p>

  <hr />

	<a href="%s">Link to product</a>

  <p>
    <strong>Q:</strong><br />
    <span style="white-space: pre-wrap;">%s</span>
  </p>

  <p>
    <strong>A:</strong><br />
    <span style="white-space: pre-wrap;">%s</span>
  </p>

  <hr />

  <p style="color: #666; font-size: 12px;">
	This mail is automated, do not reply.
  </p>
</body>
</html>
`
	newBidTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
</head>
<body style="font-family: sans-serif;">
  <h2>New bid has been placed!</h2>

  <p>
		On the product "<strong>%s</strong>"
  </p>

  <hr />

	<a href="%s">Link to product</a>

  <p>
		New bid at <strong>$%.2f</strong> placed by <strong>%s</strong>
  </p>

  <hr />

  <p style="color: #666; font-size: 12px;">
	This mail is automated, do not reply.
  </p>
</body>
</html>`
	otpEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
</head>
<body style="font-family: sans-serif;">
  <h2>One-Time Password Verification</h2>

  <p>
    You have requested to verify your identity.
  </p>

  <p>
    Please use the following one-time password (OTP) to complete the verification:
  </p>

  <hr />

  <p style="font-size: 20px; font-weight: bold; letter-spacing: 2px;">
    %s
  </p>

  <hr />

  <p>
    This code is valid for <strong>%s minutes</strong>.
    If you did not request this verification, please ignore this email.
  </p>

  <p style="color: #666; font-size: 12px;">
    This mail is automated, do not reply.
  </p>
</body>
</html>
`
)

func NewMailerService(
	config *config.Config,
	mailer *gomail.Dialer,
	productRepo *repositories.ProductRepository,
	questionRepo *repositories.QuestionRepository,
) *MailerService {
	return &MailerService{cfg: config, mailer: mailer, productRepo: productRepo, questionRepo: questionRepo}
}

func (s *MailerService) SendQuestionEmail(product *models.Product, question string, user string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		message := gomail.NewMessage()
		message.SetHeader("From", fromHeader)
		message.SetAddressHeader(
			"To",
			*product.Seller.Email,
			*product.Seller.Name,
		)
		message.SetHeader("Subject", "CherryAuctions - New Question")

		url := fmt.Sprintf("%s/products/%d", s.cfg.CORS.Origins, product.ID)

		body := fmt.Sprintf(
			questionEmailTemplate,
			product.Name,
			url,
			question,
			user,
		)

		message.SetBody("text/html", body)

		if err := s.mailer.DialAndSend(message); err != nil {
			log.Printf("failed to send question email: %v", err)
		}

		_ = ctx // gomailはcontext非対応なので存在意義は「時間の意思表示」だけ
	}()
}

// SendAnswerEmail sends an email to all current bidders of the product
func (s *MailerService) SendAnswerEmail(questionID uint) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		question, err := s.questionRepo.GetQuestionByID(ctx, questionID)
		if err != nil {
			log.Printf("failed to load question: %v", err)
			return
		}

		var bidderEmails []string
		for _, bid := range question.Product.Bids {
			if bid.User.Email != nil {
				bidderEmails = append(bidderEmails, *bid.User.Email)
			}
		}

		if len(bidderEmails) == 0 {
			fmt.Println("oh no no emails to send, what the hell!")
			return
		}

		url := fmt.Sprintf("%s/products/%d", s.cfg.CORS.Origins, question.Product.ID)
		message := gomail.NewMessage()
		message.SetHeader("From", fromHeader)

		if question.User.Email != nil {
			message.SetAddressHeader(
				"To",
				*question.User.Email,
				*question.User.Name,
			)
		}

		message.SetHeader("Bcc", bidderEmails...)
		message.SetHeader("Subject", "CherryAuctions - Question Answered")

		body := fmt.Sprintf(
			answerEmailTemplate,
			question.Product.Name,
			url,
			question.Content,
			question.Answer.String,
		)

		message.SetBody("text/html", body)

		if err := s.mailer.DialAndSend(message); err != nil {
			log.Printf("failed to send answer email: %v", err)
		}
	}()
}

// SendBidEmail sends an email to all related bidders to a certain bid,
// which includes the previous bidder, if there is.
func (s *MailerService) SendBidEmail(lastBid *models.Bid, newBid *models.Bid, product *models.Product) {
	go func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var emails []string
		if lastBid.User.Email != nil {
			emails = append(emails, *lastBid.User.Email)
		}
		emails = append(emails, *newBid.User.Email)
		emails = append(emails, *product.Seller.Email)

		url := fmt.Sprintf("%s/products/%d", s.cfg.CORS.Origins, product.ID)
		message := gomail.NewMessage()
		message.SetHeader("From", fromHeader)

		message.SetAddressHeader(
			"To",
			*newBid.User.Email,
			*newBid.User.Name,
		)

		body := fmt.Sprintf(
			newBidTemplate,
			product.Name,
			url,
			float64(newBid.Price)/100,
			*newBid.User.Name,
		)

		message.SetBody("text/html", body)

		message.SetHeader("Bcc", emails...)
		message.SetHeader("Subject", "CherryAuctions - New Bid")

		if err := s.mailer.DialAndSend(message); err != nil {
			log.Printf("failed to send bid email: %v", err)
		}
	}()
}

func (s *MailerService) SendOTPEmail(user *models.User, otp string) {
	go func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := fmt.Sprintf(otpEmailTemplate, "", "15")

		message := gomail.NewMessage()
		message.SetHeader("From", fromHeader)
		message.SetHeader("To", *user.Email, *user.Name)
		message.SetBody("text/html", body)
		message.SetHeader("Subject", "CherryAuctions - OTP Verification")

		if err := s.mailer.DialAndSend(message); err != nil {
			log.Printf("failed to send otp email: %v", err)
		}
	}()
}
