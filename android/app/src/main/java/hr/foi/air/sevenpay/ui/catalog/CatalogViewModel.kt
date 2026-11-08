package hr.foi.air.sevenpay.ui.catalog

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.ViewModelProvider.AndroidViewModelFactory.Companion.APPLICATION_KEY
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import hr.foi.air.sevenpay.SevenPayApplication
import hr.foi.air.sevenpay.core.data.CatalogRepository
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class CatalogViewModel(
    private val catalogRepository: CatalogRepository
) : ViewModel() {

    private val _state = MutableStateFlow<CatalogUiState>(CatalogUiState.Loading)
    val state: StateFlow<CatalogUiState> = _state.asStateFlow()

    init {
        loadArticles()
    }

    fun loadArticles() {
        viewModelScope.launch {
            _state.value = CatalogUiState.Loading
            _state.value = try {
                val articles = catalogRepository.getArticles(PAGE_SIZE, 0)
                if (articles.isEmpty()) {
                    CatalogUiState.Empty
                } else {
                    CatalogUiState.Success(
                        articles = articles,
                        hasMore = articles.size == PAGE_SIZE
                    )
                }
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                CatalogUiState.Error("Dohvat artikala nije uspio. Provjerite vezu i pokušajte ponovno.")
            }
        }
    }

    fun loadMore() {
        val current = _state.value as? CatalogUiState.Success ?: return
        if (!current.hasMore || current.isLoadingMore) return

        _state.value = current.copy(isLoadingMore = true, loadMoreFailed = false)

        viewModelScope.launch {
            try {
                val next = catalogRepository.getArticles(PAGE_SIZE, current.articles.size)
                _state.value = CatalogUiState.Success(
                    articles = current.articles + next,
                    hasMore = next.size == PAGE_SIZE
                )
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                _state.value = current.copy(isLoadingMore = false, loadMoreFailed = true)
            }
        }
    }

    companion object {
        private const val PAGE_SIZE = 10

        val Factory: ViewModelProvider.Factory = viewModelFactory {
            initializer {
                val application = this[APPLICATION_KEY] as SevenPayApplication
                CatalogViewModel(application.container.catalogRepository)
            }
        }
    }
}