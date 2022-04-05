-- drop table if exists mailRecipients;
create table mailRecipients (
    id serial primary key,
    email varchar(254) not null,
    usersServiceID integer,
    created timestamp not null,
    updated timestamp not null
);