-- books
CREATE TABLE IF NOT EXISTS `cg_books` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `status` ENUM('Not Read','In progress','Read') NOT NULL DEFAULT 'Not Read',
  `status_id` int NOT NULl DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now())
);
-- authors
CREATE TABLE IF NOT EXISTS `cg_authors` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `description` text,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now())
);
-- categories
CREATE TABLE IF NOT EXISTS `cg_categories` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name`  varchar(255) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now())
);
-- book authors
CREATE TABLE IF NOT EXISTS `cg_book_authors` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_id` int,
  `author_id` int,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now())
);
-- book categories
CREATE TABLE IF NOT EXISTS `cg_book_categories` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `book_id` int,
  `category_id` int,
  `created_at` datetime NOT NULL DEFAULT (now()),
  `updated_at` datetime NOT NULL DEFAULT (now())
);

CREATE INDEX `cg_books_index_0` ON `cg_books` (`title`);
CREATE INDEX `cg_books_index_1` ON `cg_books` (`id`);

CREATE INDEX `cg_authors_index_2` ON `cg_authors` (`id`);
CREATE INDEX `cg_categories_index_3` ON `cg_categories` (`id`);

ALTER TABLE `cg_book_authors` ADD FOREIGN KEY (`book_id`) REFERENCES `cg_books` (`id`) ON DELETE CASCADE;
ALTER TABLE `cg_book_authors` ADD FOREIGN KEY (`author_id`) REFERENCES `cg_authors` (`id`) ON DELETE CASCADE;

ALTER TABLE `cg_book_categories` ADD FOREIGN KEY (`book_id`) REFERENCES `cg_books` (`id`) ON DELETE CASCADE;
ALTER TABLE `cg_book_categories` ADD FOREIGN KEY (`category_id`) REFERENCES `cg_categories` (`id`) ON DELETE CASCADE;


