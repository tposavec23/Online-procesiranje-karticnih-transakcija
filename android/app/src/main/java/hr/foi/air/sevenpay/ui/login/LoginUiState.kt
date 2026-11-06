package hr.foi.air.sevenpay.ui.login

import hr.foi.air.sevenpay.core.model.User

data class LoginUiState(
    val username: String = "",
    val password: String = "",
    val isLoading: Boolean = false,
    val errorMessage: String? = null,
    val loggedInUser: User? = null
) {
    // Login request is sent only when both fields contain non-blank text.
    val canSubmit: Boolean
        get() = username.isNotBlank() && password.isNotBlank() && !isLoading
}