package com.whatsappcampus.activities;

import android.os.Bundle;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.models.User;
import com.whatsappcampus.network.ApiClient;

import java.util.HashMap;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class ProfileActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Mon Profil");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(40, 40, 40, 40);

        TextView tvInfo = new TextView(this);
        EditText etFirst = new EditText(this); etFirst.setHint("Prénom");
        EditText etLast = new EditText(this); etLast.setHint("Nom");
        EditText etBio = new EditText(this); etBio.setHint("Présentation personnelle");
        Button btnSave = new Button(this); btnSave.setText("Enregistrer");

        layout.addView(tvInfo);
        layout.addView(etFirst);
        layout.addView(etLast);
        layout.addView(etBio);
        layout.addView(btnSave);
        setContentView(layout);

        ApiClient.getService(this).getProfile().enqueue(new Callback<User>() {
            @Override
            public void onResponse(Call<User> call, Response<User> response) {
                if (response.isSuccessful() && response.body() != null) {
                    User u = response.body();
                    tvInfo.setText(u.email + " · " + u.role);
                    etFirst.setText(u.firstName);
                    etLast.setText(u.lastName);
                    etBio.setText(u.bio);
                }
            }
            @Override public void onFailure(Call<User> call, Throwable t) {}
        });

        btnSave.setOnClickListener(v -> {
            Map<String, Object> body = new HashMap<>();
            body.put("first_name", etFirst.getText().toString().trim());
            body.put("last_name", etLast.getText().toString().trim());
            body.put("bio", etBio.getText().toString().trim());
            ApiClient.getService(this).updateProfile(body).enqueue(new Callback<User>() {
                @Override
                public void onResponse(Call<User> call, Response<User> response) {
                    Toast.makeText(ProfileActivity.this, "Profil mis à jour", Toast.LENGTH_SHORT).show();
                }
                @Override public void onFailure(Call<User> call, Throwable t) {}
            });
        });
    }
}
