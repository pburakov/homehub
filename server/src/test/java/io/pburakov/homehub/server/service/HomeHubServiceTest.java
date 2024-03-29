package io.pburakov.homehub.server.service;

import static com.google.common.truth.Truth.assertThat;
import static io.pburakov.homehub.server.util.SchemaUtil.initSchema;

import io.grpc.Channel;
import io.grpc.Server;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.pburakov.homehub.schema.Ack;
import io.pburakov.homehub.schema.CheckInRequest;
import io.pburakov.homehub.schema.HomeHubGrpc;
import io.pburakov.homehub.schema.Result;
import io.pburakov.homehub.server.storage.dao.HubDao;
import io.pburakov.homehub.server.storage.model.Agent;
import java.io.File;
import java.io.IOException;
import java.util.Properties;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.core.h2.H2DatabasePlugin;
import org.jdbi.v3.sqlobject.SqlObjectPlugin;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

class HomeHubServiceTest {

  private static final String SERVER_NAME = "test-server";

  private static Channel channel;
  private static Server server;
  private static HomeHubGrpc.HomeHubBlockingStub homeHubService;
  private static HubDao hubDao;

  @BeforeAll
  static void setup() throws IOException {
    final Properties p = new Properties();
    final Jdbi jdbi =
        Jdbi.create("jdbc:h2:file:./test")
            .installPlugin(new SqlObjectPlugin())
            .installPlugin(new H2DatabasePlugin());

    initSchema(jdbi);

    hubDao = jdbi.onDemand(HubDao.class);
    server =
        InProcessServerBuilder.forName(SERVER_NAME)
            .directExecutor()
            .addService(new HomeHubService(hubDao))
            .build()
            .start();
    channel = InProcessChannelBuilder.forName(SERVER_NAME).directExecutor().build();
    homeHubService = HomeHubGrpc.newBlockingStub(channel);
  }

  @AfterAll
  static void destroy() {
    server.shutdownNow();
    final File file = new File("test.mv.db");
    if (!file.delete()) {
      throw new RuntimeException("Unable to cleanup test db file");
    }
  }

  @Test
  void TestCheckInFlows() throws InterruptedException {
    String agentId = "testagent123";
    String address = "test123";
    int ports = 4242;

    Ack response1 = givenCheckInRequest(agentId, address, ports);
    assertThat(response1).isEqualTo(Ack.newBuilder().setResult(Result.RECEIVED_NEW).build());

    // Verify record is stored
    Agent agent = hubDao.select(agentId);
    assertThat(agent.agentId()).isEqualTo(agentId);
    assertThat(agent.address()).isEqualTo(address);
    assertThat(agent.webPort()).isEqualTo(ports);
    assertThat(agent.streamPort()).isEqualTo(ports);
    assertThat(agent.sensorsPort()).isEqualTo(ports);

    // Second check-in should yield a different ack
    Ack response2 = givenCheckInRequest(agentId, address, ports);
    assertThat(response2).isEqualTo(Ack.newBuilder().setResult(Result.RECEIVED_UNCHANGED).build());

    // Check-in with new address should yield a different ack
    String newAddress = "newAddress";
    int newPorts = 4343;

    // Wait at least 1s to make sure timestamp is updated
    Thread.sleep(1000);

    Ack response3 = givenCheckInRequest(agentId, newAddress, newPorts);
    assertThat(response3).isEqualTo(Ack.newBuilder().setResult(Result.RECEIVED_UPDATED).build());

    // Updated address and port should be stored
    Agent updatedAgent = hubDao.select(agentId);
    assertThat(updatedAgent.agentId()).isEqualTo(agentId);
    assertThat(updatedAgent.address()).isEqualTo(newAddress);
    assertThat(updatedAgent.webPort()).isEqualTo(newPorts);
    assertThat(updatedAgent.streamPort()).isEqualTo(newPorts);
    assertThat(updatedAgent.sensorsPort()).isEqualTo(newPorts);
    assertThat(updatedAgent.updatedAt()).isGreaterThan(updatedAgent.createdAt());
  }

  private static Ack givenCheckInRequest(String agentId, String address, int ports) {
    return homeHubService.checkIn(
        CheckInRequest.newBuilder()
            .setAgentId(agentId)
            .setAddress(address)
            .setWebPort(ports)
            .setStreamPort(ports)
            .setSensorsPort(ports)
            .build());
  }
}
