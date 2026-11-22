package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.google.gson.JsonObject;
import com.whatsappcampus.network.ApiClient;
import com.whatsappcampus.network.NetworkConfig;
import com.whatsappcampus.network.SessionManager;

import java.util.HashMap;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class RegisterActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(48, 48, 48, 48);

        EditText etFirst = new EditText(this); etFirst.setHint("Prénom");
        EditText etLast = new EditText(this); etLast.setHint("Nom");
        EditText etEmail = new EditText(this); etEmail.setHint("Email (@ensptic.cm)");
        EditText etPass = new EditText(this); etPass.setHint("Mot de passe (min 6)");
        Button btnSubmit = new Button(this); btnSubmit.setText("Créer le compte");

        layout.addView(etFirst);
        layout.addView(etLast);
        layout.addView(etEmail);
        layout.addView(etPass);
        layout.addView(btnSubmit);
        setContentView(layout);

        btnSubmit.setOnClickListener(v -> {
            String email = etEmail.getText().toString().trim().toLowerCase();
            if (!email.endsWith(NetworkConfig.ALLOWED_DOMAIN)) {
                Toast.makeText(this, "Email @ensptic.cm obligatoire", Toast.LENGTH_SHORT).show();
                return;
            }
            Map<String, Object> body = new HashMap<>();
            body.put("first_name", etFirst.getText().toString().trim());
            body.put("last_name", etLast.getText().toString().trim());
            body.put("email", email);
            body.put("password", etPass.getText().toString());
            body.put("role", "etudiant");

            ApiClient.getService(this).register(body).enqueue(new Callback<JsonObject>() {
                @Override
                public void onResponse(Call<JsonObject> call, Response<JsonObject> response) {
                    if (response.isSuccessful() && response.body() != null) {
                        JsonObject res = response.body();
                        new SessionManager(RegisterActivity.this).saveSession(
                                res.get("token").getAsString(),
                                res.get("user_id").getAsString(),
                                email,
                                res.get("role").getAsString()
                        );
                        startActivity(new Intent(RegisterActivity.this, ConversationsActivity.class));
                        finish();
                    } else {
                        Toast.makeText(RegisterActivity.this, "Échec de l'inscription", Toast.LENGTH_SHORT).show();
                    }
                }
                @Override
                public void onFailure(Call<JsonObject> call, Throwable t) {
                    Toast.makeText(RegisterActivity.this, "Erreur réseau", Toast.LENGTH_SHORT).show();
                }
            });
        });
    }
}
