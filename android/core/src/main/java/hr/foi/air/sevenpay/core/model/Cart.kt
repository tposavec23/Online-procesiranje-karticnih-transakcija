package hr.foi.air.sevenpay.core.model

data class CartItem(
    val articleId: Int,
    val code: String,
    val name: String,
    val quatity: Int,
    val unitPrice: Long,
    val totalPrice: Long
)

data class Cart(
    val items: List<CartItem>,
    val totalPrice: Long
)


