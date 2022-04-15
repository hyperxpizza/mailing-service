insert into mailRecipients(id, email, usersServiceID, created, updated, confirmed) values (default, 'hyperxpizza@gmail.com', 1, default, default, false);
insert into mailRecipients(id, email, usersServiceID, created, updated, confirmed) values (default, 'hyperxpizza2@gmail.com', 2, default, default, false);
insert into mailRecipients(id, email, usersServiceID, created, updated, confirmed) values (default, 'hyperxpizza3@gmail.com', 3, default, default, false);

insert into mailGroups(id, groupName, created, updated) values(default, "CUSTOMERS", default, default);
insert into mailGroups(id, groupName, created, updated) values(default, "NEWSLETTER", default, default);

insert into recipientGroupMap(id, groupID, recipientID) values(default, 1, 1);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 1, 2);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 1, 3);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 2, 1);
insert into recipientGroupMap(id, groupID, recipientID) values(default, 2, 2);