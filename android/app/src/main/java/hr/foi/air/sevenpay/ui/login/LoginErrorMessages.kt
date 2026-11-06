package hr.foi.air.sevenpay.ui.login

import hr.foi.air.sevenpay.core.exception.InvalidCredentialsException
import java.io.IOException

// User-facing messages shown on the login screen, one for each type of failure.
object LoginErrorMessages {
    const val INVALID_CREDENTIALS = "Pogrešno korisničko ime ili lozinka. Pokušaj ponovno."
    const val NO_CONNECTION = "Prijava trenutno nije moguća. Provjeri internetsku vezu i pokušaj ponovno."
    const val UNKNOWN = "Došlo je do pogreške pri prijavi. Pokušaj ponovno."

    // Maps an error thrown by AccountService to the message the user should see.
    fun from(error: Throwable): String = when (error) {
        is InvalidCredentialsException -> INVALID_CREDENTIALS
        is IOException -> NO_CONNECTION
        else -> UNKNOWN
    }
}