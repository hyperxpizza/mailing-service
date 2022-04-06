-- drop table if exists mailRecipients;


create table mailGroup (
    id serial primary key,
    groupName varchar(300) not null,
    created timestamp not null,
    updated timestamp not null
);

insert into mailGroup(id, groupName) values (default, "newsletter");
insert into mailGroup(id, groupName) values (default, "customer");
insert into mailGroup(id, )

create table mailRecipients (
    id serial primary key,
    email varchar(254) not null,
    usersServiceID integer,
    created timestamp not null,
    updated timestamp not null
);

create table recipientGroupMap (
    id serial primary key,
    groupID integer foreign key references mailGroup(id),
    recipientID integer foreing key references mailRecipients(id) 
);