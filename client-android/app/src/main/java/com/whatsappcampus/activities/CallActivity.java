package com.whatsappcampus.activities;

import android.os.Bundle;
import android.widget.Button;
import android.widget.TextView;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.R;
import com.whatsappcampus.network.ApiClient;
import com.whatsappcampus.network.WebSocketManager;
import com.whatsappcampus.webrtc.WebRTCManager;

import org.json.JSONObject;

import java.util.HashMap;
import java.util.Map;

public class CallActivity extends AppCompatActivity {
    private WebRTCManager rtcManager;
    private boolean isMuted = false;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_call);

        String callId = getIntent().getStringExtra("call_id");
        String targetUserId = getIntent().getStringExtra("target_user_id");
        boolean isCaller = getIntent().getBooleanExtra("is_caller", true);

        TextView tvStatus = findViewById(R.id.tvCallStatus);
        Button btnMute = findViewById(R.id.btnMute);
        Button btnHangup = findViewById(R.id.btnHangup);

        rtcManager = new WebRTCManager(this, callId, targetUserId, new WebRTCManager.CallEvents() {
            @Override public void onConnected() { runOnUiThread(() -> tvStatus.setText("Connecté (P2P WebRTC)")); }
            @Override public void onDisconnected() { runOnUiThread(() -> finish()); }
            @Override public void onError(String message) { runOnUiThread(() -> tvStatus.setText(message)); }
        });

        WebSocketManager.getInstance().setListener(new WebSocketManager.SignalListener() {
            @Override
            public void onSignalReceived(JSONObject signal) {
                String type = signal.optString("type");
                if ("call_offer".equals(type)) {
                    rtcManager.onRemoteOfferReceived(signal.optString("sdp"));
                } else if ("call_answer".equals(type)) {
                    rtcManager.onRemoteAnswerReceived(signal.optString("sdp"));
                } else if ("ice_candidate".equals(type)) {
                    JSONObject cand = signal.optJSONObject("candidate");
                    if (cand != null) {
                        rtcManager.onRemoteIceCandidate(cand.optString("sdpMid"), cand.optInt("sdpMLineIndex"), cand.optString("candidate"));
                    }
                } else if ("call_end".equals(type)) {
                    finish();
                }
            }
            @Override public void onConnectionChanged(boolean connected) {}
        });

        if (isCaller) {
            rtcManager.createOffer();
        }

        btnMute.setOnClickListener(v -> {
            isMuted = !isMuted;
            rtcManager.setMuted(isMuted);
            btnMute.setText(isMuted ? "Micro coupé" : "Muet");
        });

        btnHangup.setOnClickListener(v -> {
            Map<String, Object> body = new HashMap<>();
            body.put("call_id", callId);
            body.put("status", "ended");
            ApiClient.getService(this).endCall(body).enqueue(new retrofit2.Callback<>() {
                @Override public void onResponse(retrofit2.Call<com.google.gson.JsonObject> c, retrofit2.Response<com.google.gson.JsonObject> r) { finish(); }
                @Override public void onFailure(retrofit2.Call<com.google.gson.JsonObject> c, Throwable t) { finish(); }
            });
        });
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        if (rtcManager != null) rtcManager.close();
    }
}
