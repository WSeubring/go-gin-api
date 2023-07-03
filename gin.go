// Start a ping pong server with gin
package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/hello/:id", GetAddress)

	router.Run(":8081")
}

func GetAddress(ctx *gin.Context) {
	type Result struct {
		ID      int    `json:"id"`
		Address string `json:"address"`
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid id provided",
		})
		return
	}

	users, err := FetchUsers()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	targetUser := FindUserByID(users, id)
	addressString := GetAddressString(targetUser)

	result := Result{
		ID:      targetUser.ID,
		Address: addressString,
	}
	ctx.JSON(http.StatusOK, result)
}

func GetAddressString(user User) string {
	return user.Address.City + " " + user.Address.Zipcode + " (" + user.Address.Geo.Lat + ", " + user.Address.Geo.Lng + ")"
}

func FindUserByID(users []User, id int) User {
	for _, user := range users {
		if user.ID == id {
			return user
		}
	}
	return User{}
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  struct {
		Street  string `json:"street"`
		Suite   string `json:"suite"`
		City    string `json:"city"`
		Zipcode string `json:"zipcode"`
		Geo     struct {
			Lat string `json:"lat"`
			Lng string `json:"lng"`
		} `json:"geo"`
	} `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
	Company struct {
		Name        string `json:"name"`
		CatchPhrase string `json:"catchPhrase"`
		Bs          string `json:"bs"`
	} `json:"company"`
}

func FetchUsers() ([]User, error) {
	url := "https://jsonplaceholder.typicode.com/users"
	resp, fetch_err := http.Get(url)
	if fetch_err != nil {
		return nil, errors.New("Error fetching users")
	}

	body, read_body_err := ioutil.ReadAll(resp.Body)
	if read_body_err != nil {
		return nil, errors.New("Error reading body")
	}

	users, err := ParseJsonToUsers(body)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func ParseJsonToUsers(body []byte) ([]User, error) {
	users := []User{}
	err := json.Unmarshal(body, &users)
	if err != nil {
		return nil, errors.New("Error unmarshalling body")
	}

	return users, nil
}
