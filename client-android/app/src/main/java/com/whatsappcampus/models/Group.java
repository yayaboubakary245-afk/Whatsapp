package com.whatsappcampus.models;

import com.google.gson.annotations.SerializedName;

public class Group {
    @SerializedName("id")
    public String id;
    @SerializedName("conversation_id")
    public String conversationId;
    @SerializedName("name")
    public String name;
    @SerializedName("description")
    public String description;
    @SerializedName("group_type")
    public String groupType;
    @SerializedName("only_admins_can_post")
    public boolean onlyAdminsCanPost;
    @SerializedName("member_count")
    public int memberCount;
    @SerializedName("created_at")
    public String createdAt;
}
