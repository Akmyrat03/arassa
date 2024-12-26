package repository

const (
	getAdminQuery = `
			SELECT 
				id, username, password
			FROM 
				admin 
			WHERE 
				username= $1 AND password=$2
		`

	signUpQuery = `
			INSERT INTO admin (username, password) 
			VALUES ($1, $2) 
			RETURNING id
		`
)
