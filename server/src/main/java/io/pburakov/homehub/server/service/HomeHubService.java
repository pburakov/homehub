package io.pburakov.homehub.server.service;

import io.grpc.stub.StreamObserver;
import io.pburakov.homehub.schema.Ack;
import io.pburakov.homehub.schema.CheckInRequest;
import io.pburakov.homehub.schema.HomeHubGrpc;
import io.pburakov.homehub.schema.Result;
import io.pburakov.homehub.server.storage.dao.HubDao;
import io.pburakov.homehub.server.storage.model.Hub;
import org.pmw.tinylog.Logger;

public class HomeHubService extends HomeHubGrpc.HomeHubImplBase {

  private HubDao hubDao;

  public HomeHubService(HubDao hubDao) {
    this.hubDao = hubDao;
  }

  @Override
  public void checkIn(CheckInRequest request, StreamObserver<Ack> responseObserver) {
    final String hubId = request.getHubId();
    Logger.info("Received check-in from a hub '{}' reporting address '{}' and ports {}, {}, {}",
        hubId.length() > 8 ? hubId.substring(0, 8).concat("...") : hubId,
        request.getAddress(),
        request.getWebPort(),
        request.getStreamPort(),
        request.getMetaPort());
    final Ack response = hubDao.inTransaction(txHubDao -> {
      final Ack.Builder responseBuilder = Ack.newBuilder();
      final Hub hub = txHubDao.select(hubId);
      if (hub != null) {
        if (equal(request, hub)) {
          responseBuilder.setResult(Result.RECEIVED_UNCHANGED);
        } else {
          responseBuilder.setResult(Result.RECEIVED_UPDATED);
          txHubDao.update(
              request.getAddress(),
              request.getWebPort(),
              request.getStreamPort(),
              request.getMetaPort(),
              hubId);
        }
      } else {
        responseBuilder.setResult(Result.RECEIVED_NEW);
        txHubDao.insert(
            hubId,
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

  private static boolean equal(CheckInRequest request, Hub entry) {
    return entry.address().equals(request.getAddress())
        && entry.webPort() == request.getWebPort()
        && entry.streamPort() == request.getStreamPort()
        && entry.metaPort() == request.getMetaPort();
  }

}
