package com.whatsappcampus.network;

/**
 * Configuration réseau du client Android vers le serveur centralisé Go.
 * Toute la logique métier passe par HTTP/REST et WebSocket.
 */
public final class NetworkConfig {
    // 10.0.2.2 pointe vers localhost de la machine hôte depuis l'émulateur Android Studio
    public static final String BASE_URL = "http://10.0.2.2:8080/";
    public static final String WS_URL = "ws://10.0.2.2:8080/api/ws";
    public static final String ALLOWED_DOMAIN = "@ensptic.cm";
    public static final String STUN_SERVER = "stun:stun.l.google.com:19302";

    private NetworkConfig() {}
}
