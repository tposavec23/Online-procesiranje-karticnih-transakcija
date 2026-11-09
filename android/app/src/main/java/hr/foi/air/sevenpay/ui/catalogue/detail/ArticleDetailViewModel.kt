package hr.foi.air.sevenpay.ui.catalogue.detail

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.ViewModelProvider.AndroidViewModelFactory.Companion.APPLICATION_KEY
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import hr.foi.air.sevenpay.SevenPayApplication
import hr.foi.air.sevenpay.core.data.CatalogueRepository
import hr.foi.air.sevenpay.core.exception.ArticleNotFoundException
import kotlinx.coroutines.CancellationException
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class ArticleDetailViewModel(
    private val articleId: Int,
    private val catalogueRepository: CatalogueRepository
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
                ArticleDetailUiState.Success(catalogueRepository.getArticle(articleId))
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

    companion object {
        fun factory(articleId: Int): ViewModelProvider.Factory = viewModelFactory {
            initializer {
                val application = this[APPLICATION_KEY] as SevenPayApplication
                ArticleDetailViewModel(articleId, application.container.catalogueRepository)
            }
        }
    }
}
