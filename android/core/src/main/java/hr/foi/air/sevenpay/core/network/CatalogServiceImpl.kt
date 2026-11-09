package hr.foi.air.sevenpay.core.service

import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.network.ApiClient
import hr.foi.air.sevenpay.core.network.CatalogApi
import retrofit2.HttpException

class CatalogServiceImpl(
    private val api: CatalogApi = ApiClient.catalogApi
) : CatalogService {

    override suspend fun getArticles(limit: Int, offset: Int): List<Article> =
        api.getArticles(limit, offset)

    override suspend fun getArticle(id: Int): Article {
        try {
            return api.getArticle(id)
        } catch (e: HttpException) {
            if (e.code() == 404) throw ArticleNotFoundException()
            throw e
        }
    }
}