package hr.foi.air.sevenpay.core.model

data class User(
    val id: Int,
    val firstName: String,
    val lastName: String,
    val username: String,
    val role: Role,
    val active: Boolean
)
