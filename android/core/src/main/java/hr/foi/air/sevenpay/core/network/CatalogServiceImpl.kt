package hr.foi.air.sevenpay.core.service

import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.network.ApiClient
import hr.foi.air.sevenpay.core.network.CatalogApi

class CatalogServiceImpl(
    private val api: CatalogApi = ApiClient.catalogApi
) : CatalogService {

    override suspend fun getArticles(): List<Article> = api.getArticles()
}