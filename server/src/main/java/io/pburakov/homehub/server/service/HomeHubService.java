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
    Logger.info("Received check-in from a hub '{}' reporting address '{}:{}'",
                request.getHubId(),
                request.getAddress(),
                request.getPort());
    final Ack response = hubDao.inTransaction(txHubDao -> {
      final Ack.Builder responseBuilder = Ack.newBuilder();
      final Hub hub = txHubDao.select(request.getHubId());
      if (hub != null) {
        if (hub.address().equals(request.getAddress()) && hub.port() == request.getPort()) {
          responseBuilder.setResult(Result.RECEIVED_UNCHANGED);
        } else {
          responseBuilder.setResult(Result.RECEIVED_UPDATED);
          txHubDao.update(request.getAddress(), request.getPort(), hub.hubId());
        }
      } else {
        responseBuilder.setResult(Result.RECEIVED_NEW);
        txHubDao.insert(request.getHubId(), request.getAddress(), request.getPort());
      }
      return responseBuilder.build();
    });
    responseObserver.onNext(response);
    responseObserver.onCompleted();
  }

}
