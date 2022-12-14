package helpers 


import(
	"os"
	"time"
	
	"log"
	"go.mongodb.org/mongo-driver/mongo"
	
	jwt "github.com/dgrijalva/jwt-go"
	"golang-hotell-app/database"
)
type signedDetails struct{
	Email				string
	First_name			string
	Last_name			string
	Uid 				string
	User_type			string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string , lastName string, userType string , uid string)(signedToken string, signedRefreshToken string ) {
	claims := &signedDetails{
		Email: email,
		First_name: firstName,
		Last_name: lastName,
		Uid: uid,
		User_type: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &signedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),

		},
	}

	token , err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken , err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil{
		log.Panic(err)
		return
	}
	return token, refreshToken
}