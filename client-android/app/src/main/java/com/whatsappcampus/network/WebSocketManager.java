package com.whatsappcampus.network;

import android.os.Handler;
import android.os.Looper;

import org.json.JSONObject;

import java.util.concurrent.TimeUnit;

import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;

/**
 * Gestionnaire WebSocket connecté au hub gorilla/websocket du serveur Go.
 * Assure la réception temps réel des messages et la signalisation WebRTC (SDP, ICE).
 */
public class WebSocketManager {

    public interface SignalListener {
        void onSignalReceived(JSONObject signal);
        void onConnectionChanged(boolean connected);
    }

    private static volatile WebSocketManager instance;
    private final OkHttpClient client;
    private final Handler mainHandler = new Handler(Looper.getMainLooper());
    private WebSocket webSocket;
    private SignalListener listener;

    private WebSocketManager() {
        this.client = new OkHttpClient.Builder()
                .pingInterval(25, TimeUnit.SECONDS)
                .build();
    }

    public static WebSocketManager getInstance() {
        if (instance == null) {
            synchronized (WebSocketManager.class) {
                if (instance == null) {
                    instance = new WebSocketManager();
                }
            }
        }
        return instance;
    }

    public void setListener(SignalListener listener) {
        this.listener = listener;
    }

    public void connect(String jwtToken) {
        if (jwtToken == null || jwtToken.isEmpty()) return;
        Request request = new Request.Builder()
                .url(NetworkConfig.WS_URL + "?token=" + jwtToken)
                .build();

        webSocket = client.newWebSocket(request, new WebSocketListener() {
            @Override
            public void onOpen(WebSocket ws, Response response) {
                if (listener != null) {
                    mainHandler.post(() -> listener.onConnectionChanged(true));
                }
            }

            @Override
            public void onMessage(WebSocket ws, String text) {
                try {
                    JSONObject obj = new JSONObject(text);
                    if (listener != null) {
                        mainHandler.post(() -> listener.onSignalReceived(obj));
                    }
                } catch (Exception ignored) {}
            }

            @Override
            public void onClosed(WebSocket ws, int code, String reason) {
                if (listener != null) {
                    mainHandler.post(() -> listener.onConnectionChanged(false));
                }
            }

            @Override
            public void onFailure(WebSocket ws, Throwable t, Response response) {
                if (listener != null) {
                    mainHandler.post(() -> listener.onConnectionChanged(false));
                }
            }
        });
    }

    public void sendSignal(JSONObject payload) {
        if (webSocket != null) {
            webSocket.send(payload.toString());
        }
    }

    public void disconnect() {
        if (webSocket != null) {
            webSocket.close(1000, "Logout");
            webSocket = null;
        }
    }
}
