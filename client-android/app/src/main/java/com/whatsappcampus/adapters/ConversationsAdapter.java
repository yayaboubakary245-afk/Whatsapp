package com.whatsappcampus.adapters;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.recyclerview.widget.RecyclerView;

import com.whatsappcampus.models.Conversation;

import java.util.List;

public class ConversationsAdapter extends RecyclerView.Adapter<ConversationsAdapter.ViewHolder> {

    public interface OnConversationClick {
        void onClick(Conversation conversation);
    }

    private final List<Conversation> items;
    private final OnConversationClick listener;

    public ConversationsAdapter(List<Conversation> items, OnConversationClick listener) {
        this.items = items;
        this.listener = listener;
    }

    @NonNull
    @Override
    public ViewHolder onCreateViewHolder(@NonNull ViewGroup parent, int viewType) {
        View v = LayoutInflater.from(parent.getContext()).inflate(android.R.layout.simple_list_item_2, parent, false);
        return new ViewHolder(v);
    }

    @Override
    public void onBindViewHolder(@NonNull ViewHolder holder, int position) {
        Conversation c = items.get(position);
        holder.title.setText(c.title + (c.peerOnline ? " (En ligne)" : ""));
        holder.subtitle.setText(c.lastMessage != null && !c.lastMessage.isEmpty() ? c.lastMessage : "Nouvelle discussion");
        holder.itemView.setOnClickListener(v -> listener.onClick(c));
    }

    @Override
    public int getItemCount() {
        return items.size();
    }

    static class ViewHolder extends RecyclerView.ViewHolder {
        final TextView title;
        final TextView subtitle;
        ViewHolder(View itemView) {
            super(itemView);
            title = itemView.findViewById(android.R.id.text1);
            subtitle = itemView.findViewById(android.R.id.text2);
        }
    }
}
