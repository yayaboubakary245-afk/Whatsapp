package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.ListView;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.models.User;
import com.whatsappcampus.network.ApiClient;

import java.util.ArrayList;
import java.util.List;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class ContactsActivity extends AppCompatActivity {
    private final List<String> rows = new ArrayList<>();
    private ArrayAdapter<String> adapter;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Contacts Campus");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(24, 24, 24, 24);

        Button btnSearch = new Button(this);
        btnSearch.setText("Rechercher des utilisateurs @ensptic.cm");
        btnSearch.setOnClickListener(v -> startActivity(new Intent(this, SearchUsersActivity.class)));

        ListView lv = new ListView(this);
        adapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_1, rows);
        lv.setAdapter(adapter);

        layout.addView(btnSearch);
        layout.addView(lv);
        setContentView(layout);
    }

    @Override
    protected void onResume() {
        super.onResume();
        ApiClient.getService(this).getContacts().enqueue(new Callback<List<User>>() {
            @Override
            public void onResponse(Call<List<User>> call, Response<List<User>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    rows.clear();
                    for (User u : response.body()) {
                        rows.add(u.getFullName() + " (" + u.email + ")");
                    }
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<User>> call, Throwable t) {}
        });
    }
}
