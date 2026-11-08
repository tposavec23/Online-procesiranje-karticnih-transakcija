package hr.foi.air.sevenpay.core.network

import hr.foi.air.sevenpay.core.exception.InvalidCredentialsException
import hr.foi.air.sevenpay.core.model.User
import hr.foi.air.sevenpay.core.service.AccountService
import retrofit2.HttpException

class AccountServiceImpl(
    private val api: AccountApi = ApiClient.accountApi
) : AccountService {

    override suspend fun login(username: String, password: String): User {
        try {
            return api.login(LoginRequest(username.trim(), password))
        } catch (e: HttpException) {
            if (e.code() == 401) throw InvalidCredentialsException()
            throw e
        }
    }
}