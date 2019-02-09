package io.pburakov.homehub.server.service;

import static com.google.common.truth.Truth.assertThat;
import static io.pburakov.homehub.server.util.SchemaUtil.initSchema;
import static org.mockito.Mockito.verify;

import io.grpc.stub.StreamObserver;
import io.pburakov.homehub.schema.Ack;
import io.pburakov.homehub.schema.CheckInRequest;
import io.pburakov.homehub.schema.Result;
import io.pburakov.homehub.server.storage.dao.HubDao;
import io.pburakov.homehub.server.storage.model.Agent;
import java.io.IOException;
import java.util.Properties;
import org.jdbi.v3.core.Jdbi;
import org.jdbi.v3.sqlobject.SqlObjectPlugin;
import org.junit.Before;
import org.junit.ClassRule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;
import org.testcontainers.containers.MySQLContainer;

@RunWith(MockitoJUnitRunner.class)
public class HomeHubServiceTest {

  @Mock
  private StreamObserver<Ack> observer;

  @ClassRule
  public static MySQLContainer mysql = new MySQLContainer();

  private HomeHubService homeHubService;
  private HubDao hubDao;

  @Before
  public void setup() throws IOException {
    final Properties p = new Properties();
    p.setProperty("user", mysql.getUsername());
    p.setProperty("password", mysql.getPassword());
    final Jdbi jdbi = Jdbi.create(mysql.getJdbcUrl(), p);
    jdbi.installPlugin(new SqlObjectPlugin());

    initSchema(jdbi);

    hubDao = jdbi.onDemand(HubDao.class);
    homeHubService = new HomeHubService(hubDao);
  }

  @Test
  public void TestCheckInFlows() throws InterruptedException {
    String agentId = "testagent123";
    String address = "test123";
    int ports = 4242;

    givenCheckInRequest(agentId, address, ports);
    verify(observer).onNext(Ack.newBuilder().setResult(Result.RECEIVED_NEW).build());

    // Verify record is stored
    Agent agent = hubDao.select(agentId);
    assertThat(agent.agentId()).isEqualTo(agentId);
    assertThat(agent.address()).isEqualTo(address);
    assertThat(agent.webPort()).isEqualTo(ports);
    assertThat(agent.streamPort()).isEqualTo(ports);
    assertThat(agent.metaPort()).isEqualTo(ports);

    // Second check-in should yield a different ack
    givenCheckInRequest(agentId, address, ports);
    verify(observer).onNext(Ack.newBuilder().setResult(Result.RECEIVED_UNCHANGED).build());

    // Check-in with new address should yield a different ack
    String newAddress = "newAddress";
    int newPorts = 4343;

    // Wait at least 1s to make sure timestamp is updated
    Thread.sleep(1000);

    givenCheckInRequest(agentId, newAddress, newPorts);
    verify(observer).onNext(Ack.newBuilder().setResult(Result.RECEIVED_UPDATED).build());

    // Updated address and port should be stored
    Agent updatedAgent = hubDao.select(agentId);
    assertThat(updatedAgent.agentId()).isEqualTo(agentId);
    assertThat(updatedAgent.address()).isEqualTo(newAddress);
    assertThat(updatedAgent.webPort()).isEqualTo(newPorts);
    assertThat(updatedAgent.streamPort()).isEqualTo(newPorts);
    assertThat(updatedAgent.metaPort()).isEqualTo(newPorts);
    assertThat(updatedAgent.updatedAt()).isGreaterThan(updatedAgent.createdAt());
  }

  private void givenCheckInRequest(String agentId, String address, int ports) {
    this.homeHubService.checkIn(
        CheckInRequest.newBuilder()
            .setAgentId(agentId)
            .setAddress(address)
            .setWebPort(ports)
            .setStreamPort(ports)
            .setMetaPort(ports)
            .build(),
        this.observer);
  }

}