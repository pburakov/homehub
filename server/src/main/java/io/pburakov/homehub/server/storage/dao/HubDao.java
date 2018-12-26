package io.pburakov.homehub.server.storage.dao;

import io.pburakov.homehub.server.storage.model.Hub;
import org.jdbi.v3.sqlobject.config.RegisterRowMapper;
import org.jdbi.v3.sqlobject.statement.SqlQuery;
import org.jdbi.v3.sqlobject.statement.SqlUpdate;
import org.jdbi.v3.sqlobject.transaction.Transactional;

public interface HubDao extends Transactional<HubDao> {

  @SqlUpdate("insert into hubs (hub_id, address, port) values (?, ?, ?)")
  void insert(String hubId, String address, int port);

  @SqlUpdate("update hubs set address = ?, port = ? where hub_id = ?")
  void update(String address, long port, String hubId);

  @SqlQuery("select * from hubs where hub_id = ?")
  @RegisterRowMapper(Hub.Mapper.class)
  Hub select(String hubId);

}
