package email

import (
	"fmt"
	//"gopkg.in/gomail.v2"
	"log"
	"os"
	"encoding/base64"
  	"io/ioutil"
	"github.com/sendgrid/sendgrid-go"
  	"github.com/sendgrid/sendgrid-go/helpers/mail"
)
  
type to struct{
	Name,Address string
}

func SendEmail(name,email,subject,message string){
	//os.Setenv("serverpassrd","")
	reciever := new(to)
	reciever.Name=name
	reciever.Address=email
	SendMail(*reciever,subject,message,"calendar1.ics")
}

func SendMail(to2 struct {Name string; Address string}, subject string, message string, filePath string) {

	m := mail.NewV3Mail()
	
	from := mail.NewEmail("Subroy", "subroto7396@gmail.com")
	content := mail.NewContent("text/html", "<p>" + message + "</p>")
	to := mail.NewEmail(to2.Name, to2.Address)
	
	m.SetFrom(from)
	m.AddContent(content)
	
	// create new *Personalization
	personalization := mail.NewPersonalization()
	personalization.AddTos(to)
	personalization.Subject = subject
   
	// add `personalization` to `m`
	m.AddPersonalizations(personalization)
	
	// read/attach .txt file
	a_txt := mail.NewAttachment()
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
	  fmt.Println(err)
	}	
	
	encoded := base64.StdEncoding.EncodeToString([]byte(dat))
  	a_txt.SetContent(encoded)
  	a_txt.SetType("text/Calendar")
  	a_txt.SetFilename("Meeting.ics")
  	a_txt.SetDisposition("attachment")
  	a_txt.SetContentID("Test Document")
	
	// add `a_txt`, `a_pdf` and `a_jpg` to `m`
	m.AddAttachment(a_txt)
	
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
  	request.Method = "POST"
  	request.Body = mail.GetRequestBody(m)
  	response, err := sendgrid.API(request)
  	if err != nil {
    	log.Println(err)
 	} else {
    	fmt.Println(response.StatusCode)
    	fmt.Println(response.Body)
    	fmt.Println(response.Headers)
  	}
	
	
	/*d := gomail.NewDialer("smtp.gmail.com", 587, "", os.Getenv("serverpassrd"))
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
		m.SetHeader("From", "")
		m.SetAddressHeader("To", to2.Address, to2.Name)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", fmt.Sprintf("Hi %s!", to2.Name) + "\n " + message)

		m.Attach(filePath)

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", to2.Address, err)
		}
		m.Reset()*/
}

