package verifyJWT

import(
        "net/http"
        "github.com/gorilla/sessions"
        verifier "github.com/okta/okta-jwt-verifier-golang"
        "fmt"
        "github.com/dgrijalva/jwt-go"

)

type Jwt struct {
        Claims map[string]interface{}
}

var sessionStore = sessions.NewCookieStore([]byte("okta-custom-login-session-store"))
var state = "ApplicationState"

func VerifyHandler(r *http.Request) (*verifier.Jwt, error) {

		session, err := sessionStore.Get(r, "okta-custom-login-session-store")
		if err != nil {
			fmt.Println(err)
	}
	var tok string
	tok = (session.Values["id_token"]).(string)
	tokenString := tok
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

	return []byte(""), nil
	})
	fmt.Println(token)
	if err != nil {
			fmt.Println("Decoding error")
	}

	Nonce := ""
	for key, val := range claims {
		if key == "nonce"{
			Nonce = val.(string)
	}
	}
	result, verificationError := verifyToken(tok,Nonce) //Token verification

	if verificationError != nil {
		fmt.Println("verf err",verificationError)
	}

	fmt.Println("result:jwt:",result)
	return result,verificationError
}

func verifyToken(t string,nonce string) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = "0oaj44d0dPhPhc3M9356"
	fmt.Println("tv",tv)
	jv := verifier.JwtVerifier{
	Issuer:  "https://dev-502722.okta.com/oauth2/default"  ,
	ClaimsToValidate: tv,
        }
        fmt.Println("jv:",jv)
        result, err := jv.New().VerifyIdToken(t)

        if err != nil {
                return nil, fmt.Errorf("err:%s", err)
        }

        if result != nil {
                fmt.Println("res",result)
                fmt.Println(result.Claims)
                return result, nil
        }

        return nil, fmt.Errorf("token could not be verified: %s", "")
}

