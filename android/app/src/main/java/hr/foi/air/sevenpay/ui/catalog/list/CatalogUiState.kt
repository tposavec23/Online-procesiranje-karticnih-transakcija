package hr.foi.air.sevenpay.ui.catalog.list

import hr.foi.air.sevenpay.core.model.Article

sealed interface CatalogUiState {
    data object Loading : CatalogUiState
    data class Success(
        val articles: List<Article>,
        val hasMore: Boolean = false,
        val isLoadingMore: Boolean = false,
        val loadMoreFailed: Boolean = false
    ) : CatalogUiState
    data object Empty : CatalogUiState
    data class Error(val message: String) : CatalogUiState
}