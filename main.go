package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gammazero/workerpool"
)

const (
	// Sender Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "support@acloset.app"

	// Recipient Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	//Recipient = "acloset.app@gmail.com"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	ConfigurationSet = "s3_delivery"

	// Subject The subject line for the email.
	Subject = Subject_kr

	Subject_kr = `지속가능한 패션🌱을 위한 에이클로젯과 중고의 결합`
	Subject_fr = `[Acloset] Présentation de la fonctionnalité "Seconde Main" pour une mode plus durable 🌱`
	Subject_es = `[Acloset] Presentamos la función "Pre-amado" para una moda más sostenible 🌱`
	Subject_en = `[Acloset] Introducing ‘Preloved’ feature for more sustainable fashion 🌱`
	Subject_ru = `[Acloset] Представляем функцию рынка для более устойчивой моды 🌱`
	// HtmlBody The HTML body for the email.
	//HtmlBody = `Amazon SES Test (AWS SDK for Go)`
	// TextBody The email body for recipients with non-HTML email clients.
	//TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// CharSet The character encoding for the email.
	CharSet = "UTF-8"
	//excludePath      = "exclude_final.csv"
	recipientCSVPath = "recipient.csv"
	htmlPath         = "test_notice.html"
)

func main() {

	content, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	htmlString := string(content)

	//Create a new session in the us-west-2 region.
	//Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)

	//Create an SES session.
	svc := ses.New(sess)

	wp := workerpool.New(20)

	// Send the email.
	for i, v := range ListRecipientString() {
		i := i
		v := v
		if strings.Contains(v, "@hanmail.net") {
			fmt.Println("hanmail.net")
			continue
		}
		//@privaterelay.appleid.com
		if strings.Contains(v, "@privaterelay.appleid.com") {
			fmt.Println("privaterelay.appleid.com")
			continue
		}
		wp.Submit(func() {
			if err := SendEmail(v, htmlString, svc); err != nil {
				fmt.Printf("Error sending email:%v, index:%v", err, i)
			}
			if err != nil {
				fmt.Printf("Error sending email:%v, index:%v", err, i)
			}
			fmt.Println(i, v)
		})
	}
	wp.StopWait()
	return
}

func SendEmail(recipient string, htmlString string, svc *ses.SES) error {
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlString),
				},
				//Text: &ses.Content{
				//	Charset: aws.String(CharSet),
				//	Data:    aws.String(TextBody),
				//},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return err
	}

	fmt.Println("Email Sent to address: ")
	fmt.Println(result)
	return nil
}

func ListRecipientString() []string {
	var result []string

	file, _ := os.Open(recipientCSVPath)
	rdr := csv.NewReader(bufio.NewReader(file))
	rows, _ := rdr.ReadAll()
	//exclude := excludeList()

	for i, row := range rows {
		for j := range row {
			//if funk.Contains(exclude, rows[i][j]) {
			//	//fmt.Println("exclude", rows[i][j])
			//	continue
			//}
			result = append(result, rows[i][j])
			//fmt.Println("!!", rows[i][j])
		}
	}

	return result
}

//
//func excludeList() []string {
//	var result []string
//
//	file, _ := os.Open(excludePath)
//	rdr := csv.NewReader(bufio.NewReader(file))
//	rows, _ := rdr.ReadAll()
//
//	for i, row := range rows {
//		for j := range row {
//			result = append(result, rows[i][j])
//			//fmt.Println("!!", rows[i][j])
//		}
//	}
//
//	return result
//}
