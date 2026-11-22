package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.EditText;
import android.widget.TextView;

import androidx.appcompat.app.AppCompatActivity;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.google.gson.JsonObject;
import com.whatsappcampus.R;
import com.whatsappcampus.adapters.MessagesAdapter;
import com.whatsappcampus.models.Message;
import com.whatsappcampus.network.ApiClient;
import com.whatsappcampus.network.SessionManager;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class ChatActivity extends AppCompatActivity {
    private final List<Message> messages = new ArrayList<>();
    private MessagesAdapter adapter;
    private String conversationId;
    private String peerUserId;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_chat);

        conversationId = getIntent().getStringExtra("conversation_id");
        peerUserId = getIntent().getStringExtra("peer_user_id");
        String title = getIntent().getStringExtra("title");

        TextView tvTitle = findViewById(R.id.tvChatTitle);
        tvTitle.setText(title != null ? title : "Discussion");

        RecyclerView rv = findViewById(R.id.rvMessages);
        rv.setLayoutManager(new LinearLayoutManager(this));
        adapter = new MessagesAdapter(messages, new SessionManager(this).getUserId());
        rv.setAdapter(adapter);

        EditText etMessage = findViewById(R.id.etMessage);
        findViewById(R.id.btnSend).setOnClickListener(v -> {
            String text = etMessage.getText().toString().trim();
            if (text.isEmpty()) return;
            Map<String, Object> body = new HashMap<>();
            body.put("conversation_id", conversationId);
            body.put("message_type", "text");
            body.put("content", text);
            ApiClient.getService(this).sendMessage(body).enqueue(new Callback<Message>() {
                @Override
                public void onResponse(Call<Message> call, Response<Message> response) {
                    if (response.isSuccessful() && response.body() != null) {
                        etMessage.setText("");
                        messages.add(response.body());
                        adapter.notifyItemInserted(messages.size() - 1);
                    }
                }
                @Override public void onFailure(Call<Message> call, Throwable t) {}
            });
        });

        findViewById(R.id.btnAudioCall).setOnClickListener(v -> {
            if (peerUserId == null || peerUserId.isEmpty()) return;
            Map<String, Object> body = new HashMap<>();
            body.put("target_user_id", peerUserId);
            body.put("conversation_id", conversationId);
            body.put("call_type", "audio");
            ApiClient.getService(this).startCall(body).enqueue(new Callback<JsonObject>() {
                @Override
                public void onResponse(Call<JsonObject> call, Response<JsonObject> response) {
                    if (response.isSuccessful() && response.body() != null) {
                        Intent i = new Intent(ChatActivity.this, CallActivity.class);
                        i.putExtra("call_id", response.body().get("call_id").getAsString());
                        i.putExtra("target_user_id", peerUserId);
                        i.putExtra("is_caller", true);
                        startActivity(i);
                    }
                }
                @Override public void onFailure(Call<JsonObject> call, Throwable t) {}
            });
        });

        loadMessages();
    }

    private void loadMessages() {
        ApiClient.getService(this).getMessages(conversationId, "").enqueue(new Callback<List<Message>>() {
            @Override
            public void onResponse(Call<List<Message>> call, Response<List<Message>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    messages.clear();
                    messages.addAll(response.body());
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<Message>> call, Throwable t) {}
        });
    }
}
