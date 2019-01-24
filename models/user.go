package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	UName    string
	FName    *string
	LName    *string
	Email    string
	Password string
	Descript *string
	Role     *string
	Banned   *int
}

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}

type AppSecrets struct {
	CookieSecret     string
	GoogleAuthID     string
	GoogleAuthKey    string
	MysqlPassword    string
	MysqlUser        string
	MysqlServicePort string
	MysqlDatabase    string
	MysqlServiceHost string
}

type Account struct {
	ID                      *int
	Username                *string
	Provider                *string
	UID                     *string
	Email                   *string
	Name                    *string
	FirstName               *string
	LastName                *string
	Company                 *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	IsAdmin                 *int
	DomainLast              *string
	IpAddressLast           *string
	Comment                 *string
	AuthenticationTokenHash *string
	IsBanned                *int
	TakedownCode            *int
	TakedownDescription     *string
	DomainId                *int
	IpAddressId             *int
}

type AccountResult struct {
	Results []Account
}

type CategoryCount struct {
	Category int
}

type FoundCategories struct {
	FoundCat []CategoryCount
}
