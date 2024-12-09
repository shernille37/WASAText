package database

import "github.com/gofrs/uuid"

func (db *appdbimpl) ListUsers(id uuid.UUID) ([]User, error) {

	var res []User

	const queryUsers = `
		SELECT u.userID, u.username, u.image
		FROM User u
		WHERE u.userID <> ?;
	`

	rows, err := db.c.Query(queryUsers, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User

		if err = rows.Scan(&u.UserID, &u.Name, &u.Image); err != nil {
			return nil, err
		}

		res = append(res, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil

}
