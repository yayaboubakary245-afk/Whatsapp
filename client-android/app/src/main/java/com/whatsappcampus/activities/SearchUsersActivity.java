package com.whatsappcampus.activities;

import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

import com.google.gson.JsonObject;
import com.whatsappcampus.models.User;
import com.whatsappcampus.network.ApiClient;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class SearchUsersActivity extends AppCompatActivity {
    private final List<User> found = new ArrayList<>();
    private final List<String> labels = new ArrayList<>();
    private ArrayAdapter<String> adapter;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Rechercher un membre");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(24, 24, 24, 24);

        EditText etQuery = new EditText(this);
        etQuery.setHint("Nom ou e-mail @ensptic.cm");
        Button btnSearch = new Button(this);
        btnSearch.setText("Rechercher");

        ListView lv = new ListView(this);
        adapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_1, labels);
        lv.setAdapter(adapter);

        lv.setOnItemClickListener((parent, view, pos, id) -> {
            User target = found.get(pos);
            Map<String, Object> body = new HashMap<>();
            body.put("contact_user_id", target.id);
            body.put("nickname", target.getFullName());
            ApiClient.getService(this).addContact(body).enqueue(new Callback<JsonObject>() {
                @Override
                public void onResponse(Call<JsonObject> call, Response<JsonObject> response) {
                    Toast.makeText(SearchUsersActivity.this, "Contact ajouté", Toast.LENGTH_SHORT).show();
                }
                @Override public void onFailure(Call<JsonObject> call, Throwable t) {}
            });
        });

        btnSearch.setOnClickListener(v -> ApiClient.getService(this).searchUsers(etQuery.getText().toString().trim()).enqueue(new Callback<List<User>>() {
            @Override
            public void onResponse(Call<List<User>> call, Response<List<User>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    found.clear();
                    labels.clear();
                    found.addAll(response.body());
                    for (User u : found) {
                        labels.add(u.getFullName() + " · " + u.role + " · " + u.email);
                    }
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<User>> call, Throwable t) {}
        }));

        layout.addView(etQuery);
        layout.addView(btnSearch);
        layout.addView(lv);
        setContentView(layout);
    }
}
