package hr.foi.air.sevenpay.core.network

data class LoginRequest(
    val username: String,
    val password: String
)