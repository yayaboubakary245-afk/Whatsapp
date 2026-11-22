package com.whatsappcampus.activities;

import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.ListView;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.models.Announcement;
import com.whatsappcampus.network.ApiClient;

import java.util.ArrayList;
import java.util.List;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class AnnouncementsActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Annonces Universitaires");

        ListView lv = new ListView(this);
        List<String> rows = new ArrayList<>();
        ArrayAdapter<String> adapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_1, rows);
        lv.setAdapter(adapter);
        setContentView(lv);

        ApiClient.getService(this).getAnnouncements().enqueue(new Callback<List<Announcement>>() {
            @Override
            public void onResponse(Call<List<Announcement>> call, Response<List<Announcement>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    rows.clear();
                    for (Announcement a : response.body()) {
                        rows.add("[" + a.category + "] " + a.title + "\n" + a.content);
                    }
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<Announcement>> call, Throwable t) {}
        });
    }
}
