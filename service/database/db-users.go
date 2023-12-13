package database

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/components"
	"github.com/dchest/uniuri"
	"github.com/julienschmidt/httprouter"
)

func (db appdbimpl) CheckCombinationIsValid(Username string, ID string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	// Prepare the SQL statement to return the row containing both the provided username and id, if there is any
	stmt, err := db.c.Prepare("SELECT COUNT(*) FROM User WHERE Username = ? AND ID = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the SQL statement to check if provided combination of username and ID is correct"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the SQL statement to check if provided combination of username and ID is correct",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while executing the SQL statement to check if provided combination of username and ID is correct"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while executing the SQL statement to check if provided combination of username and ID is correct",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}
		return &valid

	} else {
		defer rows.Close()
	}

	// Check if the returned value is exactly 1: if yes, then the user is valid
	if rows.Next() {
		var numRows string
		err := rows.Scan(&numRows)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while parsing the number of rows where username and id coincides with the ones provided"))

			error, err := json.Marshal(components.Error{
				ErrorCode:   "500",
				Description: "error while parsing the number of rows where username and id coincides with the ones provided",
			})
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			}
			_, err = w.Write([]byte(error))
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			}
			return &valid

		} else {
			valid = numRows == "1"
		}
	}

	return &valid

}

func (db appdbimpl) CheckIfUsernameExists(Username string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *bool {

	valid := false

	stmt, err := db.c.Prepare("SELECT 1 FROM User WHERE Username = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the query to check if the username provided exists or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the query to check if the username provided exists or not",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	}

	rows, err := stmt.Query(Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while executing the query to check if the username provided exists or not"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while executing the query to check if the username provided exists or not",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return &valid
	} else {
		defer rows.Close()
	}

	if rows.Next() {
		var result string
		err := rows.Scan(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while parsing the number of rows where username and id coincides with the ones provided"))

			error, err := json.Marshal(components.Error{
				ErrorCode:   "500",
				Description: "error while parsing the number of rows where username and id coincides with the ones provided",
			})
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			}
			_, err = w.Write([]byte(error))
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			}

			return &valid
		} else {
			valid = result == "1"
		}
	}

	return &valid

}

// If the user does not exist, it will be created, and an identifier is returned. If the user exists, the user identifier is returned.
func (db appdbimpl) PostUserID(Username string) (*components.ID, error) {

	// Prepare the SQL statement
	stmt, err := db.c.Prepare("SELECT ID from User WHERE Username = ?")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the id for the given user (it it exists)")
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the id for the given user (it it exists)")
	} else {
		defer rows.Close()
	}

	// Check if the username already existed
	var id string

	// If yes, just return the associated id
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the ID from the query")
		}
	} else { // If not, create a new user (and consequently a new ID for it)

		var id components.ID
		id.RandID = uniuri.NewLen(64)

		stmt, err = db.c.Prepare("INSERT INTO User (Username, ID) VALUES (?, ?)")
		if err != nil {
			return nil, fmt.Errorf("error while preparing the SQL statement to create the new user")
		}

		_, err = stmt.Query(Username, id)
		if err != nil {
			return nil, fmt.Errorf("error while performing the query to create the new user")
		}

	}

	ID := components.ID{RandID: id}
	return &ID, nil

}

func (db appdbimpl) SearchUser(Username string) (*components.UserList, error) {

	// Prepare the SQL statement for finding all the users with "uname" as substring
	stmt, err := db.c.Prepare("SELECT Username FROM User WHERE Username LIKE '%?%'")
	if err != nil {
		return nil, fmt.Errorf("error while preparing the SQL statement to obtain the list of users with the provided string as substring")
	}

	// Bind the parameters and execute the statement
	rows, err := stmt.Query(Username)
	if err != nil {
		return nil, fmt.Errorf("error while performing the query to obtain the list of users with the provided string as substring")
	} else {
		defer rows.Close()
	}

	// Instantiate the data structure that will hold the list of usernames
	var ulist components.UserList

	// Loop over the rows, and store each user id in the previously instantiated data structure
	for rows.Next() {

		// Retrieve the next username
		var user components.Username
		err = rows.Scan(&user.Uname)
		if err != nil {
			return nil, fmt.Errorf("error while extracting the username from the query")
		}

		// Insert into the returned list of usernames
		ulist.Users = append(ulist.Users, user)

	}

	// Return the list of users
	return &ulist, nil

}

func (db appdbimpl) UpdateUsername(OldUsername string, NewUsername string) error {

	stmt, err := db.c.Prepare("UPDATE User SET Username = ? WHERE Username = ?")
	if err != nil {
		return fmt.Errorf("error while preparing the SQL statement to updating the username")
	}

	rows, err := stmt.Query(NewUsername, OldUsername)
	if err != nil {
		if err != nil {
			return fmt.Errorf("error while performing the query to obtain the info about the user with the provided username")
		} else {
			defer rows.Close()
		}
	}

	return nil

}

func (db appdbimpl) GetUsernameByToken(Id string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *components.Username {

	stmt, err := db.c.Prepare("SELECT Username FROM Users WHERE ID = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while preparing the query to retrieve the username associated with the given token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while preparing the query to retrieve the username associated with the given token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return nil
	}

	rows, err := stmt.Query(Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while executing the query to retrieve the username associated with the given token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error encountered while executing the query to retrieve the username associated with the given token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return nil
	}

	var username components.Username
	if rows.Next() { // We can be sure to have one username at most since the column token is set to be unique
		err = rows.Scan(&username.Uname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(fmt.Errorf("error encountered while scanning the username associated with the given token from the returned rows"))

			error, err := json.Marshal(components.Error{
				ErrorCode:   "500",
				Description: "error encountered while scanning the username associated with the given token from the returned rows",
			})
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			}
			_, err = w.Write([]byte(error))
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			}

			return nil
		}
	}

	if len(username.Uname) == 0 { // No username associated with the provided token
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Error(fmt.Errorf("no username associated with the provided Auth token"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "400",
			Description: "no username associated with the provided Auth token",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return nil
	}

	return &username

}

func (db appdbimpl) GetOwnerUsernameOfComment(Id string, w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) *components.Username {

	stmt, err := db.c.Prepare("SELECT Author FROM Comment WHERE CommentID = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while preparing the SQL statement to retrieve the author of the provided comment"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while preparing the SQL statement to retrieve the author of the provided comment",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return nil
	}

	rows, err := stmt.Query(Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error(fmt.Errorf("error while executing the SQL statement to retrieve the author of the provided comment"))

		error, err := json.Marshal(components.Error{
			ErrorCode:   "500",
			Description: "error while executing the SQL statement to retrieve the author of the provided comment",
		})
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
		}
		_, err = w.Write([]byte(error))
		if err != nil {
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
		}

		return nil
	}

	var username components.Username
	if rows.Next() { // We can be sure to have one username at most since the column 'CommentID' is set to be unique
		err = rows.Scan(&username.Uname)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error(fmt.Errorf("error while scanning the result of the SQL query to retrieve the author of the provided comment"))

			error, err := json.Marshal(components.Error{
				ErrorCode:   "500",
				Description: "error while scanning the result of the SQL query to retrieve the author of the provided comment",
			})
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while encoding the response as JSON"))
			}
			_, err = w.Write([]byte(error))
			if err != nil {
				ctx.Logger.WithError(err).Error(fmt.Errorf("error while writing the response error in the response body"))
			}

			return nil
		}
	}

	return &username

}
