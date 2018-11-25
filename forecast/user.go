package forecast

import (
	"fmt"
)

type MeResponse struct {
	Me *Me `json:"current_user"`
}

type Me struct {
	ID int `json:"id"`
}

type UserResponse struct {
	Person *User `json:"person"`
}

type User struct {
	ID         int    `json:"id"`
	HarvestID  int    `json:"harvest_user_id"`
	Admin      bool   `json:"admin"`
	Archived   bool   `json:"archived"`
	AvatarURL  string `json:"avatar_url"`
	ColorBlind bool   `json:"color_blind"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
}

func (f *Forecast) GetMe() (*Me, error) {
	v := &MeResponse{}
	return v.Me, f.get("/whoami", nil, v)
}

func (f *Forecast) GetUser(id int) (*User, error) {
	v := &UserResponse{}
	return v.Person, f.get(fmt.Sprintf("/people/%d", id), nil, v)
}
