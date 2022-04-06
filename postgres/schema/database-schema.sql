-- drop table if exists mailRecipients;


create table mailGroup (
    id serial primary key,
    groupName varchar(300) unique not null,
    created timestamp not null,
    updated timestamp not null
);

create table mailRecipients (
    id serial primary key,
    email varchar(254) unique not null,
    usersServiceID integer,
    created timestamp not null,
    updated timestamp not null
);

create table recipientGroupMap (
    id serial primary key,
    groupID integer references mailGroup(id),
    recipientID integer references mailRecipients(id) 
);