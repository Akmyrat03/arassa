SELECT n.id AS news_id, nt.title, ct.name AS category, nt.description, n.image, n.created_at FROM news_translate AS nt 
INNER JOIN news AS n ON nt.news_id= n.id 
INNER JOIN cat_translate AS ct ON n.category_id=ct.cat_id
AND nt.lang_id=ct.lang_id WHERE nt.lang_id = '3';