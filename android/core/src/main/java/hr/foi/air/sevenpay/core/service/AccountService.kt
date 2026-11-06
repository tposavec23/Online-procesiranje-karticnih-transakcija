package hr.foi.air.sevenpay.core.service

import hr.foi.air.sevenpay.core.model.User

interface AccountService {
    suspend fun login(username: String, password: String): User
}