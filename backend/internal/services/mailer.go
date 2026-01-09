package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
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
	userRepo     *repositories.UserRepository
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
	deniedBidTemplate = `
<!DOCTYPE html>
<html>
<head>
 <meta charset="UTF-8">
</head>
<body style="font-family: sans-serif;">
 <h2>Your bid has been denied!</h2>

 <p>
 On the product "<strong>%s</strong>"
 </p>

 <hr />

	<a href="%s">Link to product</a>

 <p>
Your bid has been denied by the seller. You may no longer participate in this auction.
 </p>

 <hr />

 <p style="color: #666; font-size: 12px;">
	This mail is automated, do not reply.
 </p>
</body>
</html>`
	auctionEndedTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
</head>
<body style="font-family: sans-serif;">
  <h2>Auction has ended!</h2>

  <p>
		The product "<strong>%s</strong>" has ended with a winner!
  </p>

  <hr />

	<a href="%s">Link to product</a>

  <p>
		Winner: <strong>%s</strong> (at <strong>$%.2f</strong>)
  </p>

	<p>
		Seller: <strong>%s</strong>
	</p>

  <hr />

  <p style="color: #666; font-size: 12px;">
	This mail is automated, do not reply.
  </p>
</body>
</html>`
	auctionExpiredTemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
</head>
<body style="font-family: sans-serif;">
  <h2>Auction has expired.</h2>

  <p>
		The product "<strong>%s</strong>" has expired with no winners.
  </p>

  <hr />

	<a href="%s">Link to product</a>

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
	userRepo *repositories.UserRepository,
) *MailerService {
	return &MailerService{
		cfg:          config,
		mailer:       mailer,
		productRepo:  productRepo,
		userRepo:     userRepo,
		questionRepo: questionRepo,
	}
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

// SendOTPEmail sends an email containing an OTP code.
func (s *MailerService) SendOTPEmail(user *models.User, otp string) {
	go func() {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := fmt.Sprintf(otpEmailTemplate, otp, "15")

		message := gomail.NewMessage()
		message.SetHeader("From", fromHeader)
		message.SetAddressHeader("To", *user.Email, *user.Name)
		message.SetBody("text/html", body)
		message.SetHeader("Subject", "CherryAuctions - OTP Verification")

		if err := s.mailer.DialAndSend(message); err != nil {
			log.Printf("failed to send otp email: %v", err)
		}
	}()
}

func (s *MailerService) SendAuctionExpiredEmail(ctx context.Context, product *models.Product) {
	fmt.Println("Sending expired for", product.ID)
	url := fmt.Sprintf("%s/products/%d", s.cfg.CORS.Origins, product.ID)
	body := fmt.Sprintf(auctionExpiredTemplate, product.Name, url)

	message := gomail.NewMessage()
	message.SetHeader("From", fromHeader)
	message.SetHeader("To", *product.Seller.Email)
	message.SetHeader("Subject", "CherryAuctions - Auction Expired")
	message.SetBody("text/html", body)

	if err := s.mailer.DialAndSend(message); err != nil {
		log.Printf("failed to send expired email: %v", err)
	}

	// Mail complete! Now we update the DB, hopefully
	_, err := s.productRepo.SetProductSentEmail(ctx, product.ID)
	if err != nil {
		log.Printf("failed to mark email as sent: %v", err)
	}
}

func (s *MailerService) SendAuctionEndedEmail(ctx context.Context, product *models.Product) {
	fmt.Println("Sending ended for", product.ID)
	url := fmt.Sprintf("%s/products/%d", s.cfg.CORS.Origins, product.ID)
	body := fmt.Sprintf(
		auctionEndedTemplate,
		product.Name,
		url,
		*product.CurrentHighestBid.User.Name,
		float64(product.CurrentHighestBid.Price)/100,
		*product.Seller.Name,
	)

	message := gomail.NewMessage()
	message.SetHeader("From", fromHeader)
	message.SetHeader("To", *product.CurrentHighestBid.User.Email)
	message.SetHeader("Bcc", *product.Seller.Email)
	message.SetHeader("Subject", "CherryAuctions - Auction Ended")
	message.SetBody("text/html", body)

	if err := s.mailer.DialAndSend(message); err != nil {
		log.Printf("failed to send ended email: %v", err)
	}

	// Mail complete! Now we update the DB, hopefully
	_, err := s.productRepo.SetProductSentEmail(ctx, product.ID)
	if err != nil {
		log.Printf("failed to mark email as sent: %v", err)
	}
}

// SendEndedAuctionsEmail sends an email to all auctions ended without an email sent yet.
func (s *MailerService) SendEndedAuctionsEmail() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		fmt.Printf("Starting the sweep at %v!\n", time.Now())

		// Send all emails
		products, err := s.productRepo.GetAllExpiredProducts(ctx)
		if err != nil {
			fmt.Printf("Sweep had an error: %v\n", err)
		}
		fmt.Printf("About to sweep %d products\n", len(products))

		group, gctx := errgroup.WithContext(ctx)
		for _, product := range products {
			p := product

			group.Go(func() error {
				if product.CurrentHighestBid != nil {
					s.SendAuctionEndedEmail(gctx, &p)
				} else {
					s.SendAuctionExpiredEmail(gctx, &p)
				}
				return nil
			})
		}

		if err := group.Wait(); err != nil {
			fmt.Printf("Group finished with error: %v\n", err)
		}

		defer fmt.Printf("Finishing the sweep at %v!\n", time.Now())
	}()
}

// SendDeniedBidEmail sends an email out when there's a bidder denied from an auction.
// I honestly don't know of any structure to send to this service, but I just need the minimum
// that fits with the actual handler (since the handler didn't fetch the user who was denied).
func (s *MailerService) SendDeniedBidEmail(product *models.Product, deniedID uint) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Don't need to send to non-existent user
		user, err := s.userRepo.GetUserByID(ctx, deniedID)
		if err != nil || user.Email == nil {
			return
		}

		url := fmt.Sprintf("%s/products/%d", s.cfg.CORS.Origins, product.ID)
		body := fmt.Sprintf(deniedBidTemplate, product.Name, url)

		msg := gomail.NewMessage()
		msg.SetHeader("From", fromHeader)
		msg.SetAddressHeader("To", *user.Email, *user.Name)
		msg.SetHeader("Subject", "CherryAuctions - Bid Denied")
		msg.SetBody("text/html", body)

		if err := s.mailer.DialAndSend(msg); err != nil {
			log.Printf("failed to send bid denied email: %v", err)
		}
	}()
}
