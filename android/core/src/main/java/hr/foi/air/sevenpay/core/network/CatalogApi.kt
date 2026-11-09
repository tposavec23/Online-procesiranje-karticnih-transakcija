package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.model.Article
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Query

interface CatalogApi {
    @GET("api/articles")
    suspend fun getArticles(
        @Query("limit") limit: Int,
        @Query("offset") offset: Int
    ): List<Article>

    @GET("api/articles/{id}")
    suspend fun getArticle(@Path("id") id: Int): Article
}