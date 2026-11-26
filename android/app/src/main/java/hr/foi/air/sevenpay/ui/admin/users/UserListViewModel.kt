package hr.foi.air.sevenpay.ui.admin.users

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.ViewModelProvider.AndroidViewModelFactory.Companion.APPLICATION_KEY
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import hr.foi.air.sevenpay.SevenPayApplication
import hr.foi.air.sevenpay.core.data.AccountRepository
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class UserListViewModel(
    private val accountRepository: AccountRepository
) : ViewModel() {

    private val _state = MutableStateFlow<UserListUiState>(UserListUiState.Loading)
    val state: StateFlow<UserListUiState> = _state.asStateFlow()

    init {
        loadUsers()
    }

    fun loadUsers() {
        viewModelScope.launch {
            _state.value = UserListUiState.Loading
            _state.value = try {
                val users = accountRepository.getUsers()
                if (users.isEmpty()) UserListUiState.Empty else UserListUiState.Success(users)
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                UserListUiState.Error("Dohvat korisničkih računa nije uspio. Provjerite vezu i pokušajte ponovno.")
            }
        }
    }

    companion object {
        val Factory: ViewModelProvider.Factory = viewModelFactory {
            initializer {
                val application = this[APPLICATION_KEY] as SevenPayApplication
                UserListViewModel(application.container.accountRepository)
            }
        }
    }
}