package hr.foi.air.sevenpay.ui.catalog

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
                val articles = catalogService.getArticles()
                if (articles.isEmpty()) {
                    CatalogUiState.Empty
                } else {
                    CatalogUiState.Success(articles)
                }
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                CatalogUiState.Error("Dohvat artikala nije uspio. Provjerite vezu i pokušajte ponovno.")
            }
        }
    }
}