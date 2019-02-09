package io.pburakov.homehub.server.storage.dao;

import io.pburakov.homehub.server.storage.model.Agent;
import org.jdbi.v3.sqlobject.config.RegisterRowMapper;
import org.jdbi.v3.sqlobject.statement.SqlQuery;
import org.jdbi.v3.sqlobject.statement.SqlUpdate;
import org.jdbi.v3.sqlobject.transaction.Transactional;

public interface HubDao extends Transactional<HubDao> {

  @SqlUpdate("insert into agents "
      + "(agent_id, address, web_port, stream_port, sensors_port) "
      + "values (?, ?, ?, ?, ?)")
  void insert(String agentId, String address, int webPort, int streamPort, int metaPort);

  @SqlUpdate("update agents "
      + "set address = ?, web_port = ?, stream_port = ?, sensors_port = ? "
      + "where agent_id = ?")
  void update(String address, int webPort, int streamPort, int metaPort, String agentId);

  @SqlQuery("select * from agents where agent_id = ?")
  @RegisterRowMapper(Agent.Mapper.class)
  Agent select(String agentId);

}
