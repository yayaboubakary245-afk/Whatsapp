package com.whatsappcampus.activities;

import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.Spinner;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.models.Group;
import com.whatsappcampus.network.ApiClient;

import java.util.HashMap;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class CreateGroupActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Nouveau groupe universitaire");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(40, 40, 40, 40);

        EditText etName = new EditText(this); etName.setHint("Nom du groupe");
        EditText etDesc = new EditText(this); etDesc.setHint("Description");
        Spinner spType = new Spinner(this);
        String[] types = {"classe", "filiere", "niveau", "departement", "td", "projet", "pedagogique"};
        spType.setAdapter(new ArrayAdapter<>(this, android.R.layout.simple_spinner_dropdown_item, types));

        Button btnCreate = new Button(this);
        btnCreate.setText("Créer le groupe");
        btnCreate.setOnClickListener(v -> {
            Map<String, Object> body = new HashMap<>();
            body.put("name", etName.getText().toString().trim());
            body.put("description", etDesc.getText().toString().trim());
            body.put("group_type", spType.getSelectedItem().toString());
            ApiClient.getService(this).createGroup(body).enqueue(new Callback<Group>() {
                @Override
                public void onResponse(Call<Group> call, Response<Group> response) {
                    if (response.isSuccessful()) {
                        Toast.makeText(CreateGroupActivity.this, "Groupe créé", Toast.LENGTH_SHORT).show();
                        finish();
                    }
                }
                @Override public void onFailure(Call<Group> call, Throwable t) {}
            });
        });

        layout.addView(etName);
        layout.addView(etDesc);
        layout.addView(spType);
        layout.addView(btnCreate);
        setContentView(layout);
    }
}
