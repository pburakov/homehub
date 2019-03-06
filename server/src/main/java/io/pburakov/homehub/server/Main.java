package io.pburakov.homehub.server;

import static io.pburakov.homehub.server.util.SchemaUtil.initSchema;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.protobuf.services.ProtoReflectionService;
import io.pburakov.homehub.server.service.HomeHubService;
import io.pburakov.homehub.server.storage.dao.HubDao;
import java.io.IOException;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.core.h2.H2DatabasePlugin;
import org.jdbi.v3.sqlobject.SqlObjectPlugin;
import org.pmw.tinylog.Logger;

public class Main {

  private static final int DEFAULT_PORT = 31321;

  public static void main(String... args) throws InterruptedException, IOException {
    int port = DEFAULT_PORT;

    final Jdbi jdbi = initDb();
    final HubDao hubDao = jdbi.onDemand(HubDao.class);

    Server server = ServerBuilder.forPort(port)
        .addService(new HomeHubService(hubDao))
        .addService(ProtoReflectionService.newInstance())
        .build();

    server.start();
    Logger.info("Server started, listening on port {}", port);

    server.awaitTermination();

    Runtime.getRuntime().addShutdownHook(new Thread(() -> {
      Logger.info("Server shutdown");
      server.shutdown();
    }));
  }

  private static Jdbi initDb() {
    final Jdbi jdbi = Jdbi.create("jdbc:h2:file:./homehub")
        .installPlugin(new SqlObjectPlugin())
        .installPlugin(new H2DatabasePlugin());

    initSchema(jdbi);

    return jdbi;
  }

}
