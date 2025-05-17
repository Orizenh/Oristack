package user

import (
	"fmt"
	"net/http"
	"oristack/initializers"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	DateAdd  string `json:"dateAdd"`
	IsActive bool   `json:"isActive"`
	Token    string `json:"token"`
}

var userSetup = map[string]string{
	"payload": "id,username,date_add,is_active,token",
	"table":   "user",
}

func GetUser(wrapper *initializers.Wrapper) {
	rows, err := initializers.DB.Query("SELECT "+userSetup["payload"]+" FROM "+userSetup["table"]+" WHERE id=?", wrapper.ReturnUser())
	if err != nil {
		wrapper.Error(err.Error())
		return
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.DateAdd, &user.IsActive, &user.Token); err != nil {
			wrapper.Error(err.Error())
			return
		}
	}
	if user.ID == 0 {
		wrapper.Error("User cannot be found", http.StatusNotFound)
		return
	}
	wrapper.Render(map[string]any{
		"data": user,
	})
}

func GetUserConnect(wrapper *initializers.Wrapper) {
	if err := wrapper.WrapData("username"); err != nil {
		wrapper.Error("You have to send the username")
		return
	}
	rows, err := initializers.DB.Query("SELECT "+userSetup["payload"]+" FROM "+userSetup["table"]+" WHERE username=?", wrapper.Data["username"])
	if err != nil {
		wrapper.Error(err.Error())
		return
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.DateAdd, &user.IsActive, &user.Token); err != nil {
			wrapper.Error(err.Error())
			return
		}
	}
	if user.ID == 0 {
		wrapper.Error("User cannot be found", http.StatusNotFound)
		return
	}
	wrapper.Render(map[string]any{
		"data": user,
	})
}

func GetUserAuth(wrapper *initializers.Wrapper) (userID int, error error) {
	if err := wrapper.WrapData("token"); err != nil {
		return 0, err
	}
	token := strings.Replace(wrapper.Data["token"].(string), "Bearer ", "", -1)
	rows, err := initializers.DB.Query("SELECT "+userSetup["payload"]+" FROM "+userSetup["table"]+" WHERE token=?", token)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var user User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.DateAdd, &user.IsActive, &user.Token); err != nil {
			return 0, err
		}
	}
	if user.ID == 0 {
		return 0, fmt.Errorf("user cannot be found, token = %v", token)
	}
	return user.ID, err
}

func CreateUser(wrapper *initializers.Wrapper) {
	if err := wrapper.WrapData("username"); err != nil {
		wrapper.Error(err.Error(), http.StatusBadGateway)
		return
	}
	rows, err := initializers.DB.Query("SELECT "+userSetup["payload"]+" FROM "+userSetup["table"]+" WHERE username=?", wrapper.Data["username"])
	if err != nil {
		wrapper.Error(err.Error(), http.StatusBadGateway)
		return
	}
	defer rows.Close()
	if rows.Next() {
		wrapper.Error("Username already exist", http.StatusBadGateway)
		return
	}
	smtp, err := initializers.DB.Prepare("INSERT INTO " + userSetup["table"] + "(username,date_add,is_active,token) VALUES(?,?,?,?)")
	if err != nil {
		wrapper.Error(err.Error(), http.StatusBadGateway)
		return
	}
	defer smtp.Close()
	user := User{
		Username: wrapper.Data["username"].(string),
		DateAdd:  time.Now().UTC().Truncate(time.Second).Format(initializers.DateFormat),
		IsActive: true,
		Token:    uuid.New().String(),
	}
	result, err := smtp.Exec(user.Username, user.DateAdd, user.IsActive, user.Token)
	if err != nil {
		wrapper.Error(err.Error(), http.StatusBadGateway)
		return
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		wrapper.Error(err.Error(), http.StatusBadGateway)
		return
	}
	user.ID = int(lastInsertID)
	wrapper.Render(map[string]any{
		"data": user,
	})
}
