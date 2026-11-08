package hr.foi.air.sevenpay.ui.catalogue

import hr.foi.air.sevenpay.core.model.Article

sealed interface CatalogueUiState {
    data object Loading : CatalogueUiState
    data class Success(
        val articles: List<Article>,
        val hasMore: Boolean = false,
        val isLoadingMore: Boolean = false,
        val loadMoreFailed: Boolean = false
    ) : CatalogueUiState
    data object Empty : CatalogueUiState
    data class Error(val message: String) : CatalogueUiState
}