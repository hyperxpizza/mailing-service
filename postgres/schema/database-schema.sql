
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
    tsv tsvector,
    confirmed bool
);

-- full text search
update mailRecipients set tsv = setweight(to_tsvector(email), 'A');
create index ix_mailRecipients_tsv on mailRecipients using GIN(tsv);
--end full text search
--example search query
--select 
--  id, email, usersServiceID, created, updated, ts_headline(email, q)
--from (
--  select 
--      id, email, usersServiceID, created, updated, ts_rank(tsv, q) as rank
--  from
--      mailRecipients, plainto_tsquery($1) q
--  where 
--      tsv @@ q
--  order by
--      rank desc
--)
--order by
--  rank desc

create table recipientGroupMap (
    id serial primary key,
    groupID integer references mailGroups(id),
    recipientID integer references mailRecipients(id) on delete cascade 
);

-- get groups of a user
--select m.groupID, g.groupName, g.created, g.updated from recipientGroupMap as m join mailGroups as g on g.id=m.groupID where recipientID=1;

-- count users by group
-- select COUNT(*) from mailRecipients as r join recipientGroupMap as m on r.id = m.recipientID join mailGroups as g on m.groupID = g.id where g.groupName = '';

-- get recipients with order and pagination
--select id, email, usersServiceID, created, updated, confirmed from mailRecipients order by $1 limit $2 offset $3;


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