package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;

import androidx.appcompat.app.AppCompatActivity;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

import com.whatsappcampus.R;
import com.whatsappcampus.adapters.ConversationsAdapter;
import com.whatsappcampus.models.Conversation;
import com.whatsappcampus.network.ApiClient;
import com.whatsappcampus.network.SessionManager;
import com.whatsappcampus.network.WebSocketManager;

import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class ConversationsActivity extends AppCompatActivity {
    private final List<Conversation> list = new ArrayList<>();
    private ConversationsAdapter adapter;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_conversations);
        setTitle("WhatsApp Campus");

        RecyclerView rv = findViewById(R.id.rvConversations);
        rv.setLayoutManager(new LinearLayoutManager(this));
        adapter = new ConversationsAdapter(list, conv -> {
            Intent i = new Intent(this, ChatActivity.class);
            i.putExtra("conversation_id", conv.id);
            i.putExtra("title", conv.title);
            i.putExtra("peer_user_id", conv.peerUserId);
            startActivity(i);
        });
        rv.setAdapter(adapter);

        findViewById(R.id.btnGroups).setOnClickListener(v -> startActivity(new Intent(this, GroupsActivity.class)));
        findViewById(R.id.btnAnnouncements).setOnClickListener(v -> startActivity(new Intent(this, AnnouncementsActivity.class)));
        findViewById(R.id.btnContacts).setOnClickListener(v -> startActivity(new Intent(this, ContactsActivity.class)));

        SessionManager session = new SessionManager(this);
        WebSocketManager.getInstance().connect(session.getToken());
        WebSocketManager.getInstance().setListener(new WebSocketManager.SignalListener() {
            @Override
            public void onSignalReceived(JSONObject signal) {
                if ("incoming_call".equals(signal.optString("type"))) {
                    Intent callIntent = new Intent(ConversationsActivity.this, IncomingCallActivity.class);
                    callIntent.putExtra("call_id", signal.optString("call_id"));
                    callIntent.putExtra("from_user_id", signal.optString("from_user_id"));
                    callIntent.putExtra("call_type", signal.optString("call_type", "audio"));
                    startActivity(callIntent);
                } else {
                    loadConversations();
                }
            }
            @Override public void onConnectionChanged(boolean connected) {}
        });

        loadConversations();
    }

    private void loadConversations() {
        ApiClient.getService(this).getConversations().enqueue(new Callback<List<Conversation>>() {
            @Override
            public void onResponse(Call<List<Conversation>> call, Response<List<Conversation>> response) {
                if (response.isSuccessful() && response.body() != null) {
                    list.clear();
                    list.addAll(response.body());
                    adapter.notifyDataSetChanged();
                }
            }
            @Override public void onFailure(Call<List<Conversation>> call, Throwable t) {}
        });
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        menu.add(0, 1, 0, "Profil");
        menu.add(0, 2, 0, "Notifications");
        menu.add(0, 3, 0, "Paramètres");
        if ("admin".equals(new SessionManager(this).getRole())) {
            menu.add(0, 4, 0, "Administration");
        }
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        if (item.getItemId() == 1) startActivity(new Intent(this, ProfileActivity.class));
        if (item.getItemId() == 2) startActivity(new Intent(this, NotificationsActivity.class));
        if (item.getItemId() == 3) startActivity(new Intent(this, SettingsActivity.class));
        if (item.getItemId() == 4) startActivity(new Intent(this, AdminActivity.class));
        return super.onOptionsItemSelected(item);
    }
}
