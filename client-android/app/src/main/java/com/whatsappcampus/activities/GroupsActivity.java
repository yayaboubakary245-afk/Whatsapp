package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.ListView;

import androidx.appcompat.app.AppCompatActivity;

import com.whatsappcampus.models.Group;
import com.whatsappcampus.network.ApiClient;

import java.util.ArrayList;
import java.util.List;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class GroupsActivity extends AppCompatActivity {
    private final List<Group> groups = new ArrayList<>();
    private final List<String> labels = new ArrayList<>();
    private ArrayAdapter<String> adapter;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Groupes Universitaires");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(24, 24, 24, 24);

        Button btnCreate = new Button(this);
        btnCreate.setText("Créer un groupe (7 types)");
        btnCreate.setOnClickListener(v -> startActivity(new Intent(this, CreateGroupActivity.class)));

        ListView lv = new ListView(this);
        adapter = new ArrayAdapter<>(this, android.R.layout.simple_list_item_1, labels);
        lv.setAdapter(adapter);
        lv.setOnItemClickListener((parent, view, position, id) -> {
            Group g = groups.get(position);
            Intent i = new Intent(this, GroupInfoActivity.class);
            i.putExtra("group_name", g.name);
            i.putExtra("group_type", g.groupType);
            i.putExtra("group_desc", g.description);
            i.putExtra("conversation_id", g.conversationId);
            startActivity(i);
        });

        layout.addView(btnCreate);
        layout.addView(lv);
        setContentView(layout);
    }

    @Override
    protected void onResume() {
        super.onResume();
        ApiClient.getService(this).getGroups().enqueue(new Callback<List<Group>>() {
            @Override
            public void onResponse(Call<List<Group>> call, Response<List<Group>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    groups.clear();
                    labels.clear();
                    groups.addAll(response.body());
                    for (Group g : groups) {
                        labels.add(g.name + " [" + g.groupType + "] · " + g.memberCount + " membres");
                    }
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<Group>> call, Throwable t) {}
        });
    }
}
