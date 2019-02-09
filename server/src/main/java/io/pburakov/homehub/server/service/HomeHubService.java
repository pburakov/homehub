package io.pburakov.homehub.server.service;

import io.grpc.stub.StreamObserver;
import io.pburakov.homehub.schema.Ack;
import io.pburakov.homehub.schema.CheckInRequest;
import io.pburakov.homehub.schema.HomeHubGrpc;
import io.pburakov.homehub.schema.Result;
import io.pburakov.homehub.server.storage.dao.HubDao;
import io.pburakov.homehub.server.storage.model.Agent;
import org.pmw.tinylog.Logger;

public class HomeHubService extends HomeHubGrpc.HomeHubImplBase {

  private HubDao hubDao;

  public HomeHubService(HubDao hubDao) {
    this.hubDao = hubDao;
  }

  @Override
  public void checkIn(CheckInRequest request, StreamObserver<Ack> responseObserver) {
    final String agentId = request.getAgentId();
    Logger.info("Received check-in from agent '{}' reporting address '{}' and ports {}, {}, {}",
        agentId.length() > 8 ? agentId.substring(0, 8).concat("...") : agentId,
        request.getAddress(),
        request.getWebPort(),
        request.getStreamPort(),
        request.getMetaPort());
    final Ack response = hubDao.inTransaction(txHubDao -> {
      final Ack.Builder responseBuilder = Ack.newBuilder();
      final Agent agent = txHubDao.select(agentId);
      if (agent != null) {
        if (equal(request, agent)) {
          responseBuilder.setResult(Result.RECEIVED_UNCHANGED);
        } else {
          responseBuilder.setResult(Result.RECEIVED_UPDATED);
          txHubDao.update(
              request.getAddress(),
              request.getWebPort(),
              request.getStreamPort(),
              request.getMetaPort(),
              agentId);
        }
      } else {
        responseBuilder.setResult(Result.RECEIVED_NEW);
        txHubDao.insert(
            agentId,
            request.getAddress(),
            request.getWebPort(),
            request.getStreamPort(),
            request.getMetaPort());
      }
      return responseBuilder.build();
    });
    responseObserver.onNext(response);
    responseObserver.onCompleted();
  }

  private static boolean equal(CheckInRequest request, Agent entry) {
    return entry.address().equals(request.getAddress())
        && entry.webPort() == request.getWebPort()
        && entry.streamPort() == request.getStreamPort()
        && entry.metaPort() == request.getMetaPort();
  }

}
