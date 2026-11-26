package hr.foi.air.sevenpay.core.data

import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import hr.foi.air.sevenpay.core.exception.OutOfStockException
import hr.foi.air.sevenpay.core.model.Cart
import hr.foi.air.sevenpay.core.network.AddToCartRequest
import hr.foi.air.sevenpay.core.network.CartApi
import retrofit2.HttpException

interface CartRepository {
    suspend fun addToCart(userId: Int, articleId: Int, quantity: Int): Cart
}

class NetworkCartRepository(
    private val api: CartApi
) : CartRepository {

    override suspend fun addToCart(userId: Int, articleId: Int, quantity: Int): Cart {
        try {
            return api.addItem(userId, AddToCartRequest(articleId, quantity))
        } catch (e: HttpException) {
            when (e.code()) {
                404 -> throw ArticleNotFoundException()
                409 -> throw OutOfStockException()
            }
            throw e
        }
    }
}