CREATE TABLE `teacher` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `create_time` TIMESTAMP DEFAULT NULL,
    `update_time` TIMESTAMP DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `firstname` VARCHAR(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
    `lastname` VARCHAR(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL,
    PRIMARY KEY (`id`)
)


SELECT CONCAT(a.first_name, ' ', a.last_name) as full_name, b.title FROM `cg_authors` a
LEFT JOIN cg_book_authors bk  ON bk.author_id = a.id
WHERE bk.book_id = 30;

SELECT c.name FROM `cg_categories` c
LEFT JOIN cg_book_categories bc  ON bc.category_id = c.id
WHERE bc.book_id = 27;