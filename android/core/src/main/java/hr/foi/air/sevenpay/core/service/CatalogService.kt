package hr.foi.air.sevenpay.core.service

import hr.foi.air.sevenpay.core.model.Article

interface CatalogService {
    suspend fun getArticles(): List<Article>
}