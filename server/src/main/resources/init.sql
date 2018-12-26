create table if not exists hubs
(
	hub_id varchar(64) not null primary key,
	created_at timestamp default current_timestamp not null,
	updated_at timestamp default current_timestamp on update current_timestamp not null,
	address varchar(100) null,
	port int
);
