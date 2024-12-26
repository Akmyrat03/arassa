package repository

const (
	categoriesQuery = `
		INSERT INTO 
			categories 
		DEFAULT VALUES 
			RETURNING id
		`

	categoryTranslateQuery = `
		INSERT INTO 
			cat_translate (cat_id, lang_id, name)
		VALUES 
			($1, $2, $3)`

	deleteCategory = `
		DELETE 
		FROM categories 
		WHERE id = $1
		`

	getAllCategoriesByLangID = `
		SELECT 
			c.id, ct.name 
		FROM 
			categories AS c 
		INNER JOIN 
			cat_translate AS ct ON c.id=ct.cat_id 
		INNER JOIN 
			languages AS l ON l.id=ct.lang_id 
		WHERE lang_id = $1	
	`
)
