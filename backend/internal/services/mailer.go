package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"gopkg.in/gomail.v2"
	"luny.dev/cherryauctions/internal/models"
	"luny.dev/cherryauctions/internal/repositories"
)

type MailerService struct {
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
)

const answerEmailTemplate = `
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

func NewMailerService(
	mailer *gomail.Dialer,
	productRepo *repositories.ProductRepository,
	questionRepo *repositories.QuestionRepository,
) *MailerService {
	return &MailerService{mailer: mailer, productRepo: productRepo, questionRepo: questionRepo}
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

		body := fmt.Sprintf(
			questionEmailTemplate,
			product.Name,
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
			return
		}

		message := gomail.NewMessage()
		message.SetHeader("From", fromHeader)

		if question.Product.Seller.Email != nil {
			message.SetAddressHeader(
				"To",
				*question.Product.Seller.Email,
				*question.Product.Seller.Name,
			)
		}

		message.SetHeader("Bcc", bidderEmails...)
		message.SetHeader("Subject", "CherryAuctions - Question Answered")

		body := fmt.Sprintf(
			answerEmailTemplate,
			question.Product.Name,
			question.Content,
			question.Answer.String,
		)

		message.SetBody("text/html", body)

		if err := s.mailer.DialAndSend(message); err != nil {
			log.Printf("failed to send answer email: %v", err)
		}
	}()
}
