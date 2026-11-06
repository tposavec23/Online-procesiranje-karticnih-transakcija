package hr.foi.air.sevenpay.core.model


data class Article(
    val id: Int,
    val code: String,
    val name: String,
    val description: String,
    val price: Long,
    val availableQuantity: Int
)