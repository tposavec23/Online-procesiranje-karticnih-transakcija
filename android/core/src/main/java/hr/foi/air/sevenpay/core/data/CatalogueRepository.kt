package hr.foi.air.sevenpay.core.data

import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.network.CatalogueApi
import retrofit2.HttpException

interface CatalogueRepository {
    suspend fun getArticles(limit: Int, offset: Int): List<Article>
    suspend fun getArticle(id: Int): Article
}

class NetworkCatalogueRepository(
    private val api: CatalogueApi
) : CatalogueRepository {

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
