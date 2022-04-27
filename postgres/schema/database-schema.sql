
create table mailGroups (
    id serial primary key,
    groupName varchar(300) unique not null,
    created timestamp not null default current_timestamp,
    updated timestamp not null default current_timestamp
);

create table mailRecipients (
    id serial primary key,
    email varchar(254) unique not null,
    usersServiceID integer,
    created timestamp not null default current_timestamp,
    updated timestamp not null default current_timestamp,
    confirmed bool
);

CREATE EXTENSION pg_trgm;
CREATE INDEX ON mailRecipients USING gin (email gin_trgm_ops);

create table recipientGroupMap (
    id serial primary key,
    groupID integer references mailGroups(id),
    recipientID integer references mailRecipients(id) on delete cascade 
);


insert into mailRecipients(id, email, usersServiceID, created, updated, confirmed) values (default, 'hyperxpizza@gmail.com', 1, default, default, false);
insert into mailRecipients(id, email, usersServiceID, created, updated, confirmed) values (default, 'hyperxpizza2@gmail.com', 2, default, default, false);
insert into mailRecipients(id, email, usersServiceID, created, updated, confirmed) values (default, 'hyperxpizza3@gmail.com', 3, default, default, false);

insert into mailGroups(id, groupName, created, updated) values(default, 'CUSTOMERS', default, default);
insert into mailGroups(id, groupName, created, updated) values(default, 'NEWSLETTER', default, default);

insert into recipientGroupMap(id, groupID, recipientID) values(default, 1, 1);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 1, 2);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 1, 3);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 2, 1);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 2, 2);

