package hr.foi.air.sevenpay.ui.catalog.detail

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import hr.foi.air.sevenpay.core.service.CatalogService
import hr.foi.air.sevenpay.core.service.CatalogServiceImpl
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class ArticleDetailViewModel(
    private val articleId: Int,
    private val catalogService: CatalogService = CatalogServiceImpl()
) : ViewModel() {

    private val _state = MutableStateFlow<ArticleDetailUiState>(ArticleDetailUiState.Loading)
    val state: StateFlow<ArticleDetailUiState> = _state.asStateFlow()

    init {
        load()
    }

    fun load() {
        viewModelScope.launch {
            _state.value = ArticleDetailUiState.Loading
            _state.value = try {
                ArticleDetailUiState.Success(catalogService.getArticle(articleId))
            } catch (e: ArticleNotFoundException) {
                ArticleDetailUiState.NotFound
            } catch (e: CancellationException) {
                throw e
            } catch (e: Exception) {
                ArticleDetailUiState.Error(
                    "Dohvat detalja artikla nije uspio. Provjerite vezu i pokušajte ponovno."
                )
            }
        }
    }
}