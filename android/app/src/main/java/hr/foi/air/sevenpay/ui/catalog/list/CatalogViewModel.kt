package hr.foi.air.sevenpay.ui.catalog.list

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import hr.foi.air.sevenpay.core.service.CatalogService
import hr.foi.air.sevenpay.core.service.CatalogServiceImpl
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class CatalogViewModel(
    private val catalogService: CatalogService = CatalogServiceImpl()
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
                val articles = catalogService.getArticles(PAGE_SIZE, 0)
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
                val next = catalogService.getArticles(PAGE_SIZE, current.articles.size)
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

    private companion object {
        const val PAGE_SIZE = 10
    }
}