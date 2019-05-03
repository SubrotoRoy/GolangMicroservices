package main
 
import (
    "fmt"
    "net/http"
    "log"
    ms "microservices/email-service/email"
    _ "github.com/lib/pq"
    verify "microservices/email-service/VerifyJWT"
    "os"
    "encoding/json"
    "strconv"
    "strings"
    consulapi "github.com/hashicorp/consul/api"
)

type event struct{
    Email string `json:"email"`
    Subject string `json:"subject"`
    Description string `json:"description"`
    Location string `json:"location"`
    StartDateTime string `json:"StartDateTime"`
    EndDateTime string `json:"EndDateTime"`
}


type JwtToken struct {
	Token string `json:"token"`
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}

	registration := new(consulapi.AgentServiceRegistration)

	registration.ID = "email-service"
	registration.Name = "email-service"
	address := hostname()
	registration.Address = address
	port, err := strconv.Atoi(port()[1:len(port())]) 
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	fmt.Println("port:",port)
	consul.Agent().ServiceRegister(registration)
}

func port() string {
	p := os.Getenv("EMAIL_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8083"
	}
	return fmt.Sprintf(":%s", p)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	//hn = "52.71.115.181"
	return hn
}

func main() {
	//os.Setenv("CONSUL_HTTP_ADDR","52.71.115.181")
        registerServiceWithConsul()
	fmt.Printf("user service is up on port: %s", port())
	http.HandleFunc("/GetEvent",GetEvent)
	http.ListenAndServe(port(), nil)
    
}



//Sending an ICS as email

func lookupServiceWithConsul(serviceName string) (string, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := consul.Agent().Services()
	if err != nil {
		return "", err
	}
	srvc := services["event-service"]
	address := srvc.Address
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	res,err := verify.VerifyHandler(r)
    fmt.Println(res)
    if err == nil {
	e:= event{}
	url, err := lookupServiceWithConsul("email-service")
	fmt.Println("URL: ", url)
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Get(url + "/ListEvent")
	if err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		fmt.Fprintf(w, "Error. %s", err)
		return
	}

	fmt.Printf("Check")
	fmt.Fprintln(w,e)

    r.ParseForm()
    email:=e.Email
    name := r.FormValue("name")
    subject := r.FormValue("subject")
    message := r.FormValue("message")
    fmt.Println(email,name,subject,message)
    ms.SendEmail(name,email,subject,message)

    var file, err1 = os.Create(`calendar1.ics`)
    defer file.Close()
    fmt.Println(e.StartDateTime)
    fmt.Fprintf(file,"BEGIN:VCALENDAR\nMETHOD:PUBLISH\nVERSION:2.0\nPRODID:-//Company Name//Product//Language\nBEGIN:VEVENT")
    fmt.Fprintf(file,"\nSUMMARY:")

    fmt.Fprintf(file,e.Subject)
    fmt.Fprintf(file,"\nDTSTART:")
    fmt.Fprintf(file,e.StartDateTime)
    fmt.Fprintf(file,"\nDTEND:")
    fmt.Fprintf(file,e.EndDateTime)
    fmt.Fprintf(file,"\nDESCRIPTION:")
    fmt.Fprintf(file,e.Description)
    fmt.Fprintf(file,"\nLOCATION:")
    fmt.Fprintf(file,e.Location)
    fmt.Fprintf(file,"\nEND:VEVENT\nEND:VCALENDAR")   
    if err1 != nil {
        fmt.Println(err1)
    }
	}else{
		fmt.Println(err)
	}
}
