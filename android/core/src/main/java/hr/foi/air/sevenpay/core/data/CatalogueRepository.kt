package hr.foi.air.sevenpay.core.data

import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.network.CatalogueApi

interface CatalogueRepository {
    suspend fun getArticles(limit: Int, offset: Int): List<Article>
}

class NetworkCatalogueRepository(
    private val api: CatalogueApi
) : CatalogueRepository {

    override suspend fun getArticles(limit: Int, offset: Int): List<Article> =
        api.getArticles(limit, offset)
}
