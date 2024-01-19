package database

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"io"
	"os"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
)

func retrieveUserList(rows *sql.Rows) (*components.UserList, error) {
	var userList components.UserList
	for rows.Next() {
		var user components.User
		err := rows.Scan(&user.Username.Value, &user.Birthdate, &user.ProfilePic, &user.Name)
		if err != nil {
			return nil, err
		}

		// Open the image
		img, err := os.Open(user.ProfilePic)
		if err != nil {
			return nil, err
		}
		reader := bufio.NewReader(img)
		// Read it
		content, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		// Convert it in base64
		user.ProfilePic = base64.StdEncoding.EncodeToString(content)

		userList.Users = append(userList.Users, user)

	}

	return &userList, nil

}
