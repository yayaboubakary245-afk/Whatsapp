package com.whatsappcampus.models;

import com.google.gson.annotations.SerializedName;

public class Message {
    @SerializedName("id")
    public String id;
    @SerializedName("conversation_id")
    public String conversationId;
    @SerializedName("sender_id")
    public String senderId;
    @SerializedName("sender_name")
    public String senderName;
    @SerializedName("reply_to_message_id")
    public String replyToMessageId;
    @SerializedName("message_type")
    public String messageType;
    @SerializedName("content")
    public String content;
    @SerializedName("status")
    public String status;
    @SerializedName("created_at")
    public String createdAt;
}
