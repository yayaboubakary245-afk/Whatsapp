package com.whatsappcampus.network;

import android.content.Context;

import com.google.gson.JsonObject;
import com.whatsappcampus.models.Announcement;
import com.whatsappcampus.models.Conversation;
import com.whatsappcampus.models.Group;
import com.whatsappcampus.models.Message;
import com.whatsappcampus.models.User;

import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import okhttp3.MultipartBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import retrofit2.Call;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.http.Body;
import retrofit2.http.GET;
import retrofit2.http.Multipart;
import retrofit2.http.POST;
import retrofit2.http.PUT;
import retrofit2.http.Part;
import retrofit2.http.Path;
import retrofit2.http.Query;

public final class ApiClient {

    public interface CampusApiService {
        @GET("api/health")
        Call<JsonObject> checkHealth();

        @POST("api/register")
        Call<JsonObject> register(@Body Map<String, Object> body);

        @POST("api/login")
        Call<JsonObject> login(@Body Map<String, Object> body);

        @GET("api/profile")
        Call<User> getProfile();

        @PUT("api/profile")
        Call<User> updateProfile(@Body Map<String, Object> body);

        @GET("api/users/search")
        Call<List<User>> searchUsers(@Query("q") String query);

        @GET("api/contacts")
        Call<List<User>> getContacts();

        @POST("api/contacts")
        Call<JsonObject> addContact(@Body Map<String, Object> body);

        @GET("api/conversations")
        Call<List<Conversation>> getConversations();

        @GET("api/messages")
        Call<List<Message>> getMessages(@Query("conversation_id") String conversationId, @Query("q") String search);

        @POST("api/messages")
        Call<Message> sendMessage(@Body Map<String, Object> body);

        @GET("api/groups")
        Call<List<Group>> getGroups();

        @POST("api/groups")
        Call<Group> createGroup(@Body Map<String, Object> body);

        @Multipart
        @POST("api/attachments/upload")
        Call<JsonObject> uploadAttachment(@Part MultipartBody.Part file, @Part("message_id") RequestBody messageId);

        @GET("api/announcements")
        Call<List<Announcement>> getAnnouncements();

        @POST("api/announcements")
        Call<Announcement> publishAnnouncement(@Body Map<String, Object> body);

        @GET("api/notifications")
        Call<List<JsonObject>> getNotifications();

        @POST("api/call/start")
        Call<JsonObject> startCall(@Body Map<String, Object> body);

        @POST("api/call/accept")
        Call<JsonObject> acceptCall(@Body Map<String, Object> body);

        @POST("api/call/end")
        Call<JsonObject> endCall(@Body Map<String, Object> body);

        @POST("api/admin/users/{id}/suspend")
        Call<JsonObject> suspendUser(@Path("id") String userId, @Body Map<String, Object> body);

        @GET("api/admin/audit-logs")
        Call<List<JsonObject>> getAuditLogs();
    }

    private static volatile CampusApiService instance;

    private ApiClient() {}

    public static CampusApiService getService(Context context) {
        if (instance == null) {
            synchronized (ApiClient.class) {
                if (instance == null) {
                    SessionManager session = new SessionManager(context);
                    OkHttpClient client = new OkHttpClient.Builder()
                            .connectTimeout(15, TimeUnit.SECONDS)
                            .readTimeout(20, TimeUnit.SECONDS)
                            .addInterceptor(chain -> {
                                Request original = chain.request();
                                Request.Builder builder = original.newBuilder();
                                String token = session.getToken();
                                if (token != null && !token.isEmpty()) {
                                    builder.header("Authorization", "Bearer " + token);
                                }
                                return chain.proceed(builder.build());
                            })
                            .build();

                    Retrofit retrofit = new Retrofit.Builder()
                            .baseUrl(NetworkConfig.BASE_URL)
                            .client(client)
                            .addConverterFactory(GsonConverterFactory.create())
                            .build();

                    instance = retrofit.create(CampusApiService.class);
                }
            }
        }
        return instance;
    }
}
