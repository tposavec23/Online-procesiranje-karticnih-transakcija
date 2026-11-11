package hr.foi.air.sevenpay.core.data

import hr.foi.air.sevenpay.core.exception.InvalidCredentialsException
import hr.foi.air.sevenpay.core.model.User
import hr.foi.air.sevenpay.core.network.AccountApi
import hr.foi.air.sevenpay.core.network.LoginRequest
import retrofit2.HttpException

interface AccountRepository {
    suspend fun login(username: String, password: String): User
}

class NetworkAccountRepository(
    private val api: AccountApi
) : AccountRepository {

    override suspend fun login(username: String, password: String): User {
        try {
            return api.login(LoginRequest(username.trim(), password))
        } catch (e: HttpException) {
            if (e.code() == 401) throw InvalidCredentialsException()
            throw e
        }
    }
}
