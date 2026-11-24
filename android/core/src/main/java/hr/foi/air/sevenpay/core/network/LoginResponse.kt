package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.model.User

data class LoginResponse(
    val token: String,
    val user: User
)