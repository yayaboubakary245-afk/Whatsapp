package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.google.gson.JsonObject;
import com.whatsappcampus.R;
import com.whatsappcampus.network.ApiClient;
import com.whatsappcampus.network.NetworkConfig;
import com.whatsappcampus.network.SessionManager;
import com.whatsappcampus.network.WebSocketManager;

import java.util.HashMap;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class LoginActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_login);

        EditText etEmail = findViewById(R.id.etEmail);
        EditText etPassword = findViewById(R.id.etPassword);
        Button btnLogin = findViewById(R.id.btnLogin);
        Button btnGoRegister = findViewById(R.id.btnGoRegister);
        Button btnTestServer = findViewById(R.id.btnTestServer);
        TextView tvServerStatus = findViewById(R.id.tvServerStatus);

        btnTestServer.setOnClickListener(v -> {
            tvServerStatus.setText("Vérification du serveur Go...");
            ApiClient.getService(this).checkHealth().enqueue(new Callback<JsonObject>() {
                @Override
                public void onResponse(Call<JsonObject> call, Response<JsonObject> response) {
                    if (response.isSuccessful() && response.body() != null) {
                        boolean pg = response.body().get("postgres").getAsBoolean();
                        tvServerStatus.setText("Serveur Go OK · PostgreSQL: " + (pg ? "connecté" : "inactif"));
                    } else {
                        tvServerStatus.setText("Réponse HTTP " + response.code());
                    }
                }
                @Override
                public void onFailure(Call<JsonObject> call, Throwable t) {
                    tvServerStatus.setText("Échec connexion Go: " + t.getMessage());
                }
            });
        });

        btnGoRegister.setOnClickListener(v -> startActivity(new Intent(this, RegisterActivity.class)));

        btnLogin.setOnClickListener(v -> {
            String email = etEmail.getText().toString().trim().toLowerCase();
            String password = etPassword.getText().toString();
            if (!email.endsWith(NetworkConfig.ALLOWED_DOMAIN)) {
                Toast.makeText(this, "Utilisez votre adresse @ensptic.cm", Toast.LENGTH_SHORT).show();
                return;
            }
            Map<String, Object> body = new HashMap<>();
            body.put("email", email);
            body.put("password", password);

            ApiClient.getService(this).login(body).enqueue(new Callback<JsonObject>() {
                @Override
                public void onResponse(Call<JsonObject> call, Response<JsonObject> response) {
                    if (response.isSuccessful() && response.body() != null) {
                        JsonObject res = response.body();
                        String token = res.get("token").getAsString();
                        String userId = res.get("user_id").getAsString();
                        String role = res.get("role").getAsString();
                        new SessionManager(LoginActivity.this).saveSession(token, userId, email, role);
                        WebSocketManager.getInstance().connect(token);
                        startActivity(new Intent(LoginActivity.this, ConversationsActivity.class));
                        finish();
                    } else {
                        Toast.makeText(LoginActivity.this, "Identifiants invalides", Toast.LENGTH_SHORT).show();
                    }
                }
                @Override
                public void onFailure(Call<JsonObject> call, Throwable t) {
                    Toast.makeText(LoginActivity.this, "Serveur Go inaccessible: " + t.getMessage(), Toast.LENGTH_SHORT).show();
                }
            });
        });
    }
}
