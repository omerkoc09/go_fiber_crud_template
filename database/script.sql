--||-- 
create table if not exists `users` (
    `id` int(11) not null auto_increment,
    `name` varchar(100) not null,
    `email` varchar(100) not null,
    `password` varchar(255) not null,
    primary key (`id`)
) ENGINE=InnoDb default CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
--||--
