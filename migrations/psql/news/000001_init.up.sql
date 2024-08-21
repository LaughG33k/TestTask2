create table if not exists News  (
    Id bigserial primary key not null,
    Title varchar(255),
    Content text
);


create table if not exists NewsCategory (
    NewsId bigint not null,
    CategoryId bigint not null
);