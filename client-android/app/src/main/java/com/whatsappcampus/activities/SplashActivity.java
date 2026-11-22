package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.R;
import com.whatsappcampus.network.SessionManager;

public class SplashActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_splash);
        SessionManager session = new SessionManager(this);
        new Handler(Looper.getMainLooper()).postDelayed(() -> {
            Class<?> next = session.isLoggedIn() ? ConversationsActivity.class : LoginActivity.class;
            startActivity(new Intent(this, next));
            finish();
        }, 700);
    }
}
