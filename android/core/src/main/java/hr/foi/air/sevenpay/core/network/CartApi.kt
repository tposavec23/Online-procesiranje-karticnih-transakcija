package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.model.Cart
import retrofit2.http.Body
import retrofit2.http.Header
import retrofit2.http.POST

data class AddToCartRequest(
    val articleId: Int,
    val quantity: Int
)

interface CartApi{
    @POST("api/cart/items")
    suspend fun addItem(
        @Header("X-User-Id") userId : Int,
        @Body request: AddToCartRequest
    ) : Cart
}