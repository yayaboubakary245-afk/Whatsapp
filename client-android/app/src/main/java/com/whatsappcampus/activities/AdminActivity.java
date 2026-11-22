package com.whatsappcampus.activities;

import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.ListView;

import androidx.appcompat.app.AppCompatActivity;

import com.google.gson.JsonObject;
import com.whatsappcampus.network.ApiClient;

import java.util.ArrayList;
import java.util.List;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class AdminActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Administration SUP'PTIC");

        ListView lv = new ListView(this);
        List<String> logs = new ArrayList<>();
        ArrayAdapter<String> adapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_1, logs);
        lv.setAdapter(adapter);
        setContentView(lv);

        ApiClient.getService(this).getAuditLogs().enqueue(new Callback<List<JsonObject>>() {
            @Override
            public void onResponse(Call<List<JsonObject>> call, Response<List<JsonObject>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    logs.clear();
                    for (JsonObject entry : response.body()) {
                        logs.add(entry.get("action").getAsString() + " · " + entry.get("created_at").getAsString());
                    }
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<JsonObject>> call, Throwable t) {}
        });
    }
}
