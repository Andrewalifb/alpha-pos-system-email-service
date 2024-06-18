package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Andrewalifb/alpha-pos-system-email-service/config"
	"github.com/Andrewalifb/alpha-pos-system-email-service/dto"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	SenderName   string
	AuthEmail    string
	AuthPassword string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SenderName:   os.Getenv("CONFIG_SENDER_NAME"),
		AuthEmail:    os.Getenv("CONFIG_AUTH_EMAIL"),
		AuthPassword: os.Getenv("CONFIG_AUTH_PASSWORD"),
	}
}

func (es *EmailService) SendDigitalReceipt(req dto.DigitalReceipt, userEmail string) error {
	// Prepare the items for the email
	var emailItems string

	// iterate through receipt items for email items
	for _, item := range req.Body.Items {
		emailItem := fmt.Sprintf("<tr class='service'><td class='tableitem'><p class='itemtext'>%s</p></td><td class='tableitem'><p class='itemtext'>%d</p></td><td class='tableitem'><p class='itemtext'>%.2f</p></td><td class='tableitem'><p class='itemtext'>%.2f</p></td></tr>", item.ProductName, item.Quantity, item.Price, item.TotalPrice)
		emailItems += emailItem
	}

	subject := "Digital Receipt"
	message := fmt.Sprintf(`
	<html>
	<head>
	<style>
	#invoice-POS{
	  box-shadow: 0 0 1in -0.25in rgba(0, 0, 0, 0.5);
	  padding:2mm;
	  margin: 0 auto;
	  width: 100mm;
	  background: #FFF;
	}
	/* More CSS */
	</style>
	</head>
	<body>
	<div id="invoice-POS">
		<center id="top">
			<div class="logo"></div>
			<div class="info"> 
				<h2>%s</h2>
			</div>
		</center>
		<div id="mid">
			<div class="info">
				<h2>Contact Info</h2>
				<p> 
					Address : %s<br>
					Cashier : %s<br>
					Receipt ID : %s<br>
					Date : %s<br>
				</p>
			</div>
		</div>
		<div id="bot">
			<div id="table">
				<table>
					<tr class="tabletitle">
						<td class="item"><h2>Item</h2></td>
						<td class="Hours"><h2>Qty</h2></td>
						<td class="Rate"><h2>Price</h2></td>
						<td class="Rate"><h2>Sub Total</h2></td>
					</tr>
					%s
				</table>
			</div>
		</div>
		<div id="summary">
			<div id="table">
				<table>
					<tr class="tabletitle">
						<td></td>
						<td class="Rate"><h2>Sub Total</h2></td>
						<td class="payment"><h2>%.2f</h2></td>
					</tr>
					<tr class="tabletitle">
						<td></td>
						<td class="Rate"><h2>Discount</h2></td>
						<td class="payment"><h2>%.2f</h2></td>
					</tr>
					<tr class="tabletitle">
						<td></td>
						<td class="Rate"><h2>tax</h2></td>
						<td class="payment"><h2>%.2f</h2></td>
					</tr>
					<tr class="tabletitle">
						<td></td>
						<td class="Rate"><h2>Total</h2></td>
						<td class="payment"><h2>%.2f</h2></td>
					</tr>
					<tr class="tabletitle">
						<td></td>
						<td class="Rate"><h2>Cash</h2></td>
						<td class="payment"><h2>%.2f</h2></td>
					</tr>
					<tr class="tabletitle">
						<td></td>
						<td class="Rate"><h2>Change</h2></td>
						<td class="payment"><h2>%.2f</h2></td>
					</tr>
				</table>
			</div>
		</div>
		<div id="legalcopy">
			<p class="legal"><strong>Thank you for your business!</strong>  Payment is expected within 31 days; please process this invoice within that time. There will be a 5%% interest charge per month on late invoices. 
			</p>
		</div>
	</div>
	</body>
	</html>`, req.Header.StoreName, req.Header.StoreAddress, req.Header.CashierName, req.Header.ReceiptID, req.Header.TransactionDateTime, emailItems, req.Summary.SubTotalAmount, req.Summary.DiscountAmount, req.Summary.TaxAmount, req.Summary.TotalAmount, req.Summary.CashAmount, req.Summary.ChangeAmount)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", es.SenderName)
	mailer.SetHeader("To", userEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	config_smtp_host := os.Getenv("CONFIG_SMTP_HOST")
	config_smtp_port, _ := strconv.Atoi(os.Getenv("CONFIG_SMTP_PORT"))

	dialer := gomail.NewDialer(
		config_smtp_host,
		config_smtp_port,
		es.AuthEmail,
		es.AuthPassword,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func (es *EmailService) StartConsuming() {
	conn, err := config.ConnectToRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var req dto.DigitalReceipt
			err := json.Unmarshal(d.Body, &req)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			err = es.SendDigitalReceipt(req, req.Receiver.EmailAddress)
			if err != nil {
				log.Printf("Error sending email: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
