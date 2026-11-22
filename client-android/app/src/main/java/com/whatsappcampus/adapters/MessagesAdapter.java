package com.whatsappcampus.adapters;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import com.whatsappcampus.models.Message;

import java.util.List;

public class MessagesAdapter extends RecyclerView.Adapter<MessagesAdapter.ViewHolder> {

    private final List<Message> messages;
    private final String currentUserId;

    public MessagesAdapter(List<Message> messages, String currentUserId) {
        this.messages = messages;
        this.currentUserId = currentUserId;
    }

    @NonNull
    @Override
    public ViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View v = LayoutInflater.from(parent.getContext()).inflate(android.R.layout.simple_list_item_2, parent, false);
        return new ViewHolder(v);
    }

    @Override
    public void onBindViewHolder(@NonNull ViewHolder holder, int position) {
        Message m = messages.get(position);
        boolean mine = m.senderId != null && m.senderId.equals(currentUserId);
        String prefix = mine ? "Moi" : (m.senderName != null ? m.senderName : "Campus");
        String tick = mine ? (" [" + (m.status != null ? m.status : "sent") + "]") : "";
        holder.sender.setText(prefix + tick);
        holder.content.setText(m.content);
    }

    @Override
    public int getItemCount() {
        return messages.size();
    }

    static class ViewHolder extends RecyclerView.ViewHolder {
        final TextView sender;
        final TextView content;
        ViewHolder(View itemView) {
            super(itemView);
            sender = itemView.findViewById(android.R.id.text1);
            content = itemView.findViewById(android.R.id.text2);
        }
    }
}
