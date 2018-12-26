package io.pburakov.homehub.server.storage.model;

import io.norberg.automatter.AutoMatter;
import java.sql.ResultSet;
import java.sql.SQLException;
import org.jdbi.v3.core.mapper.RowMapper;
import org.jdbi.v3.core.statement.StatementContext;
import org.joda.time.DateTime;

@AutoMatter
public interface Hub {

  String hubId();

  String address();

  int port();

  DateTime createdAt();

  DateTime updatedAt();

  class Mapper implements RowMapper<Hub> {

    @Override
    public Hub map(ResultSet rs, StatementContext ctx) throws SQLException {
      return new HubBuilder()
          .hubId(rs.getString("hub_id"))
          .address(rs.getString("address"))
          .port(rs.getInt("port"))
          .createdAt(new DateTime(rs.getTimestamp("created_at")))
          .updatedAt(new DateTime(rs.getTimestamp("updated_at")))
          .build();
    }
  }

}
