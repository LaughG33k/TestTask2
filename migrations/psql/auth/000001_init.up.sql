CREATE EXTENSION if not exists "uuid-ossp";
CREATE EXTENSION if not exists "pgcrypto";



create table if not exists users (

    Id serial primary key,
    Uuid uuid default uuid_generate_v4() unique,
    Login varchar(30) unique not null,
    Password varchar(250) not null,
    Permission varchar(30) not null

);

create table if not exists refresh_tokens (

    Id serial primary key,
    Token varchar(300),
    Time_end_of_life bigint,
    Owner_uuid uuid not null references users(uuid) on delete cascade

);