package databases

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/openshift/sre-dashboard/models"
	"os"
	"time"
)

var (
	myHost = os.Getenv("MYSQL_SERVICE_HOST")
	myPort = os.Getenv("MYSQL_SERVICE_PORT")
	myUser = os.Getenv("MYSQL_USER")
	myPass = os.Getenv("MYSQL_PASSWORD")
	myName = os.Getenv("MYSQL_DATABASE")
)

func QueryTakedowns(dateRange int) map[string]int {
	var foundRes models.AccountResult
	var catCount = make(map[string]int)
	var bannedUsers = make(map[string][]string)
	var startRange time.Time = time.Now().AddDate(0, 0, -dateRange)
	var endRange time.Time = time.Now()

	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", myUser, myPass, myHost, myPort, myName)

	mydb, err := gorm.Open("mysql", mysqlInfo)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := mydb.Model(&models.Account{}).Where("is_banned = 1 AND created_at BETWEEN ? AND ? ", startRange, endRange).Rows()
	defer rows.Close()
	for rows.Next() {
		var account models.Account
		mydb.ScanRows(rows, &account)
		foundRes.Results = append(foundRes.Results, account)
	}

	// Populate a map whose keys are the names of the takedown categories themselves
	// e.g. spamming, phishing. The values of these keys are string slices, which we
	// will then count to get the total number of offenses per category.
	for _, item := range foundRes.Results {
		if *item.IsBanned == 1 {
			bannedUsers[models.TakedownCategory[*item.TakedownCode]] = append(bannedUsers[models.TakedownCategory[*item.TakedownCode]], *item.Username)
		}
	}

	for key, _ := range bannedUsers {
		catCount[key] = len(bannedUsers[key])
	}

	return catCount
}
