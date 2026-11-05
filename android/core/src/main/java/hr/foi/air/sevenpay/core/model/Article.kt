package hr.foi.air.sevenpay.core.model

import java.math.BigDecimal

data class Article(
    val id: Int,
    val code: String,
    val name: String,
    val description: String,
    val price: BigDecimal,
    val availableQuantity: Int
)