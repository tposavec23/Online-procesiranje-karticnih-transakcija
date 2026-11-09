package hr.foi.air.sevenpay.ui.catalogue.detail

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.viewmodel.compose.viewModel
import hr.foi.air.sevenpay.ui.theme.CreamBackground
import hr.foi.air.sevenpay.ui.theme.Raspberry

@Composable
fun ArticleDetailRoute(
    articleId: Int,
    onBack: () -> Unit,
    modifier: Modifier = Modifier
) {
    val viewModel: ArticleDetailViewModel = viewModel(
        key = "article-$articleId",
        factory = ArticleDetailViewModel.factory(articleId)
    )
    val state by viewModel.state.collectAsStateWithLifecycle()

    when (val current = state) {
        ArticleDetailUiState.Loading -> Box(
            modifier = modifier
                .fillMaxSize()
                .background(CreamBackground),
            contentAlignment = Alignment.Center
        ) {
            CircularProgressIndicator(color = Raspberry)
        }

        is ArticleDetailUiState.Success -> ArticleDetailScreen(
            article = current.article,
            onBack = onBack,
            modifier = modifier
        )

        ArticleDetailUiState.NotFound -> ArticleDetailMessage(
            message = "Artikl više ne postoji.",
            onBack = onBack,
            modifier = modifier
        )

        is ArticleDetailUiState.Error -> ArticleDetailMessage(
            message = current.message,
            onBack = onBack,
            modifier = modifier,
            onRetry = viewModel::load,
            isError = true
        )
    }
}