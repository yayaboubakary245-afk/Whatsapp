package com.whatsappcampus.models;

import com.google.gson.annotations.SerializedName;

public class User {
    @SerializedName("id")
    public String id;
    @SerializedName("user_id")
    public String userId;
    @SerializedName("email")
    public String email;
    @SerializedName("role")
    public String role;
    @SerializedName("first_name")
    public String firstName;
    @SerializedName("last_name")
    public String lastName;
    @SerializedName("bio")
    public String bio;
    @SerializedName("department")
    public String department;
    @SerializedName("program")
    public String program;
    @SerializedName("level")
    public String level;
    @SerializedName("is_online")
    public boolean isOnline;

    public String getFullName() {
        String fn = firstName != null ? firstName : "";
        String ln = lastName != null ? lastName : "";
        return (fn + " " + ln).trim();
    }
}
