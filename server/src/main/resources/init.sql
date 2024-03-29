create table if not exists agents
(
  agent_id     varchar(64)                         not null primary key,
  created_at   timestamp default current_timestamp not null,
  updated_at   timestamp default current_timestamp not null,
  address      varchar(100)                        null,
  web_port     int,
  stream_port  int,
  sensors_port int
);
