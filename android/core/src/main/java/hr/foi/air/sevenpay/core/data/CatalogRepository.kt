package hr.foi.air.sevenpay.core.data

import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.network.CatalogApi

interface CatalogRepository {
    suspend fun getArticles(limit: Int, offset: Int): List<Article>
}

class NetworkCatalogRepository(
    private val api: CatalogApi
) : CatalogRepository {

    override suspend fun getArticles(limit: Int, offset: Int): List<Article> =
        api.getArticles(limit, offset)
}
