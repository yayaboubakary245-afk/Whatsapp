package com.whatsappcampus.models;

import com.google.gson.annotations.SerializedName;

public class Conversation {
    @SerializedName("id")
    public String id;
    @SerializedName("type")
    public String type;
    @SerializedName("title")
    public String title;
    @SerializedName("peer_user_id")
    public String peerUserId;
    @SerializedName("peer_online")
    public boolean peerOnline;
    @SerializedName("last_message")
    public String lastMessage;
    @SerializedName("last_message_at")
    public String lastMessageAt;
    @SerializedName("unread_count")
    public int unreadCount;
}
