package hr.foi.air.sevenpay.ui.catalogue

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.ViewModelProvider.AndroidViewModelFactory.Companion.APPLICATION_KEY
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import hr.foi.air.sevenpay.SevenPayApplication
import hr.foi.air.sevenpay.core.data.CatalogueRepository
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class CatalogueViewModel(
    private val catalogueRepository: CatalogueRepository
) : ViewModel() {

    private val _state = MutableStateFlow<CatalogueUiState>(CatalogueUiState.Loading)
    val state: StateFlow<CatalogueUiState> = _state.asStateFlow()

    init {
        loadArticles()
    }

    fun loadArticles() {
        viewModelScope.launch {
            _state.value = CatalogueUiState.Loading
            _state.value = try {
                val articles = catalogueRepository.getArticles(PAGE_SIZE, 0)
                if (articles.isEmpty()) {
                    CatalogueUiState.Empty
                } else {
                    CatalogueUiState.Success(
                        articles = articles,
                        hasMore = articles.size == PAGE_SIZE
                    )
                }
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                CatalogueUiState.Error("Dohvat artikala nije uspio. Provjerite vezu i pokušajte ponovno.")
            }
        }
    }

    fun loadMore() {
        val current = _state.value as? CatalogueUiState.Success ?: return
        if (!current.hasMore || current.isLoadingMore) return

        _state.value = current.copy(isLoadingMore = true, loadMoreFailed = false)

        viewModelScope.launch {
            try {
                val next = catalogueRepository.getArticles(PAGE_SIZE, current.articles.size)
                _state.value = CatalogueUiState.Success(
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
                CatalogueViewModel(application.container.catalogueRepository)
            }
        }
    }
}