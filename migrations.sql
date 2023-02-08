

alter table auth_users change COLUMN is_superuser is_superuser varchar(191);
alter table auth_users change COLUMN is_staff is_staff varchar(191);
alter table auth_users change COLUMN is_active is_active varchar(191);





alter table blog_categories change COLUMN name name varchar(191);



rename table blog_general_cates to blog_general_categories;

alter table blog_general_categories change COLUMN name name varchar(191);

-- Alter table blog_general_categories drop general_id;



CREATE TABLE `blog_post_ps` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `password` varchar(191) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;





alter table blog_posts change COLUMN created_time created_time datetime DEFAULT NULL;
alter table blog_posts change COLUMN modified_time modified_time datetime DEFAULT NULL;
alter table blog_posts change COLUMN body body longtext;
alter table blog_posts change COLUMN excerpt excerpt text DEFAULT NULL;
alter table blog_posts change COLUMN cover_url cover_url varchar(191) DEFAULT NULL;
alter table blog_posts change COLUMN title_url title_url varchar(191) DEFAULT NULL;
alter table blog_posts add COLUMN is_open tinyint(1) DEFAULT NULL;




alter table blog_tags change COLUMN name name varchar(191);




CREATE TABLE `book_contents` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `book_content` longtext,
  `book_id` bigint unsigned DEFAULT NULL,
  `created_time` datetime DEFAULT NULL,
  `modified_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_book_contents_book` (`book_id`),
  CONSTRAINT `fk_book_contents_book` FOREIGN KEY (`book_id`) REFERENCES `books_lists` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



CREATE TABLE `books_lists` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `book_name` varchar(191) DEFAULT NULL,
  `book_status` varchar(191) DEFAULT NULL,
  `abstract` varchar(191) DEFAULT NULL,
  `created_time` datetime DEFAULT NULL,
  `modified_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `chicken_soups` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `sentence` varchar(191) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;









