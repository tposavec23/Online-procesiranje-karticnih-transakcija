package hr.foi.air.sevenpay.ui.catalogue.detail

import hr.foi.air.sevenpay.core.model.Article

sealed interface ArticleDetailUiState {
    data object Loading : ArticleDetailUiState
    data class Success(val article: Article) : ArticleDetailUiState
    data object NotFound : ArticleDetailUiState
    data class Error(val message: String) : ArticleDetailUiState
}