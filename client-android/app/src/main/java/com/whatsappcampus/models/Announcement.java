package com.whatsappcampus.models;

import com.google.gson.annotations.SerializedName;

public class Announcement {
    @SerializedName("id")
    public String id;
    @SerializedName("author_name")
    public String authorName;
    @SerializedName("author_role")
    public String authorRole;
    @SerializedName("title")
    public String title;
    @SerializedName("content")
    public String content;
    @SerializedName("category")
    public String category;
    @SerializedName("is_pinned")
    public boolean isPinned;
    @SerializedName("published_at")
    public String publishedAt;
}
