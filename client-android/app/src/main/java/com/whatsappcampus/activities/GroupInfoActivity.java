package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;

import androidx.appcompat.app.AppCompatActivity;

public class GroupInfoActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setTitle("Informations du groupe");

        String name = getIntent().getStringExtra("group_name");
        String type = getIntent().getStringExtra("group_type");
        String desc = getIntent().getStringExtra("group_desc");
        String convId = getIntent().getStringExtra("conversation_id");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(40, 40, 40, 40);

        TextView tvTitle = new TextView(this);
        tvTitle.setText(name + " (" + type + ")");
        tvTitle.setTextSize(20f);

        TextView tvDesc = new TextView(this);
        tvDesc.setText(desc != null ? desc : "Aucune description");

        Button btnOpenChat = new Button(this);
        btnOpenChat.setText("Ouvrir la discussion du groupe");
        btnOpenChat.setOnClickListener(v -> {
            Intent i = new Intent(this, ChatActivity.class);
            i.putExtra("conversation_id", convId);
            i.putExtra("title", name);
            startActivity(i);
        });

        layout.addView(tvTitle);
        layout.addView(tvDesc);
        layout.addView(btnOpenChat);
        setContentView(layout);
    }
}
