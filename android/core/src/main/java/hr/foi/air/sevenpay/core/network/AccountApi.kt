package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.model.User
import retrofit2.http.Body
import retrofit2.http.POST

interface AccountApi {
    @POST("api/auth/login")
    suspend fun login(@Body request: LoginRequest): User
}