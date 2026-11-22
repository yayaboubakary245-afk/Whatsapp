package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.network.NetworkConfig;
import com.whatsappcampus.network.SessionManager;
import com.whatsappcampus.network.WebSocketManager;

public class SettingsActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Paramètres");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(40, 40, 40, 40);

        SessionManager session = new SessionManager(this);
        TextView tv = new TextView(this);
        tv.setText("Compte: " + session.getEmail() + "\nServeur Go: " + NetworkConfig.BASE_URL);

        Button btnLogout = new Button(this);
        btnLogout.setText("Se déconnecter");
        btnLogout.setOnClickListener(v -> {
            WebSocketManager.getInstance().disconnect();
            session.clear();
            Intent i = new Intent(this, LoginActivity.class);
            i.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK | Intent.FLAG_ACTIVITY_CLEAR_TASK);
            startActivity(i);
        });

        layout.addView(tv);
        layout.addView(btnLogout);
        setContentView(layout);
    }
}
