package hr.foi.air.sevenpay.core.model

data class User(
    val id: Int,
    val username: String,
    val role: Role,
    val active: Boolean
)
