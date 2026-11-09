package hr.foi.air.sevenpay.ui.catalog

import hr.foi.air.sevenpay.core.model.Article

sealed interface ArticleDetailUiState {
    data object Loading : ArticleDetailUiState
    data class Success(val article: Article) : ArticleDetailUiState
}