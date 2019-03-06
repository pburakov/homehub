package io.pburakov.homehub.server.util;

import com.google.common.base.Charsets;
import com.google.common.io.Resources;
import java.io.IOException;
import java.net.URL;
import org.jdbi.v3.core.Jdbi;
import org.pmw.tinylog.Logger;

public class SchemaUtil {

  private SchemaUtil() {

  }

  public static void initSchema(Jdbi jdbi) {
    URL url = Resources.getResource("init.sql");
    try {
      String initSql = Resources.toString(url, Charsets.UTF_8);
      jdbi.withHandle(h -> h.execute(initSql));
      Logger.info("Initialized schema");
    } catch (IOException e) {
      throw new RuntimeException("Error reading init SQL script", e);
    }
  }

}
