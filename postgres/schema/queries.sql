select 
  a.id, a.email, a.usersServiceID, a.created, a.updated, a.confirmed, ts_headline(a.email, a.q)
from (
  select 
      id, email, usersServiceID, created, updated, confirmed, ts_rank(tsv, q) as rank, q
  from
      mailRecipients, plainto_tsquery('hyperxpizza') q
  where 
      tsv @@ q
  order by
      rank desc
) as a
order by
  rank desc;


select id, email, usersServiceID, created, updated, confirmed from mailRecipients, plainto_tsquery('hyper') q where tsv @@ q 


-- get groups of a user
--select m.groupID, g.groupName, g.created, g.updated from recipientGroupMap as m join mailGroups as g on g.id=m.groupID where recipientID=1;

-- count users by group
-- select COUNT(*) from mailRecipients as r join recipientGroupMap as m on r.id = m.recipientID join mailGroups as g on m.groupID = g.id where g.groupName = '';

-- get recipients with order and pagination
--select id, email, usersServiceID, created, updated, confirmed from mailRecipients order by $1 limit $2 offset $3;

-- get recipients with order and pagination where group
-- select r.id, r.email, r.usersServiceID, r.created, r.updated, r.confirmed from mailRecipients as r join recipientGroupMap as m on r.id = m.recipientID join mailGroups as g on m.groupID = g.id where g.groupName = $1 order by $2 limit $3 offset $4;

select id, email from mailRecipients where tsv @@ to_tsquery('gmail.com');