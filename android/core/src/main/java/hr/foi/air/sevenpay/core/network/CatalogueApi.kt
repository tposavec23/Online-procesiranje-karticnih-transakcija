package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.model.ImportResult
import okhttp3.MultipartBody
import okhttp3.ResponseBody
import retrofit2.http.GET
import retrofit2.http.Multipart
import retrofit2.http.POST
import retrofit2.http.Part
import retrofit2.http.Path
import retrofit2.http.Query

interface CatalogueApi {
    @GET("api/articles")
    suspend fun getArticles(
        @Query("limit") limit: Int,
        @Query("offset") offset: Int
    ): List<Article>

    @GET("api/articles/{id}")
    suspend fun getArticle(@Path("id") id: Int): Article

    @GET("api/articles/export")
    suspend fun exportArticles(): ResponseBody

    @Multipart
    @POST("api/articles/import")
    suspend fun importArticles(@Part file: MultipartBody.Part): ImportResult
}
