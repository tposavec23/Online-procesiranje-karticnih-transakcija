package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.model.Article
import retrofit2.http.GET

interface CatalogApi {
    @GET("api/articles")
    suspend fun getArticles(): List<Article>
}