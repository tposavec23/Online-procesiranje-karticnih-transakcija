package hr.foi.air.sevenpay.ui.login

import hr.foi.air.sevenpay.core.exception.InvalidCredentialsException
import hr.foi.air.sevenpay.core.model.Role
import hr.foi.air.sevenpay.core.model.User
import hr.foi.air.sevenpay.core.service.AccountService
import kotlinx.coroutines.delay

private data class SampleAccount(val user: User, val password: String)

private val sampleAccounts = listOf(
    SampleAccount(User(1, "akovac", Role.MERCHANT, active = true), "trgovac123"),
    SampleAccount(User(2, "admin", Role.ADMIN, active = true), "admin123"),
    SampleAccount(User(3, "ihorvat", Role.MERCHANT, active = false), "trgovac123")
)

class SampleAccountService : AccountService {

    override suspend fun login(username: String, password: String): User {
        delay(800) // simulates waiting for the server response

        val account = sampleAccounts.find {
            it.user.username == username.trim() && it.password == password
        }

        // Deactivated accounts get the same error as wrong credentials,
        // so the response does not reveal whether the account exists.
        if (account == null || !account.user.active) {
            throw InvalidCredentialsException()
        }
        return account.user
    }
}