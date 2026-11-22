package com.whatsappcampus.activities;

import android.content.Intent;
import android.os.Bundle;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;

import androidx.appcompat.app.AppCompatActivity;

import com.google.gson.JsonObject;
import com.whatsappcampus.network.ApiClient;

import java.util.HashMap;
import java.util.Map;

import retrofit2.Call;
import retrofit2.Callback;
import retrofit2.Response;

public class IncomingCallActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        String callId = getIntent().getStringExtra("call_id");
        String fromUserId = getIntent().getStringExtra("from_user_id");
        String callType = getIntent().getStringExtra("call_type");

        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);
        layout.setPadding(48, 96, 48, 48);

        TextView tv = new TextView(this);
        tv.setText("Appel entrant (" + callType + ") de " + fromUserId);
        tv.setTextSize(20f);

        Button btnAccept = new Button(this); btnAccept.setText("Accepter");
        Button btnReject = new Button(this); btnReject.setText("Refuser");

        layout.addView(tv);
        layout.addView(btnAccept);
        layout.addView(btnReject);
        setContentView(layout);

        btnAccept.setOnClickListener(v -> {
            Map<String, Object> body = new HashMap<>();
            body.put("call_id", callId);
            ApiClient.getService(this).acceptCall(body).enqueue(new Callback<JsonObject>() {
                @Override
                public void onResponse(Call<JsonObject> call, Response<JsonObject> response) {
                    Intent i = new Intent(IncomingCallActivity.this, CallActivity.class);
                    i.putExtra("call_id", callId);
                    i.putExtra("target_user_id", fromUserId);
                    i.putExtra("is_caller", false);
                    startActivity(i);
                    finish();
                }
                @Override public void onFailure(Call<JsonObject> call, Throwable t) { finish(); }
            });
        });

        btnReject.setOnClickListener(v -> {
            Map<String, Object> body = new HashMap<>();
            body.put("call_id", callId);
            body.put("status", "rejected");
            ApiClient.getService(this).endCall(body).enqueue(new Callback<JsonObject>() {
                @Override public void onResponse(Call<JsonObject> c, Response<JsonObject> r) { finish(); }
                @Override public void onFailure(Call<JsonObject> c, Throwable t) { finish(); }
            });
        });
    }
}
