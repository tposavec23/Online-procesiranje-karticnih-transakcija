package hr.foi.air.sevenpay.ui.admin.users

import hr.foi.air.sevenpay.core.model.User

sealed interface UserListUiState {
    data object Loading : UserListUiState
    data class Success(val users: List<User>) : UserListUiState
    data class Error(val message: String) : UserListUiState
}