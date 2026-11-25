package hr.foi.air.sevenpay.core.data

import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import hr.foi.air.sevenpay.core.model.Article
import hr.foi.air.sevenpay.core.model.ImportResult
import hr.foi.air.sevenpay.core.network.CatalogueApi
import okhttp3.MediaType
import okhttp3.MultipartBody
import okhttp3.RequestBody
import retrofit2.HttpException

interface CatalogueRepository {
    suspend fun getArticles(limit: Int, offset: Int): List<Article>
    suspend fun getArticle(id: Int): Article
    suspend fun exportArticles(): String
    suspend fun importArticles(csv: ByteArray): ImportResult
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

    override suspend fun exportArticles(): String =
        api.exportArticles().string()

    override suspend fun importArticles(csv: ByteArray): ImportResult {
        val file = MultipartBody.Part.createFormData("file", "artikli.csv", RequestBody.create(MediaType.parse("text/csv"), csv))
        return api.importArticles(file)
    }
}
