package com.whatsappcampus.network;

import android.content.Context;
import android.content.SharedPreferences;

public class SessionManager {
    private static final String PREF_NAME = "whatsapp_campus_session";
    private static final String KEY_TOKEN = "jwt_token";
    private static final String KEY_USER_ID = "user_id";
    private static final String KEY_EMAIL = "email";
    private static final String KEY_ROLE = "role";

    private final SharedPreferences prefs;

    public SessionManager(Context context) {
        this.prefs = context.getApplicationContext().getSharedPreferences(PREF_NAME, Context.MODE_PRIVATE);
    }

    public void saveSession(String token, String userId, String email, String role) {
        prefs.edit()
                .putString(KEY_TOKEN, token)
                .putString(KEY_USER_ID, userId)
                .putString(KEY_EMAIL, email)
                .putString(KEY_ROLE, role)
                .apply();
    }

    public String getToken() { return prefs.getString(KEY_TOKEN, null); }
    public String getUserId() { return prefs.getString(KEY_USER_ID, null); }
    public String getEmail() { return prefs.getString(KEY_EMAIL, ""); }
    public String getRole() { return prefs.getString(KEY_ROLE, "etudiant"); }
    public boolean isLoggedIn() { return getToken() != null && !getToken().isEmpty(); }

    public void clear() {
        prefs.edit().clear().apply();
    }
}
