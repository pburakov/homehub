package io.pburakov.homehub.server.storage.dao;

import io.pburakov.homehub.server.storage.model.Hub;
import org.jdbi.v3.sqlobject.config.RegisterRowMapper;
import org.jdbi.v3.sqlobject.statement.SqlQuery;
import org.jdbi.v3.sqlobject.statement.SqlUpdate;
import org.jdbi.v3.sqlobject.transaction.Transactional;

public interface HubDao extends Transactional<HubDao> {

  @SqlUpdate("insert into hubs "
      + "(hub_id, address, web_port, stream_port, meta_port) "
      + "values (?, ?, ?, ?, ?)")
  void insert(String hubId, String address, int webPort, int streamPort, int metaPort);

  @SqlUpdate("update hubs "
      + "set address = ?, web_port = ?, stream_port = ?, meta_port = ? "
      + "where hub_id = ?")
  void update(String address, int webPort, int streamPort, int metaPort, String hubId);

  @SqlQuery("select * from hubs where hub_id = ?")
  @RegisterRowMapper(Hub.Mapper.class)
  Hub select(String hubId);

}
