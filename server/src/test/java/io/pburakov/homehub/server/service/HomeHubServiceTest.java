package io.pburakov.homehub.server.service;

import static com.google.common.truth.Truth.assertThat;
import static io.pburakov.homehub.server.util.SchemaUtil.initSchema;
import static org.mockito.Mockito.verify;

import io.grpc.stub.StreamObserver;
import io.pburakov.homehub.schema.Ack;
import io.pburakov.homehub.schema.CheckInRequest;
import io.pburakov.homehub.schema.Result;
import io.pburakov.homehub.server.storage.dao.HubDao;
import io.pburakov.homehub.server.storage.model.Hub;
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
  public void TestCheckInFlows() {
    String hubId = "testId";
    String address = "test123";
    int port = 4242;

    givenCheckInRequest(hubId, address, port);
    verify(observer).onNext(Ack.newBuilder().setResult(Result.RECEIVED_NEW).build());

    // Verify record is stored
    Hub hub = hubDao.select(hubId);
    assertThat(hub.hubId()).isEqualTo(hubId);
    assertThat(hub.address()).isEqualTo(address);
    assertThat(hub.port()).isEqualTo(port);

    // Second check-in should yield a different ack
    givenCheckInRequest(hubId, address, port);
    verify(observer).onNext(Ack.newBuilder().setResult(Result.RECEIVED_UNCHANGED).build());

    // Check-in with new address should yield a different ack
    String newAddress = "newAddress";
    int newPort = 4343;

    givenCheckInRequest(hubId, newAddress, newPort);
    verify(observer).onNext(Ack.newBuilder().setResult(Result.RECEIVED_UPDATED).build());

    // Updated address and port should be stored
    Hub updatedHub = hubDao.select(hubId);
    assertThat(updatedHub.hubId()).isEqualTo(hubId);
    assertThat(updatedHub.address()).isEqualTo(newAddress);
    assertThat(updatedHub.port()).isEqualTo(newPort);
    assertThat(updatedHub.updatedAt()).isGreaterThan(updatedHub.createdAt());
  }

  private void givenCheckInRequest(String hubId, String address, int port) {
    this.homeHubService.checkIn(
        CheckInRequest.newBuilder()
            .setHubId(hubId)
            .setAddress(address)
            .setPort(port)
            .build(),
        this.observer);
  }

}