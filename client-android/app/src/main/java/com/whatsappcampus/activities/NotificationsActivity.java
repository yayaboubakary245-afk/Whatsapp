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

public class NotificationsActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Notifications");

        ListView lv = new ListView(this);
        List<String> items = new ArrayList<>();
        ArrayAdapter<String> adapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_1, items);
        lv.setAdapter(adapter);
        setContentView(lv);

        ApiClient.getService(this).getNotifications().enqueue(new Callback<List<JsonObject>>() {
            @Override
            public void onResponse(Call<List<JsonObject>> call, Response<List<JsonObject>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    items.clear();
                    for (JsonObject n : response.body()) {
                        items.add(n.get("title").getAsString() + ": " + n.get("body").getAsString());
                    }
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<JsonObject>> call, Throwable t) {}
        });
    }
}
