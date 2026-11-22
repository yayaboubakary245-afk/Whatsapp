package com.whatsappcampus.webrtc;

import android.content.Context;

import com.whatsappcampus.network.NetworkConfig;
import com.whatsappcampus.network.WebSocketManager;

import org.json.JSONObject;
import org.webrtc.AudioSource;
import org.webrtc.AudioTrack;
import org.webrtc.DataChannel;
import org.webrtc.IceCandidate;
import org.webrtc.MediaConstraints;
import org.webrtc.MediaStream;
import org.webrtc.PeerConnection;
import org.webrtc.PeerConnectionFactory;
import org.webrtc.RtpReceiver;
import org.webrtc.SdpObserver;
import org.webrtc.SessionDescription;

import java.util.ArrayList;
import java.util.List;

/**
 * EXCEPTION ARCHITECTURALE ENCADRÉE (C-T-08 / C-T-11 / C-T-12) :
 * WebRTCManager.java est la SEULE classe du client Android autorisée à manipuler
 * la pile org.webrtc (qui ouvre en interne des sockets UDP pour le flux média P2P).
 * Toute la signalisation (SDP / ICE / contrôle) passe par WebSocketManager -> serveur Go.
 */
public class WebRTCManager {

    public interface CallEvents {
        void onConnected();
        void onDisconnected();
        void onError(String message);
    }

    private final PeerConnectionFactory factory;
    private PeerConnection peerConnection;
    private AudioTrack localAudioTrack;
    private final String callId;
    private final String targetUserId;
    private final CallEvents events;

    public WebRTCManager(Context context, String callId, String targetUserId, CallEvents events) {
        this.callId = callId;
        this.targetUserId = targetUserId;
        this.events = events;

        PeerConnectionFactory.initialize(
                PeerConnectionFactory.InitializationOptions.builder(context.getApplicationContext())
                        .createInitializationOptions()
        );
        this.factory = PeerConnectionFactory.builder().createPeerConnectionFactory();
        initPeerConnection();
    }

    private void initPeerConnection() {
        List<PeerConnection.IceServer> iceServers = new ArrayList<>();
        iceServers.add(PeerConnection.IceServer.builder(NetworkConfig.STUN_SERVER).createIceServer());

        PeerConnection.RTCConfiguration config = new PeerConnection.RTCConfiguration(iceServers);
        peerConnection = factory.createPeerConnection(config, new PeerConnection.Observer() {
            @Override public void onSignalingChange(PeerConnection.SignalingState state) {}
            @Override public void onIceConnectionChange(PeerConnection.IceConnectionState state) {
                if (state == PeerConnection.IceConnectionState.CONNECTED && events != null) {
                    events.onConnected();
                } else if (state == PeerConnection.IceConnectionState.DISCONNECTED && events != null) {
                    events.onDisconnected();
                }
            }
            @Override public void onIceConnectionReceivingChange(boolean b) {}
            @Override public void onIceGatheringChange(PeerConnection.IceGatheringState state) {}
            @Override public void onIceCandidate(IceCandidate candidate) {
                try {
                    JSONObject sig = new JSONObject();
                    sig.put("type", "ice_candidate");
                    sig.put("call_id", callId);
                    sig.put("target_user_id", targetUserId);
                    JSONObject cand = new JSONObject();
                    cand.put("sdpMid", candidate.sdpMid);
                    cand.put("sdpMLineIndex", candidate.sdpMLineIndex);
                    cand.put("candidate", candidate.sdp);
                    sig.put("candidate", cand);
                    WebSocketManager.getInstance().sendSignal(sig);
                } catch (Exception ignored) {}
            }
            @Override public void onIceCandidatesRemoved(IceCandidate[] candidates) {}
            @Override public void onAddStream(MediaStream stream) {}
            @Override public void onRemoveStream(MediaStream stream) {}
            @Override public void onDataChannel(DataChannel dataChannel) {}
            @Override public void onRenegotiationNeeded() {}
            @Override public void onAddTrack(RtpReceiver receiver, MediaStream[] streams) {}
        });

        AudioSource audioSource = factory.createAudioSource(new MediaConstraints());
        localAudioTrack = factory.createAudioTrack("ARDAMSa0", audioSource);
        if (peerConnection != null) {
            peerConnection.addTrack(localAudioTrack);
        }
    }

    public void createOffer() {
        if (peerConnection == null) return;
        peerConnection.createOffer(new SimpleSdpObserver() {
            @Override
            public void onCreateSuccess(SessionDescription sdp) {
                peerConnection.setLocalDescription(new SimpleSdpObserver(), sdp);
                sendSdp("call_offer", sdp.description);
            }
        }, new MediaConstraints());
    }

    public void onRemoteOfferReceived(String sdpText) {
        if (peerConnection == null) return;
        SessionDescription remote = new SessionDescription(SessionDescription.Type.OFFER, sdpText);
        peerConnection.setRemoteDescription(new SimpleSdpObserver(), remote);
        peerConnection.createAnswer(new SimpleSdpObserver() {
            @Override
            public void onCreateSuccess(SessionDescription answer) {
                peerConnection.setLocalDescription(new SimpleSdpObserver(), answer);
                sendSdp("call_answer", answer.description);
            }
        }, new MediaConstraints());
    }

    public void onRemoteAnswerReceived(String sdpText) {
        if (peerConnection == null) return;
        peerConnection.setRemoteDescription(new SimpleSdpObserver(),
                new SessionDescription(SessionDescription.Type.ANSWER, sdpText));
    }

    public void onRemoteIceCandidate(String sdpMid, int sdpMLineIndex, String candidate) {
        if (peerConnection == null) return;
        peerConnection.addIceCandidate(new IceCandidate(sdpMid, sdpMLineIndex, candidate));
    }

    public void setMuted(boolean muted) {
        if (localAudioTrack != null) {
            localAudioTrack.setEnabled(!muted);
        }
    }

    public void close() {
        if (peerConnection != null) {
            peerConnection.close();
            peerConnection = null;
        }
    }

    private void sendSdp(String type, String sdp) {
        try {
            JSONObject sig = new JSONObject();
            sig.put("type", type);
            sig.put("call_id", callId);
            sig.put("target_user_id", targetUserId);
            sig.put("sdp", sdp);
            WebSocketManager.getInstance().sendSignal(sig);
        } catch (Exception ignored) {}
    }

    private static class SimpleSdpObserver implements SdpObserver {
        @Override public void onCreateSuccess(SessionDescription sdp) {}
        @Override public void onSetSuccess() {}
        @Override public void onCreateFailure(String s) {}
        @Override public void onSetFailure(String s) {}
    }
}
