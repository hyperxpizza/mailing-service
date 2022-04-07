
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
    updated timestamp not null,
    tsv tsvector
);

-- full text search
update mailRecipients set tsv = setweight(to_tsvector(email), 'A');
create index ix_mailRecipients_tsv on mailRecipients using GIN(tsv);
--end full text search
--example search query
--select 
--  id, email, usersServiceID, created, updated, ts_headling(email, q)
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
    groupID integer references mailGroup(id),
    recipientID integer references mailRecipients(id) 
);